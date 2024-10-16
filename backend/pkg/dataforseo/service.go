package dataforseo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"

	"backend/pkg/config"
)

const (
	DATAFORSEO_SERVICE_TIMEOUT = 5 * time.Second
)

var (
	ErrDataForSEOServiceGeneric  = errors.New("dataforseo service failed")
	ErrDataForSEOServiceTimedOut = errors.New("dataforseo service timed out")
)

type DataForSEOService struct {
	config   config.Config
	observer *kit.Observer
	client   *kit.HTTPClient
}

func NewDataForSEOService(observer *kit.Observer, config config.Config) *DataForSEOService {
	client := kit.NewHTTPClient(observer, kit.HTTPClientConfig{
		Timeout: DATAFORSEO_SERVICE_TIMEOUT,
		BaseURL: util.Pointer(config.DataForSEO.BaseURL),
		Headers: util.Pointer(map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": "Basic " + config.DataForSEO.APIKey,
		}),
		RaiseForStatus:   util.Pointer(true),
		AllowedRedirects: util.Pointer(0),
		DefaultRetry:     nil,
	})

	return &DataForSEOService{
		config:   config,
		observer: observer,
		client:   client,
	}
}

type postTrustpilotTaskRequest struct {
	Domain      string  `json:"domain"`
	SortBy      string  `json:"sort_by"`
	Priority    int     `json:"priority"`
	Depth       int     `json:"depth"`
	Tag         string  `json:"tag"`
	PostbackURL *string `json:"postback_url,omitempty"`
	PingbackURL string  `json:"pingback_url"`
}

type postTrustpilotTasksResponse struct {
	response[trustpilotResponseTaskResult]
}

type DataForSEOServiceCreateTrustpilotTasksParams struct {
	Domain     string
	Reviews    int
	Prioritize bool
	Identifier string
	Callback   string
}

type DataForSEOServiceCreateTrustpilotTasksResult = []Task[TrustpilotReview]

func (self *DataForSEOService) CreateTrustpilotTasks(ctx context.Context,
	params DataForSEOServiceCreateTrustpilotTasksParams) (*DataForSEOServiceCreateTrustpilotTasksResult, error) {
	if self.config.Service.Environment != kit.EnvProduction {
		self.observer.Infof(ctx, "Created %d Trustpilot tasks for '%s'", 1, params.Domain)
		return &[]Task[TrustpilotReview]{}, nil
	}

	requestBody := make([]postTrustpilotTaskRequest, 1)
	requestBody[0].Domain = params.Domain
	requestBody[0].SortBy = "recency"
	requestBody[0].Priority = 1
	if params.Prioritize {
		requestBody[0].Priority = 2
	}
	requestBody[0].Depth = params.Reviews
	requestBody[0].Tag = params.Identifier
	requestBody[0].PingbackURL = params.Callback + "?id=$id"

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx,
		"POST", "/business_data/trustpilot/reviews/task_post", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postTrustpilotTasksResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	result := make(DataForSEOServiceCreateTrustpilotTasksResult, 0)
	for _, responseTask := range responseBody.Tasks {
		if !isResponseStatusCodeOk(responseTask.StatusCode) {
			self.observer.Error(ctx,
				ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
					Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID}))

			continue
		}

		result = append(result, Task[TrustpilotReview]{
			ID:         responseTask.ID,
			Status:     responseTask.StatusCode,
			Message:    responseTask.StatusMessage,
			Cost:       responseTask.Cost,
			Identifier: responseTask.Data.Tag,
			Reviews:    make([]TrustpilotReview, 0),
		})
	}

	return &result, nil
}

type getTrustpilotTaskResponse struct {
	response[trustpilotResponseTaskResult]
}

type DataForSEOServiceGetTrustpilotTaskParams struct {
	TaskID string
}

type DataForSEOServiceGetTrustpilotTaskResult = Task[TrustpilotReview]

