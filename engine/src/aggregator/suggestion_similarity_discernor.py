from typing import List

from dspy import InputField, Module, OutputField, Prediction, Signature, backtrack_handler
from pydantic import BaseModel

from src.common import ChainOfThought
from src.config import Config


class SuggestionSimilarityDiscernor(Module):
    class Input(BaseModel):
        suggestion: str
        options: List[str]

    class Output(BaseModel):
        index: int

    class DiscernSimilarity(Signature):
        """
Discern whether suggestion A and suggestion B, that customers have about a product, are similar or not.
- Both suggestions are similar only if they are at least 80% similar.
- Customers can have similar suggestions without writing them the same way.
        """  # fmt: skip

        class Input(BaseModel):
            suggestion_a: str
            suggestion_b: str

        class Output(BaseModel):
            similar: bool

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.discern_similarity = ChainOfThought(self.DiscernSimilarity, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(
            f"{config.service.artifacts_path}/aggregator/suggestion_aggregator/suggestion_similarity_discernor.json"
        )

    def forward(self, input: Input) -> Prediction:
        for index, option in enumerate(input.options):
            if option == input.suggestion:
                return Prediction(output=self.Output(index=index))

            similar = self.discern_similarity(
                input=self.DiscernSimilarity.Input(
                    suggestion_a=input.suggestion,
                    suggestion_b=option,
                )
            ).output.similar

            if similar:
                return Prediction(output=self.Output(index=index))

        return Prediction(output=self.Output(index=-1))
