package engine

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
	ENGINE_SERVICE_TIMEOUT = 59 * time.Second
)

var (
	ErrEngineServiceGeneric  = errors.New("engine service failed")
	ErrEngineServiceTimedOut = errors.New("engine service timed out")
)

type EngineService struct {
	config   config.Config
	observer *kit.Observer
	client   *kit.HTTPClient
}

func NewEngineService(observer *kit.Observer, config config.Config) *EngineService {
	client := kit.NewHTTPClient(observer, kit.HTTPClientConfig{
		Timeout: ENGINE_SERVICE_TIMEOUT,
		BaseURL: util.Pointer(config.Engine.BaseURL),
		Headers: util.Pointer(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}),
		RaiseForStatus:   util.Pointer(true),
		AllowedRedirects: util.Pointer(0),
		DefaultRetry:     nil,
	})

	return &EngineService{
		config:   config,
		observer: observer,
		client:   client,
	}
}

type postTranslatorDetectLanguageRequest struct {
	Feedback string `json:"feedback"`
}

type postTranslatorDetectLanguageResponse struct {
	Language string `json:"language"`
	Usage    struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceDetectLanguageParams struct {
	Feedback Feedback
}

type EngineServiceDetectLanguageResult struct {
	Language string
	Usage    Usage
}

func (self *EngineService) DetectLanguage(ctx context.Context,
	params EngineServiceDetectLanguageParams) (*EngineServiceDetectLanguageResult, error) {
	requestBody := postTranslatorDetectLanguageRequest{}
	requestBody.Feedback = params.Feedback.Content

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/translator/detect-language", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postTranslatorDetectLanguageResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceDetectLanguageResult{}
	result.Language = responseBody.Language
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postTranslatorTranslateFeedbackRequest struct {
	Feedback     string `json:"feedback"`
	FromLanguage string `json:"from_language"`
	ToLanguage   string `json:"to_language"`
}

type postTranslatorTranslateFeedbackResponse struct {
	Translation string `json:"translation"`
	Usage       struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceTranslateFeedbackParams struct {
	Feedback     Feedback
	FromLanguage string
	ToLanguage   string
}

type EngineServiceTranslateFeedbackResult struct {
	Translation string
	Usage       Usage
}

func (self *EngineService) TranslateFeedback(ctx context.Context,
	params EngineServiceTranslateFeedbackParams) (*EngineServiceTranslateFeedbackResult, error) {
	requestBody := postTranslatorTranslateFeedbackRequest{}
	requestBody.Feedback = params.Feedback.Content
	requestBody.FromLanguage = params.FromLanguage
	requestBody.ToLanguage = params.ToLanguage

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/translator/translate-feedback", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postTranslatorTranslateFeedbackResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceTranslateFeedbackResult{}
	result.Translation = responseBody.Translation
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postProcessorExtractIssuesRequest struct {
	Context    string   `json:"context"`
	Categories []string `json:"categories"`
	Feedback   string   `json:"feedback"`
}

type postProcessorExtractIssuesResponse struct {
	Issues []struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Steps       []string `json:"steps"`
		Severity    string   `json:"severity"`
		Category    string   `json:"category"`
	} `json:"issues"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceExtractIssuesParams struct {
	Context    string
	Categories []string
	Feedback   Feedback
}

type EngineServiceExtractIssuesResult struct {
	Issues []Issue
	Usage  Usage
}

func (self *EngineService) ExtractIssues(ctx context.Context,
	params EngineServiceExtractIssuesParams) (*EngineServiceExtractIssuesResult, error) {
	requestBody := postProcessorExtractIssuesRequest{}
	requestBody.Context = params.Context
	requestBody.Categories = params.Categories
	requestBody.Feedback = params.Feedback.Content

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/processor/extract-issues", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postProcessorExtractIssuesResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceExtractIssuesResult{}
	result.Issues = make([]Issue, 0, len(responseBody.Issues))
	for _, issue := range responseBody.Issues {
		result.Issues = append(result.Issues, Issue{
			Title:       issue.Title,
			Description: issue.Description,
			Steps:       issue.Steps,
			Severity:    issue.Severity,
			Category:    issue.Category,
		})
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postProcessorExtractSuggestionsRequest struct {
	Context    string   `json:"context"`
	Categories []string `json:"categories"`
	Feedback   string   `json:"feedback"`
}

type postProcessorExtractSuggestionsResponse struct {
	Suggestions []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Reason      string `json:"reason"`
		Importance  string `json:"importance"`
		Category    string `json:"category"`
	} `json:"suggestions"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceExtractSuggestionsParams struct {
	Context    string
	Categories []string
	Feedback   Feedback
}

type EngineServiceExtractSuggestionsResult struct {
	Suggestions []Suggestion
	Usage       Usage
}

func (self *EngineService) ExtractSuggestions(ctx context.Context,
	params EngineServiceExtractSuggestionsParams) (*EngineServiceExtractSuggestionsResult, error) {
	requestBody := postProcessorExtractSuggestionsRequest{}
	requestBody.Context = params.Context
	requestBody.Categories = params.Categories
	requestBody.Feedback = params.Feedback.Content

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/processor/extract-suggestions", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postProcessorExtractSuggestionsResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceExtractSuggestionsResult{}
	result.Suggestions = make([]Suggestion, 0, len(responseBody.Suggestions))
	for _, suggestion := range responseBody.Suggestions {
		result.Suggestions = append(result.Suggestions, Suggestion{
			Title:       suggestion.Title,
			Description: suggestion.Description,
			Reason:      suggestion.Reason,
			Importance:  suggestion.Importance,
			Category:    suggestion.Category,
		})
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postProcessorExtractReviewRequest struct {
	Context    string   `json:"context"`
	Categories []string `json:"categories"`
	Feedback   string   `json:"feedback"`
}

type postProcessorExtractReviewResponse struct {
	Review struct {
		Content   string   `json:"content"`
		Keywords  []string `json:"keywords"`
		Sentiment string   `json:"sentiment"`
		Emotions  []string `json:"emotions"`
		Intention string   `json:"intention"`
		Category  string   `json:"category"`
	} `json:"review"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceExtractReviewParams struct {
	Context    string
	Categories []string
	Feedback   Feedback
}

type EngineServiceExtractReviewResult struct {
	Review Review
	Usage  Usage
}

func (self *EngineService) ExtractReview(ctx context.Context,
	params EngineServiceExtractReviewParams) (*EngineServiceExtractReviewResult, error) {
	requestBody := postProcessorExtractReviewRequest{}
	requestBody.Context = params.Context
	requestBody.Categories = params.Categories
	requestBody.Feedback = params.Feedback.Content

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/processor/extract-review", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postProcessorExtractReviewResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceExtractReviewResult{}
	result.Review = Review{
		Content:   responseBody.Review.Content,
		Keywords:  responseBody.Review.Keywords,
		Sentiment: responseBody.Review.Sentiment,
		Emotions:  responseBody.Review.Emotions,
		Intention: responseBody.Review.Intention,
		Category:  responseBody.Review.Category,
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postAggregatorComputeEmbeddingRequest struct {
	Text string `json:"text"`
}

type postAggregatorComputeEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
	Usage     struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceComputeEmbeddingParams struct {
	Text string
}

type EngineServiceComputeEmbeddingResult struct {
	Embedding []float32
	Usage     Usage
}

func (self *EngineService) ComputeEmbedding(ctx context.Context,
	params EngineServiceComputeEmbeddingParams) (*EngineServiceComputeEmbeddingResult, error) {
	requestBody := postAggregatorComputeEmbeddingRequest{}
	requestBody.Text = params.Text

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/aggregator/compute-embedding", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAggregatorComputeEmbeddingResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceComputeEmbeddingResult{}
	result.Embedding = responseBody.Embedding
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postAggregatorSimilarIssueRequest struct {
	Issue   string   `json:"issue"`
	Options []string `json:"options"`
}

type postAggregatorSimilarIssueResponse struct {
	Option int `json:"option"`
	Usage  struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceSimilarIssueParams struct {
	Issue   Issue
	Options []Issue
}

type EngineServiceSimilarIssueResult struct {
	Option *int
	Usage  Usage
}

func (self *EngineService) SimilarIssue(ctx context.Context,
	params EngineServiceSimilarIssueParams) (*EngineServiceSimilarIssueResult, error) {
	requestBody := postAggregatorSimilarIssueRequest{}
	requestBody.Issue = params.Issue.Description
	requestBody.Options = make([]string, 0, len(params.Options))
	for _, issue := range params.Options {
		requestBody.Options = append(requestBody.Options, issue.Description)
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/aggregator/similar-issue", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAggregatorSimilarIssueResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceSimilarIssueResult{}
	result.Option = nil
	if responseBody.Option > 0 {
		result.Option = util.Pointer(responseBody.Option - 1)
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postAggregatorMergeIssuesRequest struct {
	IssueA struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Steps       []string `json:"steps"`
	} `json:"issue_a"`
	IssueB struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Steps       []string `json:"steps"`
	} `json:"issue_b"`
}

type postAggregatorMergeIssuesResponse struct {
	Issue struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Steps       []string `json:"steps"`
	} `json:"issue"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceMergeIssuesParams struct {
	IssueA Issue
	IssueB Issue
}

type EngineServiceMergeIssuesResult struct {
	Issue Issue
	Usage Usage
}

func (self *EngineService) MergeIssues(ctx context.Context,
	params EngineServiceMergeIssuesParams) (*EngineServiceMergeIssuesResult, error) {
	requestBody := postAggregatorMergeIssuesRequest{}
	requestBody.IssueA.Title = params.IssueA.Title
	requestBody.IssueA.Description = params.IssueA.Description
	requestBody.IssueA.Steps = params.IssueA.Steps
	requestBody.IssueB.Title = params.IssueB.Title
	requestBody.IssueB.Description = params.IssueB.Description
	requestBody.IssueB.Steps = params.IssueB.Steps

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/aggregator/merge-issues", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAggregatorMergeIssuesResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceMergeIssuesResult{}
	result.Issue = Issue{
		Title:       responseBody.Issue.Title,
		Description: responseBody.Issue.Description,
		Steps:       responseBody.Issue.Steps,
		Severity:    params.IssueB.Severity,
		Category:    params.IssueB.Category,
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postAggregatorSimilarSuggestionRequest struct {
	Suggestion string   `json:"suggestion"`
	Options    []string `json:"options"`
}

type postAggregatorSimilarSuggestionResponse struct {
	Option int `json:"option"`
	Usage  struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceSimilarSuggestionParams struct {
	Suggestion Suggestion
	Options    []Suggestion
}

type EngineServiceSimilarSuggestionResult struct {
	Option *int
	Usage  Usage
}

func (self *EngineService) SimilarSuggestion(ctx context.Context,
	params EngineServiceSimilarSuggestionParams) (*EngineServiceSimilarSuggestionResult, error) {
	requestBody := postAggregatorSimilarSuggestionRequest{}
	requestBody.Suggestion = params.Suggestion.Description
	requestBody.Options = make([]string, 0, len(params.Options))
	for _, suggestion := range params.Options {
		requestBody.Options = append(requestBody.Options, suggestion.Description)
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/aggregator/similar-suggestion", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAggregatorSimilarSuggestionResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceSimilarSuggestionResult{}
	result.Option = nil
	if responseBody.Option > 0 {
		result.Option = util.Pointer(responseBody.Option - 1)
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

type postAggregatorMergeSuggestionsRequest struct {
	SuggestionA struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Reason      string `json:"reason"`
	} `json:"suggestion_a"`
	SuggestionB struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Reason      string `json:"reason"`
	} `json:"suggestion_b"`
}

type postAggregatorMergeSuggestionsResponse struct {
	Suggestion struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Reason      string `json:"reason"`
	} `json:"suggestion"`
	Usage struct {
		Input  int `json:"input"`
		Output int `json:"output"`
	} `json:"usage"`
}

type EngineServiceMergeSuggestionsParams struct {
	SuggestionA Suggestion
	SuggestionB Suggestion
}

type EngineServiceMergeSuggestionsResult struct {
	Suggestion Suggestion
	Usage      Usage
}

func (self *EngineService) MergeSuggestions(ctx context.Context,
	params EngineServiceMergeSuggestionsParams) (*EngineServiceMergeSuggestionsResult, error) {
	requestBody := postAggregatorMergeSuggestionsRequest{}
	requestBody.SuggestionA.Title = params.SuggestionA.Title
	requestBody.SuggestionA.Description = params.SuggestionA.Description
	requestBody.SuggestionA.Reason = params.SuggestionA.Reason
	requestBody.SuggestionB.Title = params.SuggestionB.Title
	requestBody.SuggestionB.Description = params.SuggestionB.Description
	requestBody.SuggestionB.Reason = params.SuggestionB.Reason

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	response, err := self.client.Request(ctx, "POST", "/aggregator/merge-suggestions", requestBodyJSON, nil)
	if err != nil {
		if kit.ErrHTTPClientTimedOut.Is(err) {
			return nil, ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}
	defer response.Body.Close()

	responseBody := postAggregatorMergeSuggestionsResponse{}

	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, ErrEngineServiceGeneric.Raise().Cause(err)
	}

	result := EngineServiceMergeSuggestionsResult{}
	result.Suggestion = Suggestion{
		Title:       responseBody.Suggestion.Title,
		Description: responseBody.Suggestion.Description,
		Reason:      responseBody.Suggestion.Reason,
		Importance:  params.SuggestionB.Importance,
		Category:    params.SuggestionB.Category,
	}
	result.Usage = Usage{
		Input:  responseBody.Usage.Input,
		Output: responseBody.Usage.Output,
	}

	return &result, nil
}

func (self *EngineService) Close(ctx context.Context) error {
	err := util.Deadline(ctx, func(exceeded <-chan struct{}) error {
		self.observer.Info(ctx, "Closing Engine service")

		err := self.client.Close(ctx)
		if err != nil {
			return ErrEngineServiceGeneric.Raise().Cause(err)
		}

		self.observer.Info(ctx, "Closed Engine service")

		return nil
	})
	if err != nil {
		if util.ErrDeadlineExceeded.Is(err) {
			return ErrEngineServiceTimedOut.Raise().Cause(err)
		}

		return err
	}

	return nil
}
