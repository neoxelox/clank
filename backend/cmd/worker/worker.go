package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/neoxelox/kit"
	kitMiddleware "github.com/neoxelox/kit/middleware"
	kitUtil "github.com/neoxelox/kit/util"

	"backend/pkg/aggregator"
	"backend/pkg/auth"
	"backend/pkg/brevo"
	"backend/pkg/collector"
	"backend/pkg/config"
	"backend/pkg/dataforseo"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/issue"
	"backend/pkg/organization"
	"backend/pkg/processor"
	"backend/pkg/product"
	"backend/pkg/review"
	"backend/pkg/scraper"
	"backend/pkg/suggestion"
	"backend/pkg/translator"
	"backend/pkg/user"
	"backend/pkg/util"
)

type Worker struct {
	Run   func(ctx context.Context) error
	Close func(ctx context.Context) error
}

func NewWorker(ctx context.Context, config config.Config) (*Worker, error) {
	retry := kit.RetryConfig{
		Attempts:     5,
		InitialDelay: 1 * time.Second,
		LimitDelay:   5 * time.Second,
	}

	level := kit.LvlInfo
	if config.Service.Environment == kit.EnvDevelopment {
		level = kit.LvlDebug
	}

	var observerSentryConfig *kit.ObserverSentryConfig
	if config.Service.Environment == kit.EnvProduction {
		observerSentryConfig = &kit.ObserverSentryConfig{
			Dsn: config.Sentry.DSN,
		}
	}

	observer, err := kit.NewObserver(ctx, kit.ObserverConfig{
		Environment: config.Service.Environment,
		Release:     config.Service.Release,
		Service:     config.Service.Name,
		Level:       level,
		Sentry:      observerSentryConfig,
	}, retry)
	if err != nil {
		return nil, err
	}

	migrator, err := kit.NewMigrator(ctx, observer, kit.MigratorConfig{
		DatabaseHost:     config.Database.Host,
		DatabasePort:     config.Database.Port,
		DatabaseSSLMode:  config.Database.SSLMode,
		DatabaseUser:     config.Database.User,
		DatabasePassword: config.Database.Password,
		DatabaseName:     config.Database.Name,
		MigrationsPath:   kitUtil.Pointer(config.Service.MigrationsPath),
	}, retry)
	if err != nil {
		return nil, err
	}

	err = migrator.Apply(ctx, config.Database.SchemaVersion)
	if err != nil {
		return nil, err
	}

	err = migrator.Assert(ctx, config.Database.SchemaVersion)
	if err != nil {
		return nil, err
	}

	err = migrator.Close(ctx)
	if err != nil {
		return nil, err
	}

	errorHandler := kit.NewErrorHandler(observer, kit.ErrorHandlerConfig{
		Environment: config.Service.Environment,
	})

	_, err = kit.NewRenderer(observer, kit.RendererConfig{
		TemplatesPath:       kitUtil.Pointer(config.Service.TemplatesPath),
		TemplateFilePattern: kitUtil.Pointer(config.Service.TemplateFilePattern),
	})
	if err != nil {
		return nil, err
	}

	_, err = kit.NewLocalizer(observer, kit.LocalizerConfig{
		DefaultLocale:     config.Service.DefaultLocale,
		LocalesPath:       kitUtil.Pointer(config.Service.LocalesPath),
		LocaleFilePattern: kitUtil.Pointer(config.Service.LocaleFilePattern),
	})
	if err != nil {
		return nil, err
	}

	worker := kit.NewWorker(observer, errorHandler, kit.WorkerConfig{
		Queues:               config.Worker.Queues,
		Concurrency:          kitUtil.Pointer(config.Worker.Concurrency),
		StrictPriority:       kitUtil.Pointer(config.Worker.StrictPriority),
		StopTimeout:          kitUtil.Pointer(config.Worker.StopTimeout),
		TimeZone:             kitUtil.Pointer(config.Service.TimeZone),
		ScheduleDefaultRetry: kitUtil.Pointer(0),
		CacheHost:            config.Cache.Host,
		CachePort:            config.Cache.Port,
		CacheSSLMode:         config.Cache.SSLMode,
		CachePassword:        config.Cache.Password,
		CacheMaxConns:        kitUtil.Pointer(config.Cache.MaxConns),
		CacheReadTimeout:     kitUtil.Pointer(config.Cache.ReadTimeout),
		CacheWriteTimeout:    kitUtil.Pointer(config.Cache.WriteTimeout),
		CacheDialTimeout:     kitUtil.Pointer(config.Cache.DialTimeout),
	})

	observerMiddleware := kitMiddleware.NewObserver(observer, kitMiddleware.ObserverConfig{})
	recoverMiddleware := kitMiddleware.NewRecover(observer, kitMiddleware.RecoverConfig{})

	worker.Use(observerMiddleware.HandleTask)
	worker.Use(recoverMiddleware.HandleTask)

	database, err := kit.NewDatabase(ctx, observer, kit.DatabaseConfig{
		Host:                  config.Database.Host,
		Port:                  config.Database.Port,
		SSLMode:               config.Database.SSLMode,
		User:                  config.Database.User,
		Password:              config.Database.Password,
		Database:              config.Database.Name,
		Service:               config.Service.Name,
		MinConns:              kitUtil.Pointer(config.Database.MinConns),
		MaxConns:              kitUtil.Pointer(config.Database.MaxConns),
		MaxConnIdleTime:       kitUtil.Pointer(config.Database.MaxConnIdleTime),
		MaxConnLifeTime:       kitUtil.Pointer(config.Database.MaxConnLifeTime),
		DialTimeout:           kitUtil.Pointer(config.Database.DialTimeout),
		StatementTimeout:      kitUtil.Pointer(config.Database.StatementTimeout),
		DefaultIsolationLevel: kitUtil.Pointer(config.Database.DefaultIsolationLevel),
	}, retry)
	if err != nil {
		return nil, err
	}

	cache, err := kit.NewCache(ctx, observer, kit.CacheConfig{
		Host:            config.Cache.Host,
		Port:            config.Cache.Port,
		SSLMode:         config.Cache.SSLMode,
		Password:        config.Cache.Password,
		MinConns:        kitUtil.Pointer(config.Cache.MinConns),
		MaxConns:        kitUtil.Pointer(config.Cache.MaxConns),
		MaxConnIdleTime: kitUtil.Pointer(config.Cache.MaxConnIdleTime),
		MaxConnLifeTime: kitUtil.Pointer(config.Cache.MaxConnLifeTime),
		ReadTimeout:     kitUtil.Pointer(config.Cache.ReadTimeout),
		WriteTimeout:    kitUtil.Pointer(config.Cache.WriteTimeout),
		DialTimeout:     kitUtil.Pointer(config.Cache.DialTimeout),
	}, retry)
	if err != nil {
		return nil, err
	}

	enqueuer := kit.NewEnqueuer(observer, kit.EnqueuerConfig{
		CacheHost:         config.Cache.Host,
		CachePort:         config.Cache.Port,
		CacheSSLMode:      config.Cache.SSLMode,
		CachePassword:     config.Cache.Password,
		CacheMaxConns:     kitUtil.Pointer(config.Cache.MaxConns),
		CacheReadTimeout:  kitUtil.Pointer(config.Cache.ReadTimeout),
		CacheWriteTimeout: kitUtil.Pointer(config.Cache.WriteTimeout),
		CacheDialTimeout:  kitUtil.Pointer(config.Cache.DialTimeout),
		TaskDefaultRetry:  kitUtil.Pointer(0),
	})

	/* REPOSITORIES  */

	invitationRepository := user.NewInvitationRepositoryImpl(observer, database, config)
	signInCodeRepository := auth.NewSignInCodeRepositoryImpl(observer, database, config)
	organizationRepository := organization.NewOrganizationRepositoryImpl(observer, database, config)
	productRepository := product.NewProductRepository(observer, database, config)
	collectorRepository := collector.NewCollectorRepository(observer, database, config)
	feedbackRepository := feedback.NewFeedbackRepository(observer, database, config)
	partialIssueRepository := issue.NewPartialIssueRepository(observer, database, config)
	partialSuggestionRepository := suggestion.NewPartialSuggestionRepository(observer, database, config)
	issueRepository := issue.NewIssueRepository(observer, database, config)
	suggestionRepository := suggestion.NewSuggestionRepository(observer, database, config)
	reviewRepository := review.NewReviewRepository(observer, database, config)

	/* SERVICES */

	engineService := engine.NewEngineService(observer, config)
	dataForSEOService := dataforseo.NewDataForSEOService(observer, config)
	brevoService := brevo.NewBrevoService(observer, config)

	/* USECASES */

	engineBreaker := engine.NewEngineBreaker(observer, cache, config)
	scraper := scraper.NewScraper(observer, config)

	trustpilotCollector := collector.NewTrustpilotCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	playStoreCollector := collector.NewPlayStoreCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	appStoreCollector := collector.NewAppStoreCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	amazonCollector := collector.NewAmazonCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	iAgoraCollector := collector.NewIAgoraCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, scraper, config)

	feedbackTranslator := translator.NewFeedbackTranslator(observer, feedbackRepository, productRepository,
		organizationRepository, enqueuer, engineService, engineBreaker, config)
	feedbackProcessor := processor.NewFeedbackProcessor(observer, database, feedbackRepository, partialIssueRepository,
		partialSuggestionRepository, reviewRepository, productRepository, organizationRepository, enqueuer,
		engineService, engineBreaker, config)
	issueAggregator := aggregator.NewIssueAggregator(observer, database, partialIssueRepository, issueRepository,
		feedbackRepository, productRepository, organizationRepository, enqueuer, engineService, engineBreaker, config)
	suggestionAggregator := aggregator.NewSuggestionAggregator(observer, database, partialSuggestionRepository,
		suggestionRepository, feedbackRepository, productRepository, organizationRepository, enqueuer, engineService,
		engineBreaker, config)

	/* TASKS */

	authTasks := auth.NewAuthTasks(observer, signInCodeRepository, config)
	userTasks := user.NewUserTasks(observer, invitationRepository, config)
	organizationTasks := organization.NewOrganizationTasks(observer, database, organizationRepository, enqueuer, config)

	/* MIDDLEWARES */

	/* REGISTRATIONS */

	worker.Register(auth.AuthTasksDeleteExpiredSignInCodes, authTasks.DeleteExpiredSignInCodes)

	worker.Register(user.UserTasksDeleteExpiredInvitations, userTasks.DeleteExpiredInvitations)

	worker.Register(organization.OrganizationTasksDowngradeEndedTrials, organizationTasks.DowngradeEndedTrials)
	worker.Register(organization.OrganizationTasksComputeUsage, organizationTasks.ComputeUsage)
	worker.Register(organization.OrganizationTasksScheduleComputeUsage, organizationTasks.ScheduleComputeUsage)

	worker.Register(collector.TrustpilotCollectorCollect, trustpilotCollector.Collect)
	worker.Register(collector.TrustpilotCollectorDispatch, trustpilotCollector.Dispatch)
	worker.Register(collector.TrustpilotCollectorSchedule, trustpilotCollector.Schedule)

	worker.Register(collector.PlayStoreCollectorCollect, playStoreCollector.Collect)
	worker.Register(collector.PlayStoreCollectorDispatch, playStoreCollector.Dispatch)
	worker.Register(collector.PlayStoreCollectorSchedule, playStoreCollector.Schedule)

	worker.Register(collector.AppStoreCollectorCollect, appStoreCollector.Collect)
	worker.Register(collector.AppStoreCollectorDispatch, appStoreCollector.Dispatch)
	worker.Register(collector.AppStoreCollectorSchedule, appStoreCollector.Schedule)

	worker.Register(collector.AmazonCollectorCollect, amazonCollector.Collect)
	worker.Register(collector.AmazonCollectorDispatch, amazonCollector.Dispatch)
	worker.Register(collector.AmazonCollectorSchedule, amazonCollector.Schedule)

	worker.Register(collector.IAgoraCollectorCollect, iAgoraCollector.Collect)
	worker.Register(collector.IAgoraCollectorSchedule, iAgoraCollector.Schedule)

	worker.Register(translator.FeedbackTranslatorTranslate, feedbackTranslator.Translate)
	worker.Register(translator.FeedbackTranslatorSchedule, feedbackTranslator.Schedule)

	worker.Register(processor.FeedbackProcessorProcess, feedbackProcessor.Process)
	worker.Register(processor.FeedbackProcessorSchedule, feedbackProcessor.Schedule)

	worker.Register(aggregator.IssueAggregatorAggregate, issueAggregator.Aggregate)
	worker.Register(aggregator.IssueAggregatorSchedule, issueAggregator.Schedule)

	worker.Register(aggregator.SuggestionAggregatorAggregate, suggestionAggregator.Aggregate)
	worker.Register(aggregator.SuggestionAggregatorSchedule, suggestionAggregator.Schedule)

	/* SCHEDULEMENTS */

	worker.Schedule(auth.AuthTasksDeleteExpiredSignInCodes, nil, "0 8 * * *", asynq.Queue("irrelevant"))                              // Every day at 08:00
	worker.Schedule(user.UserTasksDeleteExpiredInvitations, nil, "0 8 * * *", asynq.Queue("irrelevant"))                              // Every day at 08:00
	worker.Schedule(organization.OrganizationTasksDowngradeEndedTrials, nil, "0 8 * * *", asynq.Queue("critical"), asynq.MaxRetry(2)) // Every day at 08:00
	worker.Schedule(organization.OrganizationTasksScheduleComputeUsage, nil, "*/5 * * * *", asynq.Queue("critical"))                  // Every 5 minutes
	worker.Schedule(collector.TrustpilotCollectorSchedule, nil, "0 0 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))           // Every day at 00:00
	worker.Schedule(collector.PlayStoreCollectorSchedule, nil, "0 1 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))            // Every day at 01:00
	worker.Schedule(collector.AppStoreCollectorSchedule, nil, "0 2 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))             // Every day at 02:00
	worker.Schedule(collector.AmazonCollectorSchedule, nil, "0 3 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))               // Every day at 03:00
	worker.Schedule(collector.IAgoraCollectorSchedule, nil, "0 4 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))               // Every day at 04:00
	worker.Schedule(translator.FeedbackTranslatorSchedule, nil, "0 19 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))          // Every day at 19:00
	worker.Schedule(processor.FeedbackProcessorSchedule, nil, "0 20 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))            // Every day at 20:00
	worker.Schedule(aggregator.IssueAggregatorSchedule, nil, "0 22 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))             // Every day at 22:00
	worker.Schedule(aggregator.SuggestionAggregatorSchedule, nil, "0 23 * * *", asynq.MaxRetry(2), asynq.Unique(24*time.Hour))        // Every day at 23:00

	return &Worker{
		Run: func(ctx context.Context) error {
			observer.Infof(ctx, "Starting %s Worker", config.Service.Name)

			// Create a concurrent http server to satisfy health checks
			go func() {
				healthEndpoints := util.NewHealthEndpoints(observer, database, cache, config)
				err := http.ListenAndServe(fmt.Sprintf(":%d", config.Worker.HealthPort),
					http.TimeoutHandler(http.HandlerFunc(healthEndpoints.GetWorkerHealth),
						config.Service.GracefulTimeout, kit.HTTPErrServerTimeout.String()))
				if err != nil && err != http.ErrServerClosed {
					panic(fmt.Sprintf("%+v", err))
				}
			}()

			err := worker.Run(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		Close: func(ctx context.Context) error {
			err := kitUtil.Deadline(ctx, func(exceeded <-chan struct{}) error {
				observer.Infof(ctx, "Closing %s Worker", config.Service.Name)

				err := worker.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = enqueuer.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = engineService.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = dataForSEOService.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = brevoService.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = cache.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = database.Close(ctx)
				if err != nil {
					observer.Error(ctx, err)
				}

				err = observer.Close(ctx)
				if err != nil {
					fmt.Printf("%+v", err) // nolint:forbidigo
				}

				return nil
			})
			if err != nil {
				return err
			}

			return nil
		},
	}, nil
}
