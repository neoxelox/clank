package metric

import (
	"context"
	"math"

	"github.com/leporo/sqlf"
	"github.com/neoxelox/kit"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/feedback"
	"backend/pkg/issue"
	"backend/pkg/review"
	"backend/pkg/suggestion"
)

type MetricRepository struct {
	config   config.Config
	observer *kit.Observer
	database *kit.Database
}

func NewMetricRepository(observer *kit.Observer, database *kit.Database, config config.Config) *MetricRepository {
	return &MetricRepository{
		config:   config,
		observer: observer,
		database: database,
	}
}

func (self *MetricRepository) GetIssueCount(ctx context.Context,
	params IssueCountParams) (*IssueCountMetric, error) {
	var metric IssueCountMetric

	stmt := sqlf.
		Select("COUNT(*) FILTER (WHERE archived_at IS NULL)").To(&metric.ActiveIssues).
		Select("COUNT(*) FILTER (WHERE archived_at IS NOT NULL)").To(&metric.ArchivedIssues)

	if params.PeriodStartAt != nil {
		stmt.
			Select("COUNT(*) FILTER (WHERE first_seen_at >= ?)", *params.PeriodStartAt).To(&metric.NewIssues)
	}

	stmt.
		From(issue.ISSUE_MODEL_TABLE).
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	// Note that the following query+logic does not output the real number of feedbacks
	// but the number necessary for the (deduplicated) issue ratio to make sense!
	stmt = sqlf.
		Select("COUNT(*)").To(&metric.Feedbacks).
		From(feedback.FEEDBACK_MODEL_TABLE).
		LeftJoin(issue.ISSUE_FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+issue.ISSUE_FEEDBACK_MODEL_TABLE+".feedback_id").
		Where(feedback.FEEDBACK_MODEL_TABLE+".product_id = ?", params.ProductID).
		Where(issue.ISSUE_FEEDBACK_MODEL_TABLE + ".issue_id IS NULL")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err = self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	// Note that new issues are already counted either in active or archived!
	metric.Feedbacks += (metric.ActiveIssues + metric.ArchivedIssues)

	return &metric, nil
}

