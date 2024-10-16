package collector

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"slices"
	"time"

	"backend/pkg/config"
	"backend/pkg/dataforseo"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/translator"
	"backend/pkg/util"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
)

const (
	AMAZON_COLLECTOR_MAX_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE   = 10 * 10
	AMAZON_COLLECTOR_MIN_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE   = 10
	AMAZON_COLLECTOR_DAILY_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE = AMAZON_COLLECTOR_MIN_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE * 3
)

const (
	AmazonCollectorCollect  = "collector:collect-amazon-reviews"
	AmazonCollectorDispatch = "collector:dispatch-amazon-reviews"
	AmazonCollectorSchedule = "collector:schedule-amazon-reviews"
)

type AmazonCollectorSettings struct {
	CollectorSettings
	ASIN string
}

type AmazonCollectorJobdata struct {
	CollectorJobdata
	LastDispatchedAt    *time.Time
	LastDispatchedTasks []string
	Cost                float64
}

type AmazonCollector struct {
	config                 config.Config
	observer               *kit.Observer
	collectorRepository    *CollectorRepository
	productRepository      *product.ProductRepository
	organizationRepository organization.OrganizationRepository
	feedbackRepository     *feedback.FeedbackRepository
	enqueuer               *kit.Enqueuer
	dataForSEOService      *dataforseo.DataForSEOService
}

func NewAmazonCollector(observer *kit.Observer, collectorRepository *CollectorRepository,
	productRepository *product.ProductRepository, organizationRepository organization.OrganizationRepository,
	feedbackRepository *feedback.FeedbackRepository, enqueuer *kit.Enqueuer, dataForSEOService *dataforseo.DataForSEOService,
	config config.Config) *AmazonCollector {
	return &AmazonCollector{
		config:                 config,
		observer:               observer,
		collectorRepository:    collectorRepository,
		productRepository:      productRepository,
		organizationRepository: organizationRepository,
		feedbackRepository:     feedbackRepository,
		enqueuer:               enqueuer,
		dataForSEOService:      dataForSEOService,
	}
}

type AmazonCollectorCallbackRequest struct {
	Secret string `param:"secret"`
	TaskID string `query:"id"`
}

