import os

# Turn off DSPy cache before importing
os.environ["DSP_CACHEBOOL"] = "FALSE"

# Set LiteLLM production mode before importing
os.environ["LITELLM_MODE"] = "PRODUCTION"

# ruff: noqa: E402
import dspy
import sentry_sdk
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException

from src.aggregator import Aggregator, AggregatorEndpoints
from src.common import SYSTEM_PROMPT
from src.config import get_config
from src.exception_handler import ExceptionHandler
from src.health import HealthEndpoints
from src.logger_middleware import LoggerMiddleware
from src.processor import Processor, ProcessorEndpoints
from src.translator import Translator, TranslatorEndpoints

config = get_config()

if config.sentry.dsn:
    sentry_sdk.init(
        dsn=config.sentry.dsn,
        environment=config.service.environment,
        release=config.service.release,
        server_name=config.service.name,
        debug=False,
        attach_stacktrace=True,
        enable_tracing=True,
        sample_rate=1.0,
        traces_sample_rate=0.25,
        profiles_sample_rate=1.0,
    )

if config.service.environment == "dev":
    import phoenix
    from openinference.instrumentation.dspy import DSPyInstrumentor
    from opentelemetry import trace as trace_api
    from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
    from opentelemetry.sdk import trace as trace_sdk
    from opentelemetry.sdk.trace.export import SimpleSpanProcessor

    phoenix.launch_app(host="0.0.0.0", port=config.phoenix.port)
    tracer_provider = trace_sdk.TracerProvider()
    tracer_provider.add_span_processor(
        SimpleSpanProcessor(OTLPSpanExporter(endpoint=f"http://{config.phoenix.host}:{config.phoenix.port}/v1/traces"))
    )
    trace_api.set_tracer_provider(tracer_provider)
    DSPyInstrumentor().instrument()

# LM BACKENDS

backends = {}

if config.lm.openai.api_key:
    backends["openai/gpt-4o"] = {
        "model": "openai/gpt-4o",
        "api_key": config.lm.openai.api_key,
    }

if config.lm.groq.api_key:
    backends["groq/llama3-8b"] = {
        "model": "groq/llama3-8b-8192",
        "api_key": config.lm.groq.api_key,
    }

if config.lm.anyscale.api_key:
    backends["anyscale/llama3-8b"] = {
        "model": "anyscale/meta-llama/Meta-Llama-3-8B-Instruct",
        "api_key": config.lm.anyscale.api_key,
        "api_base": "https://api.endpoints.anyscale.com/v1",
    }

if config.lm.fireworks.api_key:
    backends["fireworks/llama3.1-8b"] = {
        "model": "fireworks_ai/accounts/fireworks/models/llama-v3p1-8b-instruct",
        "api_key": config.lm.fireworks.api_key,
    }

dspy.configure(
    backend=dspy.ChatBackend(
        **backends[config.lm.backend],
        params={
            "max_tokens": config.lm.max_tokens,
            "temperature": config.lm.temperature,
        },
        attempts=3,
        system_prompt=SYSTEM_PROMPT,
    ),
    trace=[],
    cache=False,
)

# API SERVER

server = FastAPI(
    debug=False,
    title="",
    summary=None,
    description="",
    version="",
    openapi_url=None,
    redirect_slashes=True,
    docs_url=None,
    redoc_url=None,
    swagger_ui_oauth2_redirect_url=None,
    swagger_ui_init_oauth=None,
)

exception_handler = ExceptionHandler(config=config)
server.add_exception_handler(Exception, exception_handler.handle)
server.add_exception_handler(HTTPException, exception_handler.handle)
server.add_exception_handler(RequestValidationError, exception_handler.handle)

logger_middleware = LoggerMiddleware(config=config)
server.middleware("http")(logger_middleware.handle)

# REPOSITORIES

# SERVICES

# USECASES

translator = Translator(config=config)
processor = Processor(config=config)
aggregator = Aggregator(config=config)

# ENDPOINTS

health_endpoints = HealthEndpoints(config=config)
translator_endpoints = TranslatorEndpoints(config=config, translator=translator)
processor_endpoints = ProcessorEndpoints(config=config, processor=processor)
aggregator_endpoints = AggregatorEndpoints(config=config, aggregator=aggregator)

# MIDDLEWARES

# ROUTES

server.get("/health")(health_endpoints.get_health)

server.post("/translator/detect-language")(translator_endpoints.post_detect_language)
server.post("/translator/translate-feedback")(translator_endpoints.post_translate_feedback)

server.post("/processor/extract-issues")(processor_endpoints.post_extract_issues)
server.post("/processor/extract-suggestions")(processor_endpoints.post_extract_suggestions)
server.post("/processor/extract-review")(processor_endpoints.post_extract_review)

server.post("/aggregator/compute-embedding")(aggregator_endpoints.post_compute_embedding)
server.post("/aggregator/similar-issue")(aggregator_endpoints.post_similar_issue)
server.post("/aggregator/merge-issues")(aggregator_endpoints.post_merge_issues)
server.post("/aggregator/similar-suggestion")(aggregator_endpoints.post_similar_suggestion)
server.post("/aggregator/merge-suggestions")(aggregator_endpoints.post_merge_suggestions)
