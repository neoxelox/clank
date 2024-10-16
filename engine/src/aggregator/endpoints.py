from typing import List

from pydantic import BaseModel

from src.common import Usage
from src.config import Config

from .aggregator import Aggregator


class AggregatorEndpoints:
    def __init__(self, config: Config, aggregator: Aggregator) -> None:
        self.config = config
        self.aggregator = aggregator

    class PostComputeEmbeddingRequest(BaseModel):
        text: str

    class PostComputeEmbeddingResponse(BaseModel):
        embedding: List[float]
        usage: Usage

    async def post_compute_embedding(self, request: PostComputeEmbeddingRequest) -> PostComputeEmbeddingResponse:
        result = self.aggregator.compute_embedding(
            params=Aggregator.ComputeEmbeddingParams(
                text=request.text,
            )
        )

        return self.PostComputeEmbeddingResponse(
            embedding=result.embedding,
            usage=result.usage,
        )

    class PostSimilarIssueRequest(BaseModel):
        issue: str
        options: List[str]

    class PostSimilarIssueResponse(BaseModel):
        option: int
        usage: Usage

    async def post_similar_issue(self, request: PostSimilarIssueRequest) -> PostSimilarIssueResponse:
        result = self.aggregator.similar_issue(
            params=Aggregator.SimilarIssueParams(
                issue=request.issue,
                options=request.options,
            )
        )

        return self.PostSimilarIssueResponse(
            option=result.option,
            usage=result.usage,
        )

    class PostMergeIssuesRequest(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue_a: Issue
        issue_b: Issue

    class PostMergeIssuesResponse(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue: Issue
        usage: Usage

    async def post_merge_issues(self, request: PostMergeIssuesRequest) -> PostMergeIssuesResponse:
        result = self.aggregator.merge_issues(
            params=Aggregator.MergeIssuesParams(
                issue_a=Aggregator.MergeIssuesParams.Issue(
                    title=request.issue_a.title,
                    description=request.issue_a.description,
                    steps=request.issue_a.steps,
                ),
                issue_b=Aggregator.MergeIssuesParams.Issue(
                    title=request.issue_b.title,
                    description=request.issue_b.description,
                    steps=request.issue_b.steps,
                ),
            )
        )

        return self.PostMergeIssuesResponse(
            issue=self.PostMergeIssuesResponse.Issue(
                title=result.issue.title,
                description=result.issue.description,
                steps=result.issue.steps,
            ),
            usage=result.usage,
        )

    class PostSimilarSuggestionRequest(BaseModel):
        suggestion: str
        options: List[str]

    class PostSimilarSuggestionResponse(BaseModel):
        option: int
        usage: Usage

    async def post_similar_suggestion(self, request: PostSimilarSuggestionRequest) -> PostSimilarSuggestionResponse:
        result = self.aggregator.similar_suggestion(
            params=Aggregator.SimilarSuggestionParams(
                suggestion=request.suggestion,
                options=request.options,
            )
        )

        return self.PostSimilarSuggestionResponse(
            option=result.option,
            usage=result.usage,
        )

    class PostMergeSuggestionsRequest(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion_a: Suggestion
        suggestion_b: Suggestion

    class PostMergeSuggestionsResponse(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion: Suggestion
        usage: Usage

    async def post_merge_suggestions(self, request: PostMergeSuggestionsRequest) -> PostMergeSuggestionsResponse:
        result = self.aggregator.merge_suggestions(
            params=Aggregator.MergeSuggestionsParams(
                suggestion_a=Aggregator.MergeSuggestionsParams.Suggestion(
                    title=request.suggestion_a.title,
                    description=request.suggestion_a.description,
                    reason=request.suggestion_a.reason,
                ),
                suggestion_b=Aggregator.MergeSuggestionsParams.Suggestion(
                    title=request.suggestion_b.title,
                    description=request.suggestion_b.description,
                    reason=request.suggestion_b.reason,
                ),
            )
        )

        return self.PostMergeSuggestionsResponse(
            suggestion=self.PostMergeSuggestionsResponse.Suggestion(
                title=result.suggestion.title,
                description=result.suggestion.description,
                reason=result.suggestion.reason,
            ),
            usage=result.usage,
        )