func (self *DataForSEOService) GetTrustpilotTask(ctx context.Context,
	params DataForSEOServiceGetTrustpilotTaskParams) (*DataForSEOServiceGetTrustpilotTaskResult, error) {
	response, err := self.client.Request(ctx,
		"GET", "/business_data/trustpilot/reviews/task_get/"+params.TaskID, nil, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := getTrustpilotTaskResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	responseTask := responseBody.Tasks[0]
	if !isResponseStatusCodeOk(responseTask.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
			Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID})
	}

	result := DataForSEOServiceGetTrustpilotTaskResult{}
	result.ID = responseTask.ID
	result.Status = responseTask.StatusCode
	result.Message = responseTask.StatusMessage
	result.Cost = responseTask.Cost
	result.Identifier = responseTask.Data.Tag
	result.Reviews = make([]TrustpilotReview, 0)
	for _, taskResult := range responseTask.Result {
		for _, reviewItem := range taskResult.Items {
			timestamp, err := time.Parse("2006-01-02 15:04:05 -07:00", reviewItem.Timestamp)
			if err != nil {
				self.observer.Error(ctx,
					ErrDataForSEOServiceGeneric.Raise().With("cannot parse review timestamp").
						Extra(map[string]any{"timestamp": reviewItem.Timestamp, "task_id": responseTask.ID}))
			}

			result.Reviews = append(result.Reviews, TrustpilotReview{
				Page:    taskResult.CheckURL,
				Link:    reviewItem.URL,
				Title:   reviewItem.Title,
				Content: reviewItem.ReviewText,
				Images:  reviewItem.ReviewImages,
				Customer: TrustpilotCustomer{
					Link:     reviewItem.UserProfile.URL,
					Name:     reviewItem.UserProfile.Name,
					Picture:  reviewItem.UserProfile.ImageURL,
					Location: reviewItem.UserProfile.Location,
					Reviews:  reviewItem.UserProfile.ReviewsCount,
				},
				Rating:    reviewItem.Rating.Value,
				Votes:     reviewItem.Rating.VotesCount,
				Verified:  reviewItem.Verified,
				Timestamp: timestamp,
			})
		}
	}

	return &result, nil
}

type postPlayStoreTaskRequest struct {
	AppID        string  `json:"app_id"`
	LocationName string  `json:"location_name"`
	LocationCode *int    `json:"location_code,omitempty"`
	LanguageName string  `json:"language_name"`
	LanguageCode *string `json:"language_code,omitempty"`
	Priority     int     `json:"priority"`
	Depth        int     `json:"depth"`
	Rating       *int    `json:"rating,omitempty"`
	SortBy       string  `json:"sort_by"`
	Tag          string  `json:"tag"`
	PostbackURL  *string `json:"postback_url,omitempty"`
	PostbackData *string `json:"postback_data,omitempty"`
	PingbackURL  string  `json:"pingback_url"`
}

type postPlayStoreTasksResponse struct {
	response[playStoreResponseTaskResult]
}

type DataForSEOServiceCreatePlayStoreTasksParams struct {
	AppID        string
	Perspectives []Perspective
	Reviews      int
	Prioritize   bool
	Identifier   string
	Callback     string
}

type DataForSEOServiceCreatePlayStoreTasksResult = []Task[PlayStoreReview]

