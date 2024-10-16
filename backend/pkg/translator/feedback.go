package translator

import (
	"context"
	"encoding/json"
	"time"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/organization"
	"backend/pkg/processor"
	"backend/pkg/product"
	"backend/pkg/util"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
)

const (
	FeedbackTranslatorTranslate = "translator:translate-feedback"
	FeedbackTranslatorSchedule  = "translator:schedule-translate-feedback"
)

type FeedbackTranslator struct {
	config                 config.Config
	observer               *kit.Observer
	feedbackRepository     *feedback.FeedbackRepository
	productRepository      *product.ProductRepository
	organizationRepository organization.OrganizationRepository
	enqueuer               *kit.Enqueuer
	engineService          *engine.EngineService
	engineBreaker          *engine.EngineBreaker
}

func NewFeedbackTranslator(observer *kit.Observer, feedbackRepository *feedback.FeedbackRepository,
	productRepository *product.ProductRepository, organizationRepository organization.OrganizationRepository,
	enqueuer *kit.Enqueuer, engineService *engine.EngineService, engineBreaker *engine.EngineBreaker,
	config config.Config) *FeedbackTranslator {
	return &FeedbackTranslator{
		config:                 config,
		observer:               observer,
		feedbackRepository:     feedbackRepository,
		productRepository:      productRepository,
		organizationRepository: organizationRepository,
		enqueuer:               enqueuer,
		engineService:          engineService,
		engineBreaker:          engineBreaker,
	}
}

type FeedbackTranslatorTranslateParams struct {
	FeedbackID string
}

func (self *FeedbackTranslator) Translate(ctx context.Context, task *asynq.Task) error {
	if self.engineBreaker.IsOpen(ctx) {
		return nil
	}

	params := FeedbackTranslatorTranslateParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	feedback, err := self.feedbackRepository.GetByID(ctx, params.FeedbackID)
	if err != nil {
		return err
	}

	if feedback == nil {
		return nil
	}

	if feedback.TranslatedAt != nil {
		return nil
	}

	product, err := self.productRepository.GetByID(ctx, feedback.ProductID)
	if err != nil {
		return err
	}

	if product == nil {
		return nil
	}

	if product.DeletedAt != nil {
		return nil
	}

	organization, err := self.organizationRepository.GetByID(ctx, product.OrganizationID)
	if err != nil {
		return err
	}

	if organization == nil {
		return nil
	}

	if organization.DeletedAt != nil {
		return nil
	}

	if organization.UsageLeft() < 1 {
		return nil
	}

	if len(feedback.Content) == 0 {
		return nil
	}

	result, err := self.engineService.DetectLanguage(ctx, engine.EngineServiceDetectLanguageParams{
		Feedback: engine.Feedback{
			Content: feedback.Content,
		},
	})
	if err != nil {
		if engine.ErrEngineServiceTimedOut.Is(err) {
			err := self.engineBreaker.Open(ctx)
			if err != nil {
				self.observer.Error(ctx, err)
			}
		}

		return err
	}

	feedback.Language = result.Language
	tokens := (result.Usage.Input + result.Usage.Output)

	if feedback.Language != product.Language {
		result, err := self.engineService.TranslateFeedback(ctx, engine.EngineServiceTranslateFeedbackParams{
			Feedback: engine.Feedback{
				Content: feedback.Content,
			},
			FromLanguage: feedback.Language,
			ToLanguage:   product.Language,
		})
		if err != nil {
			if engine.ErrEngineServiceTimedOut.Is(err) {
				err := self.engineBreaker.Open(ctx)
				if err != nil {
					self.observer.Error(ctx, err)
				}
			}

			return err
		}

		feedback.Translation = result.Translation
		tokens += (result.Usage.Input + result.Usage.Output)
	}

	feedback.Tokens += tokens
	feedback.TranslatedAt = kitUtil.Pointer(time.Now())

	err = self.feedbackRepository.UpdateTranslated(ctx, *feedback)
	if err != nil {
		return err
	}

	err = self.enqueuer.Enqueue(ctx, processor.FeedbackProcessorProcess, processor.FeedbackProcessorProcessParams{
		FeedbackID: feedback.ID,
	}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
	if err != nil {
		self.observer.Error(ctx, err)
	}

	self.observer.Infof(ctx, "Translated a feedback from %s to %s using %d tokens",
		feedback.Language, product.Language, tokens)

	return nil
}

func (self *FeedbackTranslator) Schedule(ctx context.Context, _ *asynq.Task) error {
	pagination := util.Pagination[time.Time]{
		Limit: 1000,
		From:  nil,
	}

	for {
		page, err := self.feedbackRepository.ListIDsByNotTranslated(ctx, pagination)
		if err != nil {
			return err
		}

		for _, id := range page.Items {
			err := self.enqueuer.Enqueue(ctx, FeedbackTranslatorTranslate, FeedbackTranslatorTranslateParams{
				FeedbackID: id,
			}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
			if err != nil {
				self.observer.Error(ctx, err)
			}
		}

		if page.Next == nil {
			break
		}
		pagination.From = page.Next
	}

	return nil
}
