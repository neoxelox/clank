from typing import List

from dspy import Assert, InputField, Module, OutputField, Prediction, Signature, Suggest, backtrack_handler
from pydantic import BaseModel, Field

from src.common import UNKNOWN_OPTION, ChainOfThought
from src.common.utils import StrEnum
from src.config import Config


class SuggestionGenerator(Module):
    class Input(BaseModel):
        context: str
        feedback: str

    class Output(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestions: List[Suggestion]

    class GenerateSuggestions(Signature):
        """
List valid improvement proposals, feature requests and ideas, that a customer has about a product (context is provided), from the customer's feedback.
- Suggestions that the customer did not explicitly state are invalid suggestions.
- If the customer is uncertain of a suggestion it is an invalid suggestion.
- Suggestions without reasons behind the proposals are still valid suggestions.
- Issues, concerns, complaints, reviews, opinions or preferences are invalid suggestions.
- Suggestions that come from an issue are invalid suggestions.
- Lexicographic, syntactic, spelling, grammar or any other language mistakes of the feedback's text are invalid suggestions.
- Again, a suggestion cannot be supposed to be valid if the customer did not explicitly state it.
        """  # fmt: skip

        class Input(BaseModel):
            context: str
            feedback: str

        class Output(BaseModel):
            class Suggestion(BaseModel):
                title: str = Field(
                    description="4 to 10 words, which cannot contain the words `suggestion` (or synonyms), `customer` (or synonyms) or the product's name.",
                    max_length=100,
                )
                description: str = Field(
                    description="Long, complete explanation, but without redundant information, using the feedback's original words. Must focus solely on the suggestion by depersonalizing the sentences."
                )
                reason: str = Field(
                    description=f"The customer's motivation behind the proposal of the suggestion, if any must always start with `This will`, else `{UNKNOWN_OPTION}`."
                )

            suggestions: List[Suggestion] = Field(description="If any, else `[]`.")

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.generate_suggestions = ChainOfThought(self.GenerateSuggestions, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/processor/suggestion_extractor/suggestion_generator.json")

    def forward(self, input: Input) -> Prediction:
        suggestions = self.generate_suggestions(
            input=self.GenerateSuggestions.Input(
                context=input.context,
                feedback=input.feedback,
            )
        ).output.suggestions

        return Prediction(
            output=self.Output(
                suggestions=[
                    self.Output.Suggestion(
                        title=suggestion.title,
                        description=suggestion.description,
                        reason=suggestion.reason if suggestion.reason.upper() != UNKNOWN_OPTION else "",
                    )
                    for suggestion in suggestions
                ],
            )
        )


class InfoInferrer(Module):
    class Input(BaseModel):
        context: str
        feedback: str
        suggestion: str
        categories: List[str]

    class Output(BaseModel):
        importance: str
        category: str

    class InferInfo(Signature):
        """
Infer the following information from a suggestion that an LLM extracted from the feedback of a customer.
- Discern the importance, valid options (`importances`) are provided.
- Discern the category, valid options (`categories`) are provided.
        """  # fmt: skip

        class Input(BaseModel):
            context: str
            feedback: str
            suggestion: str
            importances: List[str]
            categories: List[str]

        class Output(BaseModel):
            class Importance(StrEnum):
                CRITICAL = "CRITICAL"
                HIGH = "HIGH"
                MEDIUM = "MEDIUM"
                LOW = "LOW"

            importance: str = Field(description="The valid option that best fits.")
            category: str = Field(description=f"The valid option that best fits, if any, else `{UNKNOWN_OPTION}`.")

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.infer_info = ChainOfThought(self.InferInfo, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/processor/suggestion_extractor/info_inferrer.json")

    def forward(self, input: Input) -> Prediction:
        info = self.infer_info(
            input=self.InferInfo.Input(
                context=input.context,
                feedback=input.feedback,
                suggestion=input.suggestion,
                importances=[importance.replace("_", " ") for importance in self.InferInfo.Output.Importance.list()],
                categories=list({category.replace("_", " ") for category in input.categories + [UNKNOWN_OPTION]}),
            )
        ).output

        importance = info.importance.upper().replace(" ", "_")

        Assert(
            importance in self.InferInfo.Output.Importance.list(),
            f'Importance must be {self.InferInfo.Output.model_fields["importance"].description}! `{importance}` is NOT a valid option. Valid options are:\n'
            + "".join([f"- {option}\n" for option in self.InferInfo.Output.Importance.list()]),
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
                importance=importance,
                category=category,
            )
        )


class SuggestionExtractor(Module):
    class Input(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class Output(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str
            importance: str
            category: str

        suggestions: List[Suggestion]

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.SuggestionGenerator = SuggestionGenerator
        self.generate_suggestions = self.SuggestionGenerator(config=config)
        self.InfoInferrer = InfoInferrer
        self.infer_info = self.InfoInferrer(config=config)

    def forward(self, input: Input) -> Prediction:
        generated_suggestions = self.generate_suggestions(
            input=self.SuggestionGenerator.Input(
                context=input.context,
                feedback=input.feedback,
            )
        ).output.suggestions

        suggestions = []
        for suggestion in generated_suggestions:
            info = self.infer_info(
                input=self.InfoInferrer.Input(
                    context=input.context,
                    feedback=input.feedback,
                    suggestion=suggestion.description,
                    categories=input.categories,
                )
            ).output

            suggestions.append(
                self.Output.Suggestion(
                    title=suggestion.title,
                    description=suggestion.description,
                    reason=suggestion.reason,
                    importance=info.importance,
                    category=info.category,
                )
            )

        return Prediction(
            output=self.Output(
                suggestions=suggestions,
            )
        )
