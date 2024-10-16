from typing import Optional

import fastapi.exceptions
import starlette.exceptions
from fastapi import status
from fastapi.requests import Request
from fastapi.responses import JSONResponse

from src.config import Config


class HTTPException(Exception):
    cause: Optional[Exception]
    code: str
    status: int

    def __call__(self, cause: Exception) -> Exception:
        self.cause = cause

        return self


class ServerGenericException(HTTPException):
    code = "ERR_SERVER_GENERIC"
    status = status.HTTP_500_INTERNAL_SERVER_ERROR


class ServerUnavailableException(HTTPException):
    code = "ERR_SERVER_UNAVAILABLE"
    status = status.HTTP_503_SERVICE_UNAVAILABLE


class ServerTimeoutException(HTTPException):
    code = "ERR_SERVER_TIMEOUT"
    status = status.HTTP_504_GATEWAY_TIMEOUT


class ClientGenericException(HTTPException):
    code = "ERR_CLIENT_GENERIC"
    status = status.HTTP_400_BAD_REQUEST


class InvalidRequestException(HTTPException):
    code = "ERR_INVALID_REQUEST"
    status = status.HTTP_400_BAD_REQUEST


class NotFoundException(HTTPException):
    code = "ERR_NOT_FOUND"
    status = status.HTTP_404_NOT_FOUND


class UnauthorizedException(HTTPException):
    code = "ERR_UNAUTHORIZED"
    status = status.HTTP_401_UNAUTHORIZED


class ExceptionHandler:
    def __init__(self, config: Config) -> None:
        self.config = config

    async def handle(self, request: Request, exception: Exception) -> JSONResponse:
        if isinstance(exception, starlette.exceptions.HTTPException):
            if exception.status_code == status.HTTP_404_NOT_FOUND:
                exception = NotFoundException()(exception)
            elif exception.status_code == status.HTTP_405_METHOD_NOT_ALLOWED:
                exception = InvalidRequestException()(exception)
            elif exception.status_code == status.HTTP_413_REQUEST_ENTITY_TOO_LARGE:
                exception = InvalidRequestException()(exception)

        if isinstance(exception, fastapi.exceptions.RequestValidationError):
            exception = InvalidRequestException()(exception)

        if not isinstance(exception, HTTPException):
            exception = ServerGenericException()(exception)

        if self.config.service.environment != "dev":
            return JSONResponse({"code": exception.code}, exception.status)

        return JSONResponse({"code": exception.code, "message": str(exception.cause or "")}, exception.status)
