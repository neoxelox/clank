import json
from typing import List

from pydantic import BaseModel

from src.common import Usage
from src.common.utils import get_tokens
from src.config import Config

from .issue_extractor import IssueExtractor
from .review_extractor import ReviewExtractor
from .suggestion_extractor import SuggestionExtractor


class Processor:
    def __init__(self, config: Config) -> None:
        self.config = config

        self.issue_extractor = IssueExtractor(config=config)
        self.suggestion_extractor = SuggestionExtractor(config=config)
        self.review_extractor = ReviewExtractor(config=config)

    class ExtractIssuesParams(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class ExtractIssuesResult(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]
            severity: str
            category: str

        issues: List[Issue]
        usage: Usage

    def extract_issues(self, params: ExtractIssuesParams) -> ExtractIssuesResult:
        issues = self.issue_extractor(
            input=IssueExtractor.Input(
                context=params.context,
                categories=params.categories,
                feedback=params.feedback,
            )
        ).output.issues

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.ExtractIssuesResult.model_json_schema()) + params.model_dump_json()) + 1500
        ) * 2
        output_tokens = (sum([get_tokens(issue.model_dump_json()) for issue in issues]) + 100) * 2

        return self.ExtractIssuesResult(
            issues=[
                self.ExtractIssuesResult.Issue(
                    title=issue.title,
                    description=issue.description,
                    steps=issue.steps,
                    severity=issue.severity,
                    category=issue.category,
                )
                for issue in issues
            ],
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class ExtractSuggestionsParams(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class ExtractSuggestionsResult(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str
            importance: str
            category: str

        suggestions: List[Suggestion]
        usage: Usage

    def extract_suggestions(self, params: ExtractSuggestionsParams) -> ExtractSuggestionsResult:
        suggestions = self.suggestion_extractor(
            input=SuggestionExtractor.Input(
                context=params.context,
                categories=params.categories,
                feedback=params.feedback,
            )
        ).output.suggestions

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.ExtractSuggestionsResult.model_json_schema()) + params.model_dump_json()) + 1500
        ) * 2
        output_tokens = (sum([get_tokens(suggestion.model_dump_json()) for suggestion in suggestions]) + 100) * 2

        return self.ExtractSuggestionsResult(
            suggestions=[
                self.ExtractSuggestionsResult.Suggestion(
                    title=suggestion.title,
                    description=suggestion.description,
                    reason=suggestion.reason,
                    importance=suggestion.importance,
                    category=suggestion.category,
                )
                for suggestion in suggestions
            ],
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class ExtractReviewParams(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class ExtractReviewResult(BaseModel):
        class Review(BaseModel):
            content: str
            keywords: List[str]
            sentiment: str
            emotions: List[str]
            intention: str
            category: str

        review: Review
        usage: Usage

    def extract_review(self, params: ExtractReviewParams) -> ExtractReviewResult:
        review = self.review_extractor(
            input=ReviewExtractor.Input(
                context=params.context,
                categories=params.categories,
                feedback=params.feedback,
            )
        ).output.review

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.ExtractReviewResult.model_json_schema()) + params.model_dump_json()) + 1500
        )
        output_tokens = get_tokens(review.model_dump_json()) + 100

        return self.ExtractReviewResult(
            review=self.ExtractReviewResult.Review(
                content=review.content,
                keywords=review.keywords,
                sentiment=review.sentiment,
                emotions=review.emotions,
                intention=review.intention,
                category=review.category,
            ),
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )
