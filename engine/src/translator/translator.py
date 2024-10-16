import json

from pydantic import BaseModel

from src.common import Usage
from src.common.utils import get_tokens
from src.config import Config

from .feedback_translator import FeedbackTranslator
from .language_detector import LanguageDetector


class Translator:
    def __init__(self, config: Config) -> None:
        self.config = config

        self.language_detector = LanguageDetector(config=config)
        self.feedback_translator = FeedbackTranslator(config=config)

    class DetectLanguageParams(BaseModel):
        feedback: str

    class DetectLanguageResult(BaseModel):
        language: str
        usage: Usage

    def detect_language(self, params: DetectLanguageParams) -> DetectLanguageResult:
        language = self.language_detector(
            input=LanguageDetector.Input(
                feedback=params.feedback,
            )
        ).output.language

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.DetectLanguageResult.model_json_schema()) + params.model_dump_json()) + 1500
        )
        output_tokens = get_tokens(language) + 100

        return self.DetectLanguageResult(
            language=language,
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )

    class TranslateFeedbackParams(BaseModel):
        feedback: str
        from_language: str
        to_language: str

    class TranslateFeedbackResult(BaseModel):
        translation: str
        usage: Usage

    def translate_feedback(self, params: TranslateFeedbackParams) -> TranslateFeedbackResult:
        translation = self.feedback_translator(
            input=FeedbackTranslator.Input(
                feedback=params.feedback,
                from_language=params.from_language,
                to_language=params.to_language,
            )
        ).output.translation

        # TODO: Use real usage
        input_tokens = (
            get_tokens(json.dumps(self.TranslateFeedbackResult.model_json_schema()) + params.model_dump_json()) + 1500
        )
        output_tokens = get_tokens(translation) + 100

        return self.TranslateFeedbackResult(
            translation=translation,
            usage=Usage(
                input=input_tokens,
                output=output_tokens,
            ),
        )
