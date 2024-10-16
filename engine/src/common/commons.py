from dspy import Module, OutputField, Signature, SignatureMeta, TypedPredictor
from dspy.signatures.signature import ensure_signature
from pydantic import BaseModel


class Usage(BaseModel):
    input: int
    output: int


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


def ChainOfThought(signature: Signature | SignatureMeta, max_retries: int = 3, explain_errors: bool = False) -> Module:
    signature = ensure_signature(signature)

    output_keys = ", ".join([f"`{key.capitalize()}`" for key in signature.output_fields.keys()])

    return TypedPredictor(
        signature.prepend(
            "reasoning",
            OutputField(
                desc=f"Think step by step in order to produce the {output_keys}",
            ),
        ),
        max_retries=max_retries,
        explain_errors=explain_errors,
    )


UNKNOWN_OPTION = "UNKNOWN"
