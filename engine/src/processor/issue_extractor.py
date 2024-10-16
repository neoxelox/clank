from typing import List

from dspy import Assert, InputField, Module, OutputField, Prediction, Signature, Suggest, backtrack_handler
from pydantic import BaseModel, Field

from src.common import UNKNOWN_OPTION, ChainOfThought
from src.common.utils import StrEnum
from src.config import Config


class IssueGenerator(Module):
    class Input(BaseModel):
        context: str
        feedback: str

    class Output(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]

        issues: List[Issue]

    class GenerateIssues(Signature):
        """
List valid issues, that a customer has with a product (context is provided), from the customer's feedback.
- Issues that the customer did not explicitly state are invalid issues.
- If the customer is uncertain of an issue it is an invalid issue.
- Issues without steps to reproduce them are still valid issues.
- Suggestions, reviews, opinions or preferences are invalid issues.
- Lexicographic, syntactic, spelling, grammar or any other language mistakes of the feedback's text are invalid issues.
- Again, an issue cannot be supposed to be valid if the customer did not explicitly state it.
        """  # fmt: skip

        class Input(BaseModel):
            context: str
            feedback: str

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

            issues: List[Issue] = Field(description="If any, else `[]`.")

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.generate_issues = ChainOfThought(self.GenerateIssues, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/processor/issue_extractor/issue_generator.json")

    def forward(self, input: Input) -> Prediction:
        issues = self.generate_issues(
            input=self.GenerateIssues.Input(
                context=input.context,
                feedback=input.feedback,
            )
        ).output.issues

        return Prediction(
            output=self.Output(
                issues=[
                    self.Output.Issue(
                        title=issue.title,
                        description=issue.description,
                        steps=issue.steps,
                    )
                    for issue in issues
                ],
            )
        )


class InfoInferrer(Module):
    class Input(BaseModel):
        context: str
        feedback: str
        issue: str
        categories: List[str]

    class Output(BaseModel):
        severity: str
        category: str

    class InferInfo(Signature):
        """
Infer the following information from an issue that an LLM extracted from the feedback of a customer.
- Discern the severity, valid options (`severities`) are provided.
- Discern the category, valid options (`categories`) are provided.
        """  # fmt: skip

        class Input(BaseModel):
            context: str
            feedback: str
            issue: str
            severities: List[str]
            categories: List[str]

        class Output(BaseModel):
            class Severity(StrEnum):
                CRITICAL = "CRITICAL"
                HIGH = "HIGH"
                MEDIUM = "MEDIUM"
                LOW = "LOW"

            severity: str = Field(description="The valid option that best fits.")
            category: str = Field(description=f"The valid option that best fits, if any, else `{UNKNOWN_OPTION}`.")

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.infer_info = ChainOfThought(self.InferInfo, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/processor/issue_extractor/info_inferrer.json")

    def forward(self, input: Input) -> Prediction:
        info = self.infer_info(
            input=self.InferInfo.Input(
                context=input.context,
                feedback=input.feedback,
                issue=input.issue,
                severities=[severity.replace("_", " ") for severity in self.InferInfo.Output.Severity.list()],
                categories=list({category.replace("_", " ") for category in input.categories + [UNKNOWN_OPTION]}),
            )
        ).output

        severity = info.severity.upper().replace(" ", "_")

        Assert(
            severity in self.InferInfo.Output.Severity.list(),
            f'Severity must be {self.InferInfo.Output.model_fields["severity"].description}! `{severity}` is NOT a valid option. Valid options are:\n'
            + "".join([f"- {option}\n" for option in self.InferInfo.Output.Severity.list()]),
        )

        category = info.category.upper().replace(" ", "_")
        if not input.categories:
            category = UNKNOWN_OPTION

        Suggest(
            category in input.categories or category == UNKNOWN_OPTION,
            f'Category must be {self.InferInfo.Output.model_fields["category"].description}! `{category}` is NOT a valid option. Valid options are:\n'
            + "".join([f"- {option}\n" for option in input.categories]),
        )

        if category not in input.categories:
            category = UNKNOWN_OPTION

        return Prediction(
            output=self.Output(
                severity=severity,
                category=category,
            )
        )


class IssueExtractor(Module):
    class Input(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class Output(BaseModel):
        class Issue(BaseModel):
            title: str
            description: str
            steps: List[str]
            severity: str
            category: str

        issues: List[Issue]

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.IssueGenerator = IssueGenerator
        self.generate_issues = self.IssueGenerator(config=config)
        self.InfoInferrer = InfoInferrer
        self.infer_info = self.InfoInferrer(config=config)

    def forward(self, input: Input) -> Prediction:
        generated_issues = self.generate_issues(
            input=self.IssueGenerator.Input(
                context=input.context,
                feedback=input.feedback,
            )
        ).output.issues

        issues = []
        for issue in generated_issues:
            info = self.infer_info(
                input=self.InfoInferrer.Input(
                    context=input.context,
                    feedback=input.feedback,
                    issue=issue.description,
                    categories=input.categories,
                )
            ).output

            issues.append(
                self.Output.Issue(
                    title=issue.title,
                    description=issue.description,
                    steps=issue.steps,
                    severity=info.severity,
                    category=info.category,
                )
            )

        return Prediction(
            output=self.Output(
                issues=issues,
            )
        )
