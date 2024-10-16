package aggregator

import (
	"context"
	"encoding/json"
	"time"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/suggestion"
	"backend/pkg/util"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
)

const (
	SUGGESTION_AGGREGATOR_MAX_SIMILAR_SUGGESTIONS = 10
)

const (
	SuggestionAggregatorAggregate = "aggregator:aggregate-suggestion"
	SuggestionAggregatorSchedule  = "aggregator:schedule-aggregate-suggestion"
)

type SuggestionAggregator struct {
	config                      config.Config
	observer                    *kit.Observer
	database                    *kit.Database
	partialSuggestionRepository *suggestion.PartialSuggestionRepository
	suggestionRepository        *suggestion.SuggestionRepository
	feedbackRepository          *feedback.FeedbackRepository
	productRepository           *product.ProductRepository
	organizationRepository      organization.OrganizationRepository
	enqueuer                    *kit.Enqueuer
	engineService               *engine.EngineService
	engineBreaker               *engine.EngineBreaker
}

func NewSuggestionAggregator(observer *kit.Observer, database *kit.Database,
	partialSuggestionRepository *suggestion.PartialSuggestionRepository, suggestionRepository *suggestion.SuggestionRepository,
	feedbackRepository *feedback.FeedbackRepository, productRepository *product.ProductRepository,
	organizationRepository organization.OrganizationRepository, enqueuer *kit.Enqueuer,
	engineService *engine.EngineService, engineBreaker *engine.EngineBreaker,
	config config.Config) *SuggestionAggregator {
	return &SuggestionAggregator{
		config:                      config,
		observer:                    observer,
		database:                    database,
		partialSuggestionRepository: partialSuggestionRepository,
		suggestionRepository:        suggestionRepository,
		feedbackRepository:          feedbackRepository,
		productRepository:           productRepository,
		organizationRepository:      organizationRepository,
		enqueuer:                    enqueuer,
		engineService:               engineService,
		engineBreaker:               engineBreaker,
	}
}

func (self *SuggestionAggregator) createSuggestion(ctx context.Context, embedding []float32, partial *suggestion.PartialSuggestion,
	feedback *feedback.Feedback, product *product.Product, tokens int) error {
	_suggestion := suggestion.NewSuggestion()
	_suggestion.ID = xid.New().String()
	_suggestion.ProductID = product.ID
	_suggestion.Embedding = embedding
	_suggestion.Sources = map[string]int{feedback.Source: 1}
	_suggestion.Title = partial.Title
	_suggestion.Description = partial.Description
	_suggestion.Reason = partial.Reason
	_suggestion.Importances = map[string]int{partial.Importance: 1}
	_suggestion.Priority = suggestion.ComputePriority(_suggestion.Importances, 1)
	_suggestion.Categories = map[string]int{partial.Category: 1}
	_suggestion.Releases = map[string]int{feedback.Release: 1}
	_suggestion.Customers = 1
	_suggestion.AssigneeID = nil
	_suggestion.Quality = nil
	_suggestion.FirstSeenAt = feedback.PostedAt
	_suggestion.LastSeenAt = feedback.PostedAt
	_suggestion.CreatedAt = time.Now()
	_suggestion.ArchivedAt = nil
	_suggestion.LastAggregatedAt = nil
	_suggestion.ExportedAt = nil

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		err := self.partialSuggestionRepository.Delete(ctx, partial.ID)
		if err != nil {
			return err
		}

		_suggestion, err = self.suggestionRepository.Create(ctx, *_suggestion, *feedback)
		if err != nil {
			return err
		}

		feedback, err = self.feedbackRepository.GetByIDForUpdate(ctx, feedback.ID)
		if err != nil {
			return err
		}

		feedback.Tokens += tokens
		err = self.feedbackRepository.UpdateTokens(ctx, feedback.ID, feedback.Tokens)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	self.observer.Infof(ctx, "Created a new suggestion using %d tokens", tokens)

	return nil
}

