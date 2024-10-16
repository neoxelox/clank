package main

import (
	"runtime"
	"time"

	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"
	"golang.org/x/text/language"

	"backend/pkg/config"
)

func NewConfig() *config.Config {
	config := config.NewConfig()

	config.Service.Environment = kit.Environment(util.GetEnv("CLANK_ENVIRONMENT", "dev"))
	config.Service.Name = "api"
	config.Service.Release = util.GetEnv("CLANK_RELEASE", "wip")
	config.Service.TimeZone = *time.UTC
	config.Service.Domain = util.GetEnv("CLANK_DOMAIN", "clank.localhost")
	config.Service.GracefulTimeout = 30 * time.Second
	config.Service.DefaultLocale = language.English
	config.Service.MigrationsPath = "migrations"
	config.Service.TemplatesPath = "templates"
	config.Service.TemplateFilePattern = `^.*\.(html|txt|md)$`
	config.Service.LocalesPath = "locales"
	config.Service.LocaleFilePattern = `^.*\.(yml|yaml)$`
	config.Service.AssetsPath = "assets"
	config.Service.FilesPath = "files"

	config.Database.Host = util.GetEnv("CLANK_DATABASE_HOST", "postgres")
	config.Database.Port = util.GetEnv("CLANK_DATABASE_PORT", 5432)
	config.Database.SSLMode = util.GetEnv("CLANK_DATABASE_SSLMODE", "disable")
	config.Database.User = util.GetEnv("CLANK_DATABASE_USER", "clank")
	config.Database.Password = util.GetEnv("CLANK_DATABASE_PASSWORD", "clank")
	config.Database.Name = util.GetEnv("CLANK_DATABASE_NAME", "clank")
	config.Database.SchemaVersion = 4
	config.Database.MinConns = 1
	config.Database.MaxConns = max(4, 2*runtime.GOMAXPROCS(-1))
	config.Database.MaxConnIdleTime = 30 * time.Minute
	config.Database.MaxConnLifeTime = 1 * time.Hour
	config.Database.DialTimeout = config.Service.GracefulTimeout
	config.Database.StatementTimeout = config.Service.GracefulTimeout
	config.Database.DefaultIsolationLevel = kit.IsoLvlReadCommitted

	config.Cache.Host = util.GetEnv("CLANK_CACHE_HOST", "redis")
	config.Cache.Port = util.GetEnv("CLANK_CACHE_PORT", 6379)
	config.Cache.SSLMode = util.GetEnv("CLANK_CACHE_SSLMODE", false)
	config.Cache.Password = util.GetEnv("CLANK_CACHE_PASSWORD", "redis")
	config.Cache.MinConns = 1
	config.Cache.MaxConns = max(8, 4*runtime.GOMAXPROCS(-1))
	config.Cache.MaxConnIdleTime = 30 * time.Minute
	config.Cache.MaxConnLifeTime = 1 * time.Hour
	config.Cache.ReadTimeout = config.Service.GracefulTimeout
	config.Cache.WriteTimeout = config.Service.GracefulTimeout
	config.Cache.DialTimeout = config.Service.GracefulTimeout

	config.Server.Host = util.GetEnv("CLANK_API_HOST", "localhost")
	config.Server.Port = util.GetEnv("CLANK_API_PORT", 1111)
	config.Server.BaseURL = util.GetEnv("CLANK_API_BASE_URL", "http://api.clank.localhost")
	config.Server.Origins = util.GetEnv("CLANK_API_ORIGINS", []string{"http://localhost:1111", "http://api.clank.localhost"})
	config.Server.RequestHeaderMaxSize = 1 << 10 // 1 KB
	config.Server.RequestBodyMaxSize = 50 << 10  // 50 KB
	config.Server.RequestFileMaxSize = 2 << 20   // 2 MB
	config.Server.RequestFilePattern = `.*/file.*`
	config.Server.RequestKeepAliveTimeout = config.Service.GracefulTimeout
	config.Server.RequestReadTimeout = config.Service.GracefulTimeout
	config.Server.RequestReadHeaderTimeout = config.Service.GracefulTimeout
	config.Server.ResponseWriteTimeout = config.Service.GracefulTimeout

	config.Sentry.DSN = util.GetEnv("CLANK_BACKEND_SENTRY_DSN", "")

	config.Gilk.Port = util.GetEnv("CLANK_API_GILK_PORT", 1113)

	config.Engine.BaseURL = util.GetEnv("CLANK_ENGINE_BASE_URL", "http://engine:2222")
	config.Frontend.BaseURL = util.GetEnv("CLANK_FRONTEND_BASE_URL", "http://clank.localhost")
	config.CDN.BaseURL = util.GetEnv("CLANK_CDN_BASE_URL", "http://cdn.clank.localhost")

	config.DataForSEO.BaseURL = util.GetEnv("CLANK_DATAFORSEO_BASE_URL", "")
	config.DataForSEO.APIKey = util.GetEnv("CLANK_DATAFORSEO_API_KEY", "")
	config.DataForSEO.CallbackSecret = util.GetEnv("CLANK_DATAFORSEO_CALLBACK_SECRET", "")

	config.Brevo.APIKey = util.GetEnv("CLANK_BREVO_API_KEY", "")
	config.Brevo.SenderEmail = util.GetEnv("CLANK_BREVO_SENDER_EMAIL", "")
	config.Brevo.SenderName = util.GetEnv("CLANK_BREVO_SENDER_NAME", "")
	config.Brevo.ReplierEmail = util.GetEnv("CLANK_BREVO_REPLIER_EMAIL", "")
	config.Brevo.ReplierName = util.GetEnv("CLANK_BREVO_REPLIER_NAME", "")

	config.Auth.CryptKey = util.GetEnv("CLANK_API_AUTH_CRYPT_KEY", "")
	config.Auth.GoogleID = util.GetEnv("CLANK_API_AUTH_GOOGLE_ID", "")
	config.Auth.GoogleSecret = util.GetEnv("CLANK_API_AUTH_GOOGLE_SECRET", "")
	config.Auth.AppleID = util.GetEnv("CLANK_API_AUTH_APPLE_ID", "")
	config.Auth.AppleSecret = util.GetEnv("CLANK_API_AUTH_APPLE_SECRET", "")
	config.Auth.AmazonID = util.GetEnv("CLANK_API_AUTH_AMAZON_ID", "")
	config.Auth.AmazonSecret = util.GetEnv("CLANK_API_AUTH_AMAZON_SECRET", "")

	return config
}
