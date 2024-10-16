package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	kitMiddleware "github.com/neoxelox/kit/middleware"
	kitUtil "github.com/neoxelox/kit/util"

	"backend/pkg/auth"
	"backend/pkg/brevo"
	"backend/pkg/collector"
	"backend/pkg/config"
	"backend/pkg/dataforseo"
	"backend/pkg/engine"
	"backend/pkg/exporter"
	"backend/pkg/feedback"
	"backend/pkg/issue"
	"backend/pkg/metric"
	"backend/pkg/organization"
	"backend/pkg/product"
	"backend/pkg/review"
	"backend/pkg/suggestion"
	"backend/pkg/user"
	"backend/pkg/util"
)

type API struct {
	Run   func(ctx context.Context) error
	Close func(ctx context.Context) error
}

func NewAPI(ctx context.Context, config config.Config) (*API, error) {
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

	var observerGilkConfig *kit.ObserverGilkConfig
	if config.Service.Environment == kit.EnvDevelopment {
		observerGilkConfig = &kit.ObserverGilkConfig{
			Port: config.Gilk.Port,
		}
	}

	observer, err := kit.NewObserver(ctx, kit.ObserverConfig{
		Environment: config.Service.Environment,
		Release:     config.Service.Release,
		Service:     config.Service.Name,
		Level:       level,
		Sentry:      observerSentryConfig,
		Gilk:        observerGilkConfig,
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

	serializer := kit.NewSerializer(observer, kit.SerializerConfig{})

	binder := kit.NewBinder(observer, kit.BinderConfig{})

	renderer, err := kit.NewRenderer(observer, kit.RendererConfig{
		TemplatesPath:       kitUtil.Pointer(config.Service.TemplatesPath),
		TemplateFilePattern: kitUtil.Pointer(config.Service.TemplateFilePattern),
	})
	if err != nil {
		return nil, err
	}

	localizer, err := kit.NewLocalizer(observer, kit.LocalizerConfig{
		DefaultLocale:     config.Service.DefaultLocale,
		LocalesPath:       kitUtil.Pointer(config.Service.LocalesPath),
		LocaleFilePattern: kitUtil.Pointer(config.Service.LocaleFilePattern),
	})
	if err != nil {
		return nil, err
	}

	server := kit.NewHTTPServer(observer, serializer, binder, renderer, errorHandler, kit.HTTPServerConfig{
		Environment:              config.Service.Environment,
		Port:                     config.Server.Port,
		RequestHeaderMaxSize:     kitUtil.Pointer(config.Server.RequestHeaderMaxSize),
		RequestBodyMaxSize:       kitUtil.Pointer(config.Server.RequestBodyMaxSize),
		RequestFileMaxSize:       kitUtil.Pointer(config.Server.RequestFileMaxSize),
		RequestFilePattern:       kitUtil.Pointer(config.Server.RequestFilePattern),
		RequestKeepAliveTimeout:  kitUtil.Pointer(config.Server.RequestKeepAliveTimeout),
		RequestReadTimeout:       kitUtil.Pointer(config.Server.RequestReadTimeout),
		RequestReadHeaderTimeout: kitUtil.Pointer(config.Server.RequestReadHeaderTimeout),
		RequestIPExtractor:       kitUtil.Pointer((func(*http.Request) string)(echo.ExtractIPFromXFFHeader())),
		ResponseWriteTimeout:     kitUtil.Pointer(config.Server.ResponseWriteTimeout),
	})

	observerMiddleware := kitMiddleware.NewObserver(observer, kitMiddleware.ObserverConfig{})
	timeoutMiddleware := kitMiddleware.NewTimeout(observer, kitMiddleware.TimeoutConfig{
		// Allow timeout handler to respond before response writer is closed
		Timeout: config.Service.GracefulTimeout - (100 * time.Millisecond),
	})
	recoverMiddleware := kitMiddleware.NewRecover(observer, kitMiddleware.RecoverConfig{})
	secureMiddleware := kitMiddleware.NewSecure(observer, kitMiddleware.SecureConfig{
		CORSAllowOrigins: kitUtil.Pointer(config.Server.Origins),
		CORSAllowMethods: kitUtil.Pointer([]string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "DELETE"}),
		CORSAllowHeaders: kitUtil.Pointer([]string{"Content-Type", "X-Forwarded-For", "X-Real-IP",
			"X-Trace-Id", "sentry-trace", "baggage"}),
		CORSAllowCredentials: kitUtil.Pointer(true),
	})
	localizerMiddleware := kitMiddleware.NewLocalizer(observer, localizer, kitMiddleware.LocalizerConfig{})
	errorMiddleware := kitMiddleware.NewError(observer, kitMiddleware.ErrorConfig{})

	server.Use(observerMiddleware.HandleRequest)
	server.Use(timeoutMiddleware.Handle)
	server.Use(recoverMiddleware.HandleRequest)
	server.Use(secureMiddleware.Handle)
	server.Use(localizerMiddleware.Handle)
	server.Use(errorMiddleware.Handle)

	api := server.Default()

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

	limiter := kit.NewLimiter(observer, kit.LimiterConfig{
		CacheHost:            config.Cache.Host,
		CachePort:            config.Cache.Port,
		CacheSSLMode:         config.Cache.SSLMode,
		CachePassword:        config.Cache.Password,
		CacheMinConns:        kitUtil.Pointer(config.Cache.MinConns),
		CacheMaxConns:        kitUtil.Pointer(config.Cache.MaxConns),
		CacheMaxConnIdleTime: kitUtil.Pointer(config.Cache.MaxConnIdleTime),
		CacheMaxConnLifeTime: kitUtil.Pointer(config.Cache.MaxConnLifeTime),
		CacheReadTimeout:     kitUtil.Pointer(config.Cache.ReadTimeout),
		CacheWriteTimeout:    kitUtil.Pointer(config.Cache.WriteTimeout),
		CacheDialTimeout:     kitUtil.Pointer(config.Cache.DialTimeout),
	})

	/* REPOSITORIES  */

	userRepository := user.NewUserRepositoryImpl(observer, database, config)
	invitationRepository := user.NewInvitationRepositoryImpl(observer, database, config)
	sessionRepository := auth.NewSessionRepositoryImpl(observer, database, config)
	signInCodeRepository := auth.NewSignInCodeRepositoryImpl(observer, database, config)
	organizationRepository := organization.NewOrganizationRepositoryImpl(observer, database, config)
	productRepository := product.NewProductRepository(observer, database, config)
	collectorRepository := collector.NewCollectorRepository(observer, database, config)
	exporterRepository := exporter.NewExporterRepository(observer, database, config)
	feedbackRepository := feedback.NewFeedbackRepository(observer, database, config)
	issueRepository := issue.NewIssueRepository(observer, database, config)
	suggestionRepository := suggestion.NewSuggestionRepository(observer, database, config)
	reviewRepository := review.NewReviewRepository(observer, database, config)
	metricRepository := metric.NewMetricRepository(observer, database, config)

	/* SERVICES */

	engineService := engine.NewEngineService(observer, config)
	dataForSEOService := dataforseo.NewDataForSEOService(observer, config)
	brevoService := brevo.NewBrevoService(observer, config)

	/* USECASES */

	authVerifier := auth.NewAuthVerifier(observer, sessionRepository, userRepository, organizationRepository, config)
	authProcessor := auth.NewAuthProcessor(observer, database, signInCodeRepository, renderer, brevoService,
		authVerifier, userRepository, invitationRepository, organizationRepository, sessionRepository, config)
	trustpilotCollector := collector.NewTrustpilotCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	playStoreCollector := collector.NewPlayStoreCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	appStoreCollector := collector.NewAppStoreCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)
	amazonCollector := collector.NewAmazonCollector(observer, collectorRepository, productRepository,
		organizationRepository, feedbackRepository, enqueuer, dataForSEOService, config)

	/* ENDPOINTS */

	healthEndpoints := util.NewHealthEndpoints(observer, database, cache, config)
	authEndpoints := auth.NewAuthEndpoints(observer, authProcessor, config)
	userEndpoints := user.NewUserEndpoints(observer, database, renderer, brevoService, userRepository, invitationRepository, organizationRepository, config)
	organizationEndpoints := organization.NewOrganizationEndpoints(observer, organizationRepository, config)
	productEndpoints := product.NewProductEndpoints(observer, productRepository, config)
	collectorEndpoints := collector.NewCollectorEndpoints(observer, collectorRepository, enqueuer, config)
	exporterEndpoints := exporter.NewExporterEndpoints(observer, exporterRepository, config)
	issueEndpoints := issue.NewIssueEndpoints(observer, issueRepository, userRepository, engineService, cache, config)
	suggestionEndpoints := suggestion.NewSuggestionEndpoints(observer, suggestionRepository, userRepository, engineService, cache, config)
	reviewEndpoints := review.NewReviewEndpoints(observer, reviewRepository, cache, config)
	metricEndpoints := metric.NewMetricEndpoints(observer, metricRepository, cache, config)

	/* MIDDLEWARES */

	rateLimitMiddleware := util.NewRateLimitMiddleware(observer, limiter, config)
	authMiddlewares := auth.NewAuthMiddlewares(observer, authVerifier, sessionRepository, config)
	userMiddlewares := user.NewUserMiddlewares(observer, userRepository, invitationRepository, config)
	productMiddleware := product.NewProductMiddleware(observer, productRepository, config)
	collectorMiddleware := collector.NewCollectorMiddleware(observer, collectorRepository, config)
	exporterMiddleware := exporter.NewExporterMiddleware(observer, exporterRepository, config)
	issueMiddleware := issue.NewIssueMiddleware(observer, issueRepository, config)
	suggestionMiddleware := suggestion.NewSuggestionMiddleware(observer, suggestionRepository, config)
	reviewMiddleware := review.NewReviewMiddleware(observer, reviewRepository, config)

	/* INTERNAL ROUTES */

	internalRoutes := api.Group("/int")

	internalRoutes.GET("/health", healthEndpoints.GetServerHealth)

	/* EXTERNAL ROUTES */

	rootRoutes := api.Group("/ext")

	// NOT AUTHENTICATED

	rootRoutes.POST("/signin/:provider/start", authEndpoints.PostSignInStart, rateLimitMiddleware.Handle(3, 1*time.Minute))
	rootRoutes.POST("/signin/:provider/end", authEndpoints.PostSignInEnd, rateLimitMiddleware.Handle(5, 1*time.Minute))

	rootRoutes.GET("/callback/:secret/trustpilot", trustpilotCollector.Callback)
	rootRoutes.GET("/callback/:secret/play-store", playStoreCollector.Callback)
	rootRoutes.GET("/callback/:secret/app-store", appStoreCollector.Callback)
	rootRoutes.GET("/callback/:secret/amazon", amazonCollector.Callback)
	rootRoutes.POST("/callback/:secret/webhook", nil) // TODO
	rootRoutes.POST("/callback/:secret/widget", nil)  // TODO

	// AUTHENTICATED

	rootRoutes = rootRoutes.Group("", authMiddlewares.HandleToken)

	rootRoutes.POST("/signout", authEndpoints.PostSignOut)

	rootRoutes.GET("/user", userEndpoints.GetMe)
	rootRoutes.PUT("/user", userEndpoints.PutMe)
	rootRoutes.DELETE("/user", userEndpoints.DeleteMe)
	rootRoutes.GET("/user/settings", userEndpoints.GetMySettings)
	rootRoutes.PUT("/user/settings", userEndpoints.PutMySettings)

	rootRoutes.GET("/organization", organizationEndpoints.GetOrganization)
	rootRoutes.PUT("/organization", organizationEndpoints.PutOrganization, authMiddlewares.HandleRights)
	rootRoutes.DELETE("/organization", organizationEndpoints.DeleteOrganization, authMiddlewares.HandleRights)
	rootRoutes.GET("/organization/settings", organizationEndpoints.GetOrganizationSettings)
	rootRoutes.PUT("/organization/settings", organizationEndpoints.PutOrganizationSettings, authMiddlewares.HandleRights)
	rootRoutes.GET("/organization/usage", organizationEndpoints.GetOrganizationUsage)
	rootRoutes.GET("/organization/billing", nil, authMiddlewares.HandleRights) // TODO
	rootRoutes.PUT("/organization/billing", nil, authMiddlewares.HandleRights) // TODO

	userRoutes := rootRoutes.Group("")
	userRoutes.GET("/users", userEndpoints.ListUsers)
	userRoutes = userRoutes.Group("", userMiddlewares.HandleUser)
	userRoutes.GET("/users/:user_id", userEndpoints.GetUser)
	userRoutes.PUT("/users/:user_id", userEndpoints.PutUser, authMiddlewares.HandleRights)
	userRoutes.DELETE("/users/:user_id", userEndpoints.DeleteUser, authMiddlewares.HandleRights)

	invitationRoutes := rootRoutes.Group("")
	invitationRoutes.GET("/invitations", userEndpoints.ListInvitations)
	invitationRoutes.POST("/invitations", userEndpoints.PostInvitation, rateLimitMiddleware.Handle(3, 1*time.Minute), authMiddlewares.HandleRights)
	invitationRoutes = invitationRoutes.Group("", userMiddlewares.HandleInvitation)
	invitationRoutes.GET("/invitations/:invitation_id", userEndpoints.GetInvitation)
	invitationRoutes.DELETE("/invitations/:invitation_id", userEndpoints.DeleteInvitation, authMiddlewares.HandleRights)

	productRoutes := rootRoutes.Group("")
	productRoutes.GET("/products", productEndpoints.ListProducts)
	productRoutes.POST("/products", productEndpoints.PostProduct, authMiddlewares.HandleRights)
	productRoutes = productRoutes.Group("", productMiddleware.Handle)
	productRoutes.GET("/products/:product_id", productEndpoints.GetProduct)
	productRoutes.PUT("/products/:product_id", productEndpoints.PutProduct, authMiddlewares.HandleRights)
	productRoutes.DELETE("/products/:product_id", productEndpoints.DeleteProduct, authMiddlewares.HandleRights)
	productRoutes.GET("/products/:product_id/settings", productEndpoints.GetProductSettings)
	productRoutes.PUT("/products/:product_id/settings", productEndpoints.PutProductSettings, authMiddlewares.HandleRights)
	productRoutes.GET("/products/:product_id/usage", productEndpoints.GetProductUsage)

	collectorRoutes := productRoutes.Group("")
	collectorRoutes.GET("/products/:product_id/collectors", collectorEndpoints.ListCollectors)
	collectorRoutes.POST("/products/:product_id/collectors", collectorEndpoints.PostCollector, authMiddlewares.HandleRights)
	collectorRoutes = collectorRoutes.Group("", collectorMiddleware.Handle)
	collectorRoutes.GET("/products/:product_id/collectors/:collector_id", collectorEndpoints.GetCollector)
	collectorRoutes.PUT("/products/:product_id/collectors/:collector_id", collectorEndpoints.PutCollector, authMiddlewares.HandleRights)
	collectorRoutes.DELETE("/products/:product_id/collectors/:collector_id", collectorEndpoints.DeleteCollector, authMiddlewares.HandleRights)

	exporterRoutes := productRoutes.Group("")
	exporterRoutes.GET("/products/:product_id/exporters", exporterEndpoints.ListExporters)
	exporterRoutes.POST("/products/:product_id/exporters", exporterEndpoints.PostExporter, authMiddlewares.HandleRights)
	exporterRoutes = exporterRoutes.Group("", exporterMiddleware.Handle)
	exporterRoutes.GET("/products/:product_id/exporters/:exporter_id", exporterEndpoints.GetExporter)
	exporterRoutes.PUT("/products/:product_id/exporters/:exporter_id", exporterEndpoints.PutExporter, authMiddlewares.HandleRights)
	exporterRoutes.DELETE("/products/:product_id/exporters/:exporter_id", exporterEndpoints.DeleteExporter, authMiddlewares.HandleRights)

	issueRoutes := productRoutes.Group("")
	issueRoutes.GET("/products/:product_id/issues", issueEndpoints.ListIssues)
	issueRoutes = issueRoutes.Group("", issueMiddleware.Handle)
	issueRoutes.GET("/products/:product_id/issues/:issue_id", issueEndpoints.GetIssue)
	issueRoutes.GET("/products/:product_id/issues/:issue_id/feedbacks", issueEndpoints.ListIssueFeedbacks)
	issueRoutes.PUT("/products/:product_id/issues/:issue_id/assignee", issueEndpoints.PutIssueAssignee)
	issueRoutes.PUT("/products/:product_id/issues/:issue_id/quality", issueEndpoints.PutIssueQuality)
	issueRoutes.PUT("/products/:product_id/issues/:issue_id/archived", issueEndpoints.PutIssueArchived)

	suggestionRoutes := productRoutes.Group("")
	suggestionRoutes.GET("/products/:product_id/suggestions", suggestionEndpoints.ListSuggestions)
	suggestionRoutes = suggestionRoutes.Group("", suggestionMiddleware.Handle)
	suggestionRoutes.GET("/products/:product_id/suggestions/:suggestion_id", suggestionEndpoints.GetSuggestion)
	suggestionRoutes.GET("/products/:product_id/suggestions/:suggestion_id/feedbacks", suggestionEndpoints.ListSuggestionFeedbacks)
	suggestionRoutes.PUT("/products/:product_id/suggestions/:suggestion_id/assignee", suggestionEndpoints.PutSuggestionAssignee)
	suggestionRoutes.PUT("/products/:product_id/suggestions/:suggestion_id/quality", suggestionEndpoints.PutSuggestionQuality)
	suggestionRoutes.PUT("/products/:product_id/suggestions/:suggestion_id/archived", suggestionEndpoints.PutSuggestionArchived)

	reviewRoutes := productRoutes.Group("")
	reviewRoutes.GET("/products/:product_id/reviews", reviewEndpoints.ListReviews)
	reviewRoutes = reviewRoutes.Group("", reviewMiddleware.Handle)
	reviewRoutes.GET("/products/:product_id/reviews/:review_id", reviewEndpoints.GetReview)
	reviewRoutes.PUT("/products/:product_id/reviews/:review_id/quality", reviewEndpoints.PutReviewQuality)

	metricRoutes := productRoutes.Group("")
	metricRoutes.GET("/products/:product_id/metrics/issue-count", metricEndpoints.GetIssueCount)
	metricRoutes.GET("/products/:product_id/metrics/issue-sources", metricEndpoints.GetIssueSources)
	metricRoutes.GET("/products/:product_id/metrics/issue-severities", metricEndpoints.GetIssueSeverities)
	metricRoutes.GET("/products/:product_id/metrics/issue-categories", metricEndpoints.GetIssueCategories)
	metricRoutes.GET("/products/:product_id/metrics/issue-releases", metricEndpoints.GetIssueReleases)
	metricRoutes.GET("/products/:product_id/metrics/suggestion-count", metricEndpoints.GetSuggestionCount)
	metricRoutes.GET("/products/:product_id/metrics/suggestion-sources", metricEndpoints.GetSuggestionSources)
	metricRoutes.GET("/products/:product_id/metrics/suggestion-importances", metricEndpoints.GetSuggestionImportances)
	metricRoutes.GET("/products/:product_id/metrics/suggestion-categories", metricEndpoints.GetSuggestionCategories)
	metricRoutes.GET("/products/:product_id/metrics/suggestion-releases", metricEndpoints.GetSuggestionReleases)
	metricRoutes.GET("/products/:product_id/metrics/review-sentiments", metricEndpoints.GetReviewSentiments)
	metricRoutes.GET("/products/:product_id/metrics/review-sources", metricEndpoints.GetReviewSources)
	metricRoutes.GET("/products/:product_id/metrics/review-intentions", metricEndpoints.GetReviewIntentions)
	metricRoutes.GET("/products/:product_id/metrics/review-emotions", metricEndpoints.GetReviewEmotions)
	metricRoutes.GET("/products/:product_id/metrics/review-categories", metricEndpoints.GetReviewCategories)
	metricRoutes.GET("/products/:product_id/metrics/review-releases", metricEndpoints.GetReviewReleases)
	metricRoutes.GET("/products/:product_id/metrics/review-keywords", metricEndpoints.GetReviewKeywords)
	metricRoutes.GET("/products/:product_id/metrics/nps", metricEndpoints.GetNetPromoterScore)
	metricRoutes.GET("/products/:product_id/metrics/csat", metricEndpoints.GetCustomerSatisfactionScore)

	return &API{
		Run: func(ctx context.Context) error {
			observer.Infof(ctx, "Starting %s API", config.Service.Name)

			err := server.Run(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		Close: func(ctx context.Context) error {
			err := kitUtil.Deadline(ctx, func(exceeded <-chan struct{}) error {
				observer.Infof(ctx, "Closing %s API", config.Service.Name)

				err := server.Close(ctx)
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

				err = limiter.Close(ctx)
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
