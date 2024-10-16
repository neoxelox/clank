from dspy import InputField, Module, OutputField, Prediction, Signature, Suggest, backtrack_handler
from emoji import replace_emoji
from lingua import Language, LanguageDetectorBuilder
from pydantic import BaseModel, Field

from src.common import UNKNOWN_OPTION, ChainOfThought
from src.config import Config


class LanguageDetector(Module):
    class Input(BaseModel):
        feedback: str

    class Output(BaseModel):
        language: str

    class DetectLanguage(Signature):
        """
Detect the language from the customer's feedback.
- If the words are common in many languages including English, default to English.
- If there are lexicographic, syntactic, spelling, grammar or any other language mistakes, default to the most probable language.
        """  # fmt: skip

        class Input(BaseModel):
            feedback: str

        class Output(BaseModel):
            language: str = Field(
                description=f"The full name of the valid language that best fits, if any, else `{UNKNOWN_OPTION}`."
            )

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.LANGUAGES = [language.name.upper() for language in Language.all()]
        self.DEFAULT_LANGUAGE = "ENGLISH"
        self.MINIMUM_LENGTH = 1
        self.MINIMUM_CONFIDENCE = 0.25
        self.MINIMUM_CONFIDENCE_DISTANCE = self.MINIMUM_CONFIDENCE / 2

        self.detector = (
            LanguageDetectorBuilder.from_languages(
                Language.ENGLISH,
                Language.SPANISH,
                Language.FRENCH,
                Language.PORTUGUESE,
                Language.GERMAN,
                Language.ITALIAN,
            )
            .with_preloaded_language_models()
            .build()
        )

        self.detect_language = ChainOfThought(self.DetectLanguage, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/translator/language_detector.json")

    def forward(self, input: Input) -> Prediction:
        feedback = replace_emoji(input.feedback)

        if len(feedback.strip()) < self.MINIMUM_LENGTH:
            return Prediction(
                output=self.Output(
                    language=self.DEFAULT_LANGUAGE,
                )
            )

        confidence_values = self.detector.compute_language_confidence_values(feedback)
        most_likely = confidence_values[0]
        second_most_likely = confidence_values[1]

        if (
            most_likely.value < self.MINIMUM_CONFIDENCE
            or (most_likely.value - second_most_likely.value) < self.MINIMUM_CONFIDENCE_DISTANCE
        ):
            language = self.detect_language(
                input=self.DetectLanguage.Input(
                    feedback=feedback,
                )
            ).output.language

            language = language.upper()

            Suggest(
                language in self.LANGUAGES or language == UNKNOWN_OPTION,
                f'Language must be {self.DetectLanguage.Output.model_fields["language"].description}! `{language}` is NOT a valid language. Valid languages are:\n'
                + "".join([f"- {option}\n" for option in self.LANGUAGES]),
            )

            if language not in self.LANGUAGES:
                language = UNKNOWN_OPTION

            return Prediction(
                output=self.Output(
                    language=language,
                )
            )

        return Prediction(
            output=self.Output(
                language=most_likely.language.name.upper(),
            )
        )
