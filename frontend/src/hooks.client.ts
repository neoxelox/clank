import {
  CLANK_FRONTEND_PUBLIC_ENVIRONMENT,
  CLANK_FRONTEND_PUBLIC_RELEASE,
  CLANK_FRONTEND_PUBLIC_SENTRY_DSN,
} from "$env/static/public";
import * as entities from "$lib/entities";
import { init } from "@jill64/sentry-sveltekit-cloudflare/client";
import type { HandleClientError } from "@sveltejs/kit";

const onError = init("", {
  sentryOptions: {
    dsn: CLANK_FRONTEND_PUBLIC_SENTRY_DSN,
    environment: CLANK_FRONTEND_PUBLIC_ENVIRONMENT,
    release: CLANK_FRONTEND_PUBLIC_RELEASE,
    debug: false,
    attachStacktrace: true,
    enableTracing: true,
    sampleRate: 1.0,
    tracesSampleRate: 0.25,
    profilesSampleRate: 1.0,
    replaysSessionSampleRate: 0.1,
    replaysOnErrorSampleRate: 1.0,
    ignoreErrors: [
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_CLIENT_GENERIC],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_INVALID_REQUEST],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_NOT_FOUND],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_UNAUTHORIZED],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_RATE_LIMITED],
      entities.ApiErrorMessage[entities.ApiErrorCode.ERR_NO_PERMISSION],
    ],
  },
  enableInDevMode: !!CLANK_FRONTEND_PUBLIC_SENTRY_DSN,
});

export const handleError: HandleClientError = onError(({ message }) => {
  return { message: message };
});
