import { page } from "$app/stores";
import { CLANK_FRONTEND_PUBLIC_API_BASE_URL, CLANK_FRONTEND_PUBLIC_BASE_URL } from "$env/static/public";
import * as entities from "$lib/entities";
import * as payloads from "$lib/payloads";
import { toast } from "svelte-sonner";
import { get } from "svelte/store";

export class api {
  private static BASE_URL: string = CLANK_FRONTEND_PUBLIC_API_BASE_URL + "/ext";

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  private constructor() {}

  private static async request<Q, S>(endpoint: string, method: string, body?: Q): Promise<S> {
    let response;

    try {
      response = await fetch(this.BASE_URL + endpoint, {
        headers: { "Content-type": "application/json; charset=UTF-8" },
        method: method,
        body: body && JSON.stringify(body),
        mode: "cors",
        credentials: "include",
      });
    } catch (error) {
      toast.error("Something went wrong", { description: String(error) });
      throw new entities.ApiError(entities.ApiErrorCode.ERR_SERVER_GENERIC, error as any, 500); // eslint-disable-line @typescript-eslint/no-explicit-any
    }

    if (!response.ok) {
      const url = get(page).url;
      const path = url.pathname + url.search;
      if (response.status === 401 && !path.startsWith("/dash/signin")) {
        window.location.href = CLANK_FRONTEND_PUBLIC_BASE_URL + `/dash/signin?redirect_to=${encodeURIComponent(path)}`;
      }

      let code = entities.ApiErrorCode.ERR_SERVER_GENERIC;
      let message = response.statusText;
      const status = response.status;

      try {
        const error: payloads.ApiError = await response.json();

        if (error.code) {
          code = entities.ApiErrorCode[error.code as keyof typeof entities.ApiErrorCode];
          message = entities.ApiErrorMessage[code];
        }

        if (error.message) {
          message = error.message;
        }
        // eslint-disable-next-line no-empty
      } catch {}

      toast.error("Something went wrong", { description: message });
      throw new entities.ApiError(code, message, status);
    }

    return response.json();
  }

  public static async get<S>(endpoint: string): Promise<S> {
    return await this.request(endpoint, "GET", undefined);
  }

  public static async post<Q, S>(endpoint: string, body?: Q): Promise<S> {
    return await this.request(endpoint, "POST", body);
  }

  public static async put<Q, S>(endpoint: string, body?: Q): Promise<S> {
    return await this.request(endpoint, "PUT", body);
  }

  public static async delete<S>(endpoint: string): Promise<S> {
    return await this.request(endpoint, "DELETE", undefined);
  }
}
