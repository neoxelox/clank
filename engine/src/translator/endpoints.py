from pydantic import BaseModel

from src.common import Usage
from src.config import Config

from .translator import Translator


class TranslatorEndpoints:
    def __init__(self, config: Config, translator: Translator) -> None:
        self.config = config
        self.translator = translator

    class PostDetectLanguageRequest(BaseModel):
        feedback: str

    class PostDetectLanguageResponse(BaseModel):
        language: str
        usage: Usage

    async def post_detect_language(self, request: PostDetectLanguageRequest) -> PostDetectLanguageResponse:
        result = self.translator.detect_language(
            params=Translator.DetectLanguageParams(
                feedback=request.feedback,
            )
        )

        return self.PostDetectLanguageResponse(
            language=result.language,
            usage=result.usage,
        )

    class PostTranslateFeedbackRequest(BaseModel):
        feedback: str
        from_language: str
        to_language: str

    class PostTranslateFeedbackResponse(BaseModel):
        translation: str
        usage: Usage

    async def post_translate_feedback(self, request: PostTranslateFeedbackRequest) -> PostTranslateFeedbackResponse:
        result = self.translator.translate_feedback(
            params=Translator.TranslateFeedbackParams(
                feedback=request.feedback,
                from_language=request.from_language,
                to_language=request.to_language,
            )
        )

        return self.PostTranslateFeedbackResponse(
            translation=result.translation,
            usage=result.usage,
        )
