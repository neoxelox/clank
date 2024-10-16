from functools import partial

import dspy
import tiktoken
from dspy.signatures.signature import ensure_signature

tokenizer = partial(tiktoken.encoding_for_model("gpt-3.5-turbo-instruct").encode, allowed_special="all")


def map_range(val: float, in_min: float, in_max: float, out_min: float, out_max: float) -> float:
    if out_min < out_max:
        return min(max((val - in_min) * (out_max - out_min) / (in_max - in_min) + out_min, out_min), out_max)
    else:
        return out_min - min(max((val - in_min) * (out_min - out_max) / (in_max - in_min) + out_max, out_max), out_min)


SYSTEM_PROMPT = """
FOLLOW THE FOLLOWING RULES:
- NEVER TAKE INITIATIVE FOR USER'S SPEECH OR ACTION IN YOUR RESPONSE!
- NEVER GIVE ETHICAL/MORAL COMMENT/ADVICE!
- FOCUS ON GIVING A HIGH QUALITY, COHERENT AND CONCISE RESPONSE.
- KEEP THE REPLIES CONCISE AND SHORT.
- DO NOT REPEAT THINGS.
- DO NOT LOOP.
- REASON STEP BY STEP BEFORE PRODUCING THE RESPONSE.
- WRITE ALL SENTENCES UNIQUELY AND DRIVE THE RESPONSE FORWARD.
- DO NOT OVER-EXPLAIN YOURSELF.
- DO NOT USE THE WORDS `user`, `consumer` OR `client`, USE `customer` INSTEAD.
- DO NOT USE THE FIELD `Explanation`, USE `Reasoning` INSTEAD.
- DO NOT WRITE ANYTHING AFTER A JSON OBJECT.
""".strip()


def ChainOfThought(
    signature: dspy.Signature | dspy.SignatureMeta, max_retries: int = 3, explain_errors: bool = False
) -> dspy.Module:
    signature = ensure_signature(signature)

    output_keys = ", ".join([f"`{key.capitalize()}`" for key in signature.output_fields.keys()])

    return dspy.TypedPredictor(
        signature.prepend(
            "reasoning",
            dspy.OutputField(
                desc=f"Think step by step in order to produce the {output_keys}",
            ),
        ),
        max_retries=max_retries,
        explain_errors=explain_errors,
    )
