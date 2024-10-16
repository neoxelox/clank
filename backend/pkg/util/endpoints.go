package util

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
)

type HealthEndpoints struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
	cache    *kit.Cache
}

func NewHealthEndpoints(observer *kit.Observer, database *kit.Database, cache *kit.Cache,
	config config.Config) *HealthEndpoints {
	return &HealthEndpoints{
		config:   config,
		observer: observer,
		database: database,
		cache:    cache,
	}
}

type HealthEndpointsGetHealthResponseItem struct {
	Error   *error `json:"error"`
	Latency int64  `json:"latency"`
}

type HealthEndpointsGetHealthResponse struct {
	Database HealthEndpointsGetHealthResponseItem `json:"database"`
	Cache    HealthEndpointsGetHealthResponseItem `json:"cache"`
}

func (self *HealthEndpoints) GetServerHealth(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	response := HealthEndpointsGetHealthResponse{}

	start := time.Now()
	errD := self.database.Health(requestCtx)
	response.Database.Latency = time.Since(start).Milliseconds()
	if errD != nil {
		response.Database.Error = &errD
	} else {
		response.Database.Error = nil
	}

	start = time.Now()
	errC := self.cache.Health(requestCtx)
	response.Cache.Latency = time.Since(start).Milliseconds()
	if errC != nil {
		response.Cache.Error = &errC
	} else {
		response.Cache.Error = nil
	}

	if errD != nil || errC != nil {
		return ctx.JSON(http.StatusServiceUnavailable, &response)
	}

	return ctx.JSON(http.StatusOK, &response)
}

func (self *HealthEndpoints) GetWorkerHealth(res http.ResponseWriter, req *http.Request) {
	requestCtx := req.Context()
	response := HealthEndpointsGetHealthResponse{}

	start := time.Now()
	errD := self.database.Health(requestCtx)
	response.Database.Latency = time.Since(start).Milliseconds()
	if errD != nil {
		response.Database.Error = &errD
	} else {
		response.Database.Error = nil
	}

	start = time.Now()
	errC := self.cache.Health(requestCtx)
	response.Cache.Latency = time.Since(start).Milliseconds()
	if errC != nil {
		response.Cache.Error = &errC
	} else {
		response.Cache.Error = nil
	}

	res.Header().Set("Content-Type", "application/json")
	if errD != nil || errC != nil {
		res.WriteHeader(http.StatusServiceUnavailable)
	} else {
		res.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(res).Encode(&response) // nolint:errcheck,errchkjson
}