func (self *MetricRepository) GetIssueSources(ctx context.Context,
	params IssueSourcesParams) (*IssueSourcesMetric, error) {
	var result []struct {
		Source string `db:"key"`
		Count  int    `db:"sum"`
	}
	metric := IssueSourcesMetric{}
	metric.Sources = make(map[string]int)

	stmt := sqlf.
		Select("source.key, SUM(source.value::int)").To(&result).
		From(issue.ISSUE_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(sources) as source").
		Where("product_id = ?", params.ProductID).
		GroupBy("source.key")

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Sources[res.Source] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetIssueSeverities(ctx context.Context,
	params IssueSeveritiesParams) (*IssueSeveritiesMetric, error) {
	var result []struct {
		Severity string `db:"severity"`
		Count    int    `db:"count"`
	}
	metric := IssueSeveritiesMetric{}
	metric.Severities = make(map[string]int)

	inner := sqlf.
		Select("key as severity").
		Select("ROW_NUMBER() OVER (PARTITION BY id ORDER BY value::int DESC) AS rank").
		From(issue.ISSUE_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(severities)").
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		inner.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		inner.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	outer := sqlf.
		Select("severity, COUNT(*)").To(&result).
		From("").
		SubQuery("(", ")", inner).
		Where("rank = 1").
		GroupBy("severity")

	err := self.database.Query(ctx, outer)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Severities[res.Severity] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetIssueCategories(ctx context.Context,
	params IssueCategoriesParams) (*IssueCategoriesMetric, error) {
	var result []struct {
		Category string `db:"category"`
		Count    int    `db:"count"`
	}
	metric := IssueCategoriesMetric{}
	metric.Categories = make(map[string]int)

	inner := sqlf.
		Select("key as category").
		Select("ROW_NUMBER() OVER (PARTITION BY id ORDER BY value::int DESC) AS rank").
		From(issue.ISSUE_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(categories)").
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		inner.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		inner.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	outer := sqlf.
		Select("category, COUNT(*)").To(&result).
		From("").
		SubQuery("(", ")", inner).
		Where("rank = 1").
		GroupBy("category")

	err := self.database.Query(ctx, outer)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Categories[res.Category] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetIssueReleases(ctx context.Context,
	params IssueReleasesParams) (*IssueReleasesMetric, error) {
	var result []struct {
		Release string `db:"key"`
		Count   int    `db:"sum"`
	}
	metric := IssueReleasesMetric{}
	metric.Releases = make(map[string]int)

	stmt := sqlf.
		Select("release.key, SUM(release.value::int)").To(&result).
		From(issue.ISSUE_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(releases) as release").
		Where("product_id = ?", params.ProductID).
		GroupBy("release.key")

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Releases[res.Release] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetSuggestionCount(ctx context.Context,
	params SuggestionCountParams) (*SuggestionCountMetric, error) {
	var metric SuggestionCountMetric

	stmt := sqlf.
		Select("COUNT(*) FILTER (WHERE archived_at IS NULL)").To(&metric.ActiveSuggestions).
		Select("COUNT(*) FILTER (WHERE archived_at IS NOT NULL)").To(&metric.ArchivedSuggestions)

	if params.PeriodStartAt != nil {
		stmt.
			Select("COUNT(*) FILTER (WHERE first_seen_at >= ?)", *params.PeriodStartAt).To(&metric.NewSuggestions)
	}

	stmt.
		From(suggestion.SUGGESTION_MODEL_TABLE).
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	// Note that the following query+logic does not output the real number of feedbacks
	// but the number necessary for the (deduplicated) suggestion ratio to make sense!
	stmt = sqlf.
		Select("COUNT(*)").To(&metric.Feedbacks).
		From(feedback.FEEDBACK_MODEL_TABLE).
		LeftJoin(suggestion.SUGGESTION_FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+suggestion.SUGGESTION_FEEDBACK_MODEL_TABLE+".feedback_id").
		Where(feedback.FEEDBACK_MODEL_TABLE+".product_id = ?", params.ProductID).
		Where(suggestion.SUGGESTION_FEEDBACK_MODEL_TABLE + ".suggestion_id IS NULL")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err = self.database.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	// Note that new suggestions are already counted either in active or archived!
	metric.Feedbacks += (metric.ActiveSuggestions + metric.ArchivedSuggestions)

	return &metric, nil
}

func (self *MetricRepository) GetSuggestionSources(ctx context.Context,
	params SuggestionSourcesParams) (*SuggestionSourcesMetric, error) {
	var result []struct {
		Source string `db:"key"`
		Count  int    `db:"sum"`
	}
	metric := SuggestionSourcesMetric{}
	metric.Sources = make(map[string]int)

	stmt := sqlf.
		Select("source.key, SUM(source.value::int)").To(&result).
		From(suggestion.SUGGESTION_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(sources) as source").
		Where("product_id = ?", params.ProductID).
		GroupBy("source.key")

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Sources[res.Source] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetSuggestionImportances(ctx context.Context,
	params SuggestionImportancesParams) (*SuggestionImportancesMetric, error) {
	var result []struct {
		Importance string `db:"importance"`
		Count      int    `db:"count"`
	}
	metric := SuggestionImportancesMetric{}
	metric.Importances = make(map[string]int)

	inner := sqlf.
		Select("key as importance").
		Select("ROW_NUMBER() OVER (PARTITION BY id ORDER BY value::int DESC) AS rank").
		From(suggestion.SUGGESTION_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(importances)").
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		inner.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		inner.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	outer := sqlf.
		Select("importance, COUNT(*)").To(&result).
		From("").
		SubQuery("(", ")", inner).
		Where("rank = 1").
		GroupBy("importance")

	err := self.database.Query(ctx, outer)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Importances[res.Importance] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetSuggestionCategories(ctx context.Context,
	params SuggestionCategoriesParams) (*SuggestionCategoriesMetric, error) {
	var result []struct {
		Category string `db:"category"`
		Count    int    `db:"count"`
	}
	metric := SuggestionCategoriesMetric{}
	metric.Categories = make(map[string]int)

	inner := sqlf.
		Select("key as category").
		Select("ROW_NUMBER() OVER (PARTITION BY id ORDER BY value::int DESC) AS rank").
		From(suggestion.SUGGESTION_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(categories)").
		Where("product_id = ?", params.ProductID)

	if params.PeriodStartAt != nil {
		inner.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		inner.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	outer := sqlf.
		Select("category, COUNT(*)").To(&result).
		From("").
		SubQuery("(", ")", inner).
		Where("rank = 1").
		GroupBy("category")

	err := self.database.Query(ctx, outer)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Categories[res.Category] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetSuggestionReleases(ctx context.Context,
	params SuggestionReleasesParams) (*SuggestionReleasesMetric, error) {
	var result []struct {
		Release string `db:"key"`
		Count   int    `db:"sum"`
	}
	metric := SuggestionReleasesMetric{}
	metric.Releases = make(map[string]int)

	stmt := sqlf.
		Select("release.key, SUM(release.value::int)").To(&result).
		From(suggestion.SUGGESTION_MODEL_TABLE).
		Clause(", LATERAL jsonb_each_text(releases) as release").
		Where("product_id = ?", params.ProductID).
		GroupBy("release.key")

	if params.PeriodStartAt != nil {
		stmt.
			Where("last_seen_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("first_seen_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Releases[res.Release] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewSentiments(ctx context.Context,
	params ReviewSentimentsParams) (*ReviewSentimentsMetric, error) {
	var result []struct {
		Sentiment string `db:"sentiment"`
		Count     int    `db:"count"`
	}
	metric := ReviewSentimentsMetric{}
	metric.Sentiments = make(map[string]int)

	stmt := sqlf.
		Select(review.REVIEW_MODEL_TABLE+".sentiment, COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy(review.REVIEW_MODEL_TABLE + ".sentiment")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Sentiments[res.Sentiment] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewSources(ctx context.Context,
	params ReviewSourcesParams) (*ReviewSourcesMetric, error) {
	var result []struct {
		Source string `db:"source"`
		Count  int    `db:"count"`
	}
	metric := ReviewSourcesMetric{}
	metric.Sources = make(map[string]int)

	stmt := sqlf.
		Select("source, COUNT(*)").To(&result).
		From(feedback.FEEDBACK_MODEL_TABLE).
		Where("product_id = ?", params.ProductID).
		GroupBy("source")

	if params.PeriodStartAt != nil {
		stmt.
			Where("posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Sources[res.Source] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewIntentions(ctx context.Context,
	params ReviewIntentionsParams) (*ReviewIntentionsMetric, error) {
	var result []struct {
		Intention string `db:"intention"`
		Count     int    `db:"count"`
	}
	metric := ReviewIntentionsMetric{}
	metric.Intentions = make(map[string]int)

	stmt := sqlf.
		Select(review.REVIEW_MODEL_TABLE+".intention, COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy(review.REVIEW_MODEL_TABLE + ".intention")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Intentions[res.Intention] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewEmotions(ctx context.Context,
	params ReviewEmotionsParams) (*ReviewEmotionsMetric, error) {
	var result []struct {
		Emotion string `db:"unnest"`
		Count   int    `db:"count"`
	}
	metric := ReviewEmotionsMetric{}
	metric.Emotions = make(map[string]int)

	stmt := sqlf.
		Select("UNNEST("+review.REVIEW_MODEL_TABLE+".emotions), COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy("UNNEST(" + review.REVIEW_MODEL_TABLE + ".emotions)")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Emotions[res.Emotion] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewCategories(ctx context.Context,
	params ReviewCategoriesParams) (*ReviewCategoriesMetric, error) {
	var result []struct {
		Category string `db:"category"`
		Count    int    `db:"count"`
	}
	metric := ReviewCategoriesMetric{}
	metric.Categories = make(map[string]int)

	stmt := sqlf.
		Select(review.REVIEW_MODEL_TABLE+".category, COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy(review.REVIEW_MODEL_TABLE + ".category")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Categories[res.Category] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewReleases(ctx context.Context,
	params ReviewReleasesParams) (*ReviewReleasesMetric, error) {
	var result []struct {
		Release string `db:"release"`
		Count   int    `db:"count"`
	}
	metric := ReviewReleasesMetric{}
	metric.Releases = make(map[string]int)

	stmt := sqlf.
		Select("release, COUNT(*)").To(&result).
		From(feedback.FEEDBACK_MODEL_TABLE).
		Where("product_id = ?", params.ProductID).
		GroupBy("release")

	if params.PeriodStartAt != nil {
		stmt.
			Where("posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where("posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	for _, res := range result {
		metric.Releases[res.Release] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetReviewKeywords(ctx context.Context,
	params ReviewKeywordsParams) (*ReviewKeywordsMetric, error) {
	var result []struct {
		Keyword string `db:"unnest"`
		Count   int    `db:"count"`
	}
	metric := ReviewKeywordsMetric{}
	metric.Positive = make(map[string]int)

	stmt := sqlf.
		Select("UNNEST("+review.REVIEW_MODEL_TABLE+".keywords), COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		Where(review.REVIEW_MODEL_TABLE + ".sentiment = 'POSITIVE'").
		GroupBy("UNNEST(" + review.REVIEW_MODEL_TABLE + ".keywords)").
		OrderBy("COUNT(*) DESC").
		Limit(METRIC_REVIEW_KEYWORDS_LIMIT)

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil && !kit.ErrDatabaseNoRows.Is(err) {
		return nil, err
	}

	for _, res := range result {
		metric.Positive[res.Keyword] = res.Count
	}

	result = []struct {
		Keyword string `db:"unnest"`
		Count   int    `db:"count"`
	}{}
	metric.Neutral = make(map[string]int)

	stmt = sqlf.
		Select("UNNEST("+review.REVIEW_MODEL_TABLE+".keywords), COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		Where(review.REVIEW_MODEL_TABLE + ".sentiment = 'NEUTRAL'").
		GroupBy("UNNEST(" + review.REVIEW_MODEL_TABLE + ".keywords)").
		OrderBy("COUNT(*) DESC").
		Limit(METRIC_REVIEW_KEYWORDS_LIMIT)

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err = self.database.Query(ctx, stmt)
	if err != nil && !kit.ErrDatabaseNoRows.Is(err) {
		return nil, err
	}

	for _, res := range result {
		metric.Neutral[res.Keyword] = res.Count
	}

	result = []struct {
		Keyword string `db:"unnest"`
		Count   int    `db:"count"`
	}{}
	metric.Negative = make(map[string]int)

	stmt = sqlf.
		Select("UNNEST("+review.REVIEW_MODEL_TABLE+".keywords), COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		Where(review.REVIEW_MODEL_TABLE + ".sentiment = 'NEGATIVE'").
		GroupBy("UNNEST(" + review.REVIEW_MODEL_TABLE + ".keywords)").
		OrderBy("COUNT(*) DESC").
		Limit(METRIC_REVIEW_KEYWORDS_LIMIT)

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err = self.database.Query(ctx, stmt)
	if err != nil && !kit.ErrDatabaseNoRows.Is(err) {
		return nil, err
	}

	for _, res := range result {
		metric.Negative[res.Keyword] = res.Count
	}

	return &metric, nil
}

func (self *MetricRepository) GetNetPromoterScore(ctx context.Context,
	params NetPromoterScoreParams) (*NetPromoterScoreMetric, error) {
	var result []struct {
		Intention string `db:"intention"`
		Count     int    `db:"count"`
	}
	metric := NetPromoterScoreMetric{}
	metric.Score = 0

	stmt := sqlf.
		Select(review.REVIEW_MODEL_TABLE+".intention, COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy(review.REVIEW_MODEL_TABLE + ".intention")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	promoters := 0
	detractors := 0
	customers := 0
	for _, res := range result {
		switch res.Intention {
		case review.ReviewIntentionRetainAndRecommend:
			promoters += res.Count
		case review.ReviewIntentionChurn:
			detractors += res.Count
		case review.ReviewIntentionChurnAndDiscourage:
			detractors += res.Count
		}

		customers += res.Count
	}

	if customers == 0 {
		return &metric, nil
	}

	// NPS can be negative
	metric.Score = math.Round((float64(promoters-detractors) / float64(customers)) * 100)

	return &metric, nil
}

func (self *MetricRepository) GetCustomerSatisfactionScore(ctx context.Context,
	params CustomerSatisfactionScoreParams) (*CustomerSatisfactionScoreMetric, error) {
	var result []struct {
		Sentiment string `db:"sentiment"`
		Intention string `db:"intention"`
		Count     int    `db:"count"`
	}
	metric := CustomerSatisfactionScoreMetric{}
	metric.Score = 0

	stmt := sqlf.
		Select(review.REVIEW_MODEL_TABLE+".sentiment, "+review.REVIEW_MODEL_TABLE+".intention, COUNT(*)").To(&result).
		From(review.REVIEW_MODEL_TABLE).
		Join(feedback.FEEDBACK_MODEL_TABLE,
			feedback.FEEDBACK_MODEL_TABLE+".id = "+review.REVIEW_MODEL_TABLE+".feedback_id").
		Where(review.REVIEW_MODEL_TABLE+".product_id = ?", params.ProductID).
		GroupBy(review.REVIEW_MODEL_TABLE + ".sentiment").
		GroupBy(review.REVIEW_MODEL_TABLE + ".intention")

	if params.PeriodStartAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at >= ?", *params.PeriodStartAt)
	}

	if params.PeriodEndAt != nil {
		stmt.
			Where(feedback.FEEDBACK_MODEL_TABLE+".posted_at <= ?", *params.PeriodEndAt)
	}

	err := self.database.Query(ctx, stmt)
	if err != nil {
		if kit.ErrDatabaseNoRows.Is(err) {
			return &metric, nil
		}

		return nil, err
	}

	points := 0
	customers := 0
	for _, res := range result {
		switch res.Sentiment {
		case review.ReviewSentimentPositive:
			points += res.Count * 1
		case review.ReviewSentimentNeutral:
			points += res.Count * 0
		case review.ReviewSentimentNegative:
			points += res.Count * -1
		}

		switch res.Intention {
		case review.ReviewIntentionRetainAndRecommend:
			points += res.Count * 2
		case review.ReviewIntentionRetain:
			points += res.Count * 1
		case engine.OPTION_UNKNOWN:
			points += res.Count * 0
		case review.ReviewIntentionChurn:
			points += res.Count * -1
		case review.ReviewIntentionChurnAndDiscourage:
			points += res.Count * -2
		}

		customers += res.Count
	}

	if customers == 0 {
		return &metric, nil
	}

	// CSAT cannot be negative
	metric.Score = math.Round((((float64(points) / float64(customers*3)) * 100) + 100) / 2)

	return &metric, nil
}