func (self *DataForSEOService) CreatePlayStoreTasks(ctx context.Context,
	params DataForSEOServiceCreatePlayStoreTasksParams) (*DataForSEOServiceCreatePlayStoreTasksResult, error) {
	if self.config.Service.Environment != kit.EnvProduction {
		self.observer.Infof(ctx, "Created %d Play Store tasks for '%s'", len(params.Perspectives), params.AppID)
		return &[]Task[PlayStoreReview]{}, nil
	}

	requestBody := make([]postPlayStoreTaskRequest, 0, len(params.Perspectives))
	priority := 1
	if params.Prioritize {
		priority = 2
	}
	for _, perspective := range params.Perspectives {
		requestBody = append(requestBody, postPlayStoreTaskRequest{
			AppID:        params.AppID,
			LocationName: perspective.Location,
			LanguageName: perspective.Language,
			Priority:     priority,
			Depth:        params.Reviews,
			SortBy:       "newest",
			Tag:          params.Identifier,
			PingbackURL:  params.Callback + "?id=$id",
		})
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx,
		"POST", "/app_data/google/app_reviews/task_post", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postPlayStoreTasksResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	result := make(DataForSEOServiceCreatePlayStoreTasksResult, 0)
	for _, responseTask := range responseBody.Tasks {
		if !isResponseStatusCodeOk(responseTask.StatusCode) {
			self.observer.Error(ctx,
				ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
					Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID}))

			continue
		}

		result = append(result, Task[PlayStoreReview]{
			ID:         responseTask.ID,
			Status:     responseTask.StatusCode,
			Message:    responseTask.StatusMessage,
			Cost:       responseTask.Cost,
			Identifier: responseTask.Data.Tag,
			Reviews:    make([]PlayStoreReview, 0),
		})
	}

	return &result, nil
}

type getPlayStoreTaskResponse struct {
	response[playStoreResponseTaskResult]
}

type DataForSEOServiceGetPlayStoreTaskParams struct {
	TaskID string
}

type DataForSEOServiceGetPlayStoreTaskResult = Task[PlayStoreReview]

func (self *DataForSEOService) GetPlayStoreTask(ctx context.Context,
	params DataForSEOServiceGetPlayStoreTaskParams) (*DataForSEOServiceGetPlayStoreTaskResult, error) {
	response, err := self.client.Request(ctx,
		"GET", "/app_data/google/app_reviews/task_get/advanced/"+params.TaskID, nil, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := getPlayStoreTaskResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	responseTask := responseBody.Tasks[0]
	if !isResponseStatusCodeOk(responseTask.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
			Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID})
	}

	result := DataForSEOServiceGetPlayStoreTaskResult{}
	result.ID = responseTask.ID
	result.Status = responseTask.StatusCode
	result.Message = responseTask.StatusMessage
	result.Cost = responseTask.Cost
	result.Identifier = responseTask.Data.Tag
	result.Reviews = make([]PlayStoreReview, 0)
	for _, taskResult := range responseTask.Result {
		for _, reviewItem := range taskResult.Items {
			timestamp, err := time.Parse("2006-01-02 15:04:05 -07:00", reviewItem.Timestamp)
			if err != nil {
				self.observer.Error(ctx,
					ErrDataForSEOServiceGeneric.Raise().With("cannot parse review timestamp").
						Extra(map[string]any{"timestamp": reviewItem.Timestamp, "task_id": responseTask.ID}))
			}

			result.Reviews = append(result.Reviews, PlayStoreReview{
				ID:      reviewItem.ID,
				Page:    taskResult.CheckURL,
				Title:   reviewItem.Title,
				Content: reviewItem.ReviewText,
				Customer: PlayStoreCustomer{
					Name:    reviewItem.UserProfile.ProfileName,
					Picture: reviewItem.UserProfile.ProfileImageURL,
				},
				Rating:    reviewItem.Rating.Value,
				Votes:     reviewItem.HelpfulCount,
				Release:   reviewItem.Version,
				Timestamp: timestamp,
			})
		}
	}

	return &result, nil
}

type postAppStoreTaskRequest struct {
	AppID        string  `json:"app_id"`
	LocationName string  `json:"location_name"`
	LocationCode *int    `json:"location_code,omitempty"`
	LanguageName string  `json:"language_name"`
	LanguageCode *string `json:"language_code,omitempty"`
	Priority     int     `json:"priority"`
	Depth        int     `json:"depth"`
	SortBy       string  `json:"sort_by"`
	Tag          string  `json:"tag"`
	PostbackURL  *string `json:"postback_url,omitempty"`
	PostbackData *string `json:"postback_data,omitempty"`
	PingbackURL  string  `json:"pingback_url"`
}

