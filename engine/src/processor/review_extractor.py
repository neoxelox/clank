from typing import List

from dspy import Assert, InputField, Module, OutputField, Prediction, Signature, Suggest, backtrack_handler
from emoji import emoji_count
from pydantic import BaseModel, Field

from src.common import UNKNOWN_OPTION, ChainOfThought
from src.common.utils import StrEnum
from src.config import Config


class ReviewExtractor(Module):
    class Input(BaseModel):
        context: str
        categories: List[str]
        feedback: str

    class Output(BaseModel):
        class Review(BaseModel):
            content: str
            keywords: List[str]
            sentiment: str
            emotions: List[str]
            intention: str
            category: str

        review: Review

    class InferInfo(Signature):
        """
Infer the following information from the customer's feedback of a product (context is provided).
- List the most important keywords, following the following rules:
    - Limit each keyword to 3 words maximum.
    - Only include keywords the customer explicitly stated.
    - Do not include emojis in the keywords.
    - Do not include the name of the product in the keywords.
- Discern the sentiment, valid options (`sentiments`) are provided.
- Discern the emotions, valid options (`emotions`) are provided.
- Discern the intention, valid options (`intentions`) are provided, following the following rules:
    - To `retain` means to have the intention to buy again, renew and/or recommend the product.
    - To `churn` means to have the intention to return, refund, cancel and/or discourage the product. Critical issues also cause customer churn.
    - To `recommend` cannot be assumed if the customer did not explicitly state it (except if synonyms were used).
    - To `discourage` is very likely if the customer has the intention to churn.
- Discern the category, valid options (`categories`) are provided.
        """  # fmt: skip

        class Input(BaseModel):
            context: str
            feedback: str
            sentiments: List[str]
            emotions: List[str]
            intentions: List[str]
            categories: List[str]

        class Output(BaseModel):
            class Sentiment(StrEnum):
                POSITIVE = "POSITIVE"
                NEUTRAL = "NEUTRAL"
                NEGATIVE = "NEGATIVE"

            class Emotion(StrEnum):
                TRUST = "TRUST"
                ACCEPTANCE = "ACCEPTANCE"
                FEAR = "FEAR"
                APPREHENSION = "APPREHENSION"
                SURPRISE = "SURPRISE"
                DISTRACTION = "DISTRACTION"
                SADNESS = "SADNESS"
                PENSIVENESS = "PENSIVENESS"
                DISGUST = "DISGUST"
                BOREDOM = "BOREDOM"
                ANGER = "ANGER"
                ANNOYANCE = "ANNOYANCE"
                ANTICIPATION = "ANTICIPATION"
                INTEREST = "INTEREST"
                JOY = "JOY"
                SERENITY = "SERENITY"

            class Intention(StrEnum):
                RETAIN = "RETAIN"
                CHURN = "CHURN"
                RETAIN_AND_RECOMMEND = "RETAIN_AND_RECOMMEND"
                CHURN_AND_DISCOURAGE = "CHURN_AND_DISCOURAGE"

            keywords: List[str] = Field(description="If any, else `[]`.", max_items=10)
            sentiment: str = Field(description="The valid option that best fits.")
            emotions: List[str] = Field(description="The valid options that best fit, if any, else `[]`.", max_items=4)
            intention: str = Field(description=f"The valid option that best fits, if any, else `{UNKNOWN_OPTION}`.")
            category: str = Field(description=f"The valid option that best fits, if any, else `{UNKNOWN_OPTION}`.")

        input: Input = InputField()
        output: Output = OutputField()

    def __init__(self, config: Config) -> None:
        super().__init__()

        self.infer_info = ChainOfThought(self.InferInfo, max_retries=3, explain_errors=False)

        self.activate_assertions(handler=backtrack_handler, max_backtracks=3)
        self.load(f"{config.service.artifacts_path}/processor/review_extractor.json")

    def forward(self, input: Input) -> Prediction:
        info = self.infer_info(
            input=self.InferInfo.Input(
                context=input.context,
                feedback=input.feedback,
                sentiments=[sentiment.replace("_", " ") for sentiment in self.InferInfo.Output.Sentiment.list()],
                emotions=[emotion.replace("_", " ") for emotion in self.InferInfo.Output.Emotion.list()],
                intentions=[
                    intention.replace("_", " ")
                    for intention in self.InferInfo.Output.Intention.list() + [UNKNOWN_OPTION]
                ],
                categories=list({category.replace("_", " ") for category in input.categories + [UNKNOWN_OPTION]}),
            )
        ).output

        keywords = [keyword.lower() for keyword in info.keywords if keyword]

        # TODO: Decide whether to keep this rule because it fails a lot!
        # inexistent_keywords = list(filter(lambda keyword: keyword not in input.feedback.lower(), keywords))
        # Suggest(
        #     not inexistent_keywords,
        #     "All keywords must be included in the customer's feedback! Keywords not included:\n"
        #     + "".join([f"- {keyword}\n" for keyword in inexistent_keywords]),
        # )

        # keywords = list(set(keywords) - set(inexistent_keywords))

        long_keywords = list(filter(lambda keyword: len(keyword.split()) > 3, keywords))
        Suggest(
            not long_keywords,
            "Each keyword must be 3 words maximum! Keywords too long:\n"
            + "".join([f"- {keyword}\n" for keyword in long_keywords]),
        )

        keywords = list(set(keywords) - set(long_keywords))

        emoji_keywords = list(filter(lambda keyword: emoji_count(keyword) > 0, keywords))
        Suggest(
            not emoji_keywords,
            "Keywords cannot include emojis! Keywords with emojis:\n"
            + "".join([f"- {keyword}\n" for keyword in emoji_keywords]),
        )

        keywords = list(set(keywords) - set(emoji_keywords))

        sentiment = info.sentiment.upper().replace(" ", "_")

        Assert(
            sentiment in self.InferInfo.Output.Sentiment.list(),
            f'Sentiment must be {self.InferInfo.Output.model_fields["sentiment"].description}! `{sentiment}` is NOT a valid option. Valid options are:\n'
            + "".join([f"- {option}\n" for option in self.InferInfo.Output.Sentiment.list()]),
        )

        emotions = [emotion.upper().replace(" ", "_") for emotion in info.emotions]

        invalid_emotions = list(filter(lambda emotion: emotion not in self.InferInfo.Output.Emotion.list(), emotions))
        Assert(
            not invalid_emotions,
            f'Emotions must be {self.InferInfo.Output.model_fields["emotions"].description}! Invalid options:\n'
            + "".join([f"- {emotion}\n" for emotion in invalid_emotions])
            + "Valid options are:\n"
            + "".join([f"- {option}\n" for option in self.InferInfo.Output.Emotion.list()]),
        )

        emotions = list(set(emotions))

        intention = info.intention.upper().replace(" ", "_")

        Suggest(
            intention in self.InferInfo.Output.Intention.list() or intention == UNKNOWN_OPTION,
            f'Intention must be {self.InferInfo.Output.model_fields["intention"].description}! `{intention}` is NOT a valid option. Valid options are:\n'
            + "".join([f"- {option}\n" for option in self.InferInfo.Output.Intention.list()]),
        )

        if intention not in self.InferInfo.Output.Intention.list():
            intention = UNKNOWN_OPTION

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
                review=self.Output.Review(
                    content=input.feedback,
                    keywords=keywords,
                    sentiment=sentiment,
                    emotions=emotions,
                    intention=intention,
                    category=category,
                ),
            )
        )
