package processor

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"backend/pkg/aggregator"
	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/issue"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/review"
	"backend/pkg/suggestion"
	"backend/pkg/util"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
)

const (
	FeedbackProcessorProcess  = "processor:process-feedback"
	FeedbackProcessorSchedule = "processor:schedule-process-feedback"
)

type FeedbackProcessor struct {
	config                      config.Config
	observer                    *kit.Observer
	database                    *kit.Database
	feedbackRepository          *feedback.FeedbackRepository
	partialIssueRepository      *issue.PartialIssueRepository
	partialSuggestionRepository *suggestion.PartialSuggestionRepository
	reviewRepository            *review.ReviewRepository
	productRepository           *product.ProductRepository
	organizationRepository      organization.OrganizationRepository
	enqueuer                    *kit.Enqueuer
	engineService               *engine.EngineService
	engineBreaker               *engine.EngineBreaker
}

func NewFeedbackProcessor(observer *kit.Observer, database *kit.Database,
	feedbackRepository *feedback.FeedbackRepository, partialIssueRepository *issue.PartialIssueRepository,
	partialSuggestionRepository *suggestion.PartialSuggestionRepository,
	reviewRepository *review.ReviewRepository, productRepository *product.ProductRepository,
	organizationRepository organization.OrganizationRepository, enqueuer *kit.Enqueuer,
	engineService *engine.EngineService, engineBreaker *engine.EngineBreaker,
	config config.Config) *FeedbackProcessor {
	return &FeedbackProcessor{
		config:                      config,
		observer:                    observer,
		database:                    database,
		feedbackRepository:          feedbackRepository,
		partialIssueRepository:      partialIssueRepository,
		partialSuggestionRepository: partialSuggestionRepository,
		reviewRepository:            reviewRepository,
		productRepository:           productRepository,
		organizationRepository:      organizationRepository,
		enqueuer:                    enqueuer,
		engineService:               engineService,
		engineBreaker:               engineBreaker,
	}
}

type FeedbackProcessorProcessParams struct {
	FeedbackID string
}

