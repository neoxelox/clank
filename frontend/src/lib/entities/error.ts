export enum ApiErrorCode {
  ERR_SERVER_GENERIC = 0,
  ERR_SERVER_UNAVAILABLE,
  ERR_SERVER_TIMEOUT,
  ERR_CLIENT_GENERIC,
  ERR_INVALID_REQUEST,
  ERR_NOT_FOUND,
  ERR_UNAUTHORIZED,
  ERR_RATE_LIMITED,
  ERR_NO_PERMISSION,
}

export const ApiErrorMessage: Record<ApiErrorCode | number, string> = {
  [ApiErrorCode.ERR_SERVER_GENERIC]:
    "An error occurred, please try again later. If the issue persists, contact our support team.",
  [ApiErrorCode.ERR_SERVER_UNAVAILABLE]:
    "An error occurred, please try again later. If the issue persists, contact our support team.",
  [ApiErrorCode.ERR_SERVER_TIMEOUT]:
    "An error occurred, please try again later. If the issue persists, contact our support team.",
  [ApiErrorCode.ERR_CLIENT_GENERIC]:
    "An error occurred, please try again later. If the issue persists, contact our support team.",
  [ApiErrorCode.ERR_INVALID_REQUEST]: "Invalid request, please check your input and try again.",
  [ApiErrorCode.ERR_NOT_FOUND]: "Resource not found, the requested resource does not exist.",
  [ApiErrorCode.ERR_UNAUTHORIZED]: "Authentication failed, please sign in or check your credentials.",
  [ApiErrorCode.ERR_RATE_LIMITED]: "Too many requests, please calm down and try again later.",
  [ApiErrorCode.ERR_NO_PERMISSION]:
    "Access denied, you don't have permission to view this resource or perform this action.",
};

export class ApiError extends Error {
  public code: ApiErrorCode;
  public status: number;

  public constructor(code: ApiErrorCode, message: string, status: number) {
    super(message);
    this.code = code;
    this.status = status;
    Object.setPrototypeOf(this, ApiError.prototype);
  }
}