func (self *SuggestionAggregator) mergeSuggestions(ctx context.Context, partial *suggestion.PartialSuggestion,
	_suggestion *suggestion.Suggestion, feedback *feedback.Feedback, product *product.Product, tokens int) error {
	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		_suggestion, err := self.suggestionRepository.GetByIDForUpdate(ctx, _suggestion.ID)
		if err != nil {
			return err
		}

		msResult, err := self.engineService.MergeSuggestions(ctx, engine.EngineServiceMergeSuggestionsParams{
			SuggestionA: engine.Suggestion{
				Title:       partial.Title,
				Description: partial.Description,
				Reason:      partial.Reason,
			},
			SuggestionB: engine.Suggestion{
				Title:       _suggestion.Title,
				Description: _suggestion.Description,
				Reason:      _suggestion.Reason,
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

		_suggestion.Title = msResult.Suggestion.Title
		_suggestion.Description = msResult.Suggestion.Description
		_suggestion.Reason = msResult.Suggestion.Reason
		tokens += (msResult.Usage.Input + msResult.Usage.Output)

		ceResult, err := self.engineService.ComputeEmbedding(ctx, engine.EngineServiceComputeEmbeddingParams{
			Text: _suggestion.Description,
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

		_suggestion.Embedding = ceResult.Embedding
		tokens += (ceResult.Usage.Input + ceResult.Usage.Output)

		_suggestion.Sources[feedback.Source]++
		_suggestion.Importances[partial.Importance]++
		_suggestion.Priority = suggestion.ComputePriority(_suggestion.Importances, _suggestion.Customers+1)
		_suggestion.Categories[partial.Category]++
		_suggestion.Releases[feedback.Release]++
		_suggestion.Customers++
		if feedback.PostedAt.Before(_suggestion.FirstSeenAt) {
			_suggestion.FirstSeenAt = feedback.PostedAt
		}
		if feedback.PostedAt.After(_suggestion.LastSeenAt) {
			_suggestion.LastSeenAt = feedback.PostedAt
		}
		_suggestion.LastAggregatedAt = kitUtil.Pointer(time.Now())

		err = self.suggestionRepository.UpdateAggregated(ctx, *_suggestion, *feedback)
		if err != nil {
			return err
		}

		err = self.partialSuggestionRepository.Delete(ctx, partial.ID)
		if err != nil {
			return err
		}

		feedback, err = self.feedbackRepository.GetByIDForUpdate(ctx, feedback.ID)
		if err != nil {
			return err
		}

		feedback.Tokens += tokens
		err = self.feedbackRepository.UpdateTokens(ctx, feedback.ID, feedback.Tokens)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if kit.ErrDatabaseIntegrityViolation.In(err) {
			// There are two or more suggestions from the same feedback trying to merge themselves.
			// This should not be possible, but protect from retry-loop by deleting the partial.
			self.observer.Error(ctx, err)

			err = self.partialSuggestionRepository.Delete(ctx, partial.ID)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	self.observer.Infof(ctx, "Merged 2 suggestions using %d tokens", tokens)

	return nil
}

type SuggestionAggregatorAggregateParams struct {
	PartialID string
}

func (self *SuggestionAggregator) Aggregate(ctx context.Context, task *asynq.Task) error {
	if self.engineBreaker.IsOpen(ctx) {
		return nil
	}

	params := SuggestionAggregatorAggregateParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	partial, err := self.partialSuggestionRepository.GetByID(ctx, params.PartialID)
	if err != nil {
		return err
	}

	if partial == nil {
		return nil
	}

	feedback, err := self.feedbackRepository.GetByID(ctx, partial.FeedbackID)
	if err != nil {
		return err
	}

	if feedback == nil {
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

	if len(partial.Description) == 0 {
		return nil
	}

	ceResult, err := self.engineService.ComputeEmbedding(ctx, engine.EngineServiceComputeEmbeddingParams{
		Text: partial.Description,
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

	tokens := (ceResult.Usage.Input + ceResult.Usage.Output)

	suggestions, err := self.suggestionRepository.ListByEmbeddingAndProductID(ctx, ceResult.Embedding,
		suggestion.SUGGESTION_SIMILAR_THRESHOLD, SUGGESTION_AGGREGATOR_MAX_SIMILAR_SUGGESTIONS, product.ID)
	if err != nil {
		return err
	}

	if len(suggestions) == 0 {
		return self.createSuggestion(ctx, ceResult.Embedding, partial, feedback, product, tokens)
	}

	options := make([]engine.Suggestion, 0, len(suggestions))
	for _, suggestion := range suggestions {
		options = append(options, engine.Suggestion{
			Title:       suggestion.Title,
			Description: suggestion.Description,
			Reason:      suggestion.Reason,
		})
	}

	ssResult, err := self.engineService.SimilarSuggestion(ctx, engine.EngineServiceSimilarSuggestionParams{
		Suggestion: engine.Suggestion{
			Title:       partial.Title,
			Description: partial.Description,
			Reason:      partial.Reason,
		},
		Options: options,
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

	tokens += (ssResult.Usage.Input + ssResult.Usage.Output)

	if ssResult.Option == nil {
		return self.createSuggestion(ctx, ceResult.Embedding, partial, feedback, product, tokens)
	}

	return self.mergeSuggestions(ctx, partial, &suggestions[*ssResult.Option], feedback, product, tokens)
}

func (self *SuggestionAggregator) Schedule(ctx context.Context, _ *asynq.Task) error {
	pagination := util.Pagination[time.Time]{
		Limit: 1000,
		From:  nil,
	}

	for {
		page, err := self.partialSuggestionRepository.ListIDsByCreatedAt(ctx, pagination)
		if err != nil {
			return err
		}

		for _, id := range page.Items {
			err := self.enqueuer.Enqueue(ctx, SuggestionAggregatorAggregate, SuggestionAggregatorAggregateParams{
				PartialID: id,
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