func (self *FeedbackProcessor) Process(ctx context.Context, task *asynq.Task) error {
	if self.engineBreaker.IsOpen(ctx) {
		return nil
	}

	params := FeedbackProcessorProcessParams{}

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

	if feedback.ProcessedAt != nil {
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

	content := feedback.Content
	if feedback.Language != product.Language {
		content = feedback.Translation
	}

	if len(content) == 0 {
		return nil
	}

	var group sync.WaitGroup

	var eiResult *engine.EngineServiceExtractIssuesResult
	var eiErr error
	group.Add(1)
	go func() {
		defer group.Done()
		eiResult, eiErr = self.engineService.ExtractIssues(ctx, engine.EngineServiceExtractIssuesParams{
			Context:    product.Context,
			Categories: product.Categories,
			Feedback: engine.Feedback{
				Content: content,
			},
		})
		if eiErr != nil && engine.ErrEngineServiceTimedOut.Is(eiErr) {
			err := self.engineBreaker.Open(ctx)
			if err != nil {
				self.observer.Error(ctx, err)
			}
		}
	}()

	var esResult *engine.EngineServiceExtractSuggestionsResult
	var esErr error
	group.Add(1)
	go func() {
		defer group.Done()
		esResult, esErr = self.engineService.ExtractSuggestions(ctx, engine.EngineServiceExtractSuggestionsParams{
			Context:    product.Context,
			Categories: product.Categories,
			Feedback: engine.Feedback{
				Content: content,
			},
		})
		if esErr != nil && engine.ErrEngineServiceTimedOut.Is(esErr) {
			err := self.engineBreaker.Open(ctx)
			if err != nil {
				self.observer.Error(ctx, err)
			}
		}
	}()

	var erResult *engine.EngineServiceExtractReviewResult
	var erErr error
	group.Add(1)
	go func() {
		defer group.Done()
		erResult, erErr = self.engineService.ExtractReview(ctx, engine.EngineServiceExtractReviewParams{
			Context:    product.Context,
			Categories: product.Categories,
			Feedback: engine.Feedback{
				Content: content,
			},
		})
		if erErr != nil && engine.ErrEngineServiceTimedOut.Is(erErr) {
			err := self.engineBreaker.Open(ctx)
			if err != nil {
				self.observer.Error(ctx, err)
			}
		}
	}()

	group.Wait()

	if eiErr != nil {
		return eiErr
	} else if esErr != nil {
		return esErr
	} else if erErr != nil {
		return erErr
	}

	now := time.Now()

	var issues []issue.PartialIssue
	for _, resIssue := range eiResult.Issues {
		issue := issue.NewPartialIssue()
		issue.ID = xid.New().String()
		issue.FeedbackID = feedback.ID
		issue.Title = resIssue.Title
		issue.Description = resIssue.Description
		issue.Steps = resIssue.Steps
		issue.Severity = resIssue.Severity
		issue.Category = resIssue.Category
		issue.CreatedAt = now

		issues = append(issues, *issue)
	}

	var suggestions []suggestion.PartialSuggestion
	for _, resSuggestion := range esResult.Suggestions {
		suggestion := suggestion.NewPartialSuggestion()
		suggestion.ID = xid.New().String()
		suggestion.FeedbackID = feedback.ID
		suggestion.Title = resSuggestion.Title
		suggestion.Description = resSuggestion.Description
		suggestion.Reason = resSuggestion.Reason
		suggestion.Importance = resSuggestion.Importance
		suggestion.Category = resSuggestion.Category
		suggestion.CreatedAt = now

		suggestions = append(suggestions, *suggestion)
	}

	review := review.NewReview()
	review.ID = xid.New().String()
	review.ProductID = product.ID
	review.Feedback = *feedback
	review.Keywords = erResult.Review.Keywords
	review.Sentiment = erResult.Review.Sentiment
	review.Emotions = erResult.Review.Emotions
	review.Intention = erResult.Review.Intention
	review.Category = erResult.Review.Category
	review.Quality = nil
	review.CreatedAt = now
	review.ExportedAt = nil

	tokens := eiResult.Usage.Input + eiResult.Usage.Output +
		esResult.Usage.Input + esResult.Usage.Output +
		erResult.Usage.Input + erResult.Usage.Output

	feedback.Tokens += tokens
	feedback.ProcessedAt = kitUtil.Pointer(time.Now())

	err = self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		if len(issues) > 0 {
			err := self.partialIssueRepository.BulkCreate(ctx, issues)
			if err != nil {
				return err
			}
		}

		if len(suggestions) > 0 {
			err := self.partialSuggestionRepository.BulkCreate(ctx, suggestions)
			if err != nil {
				return err
			}
		}

		review, err = self.reviewRepository.Create(ctx, *review)
		if err != nil {
			return err
		}

		err = self.feedbackRepository.UpdateProcessed(ctx, *feedback)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, partial := range issues {
		err := self.enqueuer.Enqueue(ctx, aggregator.IssueAggregatorAggregate,
			aggregator.IssueAggregatorAggregateParams{
				PartialID: partial.ID,
			}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	for _, partial := range suggestions {
		err := self.enqueuer.Enqueue(ctx, aggregator.SuggestionAggregatorAggregate,
			aggregator.SuggestionAggregatorAggregateParams{
				PartialID: partial.ID,
			}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	self.observer.Infof(ctx, "Processed a feedback with %d issues, %d suggestions and 1 review using %d tokens",
		len(issues), len(suggestions), tokens)

	return nil
}

func (self *FeedbackProcessor) Schedule(ctx context.Context, _ *asynq.Task) error {
	pagination := util.Pagination[time.Time]{
		Limit: 1000,
		From:  nil,
	}

	for {
		page, err := self.feedbackRepository.ListIDsByNotProcessed(ctx, pagination)
		if err != nil {
			return err
		}

		for _, id := range page.Items {
			err := self.enqueuer.Enqueue(ctx, FeedbackProcessorProcess, FeedbackProcessorProcessParams{
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
