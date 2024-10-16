import * as entities from "$lib/entities";
import dayjs from "$lib/utils/datetime";

export interface BaseMetric {}

export type IssueCountMetric = BaseMetric & {
  active_issues: number;
  archived_issues: number;
  new_issues: number;
  feedbacks: number;
};

export function toIssueCountMetric(metric: IssueCountMetric): entities.IssueCountMetric {
  return {
    activeIssues: metric.active_issues,
    archivedIssues: metric.archived_issues,
    newIssues: metric.new_issues,
    feedbacks: metric.feedbacks,
  };
}

export type IssueSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export function toIssueSourcesMetric(metric: IssueSourcesMetric): entities.IssueSourcesMetric {
  return {
    sources: metric.sources,
  };
}

export type IssueSeveritiesMetric = BaseMetric & {
  severities: Record<string, number>;
};

export function toIssueSeveritiesMetric(metric: IssueSeveritiesMetric): entities.IssueSeveritiesMetric {
  return {
    severities: metric.severities,
  };
}

export type IssueCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export function toIssueCategoriesMetric(metric: IssueCategoriesMetric): entities.IssueCategoriesMetric {
  return {
    categories: metric.categories,
  };
}

export type IssueReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export function toIssueReleasesMetric(metric: IssueReleasesMetric): entities.IssueReleasesMetric {
  return {
    releases: metric.releases,
  };
}

export type SuggestionCountMetric = BaseMetric & {
  active_suggestions: number;
  archived_suggestions: number;
  new_suggestions: number;
  feedbacks: number;
};

export function toSuggestionCountMetric(metric: SuggestionCountMetric): entities.SuggestionCountMetric {
  return {
    activeSuggestions: metric.active_suggestions,
    archivedSuggestions: metric.archived_suggestions,
    newSuggestions: metric.new_suggestions,
    feedbacks: metric.feedbacks,
  };
}

export type SuggestionSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export function toSuggestionSourcesMetric(metric: SuggestionSourcesMetric): entities.SuggestionSourcesMetric {
  return {
    sources: metric.sources,
  };
}

export type SuggestionImportancesMetric = BaseMetric & {
  importances: Record<string, number>;
};

export function toSuggestionImportancesMetric(
  metric: SuggestionImportancesMetric,
): entities.SuggestionImportancesMetric {
  return {
    importances: metric.importances,
  };
}

export type SuggestionCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export function toSuggestionCategoriesMetric(metric: SuggestionCategoriesMetric): entities.SuggestionCategoriesMetric {
  return {
    categories: metric.categories,
  };
}

export type SuggestionReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export function toSuggestionReleasesMetric(metric: SuggestionReleasesMetric): entities.SuggestionReleasesMetric {
  return {
    releases: metric.releases,
  };
}

export type ReviewSentimentsMetric = BaseMetric & {
  sentiments: Record<string, number>;
};

export function toReviewSentimentsMetric(metric: ReviewSentimentsMetric): entities.ReviewSentimentsMetric {
  return {
    sentiments: metric.sentiments,
  };
}

export type ReviewSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export function toReviewSourcesMetric(metric: ReviewSourcesMetric): entities.ReviewSourcesMetric {
  return {
    sources: metric.sources,
  };
}

export type ReviewIntentionsMetric = BaseMetric & {
  intentions: Record<string, number>;
};

export function toReviewIntentionsMetric(metric: ReviewIntentionsMetric): entities.ReviewIntentionsMetric {
  return {
    intentions: metric.intentions,
  };
}

export type ReviewEmotionsMetric = BaseMetric & {
  emotions: Record<string, number>;
};

export function toReviewEmotionsMetric(metric: ReviewEmotionsMetric): entities.ReviewEmotionsMetric {
  return {
    emotions: metric.emotions,
  };
}

export type ReviewCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export function toReviewCategoriesMetric(metric: ReviewCategoriesMetric): entities.ReviewCategoriesMetric {
  return {
    categories: metric.categories,
  };
}

export type ReviewReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export function toReviewReleasesMetric(metric: ReviewReleasesMetric): entities.ReviewReleasesMetric {
  return {
    releases: metric.releases,
  };
}

export type ReviewKeywordsMetric = BaseMetric & {
  positive: Record<string, number>;
  neutral: Record<string, number>;
  negative: Record<string, number>;
};

export function toReviewKeywordsMetric(metric: ReviewKeywordsMetric): entities.ReviewKeywordsMetric {
  return {
    positive: metric.positive,
    neutral: metric.neutral,
    negative: metric.negative,
  };
}

export type NetPromoterScoreMetric = BaseMetric & {
  score: number;
};

export function toNetPromoterScoreMetric(metric: NetPromoterScoreMetric): entities.NetPromoterScoreMetric {
  return {
    score: metric.score,
  };
}

export type CustomerSatisfactionScoreMetric = BaseMetric & {
  score: number;
};

export function toCustomerSatisfactionScoreMetric(
  metric: CustomerSatisfactionScoreMetric,
): entities.CustomerSatisfactionScoreMetric {
  return {
    score: metric.score,
  };
}

