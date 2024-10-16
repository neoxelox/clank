from dspy import Assert, InputField, Module, OutputField, Prediction, Signature, backtrack_handler
from pydantic import BaseModel

from src.common import ChainOfThought
from src.config import Config


class FeedbackTranslator(Module):
    class Input(BaseModel):
        feedback: str
        from_language: str
        to_language: str

    class Output(BaseModel):
        translation: str

    class TranslateFeedback(Signature):
        """
Translate the customer's feedback from a language to a language.
Maintain the feedback's:
- Style
- Format (including newlines and tabs)
- Emphasis
- Emojis
- Punctuation
- Names
- Measures
- Units
- Dates (use the translated format)
        """  # fmt: skip

        class Input(BaseModel):
            feedback: str
            from_language: str
            to_language: str

        class Output(BaseModel):
            translation: str

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.MAXIMUM_DIFFERENCE = 50

        self.translate_feedback = ChainOfThought(self.TranslateFeedback, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/translator/feedback_translator.json")

    def forward(self, input: Input) -> Prediction:
        if input.from_language == input.to_language:
            return Prediction(
                output=self.Output(
                    translation=input.feedback,
                )
            )

        translation = self.translate_feedback(
            input=self.TranslateFeedback.Input(
                feedback=input.feedback,
                from_language=input.from_language,
                to_language=input.to_language,
            )
        ).output.translation

        Assert(translation != "", "Translation cannot be empty!")

        return Prediction(
            output=self.Output(
                translation=translation,
            )
        )
