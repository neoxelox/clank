from dspy import InputField, Module, OutputField, Prediction, Signature, backtrack_handler
from pydantic import BaseModel, Field

from src.common import UNKNOWN_OPTION, ChainOfThought
from src.config import Config


class SuggestionMerger(Module):
    class Input(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion_a: Suggestion
        suggestion_b: Suggestion

    class Output(BaseModel):
        class Suggestion(BaseModel):
            title: str
            description: str
            reason: str

        suggestion: Suggestion

    class MergeSuggestions(Signature):
        """
Merge, coherently, suggestion A and suggestion B, that customers have about a product, into a single suggestion.
- Maintain the core idea, context and nuances of both suggestions.
- Do not create information that is not present in any of the suggestions.
        """  # fmt: skip

        class Input(BaseModel):
            class Suggestion(BaseModel):
                title: str
                description: str
                reason: str

            suggestion_a: Suggestion
            suggestion_b: Suggestion

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

            suggestion: Suggestion

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.merge_suggestions = ChainOfThought(self.MergeSuggestions, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/aggregator/suggestion_aggregator/suggestion_merger.json")

    def forward(self, input: Input) -> Prediction:
        if input.suggestion_a == input.suggestion_b:
            return Prediction(
                output=self.Output(
                    suggestion=self.Output.Suggestion(
                        title=input.suggestion_a.title,
                        description=input.suggestion_a.description,
                        reason=input.suggestion_a.reason,
                    ),
                )
            )

        suggestion = self.merge_suggestions(
            input=self.MergeSuggestions.Input(
                suggestion_a=self.MergeSuggestions.Input.Suggestion(
                    title=input.suggestion_a.title,
                    description=input.suggestion_a.description,
                    reason=input.suggestion_a.reason,
                ),
                suggestion_b=self.MergeSuggestions.Input.Suggestion(
                    title=input.suggestion_b.title,
                    description=input.suggestion_b.description,
                    reason=input.suggestion_b.reason,
                ),
            )
        ).output.suggestion

        return Prediction(
            output=self.Output(
                suggestion=self.Output.Suggestion(
                    title=suggestion.title,
                    description=suggestion.description,
                    reason=suggestion.reason if suggestion.reason.upper() != UNKNOWN_OPTION else "",
                ),
            )
        )
