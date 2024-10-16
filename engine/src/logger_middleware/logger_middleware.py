import json
import logging
import time
from typing import Callable

from fastapi import Request, Response

from src.config import Config


class LoggerMiddleware:
    def __init__(self, config: Config) -> None:
        self.config = config

        level = logging.INFO
        if self.config.service.environment == "dev":
            level = logging.DEBUG

        logging.basicConfig(level=level, force=True)

        for logger in [
            "uvicorn.error",
            "uvicorn.access",
            "uvicorn",
            "fastapi",
            "httpx",
            "httpcore",
            "urllib3.connectionpool",
            "aiosqlite",
        ]:
            logging.getLogger(logger).setLevel(level=logging.CRITICAL)

    async def handle(self, request: Request, call_next: Callable) -> Response:
        start = time.time()

        response = await call_next(request)

        stop = time.time()

        logging.info(
            json.dumps(
                {
                    "host": request.url.hostname,
                    "method": request.method,
                    "path": request.url.path,
                    "status": response.status_code,
                    "ip_address": request.client.host,
                    "latency": round((stop - start) * 1000, 2),
                    "timestamp": int(round(stop, 2)),
                    # "trace_id": "TODO",
                }
            )
        )

        return response
