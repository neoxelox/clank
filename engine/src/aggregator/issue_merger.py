from typing import List

from dspy import InputField, Module, OutputField, Prediction, Signature, Suggest, backtrack_handler
from pydantic import BaseModel, Field

from src.common import ChainOfThought
from src.config import Config


class IssueMerger(Module):
    class Input(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue_a: Issue
        issue_b: Issue

    class Output(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issue: Issue

    class MergeIssues(Signature):
        """
Merge, coherently, issue A and issue B, that customers have with a product, into a single issue.
- Maintain the core problem, context and nuances of both issues.
- Do not create information that is not present in any of the issues.
        """  # fmt: skip

        class Input(BaseModel):
            class Issue(BaseModel):
                title: str
                description: str
                steps: List[str]

            issue_a: Issue
            issue_b: Issue

        class Output(BaseModel):
            class Issue(BaseModel):
                title: str = Field(
                    description="4 to 10 words, which cannot contain the words `issue` (or synonyms), `customer` (or synonyms) or the product's name.",
                    max_length=100,
                )
                description: str = Field(
                    description="Long, complete explanation, but without redundant information, using the feedback's original words. Must focus solely on the issue by depersonalizing the sentences."
                )
                steps: List[str] = Field(
                    description="Precise steps, but very concise, if any, to be able to reproduce the issue, else `[]`.",
                    max_items=5,
                )

            issue: Issue

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.merge_issues = ChainOfThought(self.MergeIssues, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/aggregator/issue_aggregator/issue_merger.json")

    def forward(self, input: Input) -> Prediction:
        if input.issue_a == input.issue_b:
            return Prediction(
                output=self.Output(
                    issue=self.Output.Issue(
                        title=input.issue_a.title,
                        description=input.issue_a.description,
                        steps=input.issue_a.steps,
                    ),
                )
            )

        issue = self.merge_issues(
            input=self.MergeIssues.Input(
                issue_a=self.MergeIssues.Input.Issue(
                    title=input.issue_a.title,
                    description=input.issue_a.description,
                    steps=input.issue_a.steps,
                ),
                issue_b=self.MergeIssues.Input.Issue(
                    title=input.issue_b.title,
                    description=input.issue_b.description,
                    steps=input.issue_b.steps,
                ),
            )
        ).output.issue

        Suggest(
            len(issue.steps) <= len(input.issue_a.steps) + len(input.issue_b.steps),
            f"The merged issue's `steps to reproduce` ({len(issue.steps)}) cannot be longer than the sum of the `steps to reproduce` of the original issues ({len(input.issue_a.steps) + len(input.issue_b.steps)})!",
        )

        return Prediction(
            output=self.Output(
                issue=self.Output.Issue(
                    title=issue.title,
                    description=issue.description,
                    steps=issue.steps,
                ),
            )
        )