type postAppStoreTasksResponse struct {
	response[appStoreResponseTaskResult]
}

type DataForSEOServiceCreateAppStoreTasksParams struct {
	AppID        string
	Perspectives []Perspective
	Reviews      int
	Prioritize   bool
	Identifier   string
	Callback     string
}

type DataForSEOServiceCreateAppStoreTasksResult = []Task[AppStoreReview]

func (self *DataForSEOService) CreateAppStoreTasks(ctx context.Context,
	params DataForSEOServiceCreateAppStoreTasksParams) (*DataForSEOServiceCreateAppStoreTasksResult, error) {
	if self.config.Service.Environment != kit.EnvProduction {
		self.observer.Infof(ctx, "Created %d App Store tasks for '%s'", len(params.Perspectives), params.AppID)
		return &[]Task[AppStoreReview]{}, nil
	}

	requestBody := make([]postAppStoreTaskRequest, 0, len(params.Perspectives))
	priority := 1
	if params.Prioritize {
		priority = 2
	}
	for _, perspective := range params.Perspectives {
		requestBody = append(requestBody, postAppStoreTaskRequest{
			AppID:        params.AppID,
			LocationName: perspective.Location,
			LanguageName: perspective.Language,
			Priority:     priority,
			Depth:        params.Reviews,
			SortBy:       "most_recent",
			Tag:          params.Identifier,
			PingbackURL:  params.Callback + "?id=$id",
		})
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx,
		"POST", "/app_data/apple/app_reviews/task_post", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAppStoreTasksResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	result := make(DataForSEOServiceCreateAppStoreTasksResult, 0)
	for _, responseTask := range responseBody.Tasks {
		if !isResponseStatusCodeOk(responseTask.StatusCode) {
			self.observer.Error(ctx,
				ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
					Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID}))

			continue
		}

		result = append(result, Task[AppStoreReview]{
			ID:         responseTask.ID,
			Status:     responseTask.StatusCode,
			Message:    responseTask.StatusMessage,
			Cost:       responseTask.Cost,
			Identifier: responseTask.Data.Tag,
			Reviews:    make([]AppStoreReview, 0),
		})
	}

	return &result, nil
}

type getAppStoreTaskResponse struct {
	response[appStoreResponseTaskResult]
}

type DataForSEOServiceGetAppStoreTaskParams struct {
	TaskID string
}

type DataForSEOServiceGetAppStoreTaskResult = Task[AppStoreReview]

func (self *DataForSEOService) GetAppStoreTask(ctx context.Context,
	params DataForSEOServiceGetAppStoreTaskParams) (*DataForSEOServiceGetAppStoreTaskResult, error) {
	response, err := self.client.Request(ctx,
		"GET", "/app_data/apple/app_reviews/task_get/advanced/"+params.TaskID, nil, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := getAppStoreTaskResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	responseTask := responseBody.Tasks[0]
	if !isResponseStatusCodeOk(responseTask.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
			Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID})
	}

	result := DataForSEOServiceGetAppStoreTaskResult{}
	result.ID = responseTask.ID
	result.Status = responseTask.StatusCode
	result.Message = responseTask.StatusMessage
	result.Cost = responseTask.Cost
	result.Identifier = responseTask.Data.Tag
	result.Reviews = make([]AppStoreReview, 0)
	for _, taskResult := range responseTask.Result {
		for _, reviewItem := range taskResult.Items {
			timestamp, err := time.Parse("2006-01-02 15:04:05 -07:00", reviewItem.Timestamp)
			if err != nil {
				self.observer.Error(ctx,
					ErrDataForSEOServiceGeneric.Raise().With("cannot parse review timestamp").
						Extra(map[string]any{"timestamp": reviewItem.Timestamp, "task_id": responseTask.ID}))
			}

			result.Reviews = append(result.Reviews, AppStoreReview{
				ID:      reviewItem.ID,
				Page:    taskResult.CheckURL,
				Title:   reviewItem.Title,
				Content: reviewItem.ReviewText,
				Customer: AppStoreCustomer{
					Name:    reviewItem.UserProfile.ProfileName,
					Picture: reviewItem.UserProfile.ProfileImageURL,
				},
				Rating:    reviewItem.Rating.Value,
				Release:   reviewItem.Version,
				Timestamp: timestamp,
			})
		}
	}

	return &result, nil
}

