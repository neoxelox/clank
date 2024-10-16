from typing import List

from pydantic import BaseModel

from src.common import Usage
from src.config import Config

from .processor import Processor


class ProcessorEndpoints:
    def __init__(self, config: Config, processor: Processor) -> None:
        self.config = config
        self.processor = processor

    class PostExtractIssuesRequest(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class PostExtractIssuesResponse(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]
            severity: str
            category: str

        issues: List[Issue]
        usage: Usage

    async def post_extract_issues(self, request: PostExtractIssuesRequest) -> PostExtractIssuesResponse:
        result = self.processor.extract_issues(
            params=Processor.ExtractIssuesParams(
                context=request.context,
                categories=request.categories,
                feedback=request.feedback,
            )
        )

        return self.PostExtractIssuesResponse(
            issues=[
                self.PostExtractIssuesResponse.Issue(
                    title=issue.title,
                    description=issue.description,
                    steps=issue.steps,
                    severity=issue.severity,
                    category=issue.category,
                )
                for issue in result.issues
            ],
            usage=result.usage,
        )

    class PostExtractSuggestionsRequest(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class PostExtractSuggestionsResponse(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str
            importance: str
            category: str

        suggestions: List[Suggestion]
        usage: Usage

    async def post_extract_suggestions(self, request: PostExtractSuggestionsRequest) -> PostExtractSuggestionsResponse:
        result = self.processor.extract_suggestions(
            params=Processor.ExtractSuggestionsParams(
                context=request.context,
                categories=request.categories,
                feedback=request.feedback,
            )
        )

        return self.PostExtractSuggestionsResponse(
            suggestions=[
                self.PostExtractSuggestionsResponse.Suggestion(
                    title=suggestion.title,
                    description=suggestion.description,
                    reason=suggestion.reason,
                    importance=suggestion.importance,
                    category=suggestion.category,
                )
                for suggestion in result.suggestions
            ],
            usage=result.usage,
        )

    class PostExtractReviewRequest(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class PostExtractReviewResponse(BaseModel):
        class Review(BaseModel):
            content: str
            keywords: List[str]
            sentiment: str
            emotions: List[str]
            intention: str
            category: str

        review: Review
        usage: Usage

    async def post_extract_review(self, request: PostExtractReviewRequest) -> PostExtractReviewResponse:
        result = self.processor.extract_review(
            params=Processor.ExtractReviewParams(
                context=request.context,
                categories=request.categories,
                feedback=request.feedback,
            )
        )

        return self.PostExtractReviewResponse(
            review=self.PostExtractReviewResponse.Review(
                content=result.review.content,
                keywords=result.review.keywords,
                sentiment=result.review.sentiment,
                emotions=result.review.emotions,
                intention=result.review.intention,
                category=result.review.category,
            ),
            usage=result.usage,
        )