export type GetBaseMetricRequest = {
  period_start_at?: Date;
  period_end_at?: Date;
};

export type GetBaseMetricResponse = BaseMetric;

export type GetIssueCountMetricRequest = GetBaseMetricRequest;

export function toGetIssueCountMetricQuery(request: GetIssueCountMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetIssueCountMetricResponse = IssueCountMetric;

export type GetIssueSourcesMetricRequest = GetBaseMetricRequest;

export function toGetIssueSourcesMetricQuery(request: GetIssueSourcesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetIssueSourcesMetricResponse = IssueSourcesMetric;

export type GetIssueSeveritiesMetricRequest = GetBaseMetricRequest;

export function toGetIssueSeveritiesMetricQuery(request: GetIssueSeveritiesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetIssueSeveritiesMetricResponse = IssueSeveritiesMetric;

export type GetIssueCategoriesMetricRequest = GetBaseMetricRequest;

export function toGetIssueCategoriesMetricQuery(request: GetIssueCategoriesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetIssueCategoriesMetricResponse = IssueCategoriesMetric;

export type GetIssueReleasesMetricRequest = GetBaseMetricRequest;

export function toGetIssueReleasesMetricQuery(request: GetIssueReleasesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetIssueReleasesMetricResponse = IssueReleasesMetric;

export type GetSuggestionCountMetricRequest = GetBaseMetricRequest;

export function toGetSuggestionCountMetricQuery(request: GetSuggestionCountMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetSuggestionCountMetricResponse = SuggestionCountMetric;

export type GetSuggestionSourcesMetricRequest = GetBaseMetricRequest;

export function toGetSuggestionSourcesMetricQuery(request: GetSuggestionSourcesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetSuggestionSourcesMetricResponse = SuggestionSourcesMetric;

export type GetSuggestionImportancesMetricRequest = GetBaseMetricRequest;

export function toGetSuggestionImportancesMetricQuery(request: GetSuggestionImportancesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetSuggestionImportancesMetricResponse = SuggestionImportancesMetric;

export type GetSuggestionCategoriesMetricRequest = GetBaseMetricRequest;

export function toGetSuggestionCategoriesMetricQuery(request: GetSuggestionCategoriesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetSuggestionCategoriesMetricResponse = SuggestionCategoriesMetric;

export type GetSuggestionReleasesMetricRequest = GetBaseMetricRequest;

export function toGetSuggestionReleasesMetricQuery(request: GetSuggestionReleasesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetSuggestionReleasesMetricResponse = SuggestionReleasesMetric;

export type GetReviewSentimentsMetricRequest = GetBaseMetricRequest;

export function toGetReviewSentimentsMetricQuery(request: GetReviewSentimentsMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewSentimentsMetricResponse = ReviewSentimentsMetric;

export type GetReviewSourcesMetricRequest = GetBaseMetricRequest;

export function toGetReviewSourcesMetricQuery(request: GetReviewSourcesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewSourcesMetricResponse = ReviewSourcesMetric;

export type GetReviewIntentionsMetricRequest = GetBaseMetricRequest;

export function toGetReviewIntentionsMetricQuery(request: GetReviewIntentionsMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewIntentionsMetricResponse = ReviewIntentionsMetric;

export type GetReviewEmotionsMetricRequest = GetBaseMetricRequest;

export function toGetReviewEmotionsMetricQuery(request: GetReviewEmotionsMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewEmotionsMetricResponse = ReviewEmotionsMetric;

export type GetReviewCategoriesMetricRequest = GetBaseMetricRequest;

export function toGetReviewCategoriesMetricQuery(request: GetReviewCategoriesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewCategoriesMetricResponse = ReviewCategoriesMetric;

export type GetReviewReleasesMetricRequest = GetBaseMetricRequest;

export function toGetReviewReleasesMetricQuery(request: GetReviewReleasesMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewReleasesMetricResponse = ReviewReleasesMetric;

export type GetReviewKeywordsMetricRequest = GetBaseMetricRequest;

export function toGetReviewKeywordsMetricQuery(request: GetReviewKeywordsMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetReviewKeywordsMetricResponse = ReviewKeywordsMetric;

export type GetNetPromoterScoreMetricRequest = GetBaseMetricRequest;

export function toGetNetPromoterScoreMetricQuery(request: GetNetPromoterScoreMetricRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetNetPromoterScoreMetricResponse = NetPromoterScoreMetric;

export type GetCustomerSatisfactionScoreMetricRequest = GetBaseMetricRequest;

export function toGetCustomerSatisfactionScoreMetricQuery(
  request: GetCustomerSatisfactionScoreMetricRequest,
): URLSearchParams {
  const params = new URLSearchParams();

  if (request.period_start_at) params.set("period_start_at", dayjs(request.period_start_at).toISOString());

  if (request.period_end_at) params.set("period_end_at", dayjs(request.period_end_at).toISOString());

  return params;
}

export type GetCustomerSatisfactionScoreMetricResponse = CustomerSatisfactionScoreMetric;