type postAmazonTaskRequest struct {
	ASIN               string  `json:"asin"`
	Priority           int     `json:"priority"`
	LocationName       string  `json:"location_name"`
	LocationCode       *int    `json:"location_code,omitempty"`
	LocationCoordinate *string `json:"location_coordinate,omitempty"`
	LanguageName       string  `json:"language_name"`
	LanguageCode       *string `json:"language_code,omitempty"`
	SEDomain           *string `json:"se_domain,omitempty"`
	Depth              int     `json:"depth"`
	SortBy             string  `json:"sort_by"`
	ReviewerType       *string `json:"reviewer_type,omitempty"`
	FilterByStar       *string `json:"filter_by_star,omitempty"`
	FilterByKeyword    *string `json:"filter_by_keyword,omitempty"`
	MediaType          *string `json:"media_type,omitempty"`
	FormatType         *string `json:"format_type,omitempty"`
	Tag                string  `json:"tag"`
	PostbackURL        *string `json:"postback_url,omitempty"`
	PostbackData       *string `json:"postback_data,omitempty"`
	PingbackURL        string  `json:"pingback_url"`
}

type postAmazonTasksResponse struct {
	response[amazonResponseTaskResult]
}

type DataForSEOServiceCreateAmazonTasksParams struct {
	ASIN         string
	Perspectives []Perspective
	Reviews      int
	Prioritize   bool
	Identifier   string
	Callback     string
}

type DataForSEOServiceCreateAmazonTasksResult = []Task[AmazonReview]

func (self *DataForSEOService) CreateAmazonTasks(ctx context.Context,
	params DataForSEOServiceCreateAmazonTasksParams) (*DataForSEOServiceCreateAmazonTasksResult, error) {
	if self.config.Service.Environment != kit.EnvProduction {
		self.observer.Infof(ctx, "Created %d Amazon tasks for '%s'", len(params.Perspectives), params.ASIN)
		return &[]Task[AmazonReview]{}, nil
	}

	requestBody := make([]postAmazonTaskRequest, 0, len(params.Perspectives))
	priority := 1
	if params.Prioritize {
		priority = 2
	}
	for _, perspective := range params.Perspectives {
		requestBody = append(requestBody, postAmazonTaskRequest{
			ASIN:         params.ASIN,
			Priority:     priority,
			LocationName: perspective.Location,
			LanguageName: perspective.Language,
			Depth:        params.Reviews,
			SortBy:       "recent",
			Tag:          params.Identifier,
			PingbackURL:  params.Callback + "?id=$id",
		})
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx,
		"POST", "/merchant/amazon/reviews/task_post", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAmazonTasksResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	result := make(DataForSEOServiceCreateAmazonTasksResult, 0)
	for _, responseTask := range responseBody.Tasks {
		if !isResponseStatusCodeOk(responseTask.StatusCode) {
			self.observer.Error(ctx,
				ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
					Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID}))

			continue
		}

		result = append(result, Task[AmazonReview]{
			ID:         responseTask.ID,
			Status:     responseTask.StatusCode,
			Message:    responseTask.StatusMessage,
			Cost:       responseTask.Cost,
			Identifier: responseTask.Data.Tag,
			Reviews:    make([]AmazonReview, 0),
		})
	}

	return &result, nil
}