func (self *AmazonCollector) Callback(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	request := AmazonCollectorCallbackRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Secret != self.config.DataForSEO.CallbackSecret {
		return kit.HTTPErrUnauthorized
	}

	err = self.enqueuer.Enqueue(requestCtx, AmazonCollectorCollect, AmazonCollectorCollectParams{
		TaskID: kitUtil.Pointer(request.TaskID),
	}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

func (self *AmazonCollector) getCollectorProductAndOrganization(ctx context.Context,
	collectorID string) (*Collector, *product.Product, *organization.Organization, error) {
	collector, err := self.collectorRepository.GetByID(ctx, collectorID)
	if err != nil {
		return nil, nil, nil, err
	}

	if collector == nil {
		return nil, nil, nil, nil
	}

	if collector.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	product, err := self.productRepository.GetByID(ctx, collector.ProductID)
	if err != nil {
		return nil, nil, nil, err
	}

	if product == nil {
		return nil, nil, nil, nil
	}

	if product.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	organization, err := self.organizationRepository.GetByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, nil, nil, err
	}

	if organization == nil {
		return nil, nil, nil, nil
	}

	if organization.DeletedAt != nil {
		return nil, nil, nil, nil
	}

	return collector, product, organization, nil
}

func (self *AmazonCollector) saveAndEnqueue(ctx context.Context, feedbacks []feedback.Feedback) (int, error) {
	newFeedbacks, err := self.feedbackRepository.BulkCreate(ctx, feedbacks)
	if err != nil {
		return 0, err
	}

	for _, feedback := range feedbacks {
		err := self.enqueuer.Enqueue(ctx, translator.FeedbackTranslatorTranslate,
			translator.FeedbackTranslatorTranslateParams{
				FeedbackID: feedback.ID,
			}, asynq.MaxRetry(2), asynq.Unique(12*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	return newFeedbacks, nil
}

type AmazonCollectorCollectParams struct {
	TaskID      *string
	CollectorID *string
}

func (self *AmazonCollector) Collect(ctx context.Context, task *asynq.Task) error {
	params := AmazonCollectorCollectParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	var collector *Collector
	var jobdata AmazonCollectorJobdata
	var product *product.Product
	var organization *organization.Organization

	var taskIDs []string
	if params.TaskID != nil {
		taskIDs = append(taskIDs, *params.TaskID)
	} else if params.CollectorID != nil {
		collector, product, organization, err = self.getCollectorProductAndOrganization(ctx, *params.CollectorID)
		if err != nil {
			return err
		} else if collector == nil || product == nil || organization == nil {
			return nil
		}
		jobdata = collector.Jobdata.(AmazonCollectorJobdata)
		taskIDs = append(taskIDs, jobdata.LastDispatchedTasks...)
	}

	totalFeedbacks := 0
	newFeedbacks := 0
	feedbacks := []feedback.Feedback{}
	for _, taskID := range taskIDs {
		task, err := self.dataForSEOService.GetAmazonTask(ctx,
			dataforseo.DataForSEOServiceGetAmazonTaskParams{
				TaskID: taskID,
			})
		if err != nil {
			self.observer.Error(ctx, err)
			continue
		}

		if collector == nil || task.Identifier != collector.ID {
			collector, product, organization, err = self.getCollectorProductAndOrganization(ctx, task.Identifier)
			if err != nil {
				self.observer.Error(ctx, err)
				continue
			} else if collector == nil || product == nil || organization == nil {
				self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().
					With("DataForSEO Amazon task without collector, product or organization attached").
					Extra(map[string]any{"task_id": task.ID}))
				continue
			}

			jobdata = collector.Jobdata.(AmazonCollectorJobdata)
		}

		now := time.Now()

		for _, review := range task.Reviews {
			content := feedback.CleanContent(review.Title, review.Content)
			if len(content) == 0 {
				continue
			}
			link := review.Page
			if review.Link != nil {
				link = *review.Link
			}
			hash := feedback.ComputeHash(feedback.FeedbackSourceAmazon, review.Customer.Name, content)

			_feedback := feedback.NewFeedback()
			_feedback.ID = xid.New().String()
			_feedback.ProductID = product.ID
			_feedback.Hash = hash
			_feedback.Source = feedback.FeedbackSourceAmazon
			_feedback.Customer.Email = nil
			_feedback.Customer.Name = review.Customer.Name
			_feedback.Customer.Picture = review.Customer.Picture
			_feedback.Customer.Location = review.Customer.Location
			_feedback.Customer.Verified = nil
			_feedback.Customer.Reviews = review.Customer.Reviews
			_feedback.Customer.Link = kitUtil.Pointer(review.Customer.Link)
			_feedback.Content = content
			_feedback.Language = engine.OPTION_UNKNOWN
			_feedback.Translation = ""
			_feedback.Release = engine.OPTION_UNKNOWN
			_feedback.Metadata.Rating = kitUtil.Pointer(review.Rating)
			_feedback.Metadata.Media = kitUtil.Pointer(slices.Concat(review.Images, review.Videos))
			_feedback.Metadata.Verified = kitUtil.Pointer(review.Verified)
			_feedback.Metadata.Votes = kitUtil.Pointer(review.Votes)
			_feedback.Metadata.Link = kitUtil.Pointer(link)
			_feedback.Tokens = 0
			_feedback.PostedAt = review.Timestamp
			_feedback.CollectedAt = now
			_feedback.TranslatedAt = nil
			_feedback.ProcessedAt = nil

			feedbacks = append(feedbacks, *_feedback)

			if len(feedbacks) == 1000 {
				_newFeedbacks, err := self.saveAndEnqueue(ctx, feedbacks)
				if err != nil {
					return err
				}

				totalFeedbacks += 1000
				newFeedbacks += _newFeedbacks
				feedbacks = []feedback.Feedback{}
			}
		}

		jobdata.LastDispatchedTasks = util.Filter(jobdata.LastDispatchedTasks, func(dispatched string) bool {
			return dispatched != task.ID
		})

		// We could lose some reviews if the jobdata update isn't in a transaction,
		// but then how to bulk insert in batches in a performant way?
		collector.Jobdata = jobdata
		err = self.collectorRepository.UpdateJobdata(ctx, *collector)
		if err != nil {
			return err
		}
	}

	if len(feedbacks) > 0 {
		_newFeedbacks, err := self.saveAndEnqueue(ctx, feedbacks)
		if err != nil {
			return err
		}

		totalFeedbacks += len(feedbacks)
		newFeedbacks += _newFeedbacks
	}

	self.observer.Infof(ctx,
		"Collected %d DataForSEO Amazon reviews of which %d were duplicated",
		totalFeedbacks, totalFeedbacks-newFeedbacks)

	return nil
}

type AmazonCollectorDispatchParams struct {
	CollectorID string
}

func (self *AmazonCollector) Dispatch(ctx context.Context, task *asynq.Task) error {
	params := AmazonCollectorDispatchParams{}

	err := json.Unmarshal(task.Payload(), &params)
	if err != nil {
		self.observer.Error(ctx, kit.ErrWorkerGeneric.Raise().Cause(err))
		return nil
	}

	collector, product, organization, err := self.getCollectorProductAndOrganization(ctx, params.CollectorID)
	if err != nil {
		return err
	} else if collector == nil || product == nil || organization == nil {
		return nil
	}

	settings := collector.Settings.(AmazonCollectorSettings)
	jobdata := collector.Jobdata.(AmazonCollectorJobdata)

	reviewsPerPerspective := AMAZON_COLLECTOR_MAX_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE
	prioritize := true
	if jobdata.LastDispatchedAt != nil {
		days := int(time.Since(*jobdata.LastDispatchedAt).Hours() / 24)
		reviewsPerPerspective = min(AMAZON_COLLECTOR_MAX_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE, AMAZON_COLLECTOR_DAILY_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE*days)
		prioritize = false
	}

	reviews := reviewsPerPerspective * len(dataforseo.AmazonPerspectives)
	reviews = min(organization.UsageLeft(), reviews)
	if reviews <= 0 {
		return nil
	}

	// This is not very beautiful
	reviewsPerPerspective = int(math.Ceil(float64(
		float64(reviews)/float64(len(dataforseo.AmazonPerspectives)))/float64(AMAZON_COLLECTOR_MIN_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE))) *
		AMAZON_COLLECTOR_MIN_REVIEWS_TO_DISPATCH_PER_PERSPECTIVE
	reviews = reviewsPerPerspective * len(dataforseo.AmazonPerspectives)

	tasks, err := self.dataForSEOService.CreateAmazonTasks(ctx,
		dataforseo.DataForSEOServiceCreateAmazonTasksParams{
			ASIN:         settings.ASIN,
			Perspectives: dataforseo.AmazonPerspectives,
			Reviews:      reviewsPerPerspective,
			Prioritize:   prioritize,
			Identifier:   collector.ID,
			Callback: self.config.Server.BaseURL + "/ext/callback/" +
				url.QueryEscape(self.config.DataForSEO.CallbackSecret) + "/amazon",
		})
	if err != nil {
		return err
	}

	jobdata.LastDispatchedAt = kitUtil.Pointer(time.Now())
	cost := 0.0
	for _, task := range *tasks {
		jobdata.LastDispatchedTasks = append(jobdata.LastDispatchedTasks, task.ID)
		cost += task.Cost
	}
	jobdata.Cost += cost

	collector.Jobdata = jobdata
	err = self.collectorRepository.UpdateJobdata(ctx, *collector)
	if err != nil {
		return err
	}

	self.observer.Infof(ctx,
		"Dispatched %d DataForSEO Amazon tasks with a total of %d reviews and %.4f cost", len(*tasks), reviews, cost)

	return nil
}

func (self *AmazonCollector) Schedule(ctx context.Context, _ *asynq.Task) error {
	ids, err := self.collectorRepository.ListIDsByTypeNotDeleted(ctx, CollectorTypeAmazon)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err := self.enqueuer.Enqueue(ctx, AmazonCollectorDispatch, AmazonCollectorDispatchParams{
			CollectorID: id,
		}, asynq.MaxRetry(2), asynq.Unique(24*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}

		err = self.enqueuer.Enqueue(ctx, AmazonCollectorCollect, AmazonCollectorCollectParams{
			CollectorID: kitUtil.Pointer(id),
		}, asynq.MaxRetry(2), asynq.ProcessIn(12*time.Hour))
		if err != nil {
			self.observer.Error(ctx, err)
		}
	}

	return nil
}
