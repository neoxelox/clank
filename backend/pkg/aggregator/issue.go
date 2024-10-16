package aggregator

import (
	"context"
	"encoding/json"
	"time"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/issue"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/util"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
)

const (
	ISSUE_AGGREGATOR_MAX_SIMILAR_ISSUES = 10
)

const (
	IssueAggregatorAggregate = "aggregator:aggregate-issue"
	IssueAggregatorSchedule  = "aggregator:schedule-aggregate-issue"
)

type IssueAggregator struct {
	config                 config.Config
	observer               *kit.Observer
	database               *kit.Database
	partialIssueRepository *issue.PartialIssueRepository
	issueRepository        *issue.IssueRepository
	feedbackRepository     *feedback.FeedbackRepository
	productRepository      *product.ProductRepository
	organizationRepository organization.OrganizationRepository
	enqueuer               *kit.Enqueuer
	engineService          *engine.EngineService
	engineBreaker          *engine.EngineBreaker
}

func NewIssueAggregator(observer *kit.Observer, database *kit.Database,
	partialIssueRepository *issue.PartialIssueRepository, issueRepository *issue.IssueRepository,
	feedbackRepository *feedback.FeedbackRepository, productRepository *product.ProductRepository,
	organizationRepository organization.OrganizationRepository, enqueuer *kit.Enqueuer,
	engineService *engine.EngineService, engineBreaker *engine.EngineBreaker,
	config config.Config) *IssueAggregator {
	return &IssueAggregator{
		config:                 config,
		observer:               observer,
		database:               database,
		partialIssueRepository: partialIssueRepository,
		issueRepository:        issueRepository,
		feedbackRepository:     feedbackRepository,
		productRepository:      productRepository,
		organizationRepository: organizationRepository,
		enqueuer:               enqueuer,
		engineService:          engineService,
		engineBreaker:          engineBreaker,
	}
}

func (self *IssueAggregator) createIssue(ctx context.Context, embedding []float32, partial *issue.PartialIssue,
	feedback *feedback.Feedback, product *product.Product, tokens int) error {
	_issue := issue.NewIssue()
	_issue.ID = xid.New().String()
	_issue.ProductID = product.ID
	_issue.Embedding = embedding
	_issue.Sources = map[string]int{feedback.Source: 1}
	_issue.Title = partial.Title
	_issue.Description = partial.Description
	_issue.Steps = partial.Steps
	_issue.Severities = map[string]int{partial.Severity: 1}
	_issue.Priority = issue.ComputePriority(_issue.Severities, 1)
	_issue.Categories = map[string]int{partial.Category: 1}
	_issue.Releases = map[string]int{feedback.Release: 1}
	_issue.Customers = 1
	_issue.AssigneeID = nil
	_issue.Quality = nil
	_issue.FirstSeenAt = feedback.PostedAt
	_issue.LastSeenAt = feedback.PostedAt
	_issue.CreatedAt = time.Now()
	_issue.ArchivedAt = nil
	_issue.LastAggregatedAt = nil
	_issue.ExportedAt = nil

	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		err := self.partialIssueRepository.Delete(ctx, partial.ID)
		if err != nil {
			return err
		}

		_issue, err = self.issueRepository.Create(ctx, *_issue, *feedback)
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

	self.observer.Infof(ctx, "Created a new issue using %d tokens", tokens)

	return nil
}