type getAmazonTaskResponse struct {
	response[amazonResponseTaskResult]
}

type DataForSEOServiceGetAmazonTaskParams struct {
	TaskID string
}

type DataForSEOServiceGetAmazonTaskResult = Task[AmazonReview]

func (self *DataForSEOService) GetAmazonTask(ctx context.Context,
	params DataForSEOServiceGetAmazonTaskParams) (*DataForSEOServiceGetAmazonTaskResult, error) {
	response, err := self.client.Request(ctx,
		"GET", "/merchant/amazon/reviews/task_get/advanced/"+params.TaskID, nil, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := getAmazonTaskResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrDataForSEOServiceGeneric.Raise().Cause(err)
	}

	if !isResponseStatusCodeOk(responseBody.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseBody.StatusMessage).
			Extra(map[string]any{"status_code": responseBody.StatusCode})
	}

	responseTask := responseBody.Tasks[0]
	if !isResponseStatusCodeOk(responseTask.StatusCode) {
		return nil, ErrDataForSEOServiceGeneric.Raise().With(responseTask.StatusMessage).
			Extra(map[string]any{"status_code": responseTask.StatusCode, "task_id": responseTask.ID})
	}

	result := DataForSEOServiceGetAmazonTaskResult{}
	result.ID = responseTask.ID
	result.Status = responseTask.StatusCode
	result.Message = responseTask.StatusMessage
	result.Cost = responseTask.Cost
	result.Identifier = responseTask.Data.Tag
	result.Reviews = make([]AmazonReview, 0)
	for _, taskResult := range responseTask.Result {
		for _, reviewItem := range taskResult.Items {
			timestamp, err := time.Parse("2006-01-02 15:04:05 -07:00", reviewItem.PublicationDate)
			if err != nil {
				self.observer.Error(ctx,
					ErrDataForSEOServiceGeneric.Raise().With("cannot parse review timestamp").
						Extra(map[string]any{"timestamp": reviewItem.PublicationDate, "task_id": responseTask.ID}))
			}

			images := make([]string, 0, len(reviewItem.Images))
			for _, reviewImage := range reviewItem.Images {
				images = append(images, reviewImage.ImageURL)
			}

			videos := make([]string, 0, len(reviewItem.Videos))
			for _, reviewVideo := range reviewItem.Videos {
				videos = append(videos, reviewVideo.Source)
			}

			result.Reviews = append(result.Reviews, AmazonReview{
				Page:     taskResult.CheckURL,
				Link:     reviewItem.URL,
				Title:    reviewItem.Title,
				Subtitle: reviewItem.Subtitle,
				Content:  reviewItem.ReviewText,
				Images:   images,
				Videos:   videos,
				Customer: AmazonCustomer{
					Link:     reviewItem.UserProfile.URL,
					Name:     reviewItem.UserProfile.Name,
					Picture:  reviewItem.UserProfile.Avatar,
					Location: reviewItem.UserProfile.Locations,
					Reviews:  reviewItem.UserProfile.ReviewsCount,
				},
				Rating:    reviewItem.Rating.Value,
				Votes:     reviewItem.HelpfulVotes,
				Verified:  reviewItem.Verified,
				Timestamp: timestamp,
			})
		}
	}

	return &result, nil
}

func (self *DataForSEOService) Close(ctx context.Context) error {
	err := util.Deadline(ctx, func(exceeded <-chan struct{}) error {
		self.observer.Info(ctx, "Closing DataForSEO service")

		err := self.client.Close(ctx)
		if err != nil {
			return ErrDataForSEOServiceGeneric.Raise().Cause(err)
		}

		self.observer.Info(ctx, "Closed DataForSEO service")

		return nil
	})
	if err != nil {
		if util.ErrDeadlineExceeded.Is(err) {
			return ErrDataForSEOServiceTimedOut.Raise().Cause(err)
		}

		return err
	}

	return nil
}
