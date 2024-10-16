import os
from enum import Enum
from functools import partial
from typing import List, TypeVar

from tiktoken import encoding_for_model


class StrEnum(str, Enum):
    def __str__(self) -> str:
        return str(self.value)

    @classmethod
    def list(cls) -> List[str]:
        return list(map(lambda v: v.value, cls))


T = TypeVar("T", str, int, bool, List[str])


def get_env(key: str, default: T) -> T:
    value = os.getenv(key)
    if not value:
        return default

    if isinstance(default, str):
        return value

    elif isinstance(default, int):
        return int(value)

    elif isinstance(default, bool):
        return bool(value)

    elif isinstance(default, list):
        return value.split(",")

    raise TypeError()


tokenizer = partial(encoding_for_model("gpt-3.5-turbo-instruct").encode, allowed_special="all")


def get_tokens(text: str) -> int:
    return len(tokenizer(text))


def map_range(val: float, in_min: float, in_max: float, out_min: float, out_max: float) -> float:
    if out_min < out_max:
        return min(max((val - in_min) * (out_max - out_min) / (in_max - in_min) + out_min, out_min), out_max)
    else:
        return out_min - min(max((val - in_min) * (out_min - out_max) / (in_max - in_min) + out_max, out_max), out_min)
