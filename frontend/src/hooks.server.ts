import { CLANK_ENVIRONMENT, CLANK_FRONTEND_SENTRY_DSN, CLANK_RELEASE } from "$env/static/private";
import * as entities from "$lib/entities";
import { init } from "@jill64/sentry-sveltekit-cloudflare/server";
import type { Handle, HandleServerError } from "@sveltejs/kit";

const { onHandle, onError } = init("", {
  toucanOptions: {
    dsn: CLANK_FRONTEND_SENTRY_DSN,
    environment: CLANK_ENVIRONMENT,
    release: CLANK_RELEASE,
    debug: false,
    attachStacktrace: true,
    enableTracing: true,
    sampleRate: 1.0,
    tracesSampleRate: 0.25,
    ignoreErrors: [
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_CLIENT_GENERIC],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_INVALID_REQUEST],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_NOT_FOUND],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_UNAUTHORIZED],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_RATE_LIMITED],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_NO_PERMISSION],
    ],
  },
  enableInDevMode: !!CLANK_FRONTEND_SENTRY_DSN,
});

export const handle: Handle = onHandle();

export const handleError: HandleServerError = onError(({ message }) => {
  return { message: message };
});
