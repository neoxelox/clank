from dataclasses import dataclass
from enum import Enum
from typing import List


class StrEnum(str, Enum):
    def __str__(self) -> str:
        return str(self.value)

    @classmethod
    def list(cls) -> List[str]:
        return list(map(lambda v: v.value, cls))


UNKNOWN_OPTION = "UNKNOWN"


class Source(StrEnum):
    TRUSTPILOT = "TRUSTPILOT"
    APP_STORE = "APP_STORE"
    PLAY_STORE = "PLAY_STORE"
    AMAZON = "AMAZON"


@dataclass
class Feedback:
    id: str
    hash: str
    context: str
    categories: List[str]
    source: str
    customer: str
    content: str


class Severity(StrEnum):
    CRITICAL = "CRITICAL"
    HIGH = "HIGH"
    MEDIUM = "MEDIUM"
    LOW = "LOW"


@dataclass
class Issue:
    id: str
    embedding: List[float]
    title: str  # (max100)
    description: str
    steps: List[str]
    severity: str
    category: str  # (optUNKNOWN) User-defined LLM-inferred


class Importance(StrEnum):
    CRITICAL = "CRITICAL"
    HIGH = "HIGH"
    MEDIUM = "MEDIUM"
    LOW = "LOW"


@dataclass
class Suggestion:
    id: str
    embedding: List[float]
    title: str  # (max100)
    description: str
    reason: str
    importance: str
    category: str  # (optUNKNOWN) User-defined LLM-inferred


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
    RETAIN = "RETAIN"  # BUY_AGAIN / RENEW / RECOMMEND
    CHURN = "CHURN"  # RETURN / REFUND / CANCEL / DISCOURAGE
    RETAIN_AND_RECOMMEND = "RETAIN_AND_RECOMMEND"
    CHURN_AND_DISCOURAGE = "CHURN_AND_DISCOURAGE"


@dataclass
class Review:
    id: str
    content: str
    keywords: List[str]  # max(10)
    sentiment: str
    emotions: List[str]  # (max4)
    intention: str  # (optUNKNOWN)
    category: str  # (optUNKNOWN) User-defined LLM-inferred


@dataclass
class LabeledFeedback(Feedback):
    language: str  # TODO: Sure?: ENGLISH | SPANISH | FRENCH | ITALIAN | GERMAN | PORTUGUESE | UNKNOWN
    translation: str
    issues: List[Issue]
    suggestions: List[Suggestion]
    review: Review
