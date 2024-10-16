export type BaseParams = {
  productID: string;
  periodStartAt?: Date;
  periodEndAt?: Date;
};

export interface BaseMetric {}

export type IssueCountParams = BaseParams;

export type IssueCountMetric = BaseMetric & {
  activeIssues: number;
  archivedIssues: number;
  newIssues: number;
  feedbacks: number;
};

export type IssueTrendsParams = BaseParams;

export type IssueTrendsMetric = BaseMetric & {
  source?: string;
  severity?: string;
  category?: string;
  release?: string;
};

export type IssueSourcesParams = BaseParams;

export type IssueSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export type IssueSeveritiesParams = BaseParams;

export type IssueSeveritiesMetric = BaseMetric & {
  severities: Record<string, number>;
};

export type IssueCategoriesParams = BaseParams;

export type IssueCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export type IssueReleasesParams = BaseParams;

export type IssueReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export type SuggestionCountParams = BaseParams;

export type SuggestionCountMetric = BaseMetric & {
  activeSuggestions: number;
  archivedSuggestions: number;
  newSuggestions: number;
  feedbacks: number;
};

export type SuggestionTrendsParams = BaseParams;

export type SuggestionTrendsMetric = BaseMetric & {
  source?: string;
  importance?: string;
  category?: string;
  release?: string;
};

export type SuggestionSourcesParams = BaseParams;

export type SuggestionSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export type SuggestionImportancesParams = BaseParams;

export type SuggestionImportancesMetric = BaseMetric & {
  importances: Record<string, number>;
};

export type SuggestionCategoriesParams = BaseParams;

export type SuggestionCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export type SuggestionReleasesParams = BaseParams;

export type SuggestionReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export type ReviewSentimentsParams = BaseParams;

export type ReviewSentimentsMetric = BaseMetric & {
  sentiments: Record<string, number>;
};

export type ReviewSourcesParams = BaseParams;

export type ReviewSourcesMetric = BaseMetric & {
  sources: Record<string, number>;
};

export type ReviewIntentionsParams = BaseParams;

export type ReviewIntentionsMetric = BaseMetric & {
  intentions: Record<string, number>;
};

export type ReviewEmotionsParams = BaseParams;

export type ReviewEmotionsMetric = BaseMetric & {
  emotions: Record<string, number>;
};

export type ReviewCategoriesParams = BaseParams;

export type ReviewCategoriesMetric = BaseMetric & {
  categories: Record<string, number>;
};

export type ReviewReleasesParams = BaseParams;

export type ReviewReleasesMetric = BaseMetric & {
  releases: Record<string, number>;
};

export type ReviewKeywordsParams = BaseParams;

export type ReviewKeywordsMetric = BaseMetric & {
  positive: Record<string, number>;
  neutral: Record<string, number>;
  negative: Record<string, number>;
};

export type NetPromoterScoreParams = BaseParams;

export type NetPromoterScoreMetric = BaseMetric & {
  score: number;
};

export type CustomerSatisfactionScoreParams = BaseParams;

export type CustomerSatisfactionScoreMetric = BaseMetric & {
  score: number;
};