func (self *IssueAggregator) mergeIssues(ctx context.Context, partial *issue.PartialIssue,
	_issue *issue.Issue, feedback *feedback.Feedback, product *product.Product, tokens int) error {
	err := self.database.Transaction(ctx, nil, func(ctx context.Context) error {
		_issue, err := self.issueRepository.GetByIDForUpdate(ctx, _issue.ID)
		if err != nil {
			return err
		}

		miResult, err := self.engineService.MergeIssues(ctx, engine.EngineServiceMergeIssuesParams{
			IssueA: engine.Issue{
				Title:       partial.Title,
				Description: partial.Description,
				Steps:       partial.Steps,
			},
			IssueB: engine.Issue{
				Title:       _issue.Title,
				Description: _issue.Description,
				Steps:       _issue.Steps,
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

		_issue.Title = miResult.Issue.Title
		_issue.Description = miResult.Issue.Description
		_issue.Steps = miResult.Issue.Steps
		tokens += (miResult.Usage.Input + miResult.Usage.Output)

		ceResult, err := self.engineService.ComputeEmbedding(ctx, engine.EngineServiceComputeEmbeddingParams{
			Text: _issue.Description,
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

		_issue.Embedding = ceResult.Embedding
		tokens += (ceResult.Usage.Input + ceResult.Usage.Output)

		_issue.Sources[feedback.Source]++
		_issue.Severities[partial.Severity]++
		_issue.Priority = issue.ComputePriority(_issue.Severities, _issue.Customers+1)
		_issue.Categories[partial.Category]++
		_issue.Releases[feedback.Release]++
		_issue.Customers++
		if feedback.PostedAt.Before(_issue.FirstSeenAt) {
			_issue.FirstSeenAt = feedback.PostedAt
		}
		if feedback.PostedAt.After(_issue.LastSeenAt) {
			_issue.LastSeenAt = feedback.PostedAt
		}
		_issue.LastAggregatedAt = kitUtil.Pointer(time.Now())

		err = self.issueRepository.UpdateAggregated(ctx, *_issue, *feedback)
		if err != nil {
			return err
		}

		err = self.partialIssueRepository.Delete(ctx, partial.ID)
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
			// There are two or more issues from the same feedback trying to merge themselves.
			// This should not be possible, but protect from retry-loop by deleting the partial.
			self.observer.Error(ctx, err)

			err = self.partialIssueRepository.Delete(ctx, partial.ID)
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	self.observer.Infof(ctx, "Merged 2 issues using %d tokens", tokens)

	return nil
}

type IssueAggregatorAggregateParams struct {
	PartialID string
}

func (self *IssueAggregator) Aggregate(ctx context.Context, task *asynq.Task) error {
	if self.engineBreaker.IsOpen(ctx) {
		return nil
	}

	params := IssueAggregatorAggregateParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	partial, err := self.partialIssueRepository.GetByID(ctx, params.PartialID)
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

	issues, err := self.issueRepository.ListByEmbeddingAndProductID(ctx, ceResult.Embedding,
		issue.ISSUE_SIMILAR_THRESHOLD, ISSUE_AGGREGATOR_MAX_SIMILAR_ISSUES, product.ID)
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return self.createIssue(ctx, ceResult.Embedding, partial, feedback, product, tokens)
	}

	options := make([]engine.Issue, 0, len(issues))
	for _, issue := range issues {
		options = append(options, engine.Issue{
			Title:       issue.Title,
			Description: issue.Description,
			Steps:       issue.Steps,
		})
	}

	siResult, err := self.engineService.SimilarIssue(ctx, engine.EngineServiceSimilarIssueParams{
		Issue: engine.Issue{
			Title:       partial.Title,
			Description: partial.Description,
			Steps:       partial.Steps,
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

	tokens += (siResult.Usage.Input + siResult.Usage.Output)

	if siResult.Option == nil {
		return self.createIssue(ctx, ceResult.Embedding, partial, feedback, product, tokens)
	}

	return self.mergeIssues(ctx, partial, &issues[*siResult.Option], feedback, product, tokens)
}

func (self *IssueAggregator) Schedule(ctx context.Context, _ *asynq.Task) error {
	pagination := util.Pagination[time.Time]{
		Limit: 1000,
		From:  nil,
	}

	for {
		page, err := self.partialIssueRepository.ListIDsByCreatedAt(ctx, pagination)
		if err != nil {
			return err
		}

		for _, id := range page.Items {
			err := self.enqueuer.Enqueue(ctx, IssueAggregatorAggregate, IssueAggregatorAggregateParams{
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
