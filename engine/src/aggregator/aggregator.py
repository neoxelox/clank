import json
from functools import partial
from typing import List

from flashrank import Ranker, RerankRequest
from openai import OpenAI
from pydantic import BaseModel

from src.common import Usage
from src.common.utils import get_tokens
from src.config import Config

from .issue_merger import IssueMerger
from .issue_similarity_discernor import IssueSimilarityDiscernor
from .suggestion_merger import SuggestionMerger
from .suggestion_similarity_discernor import SuggestionSimilarityDiscernor


class Aggregator:
    def __init__(self, config: Config) -> None:
        self.config = config

        self.embedder = partial(
            OpenAI(api_key=config.lm.openai.api_key).embeddings.create,
            model="text-embedding-3-small",
            encoding_format="float",
        )

        self.ranker = Ranker(
            model_name="rank-T5-flan", cache_dir=f"{config.service.resources_path}/rank-T5-flan"
        ).rerank

        self.issue_similarity_discernor = IssueSimilarityDiscernor(config=config)
        self.issue_merger = IssueMerger(config=config)

        self.suggestion_similarity_discernor = SuggestionSimilarityDiscernor(config=config)
        self.suggestion_merger = SuggestionMerger(config=config)

    class ComputeEmbeddingParams(BaseModel):
        text: str

    class ComputeEmbeddingResult(BaseModel):
        embedding: List[float]
        usage: Usage

    def compute_embedding(self, params: ComputeEmbeddingParams) -> ComputeEmbeddingResult:
        embedding = self.embedder(input=params.text).data[0].embedding

        # TODO: Use real usage
        input_tokens = get_tokens(params.text)
        output_tokens = 0

        return self.ComputeEmbeddingResult(
            embedding=embedding,
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class SimilarIssueParams(BaseModel):
        issue: str
        options: List[str]

    class SimilarIssueResult(BaseModel):
        option: int
        usage: Usage

    def similar_issue(self, params: SimilarIssueParams) -> SimilarIssueResult:
        if not params.options:
            return self.SimilarIssueResult(
                option=0,
                usage=Usage(
                    input=0,
                    output=0,
                ),
            )

        similar_issues = self.ranker(
            RerankRequest(
                query=params.issue,
                passages=[{"index": index, "text": issue} for index, issue in enumerate(params.options)],
            )
        )

        similar_issues = [(issue["index"], issue["text"]) for issue in similar_issues if float(issue["score"]) >= 0.30][
            :3
        ]

        if not similar_issues:
            return self.SimilarIssueResult(
                option=0,
                usage=Usage(
                    input=0,
                    output=0,
                ),
            )

        index = self.issue_similarity_discernor(
            input=IssueSimilarityDiscernor.Input(
                issue=params.issue,
                options=[issue[1] for issue in similar_issues],
            )
        ).output.index

        if index >= 0:
            index = similar_issues[index][0]

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.SimilarIssueResult.model_json_schema()) + params.model_dump_json()) + 1500
        ) * len(similar_issues)
        output_tokens = (get_tokens(str(index)) + 100) * len(similar_issues)

        return self.SimilarIssueResult(
            option=index + 1,
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class MergeIssuesParams(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue_a: Issue
        issue_b: Issue

    class MergeIssuesResult(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue: Issue
        usage: Usage

    def merge_issues(self, params: MergeIssuesParams) -> MergeIssuesResult:
        issue = self.issue_merger(
            input=IssueMerger.Input(
                issue_a=IssueMerger.Input.Issue(
                    title=params.issue_a.title,
                    description=params.issue_a.description,
                    steps=params.issue_a.steps,
                ),
                issue_b=IssueMerger.Input.Issue(
                    title=params.issue_b.title,
                    description=params.issue_b.description,
                    steps=params.issue_b.steps,
                ),
            )
        ).output.issue

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.MergeIssuesResult.model_json_schema()) + params.model_dump_json()) + 1500
        )
        output_tokens = get_tokens(issue.model_dump_json()) + 100

        return self.MergeIssuesResult(
            issue=self.MergeIssuesResult.Issue(
                title=issue.title,
                description=issue.description,
                steps=issue.steps,
            ),
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class SimilarSuggestionParams(BaseModel):
        suggestion: str
        options: List[str]

    class SimilarSuggestionResult(BaseModel):
        option: int
        usage: Usage

    def similar_suggestion(self, params: SimilarSuggestionParams) -> SimilarSuggestionResult:
        if not params.options:
            return self.SimilarSuggestionResult(
                option=0,
                usage=Usage(
                    input=0,
                    output=0,
                ),
            )

        similar_suggestions = self.ranker(
            RerankRequest(
                query=params.suggestion,
                passages=[{"index": index, "text": suggestion} for index, suggestion in enumerate(params.options)],
            )
        )

        similar_suggestions = [
            (suggestion["index"], suggestion["text"])
            for suggestion in similar_suggestions
            if float(suggestion["score"]) >= 0.30
        ][:3]

        if not similar_suggestions:
            return self.SimilarSuggestionResult(
                option=0,
                usage=Usage(
                    input=0,
                    output=0,
                ),
            )

        index = self.suggestion_similarity_discernor(
            input=SuggestionSimilarityDiscernor.Input(
                suggestion=params.suggestion,
                options=[suggestion[1] for suggestion in similar_suggestions],
            )
        ).output.index

        if index >= 0:
            index = similar_suggestions[index][0]

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.SimilarSuggestionResult.model_json_schema()) + params.model_dump_json()) + 1500
        ) * len(similar_suggestions)
        output_tokens = (get_tokens(str(index)) + 100) * len(similar_suggestions)

        return self.SimilarSuggestionResult(
            option=index + 1,
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class MergeSuggestionsParams(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion_a: Suggestion
        suggestion_b: Suggestion

    class MergeSuggestionsResult(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion: Suggestion
        usage: Usage

    def merge_suggestions(self, params: MergeSuggestionsParams) -> MergeSuggestionsResult:
        suggestion = self.suggestion_merger(
            input=SuggestionMerger.Input(
                suggestion_a=SuggestionMerger.Input.Suggestion(
                    title=params.suggestion_a.title,
                    description=params.suggestion_a.description,
                    reason=params.suggestion_a.reason,
                ),
                suggestion_b=SuggestionMerger.Input.Suggestion(
                    title=params.suggestion_b.title,
                    description=params.suggestion_b.description,
                    reason=params.suggestion_b.reason,
                ),
            )
        ).output.suggestion

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.MergeSuggestionsResult.model_json_schema()) + params.model_dump_json()) + 1500
        )
        output_tokens = get_tokens(suggestion.model_dump_json()) + 100

        return self.MergeSuggestionsResult(
            suggestion=self.MergeSuggestionsResult.Suggestion(
                title=suggestion.title,
                description=suggestion.description,
                reason=suggestion.reason,
            ),
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )
