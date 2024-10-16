import * as entities from "$lib/entities";
import dayjs from "$lib/utils/datetime";
import type { Feedback } from "./feedback";

export type Issue = {
  id: string;
  product_id: string;
  sources: Record<string, number>;
  title: string;
  description: string;
  steps: string[];
  severities: Record<string, number>;
  priority: number;
  categories: Record<string, number>;
  releases: Record<string, number>;
  customers: number;
  assignee_id?: string;
  quality?: number;
  first_seen_at: Date;
  last_seen_at: Date;
  archived_at?: Date;
};

export function toIssue(issue: Issue): entities.Issue {
  return {
    id: issue.id,
    productID: issue.product_id,
    sources: issue.sources,
    title: issue.title,
    description: issue.description,
    steps: issue.steps,
    severities: issue.severities,
    priority: issue.priority,
    categories: issue.categories,
    releases: issue.releases,
    customers: issue.customers,
    assigneeID: issue.assignee_id,
    quality: issue.quality,
    firstSeenAt: issue.first_seen_at,
    lastSeenAt: issue.last_seen_at,
    archivedAt: issue.archived_at,
  };
}

export const ISSUE_FILTER_UNASSIGNED: string = "UNASSIGNED";

export enum IssueFilterStatus {
  ACTIVE = "ACTIVE",
  REGRESSED = "REGRESSED",
  ARCHIVED = "ARCHIVED",
  UNARCHIVED = "UNARCHIVED",
}

export type ListIssuesRequest = {
  filters: {
    content?: string;
    sources?: string[];
    severities?: string[];
    releases?: string[];
    categories?: string[];
    assignees?: string[];
    status?: string;
    first_seen_start_at?: Date;
    first_seen_end_at?: Date;
    last_seen_start_at?: Date;
    last_seen_end_at?: Date;
  };
  orders: {
    relevance?: string;
  };
  pagination: {
    limit?: number;
    from?: string;
  };
};

export function fromListIssuesQuery(params: URLSearchParams): ListIssuesRequest {
  let content = undefined;
  if (params.has("content")) content = params.get("content");

  let sources = undefined;
  if (params.has("sources")) sources = params.getAll("sources");

  let severities = undefined;
  if (params.has("severities")) severities = params.getAll("severities");

  let releases = undefined;
  if (params.has("releases")) releases = params.getAll("releases");

  let categories = undefined;
  if (params.has("categories")) categories = params.getAll("categories");

  let assignees = undefined;
  if (params.has("assignees")) assignees = params.getAll("assignees");

  let status = undefined;
  if (params.has("status")) status = params.get("status").toUpperCase();

  let firstSeenStartAt = undefined;
  if (params.has("first_seen_start_at")) firstSeenStartAt = dayjs(params.get("first_seen_start_at")).toDate();

  let firstSeenEndAt = undefined;
  if (params.has("first_seen_end_at")) firstSeenEndAt = dayjs(params.get("first_seen_end_at")).toDate();

  let lastSeenStartAt = undefined;
  if (params.has("last_seen_start_at")) lastSeenStartAt = dayjs(params.get("last_seen_start_at")).toDate();

  let lastSeenEndAt = undefined;
  if (params.has("last_seen_end_at")) lastSeenEndAt = dayjs(params.get("last_seen_end_at")).toDate();

  let relevance = undefined;
  if (params.has("relevance")) relevance = params.get("relevance").toUpperCase();

  let limit = undefined;
  if (params.has("limit")) limit = parseInt(params.get("limit"));

  let from = undefined;
  if (params.has("from")) from = params.get("from");

  return {
    filters: {
      content: content,
      sources: sources,
      severities: severities,
      releases: releases,
      categories: categories,
      assignees: assignees,
      status: status,
      first_seen_start_at: firstSeenStartAt,
      first_seen_end_at: firstSeenEndAt,
      last_seen_start_at: lastSeenStartAt,
      last_seen_end_at: lastSeenEndAt,
    },
    orders: {
      relevance: relevance,
    },
    pagination: {
      limit: limit,
      from: from,
    },
  };
}

export function toListIssuesQuery(request: ListIssuesRequest): URLSearchParams {
  const params = new URLSearchParams();

  if (request.filters.content && request.filters.content.length >= 5 && request.filters.content.length <= 250)
    params.set("content", request.filters.content);

  if (request.filters.sources) request.filters.sources.forEach((source) => params.append("sources", source));

  if (request.filters.severities)
    request.filters.severities.forEach((severity) => params.append("severities", severity));

  if (request.filters.releases) request.filters.releases.forEach((release) => params.append("releases", release));

  if (request.filters.categories)
    request.filters.categories.forEach((category) => params.append("categories", category));

  if (request.filters.assignees) request.filters.assignees.forEach((assignee) => params.append("assignees", assignee));

  if (request.filters.status) params.set("status", request.filters.status.toUpperCase());

  if (request.filters.first_seen_start_at)
    params.set("first_seen_start_at", dayjs(request.filters.first_seen_start_at).toISOString());

  if (request.filters.first_seen_end_at)
    params.set("first_seen_end_at", dayjs(request.filters.first_seen_end_at).toISOString());

  if (request.filters.last_seen_start_at)
    params.set("last_seen_start_at", dayjs(request.filters.last_seen_start_at).toISOString());

  if (request.filters.last_seen_end_at)
    params.set("last_seen_end_at", dayjs(request.filters.last_seen_end_at).toISOString());

  if (request.orders.relevance) params.set("relevance", request.orders.relevance.toUpperCase());

  if (request.pagination.limit && request.pagination.limit <= 100)
    params.set("limit", request.pagination.limit.toString());

  if (request.pagination.from) params.set("from", request.pagination.from);

  return params;
}

export type ListIssuesResponse = {
  issues: Issue[];
  next?: string;
};

export type GetIssueResponse = Issue;

export type ListIssueFeedbacksRequest = {
  from?: string;
};

export type ListIssueFeedbacksResponse = {
  feedbacks: Feedback[];
  next?: string;
};

export type PutIssueAssigneeRequest = {
  assignee_id?: string;
};

export type PutIssueAssigneeResponse = {
  assignee_id?: string;
};

export type PutIssueArchivedRequest = {
  archived: boolean;
};

export type PutIssueArchivedResponse = {
  archived: boolean;
};
