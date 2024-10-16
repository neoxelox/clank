from pydantic import BaseModel

from src.common.utils import get_env


class Config(BaseModel):
    class Service(BaseModel):
        environment: str
        name: str
        release: str
        timezone: str
        graceful_timeout: int
        resources_path: str
        artifacts_path: str

    class Server(BaseModel):
        host: str
        port: int

    class LM(BaseModel):
        backend: str
        max_tokens: int
        temperature: float

        class OpenAI(BaseModel):
            api_key: str

        class Groq(BaseModel):
            api_key: str

        class Anyscale(BaseModel):
            api_key: str

        class Fireworks(BaseModel):
            api_key: str

        openai: OpenAI
        groq: Groq
        anyscale: Anyscale
        fireworks: Fireworks

    class Phoenix(BaseModel):
        host: str
        port: int

    class Sentry(BaseModel):
        dsn: str

    service: Service
    server: Server
    lm: LM
    phoenix: Phoenix
    sentry: Sentry


def get_config() -> Config:
    return Config(
        service=Config.Service(
            environment=get_env("CLANK_ENVIRONMENT", "dev"),
            name="engine",
            release=get_env("CLANK_RELEASE", "wip"),
            timezone=get_env("CLANK_TIMEZONE", "UTC"),
            graceful_timeout=30,
            resources_path="resources",
            artifacts_path="artifacts",
        ),
        server=Config.Server(
            host=get_env("CLANK_ENGINE_HOST", "localhost"),
            port=get_env("CLANK_ENGINE_PORT", 2222),
        ),
        lm=Config.LM(
            backend=get_env("CLANK_ENGINE_LM_BACKEND", ""),
            max_tokens=1024,
            temperature=0.7,
            openai=Config.LM.OpenAI(
                api_key=get_env("CLANK_OPENAI_API_KEY", ""),
            ),
            groq=Config.LM.Groq(
                api_key=get_env("CLANK_GROQ_API_KEY", ""),
            ),
            anyscale=Config.LM.Anyscale(
                api_key=get_env("CLANK_ANYSCALE_API_KEY", ""),
            ),
            fireworks=Config.LM.Fireworks(
                api_key=get_env("CLANK_FIREWORKS_API_KEY", ""),
            ),
        ),
        phoenix=Config.Phoenix(
            host=get_env("CLANK_ENGINE_PHOENIX_HOST", "localhost"),
            port=get_env("CLANK_ENGINE_PHOENIX_PORT", 2223),
        ),
        sentry=Config.Sentry(
            dsn=get_env("CLANK_ENGINE_SENTRY_DSN", ""),
        ),
    )
