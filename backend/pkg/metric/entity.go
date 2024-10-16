package metric

import "time"

type Params struct {
	ProductID     string
	PeriodStartAt *time.Time
	PeriodEndAt   *time.Time
}

type Metric struct {
}

type IssueCountParams struct {
	Params
}

type IssueCountMetric struct {
	Metric
	ActiveIssues   int
	ArchivedIssues int
	NewIssues      int
	Feedbacks      int
}

type IssueSourcesParams struct {
	Params
}

type IssueSourcesMetric struct {
	Metric
	Sources map[string]int
}

type IssueSeveritiesParams struct {
	Params
}

type IssueSeveritiesMetric struct {
	Metric
	Severities map[string]int
}

type IssueCategoriesParams struct {
	Params
}

type IssueCategoriesMetric struct {
	Metric
	Categories map[string]int
}

type IssueReleasesParams struct {
	Params
}

type IssueReleasesMetric struct {
	Metric
	Releases map[string]int
}

type SuggestionCountParams struct {
	Params
}

type SuggestionCountMetric struct {
	Metric
	ActiveSuggestions   int
	ArchivedSuggestions int
	NewSuggestions      int
	Feedbacks           int
}

type SuggestionSourcesParams struct {
	Params
}

type SuggestionSourcesMetric struct {
	Metric
	Sources map[string]int
}

type SuggestionImportancesParams struct {
	Params
}

type SuggestionImportancesMetric struct {
	Metric
	Importances map[string]int
}

type SuggestionCategoriesParams struct {
	Params
}

type SuggestionCategoriesMetric struct {
	Metric
	Categories map[string]int
}

type SuggestionReleasesParams struct {
	Params
}

type SuggestionReleasesMetric struct {
	Metric
	Releases map[string]int
}

type ReviewSentimentsParams struct {
	Params
}

type ReviewSentimentsMetric struct {
	Metric
	Sentiments map[string]int
}

type ReviewSourcesParams struct {
	Params
}

type ReviewSourcesMetric struct {
	Metric
	Sources map[string]int
}

type ReviewIntentionsParams struct {
	Params
}

type ReviewIntentionsMetric struct {
	Metric
	Intentions map[string]int
}

type ReviewEmotionsParams struct {
	Params
}

type ReviewEmotionsMetric struct {
	Metric
	Emotions map[string]int
}

type ReviewCategoriesParams struct {
	Params
}

type ReviewCategoriesMetric struct {
	Metric
	Categories map[string]int
}

type ReviewReleasesParams struct {
	Params
}

type ReviewReleasesMetric struct {
	Metric
	Releases map[string]int
}

const (
	METRIC_REVIEW_KEYWORDS_LIMIT = 50
)

type ReviewKeywordsParams struct {
	Params
}

type ReviewKeywordsMetric struct {
	Metric
	Positive map[string]int
	Neutral  map[string]int
	Negative map[string]int
}

type NetPromoterScoreParams struct {
	Params
}

type NetPromoterScoreMetric struct {
	Metric
	Score float64
}

type CustomerSatisfactionScoreParams struct {
	Params
}

type CustomerSatisfactionScoreMetric struct {
	Metric
	Score float64
}
