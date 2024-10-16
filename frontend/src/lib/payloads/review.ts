import * as entities from "$lib/entities";
import dayjs from "$lib/utils/datetime";
import type { Feedback } from "./feedback";
import { toFeedback } from "./feedback";

export type Review = {
  id: string;
  product_id: string;
  feedback: Feedback;
  keywords: string[];
  sentiment: string;
  emotions: string[];
  intention: string;
  category: string;
  quality?: number;
};

export function toReview(review: Review): entities.Review {
  return {
    id: review.id,
    productID: review.product_id,
    feedback: toFeedback(review.feedback),
    keywords: review.keywords,
    sentiment: review.sentiment,
    emotions: review.emotions,
    intention: review.intention,
    category: review.category,
    quality: review.quality,
  };
}

export type ListReviewsRequest = {
  filters: {
    sources?: string[];
    releases?: string[];
    categories?: string[];
    keywords?: string[];
    sentiments?: string[];
    emotions?: string[];
    intentions?: string[];
    languages?: string[];
    seen_start_at?: Date;
    seen_end_at?: Date;
  };
  orders: {
    recency?: string;
  };
  pagination: {
    limit?: number;
    from?: string;
  };
};

export function fromListReviewsQuery(params: URLSearchParams): ListReviewsRequest {
  let sources = undefined;
  if (params.has("sources")) sources = params.getAll("sources");

  let releases = undefined;
  if (params.has("releases")) releases = params.getAll("releases");

  let categories = undefined;
  if (params.has("categories")) categories = params.getAll("categories");

  let keywords = undefined;
  if (params.has("keywords")) keywords = params.getAll("keywords");

  let sentiments = undefined;
  if (params.has("sentiments")) sentiments = params.getAll("sentiments");

  let emotions = undefined;
  if (params.has("emotions")) emotions = params.getAll("emotions");

  let intentions = undefined;
  if (params.has("intentions")) intentions = params.getAll("intentions");

  let languages = undefined;
  if (params.has("languages")) languages = params.getAll("languages");

  let seenStartAt = undefined;
  if (params.has("seen_start_at")) seenStartAt = dayjs(params.get("seen_start_at")).toDate();

  let seenEndAt = undefined;
  if (params.has("seen_end_at")) seenEndAt = dayjs(params.get("seen_end_at")).toDate();

  let recency = undefined;
  if (params.has("recency")) recency = params.get("recency").toUpperCase();

  let limit = undefined;
  if (params.has("limit")) limit = parseInt(params.get("limit"));

  let from = undefined;
  if (params.has("from")) from = params.get("from");

  return {
    filters: {
      sources: sources,
      releases: releases,
      categories: categories,
      keywords: keywords,
      sentiments: sentiments,
      emotions: emotions,
      intentions: intentions,
      languages: languages,
      seen_start_at: seenStartAt,
      seen_end_at: seenEndAt,
    },
    orders: {
      recency: recency,
    },
    pagination: {
      limit: limit,
      from: from,
    },
  };
}

export function toListReviewsQuery(request: ListReviewsRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.filters.sources) request.filters.sources.forEach((source) => params.append("sources", source));

  if (request.filters.releases) request.filters.releases.forEach((release) => params.append("releases", release));

  if (request.filters.categories)
    request.filters.categories.forEach((category) => params.append("categories", category));

  if (request.filters.keywords) request.filters.keywords.forEach((keyword) => params.append("keywords", keyword));

  if (request.filters.sentiments)
    request.filters.sentiments.forEach((sentiment) => params.append("sentiments", sentiment));

  if (request.filters.emotions) request.filters.emotions.forEach((emotion) => params.append("emotions", emotion));

  if (request.filters.intentions)
    request.filters.intentions.forEach((intention) => params.append("intentions", intention));

  if (request.filters.languages) request.filters.languages.forEach((language) => params.append("languages", language));

  if (request.filters.seen_start_at) params.set("seen_start_at", dayjs(request.filters.seen_start_at).toISOString());

  if (request.filters.seen_end_at) params.set("seen_end_at", dayjs(request.filters.seen_end_at).toISOString());

  if (request.orders.recency) params.set("recency", request.orders.recency.toUpperCase());

  if (request.pagination.limit && request.pagination.limit <= 100)
    params.set("limit", request.pagination.limit.toString());

  if (request.pagination.from) params.set("from", request.pagination.from);

  return params;
}

export type ListReviewsResponse = {
  reviews: Review[];
  next?: string;
};

export type GetReviewResponse = Review;
