import createClient, { type Middleware } from "openapi-fetch";
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/lib/access-token";
import { refreshRequest } from "./auth";
import type { components, paths } from "./schema";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const RETRY_HEADER = "x-musclead-retried";
const AUTH_EXPIRED_EVENT = "musclead:auth-expired";

const credentialsFetch: typeof fetch = (input, init) =>
  fetch(input, { ...init, credentials: "include" });

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    const token = getAccessToken();
    if (token) request.headers.set("Authorization", `Bearer ${token}`);
    return request;
  },
  async onResponse({ request, response }) {
    if (response.status !== 401) return undefined;
    if (request.headers.get(RETRY_HEADER)) return undefined;

    const tokens = await refreshRequest();
    if (!tokens) {
      clearAccessToken();
      if (typeof window !== "undefined") {
        window.dispatchEvent(new Event(AUTH_EXPIRED_EVENT));
      }
      return undefined;
    }
    setAccessToken(tokens.access_token);

    const retry = new Request(request);
    retry.headers.set("Authorization", `Bearer ${tokens.access_token}`);
    retry.headers.set(RETRY_HEADER, "1");
    return credentialsFetch(retry);
  },
};

export const apiClient = createClient<paths>({
  baseUrl: BASE_URL,
  fetch: credentialsFetch,
});
apiClient.use(authMiddleware);

type Schemas = components["schemas"];
export type UserDTO =
  Schemas["github_com_Watari995_musclead_internal_user_dto.UserDTO"];
export type MealDTO =
  Schemas["github_com_Watari995_musclead_internal_meal_dto.MealDTO"];
export type RegisterRequest =
  Schemas["internal_user_internal_handler.RegisterRequest"];
export type RegisterResponse =
  Schemas["internal_user_internal_handler.RegisterResponse"];
export type RecordMealRequest =
  Schemas["internal_meal_internal_handler.RecordMealRequest"];
export type ListMealsResponse =
  Schemas["internal_meal_internal_handler.ListMealsResponse"];
export type ErrorResponse =
  Schemas["github_com_Watari995_musclead_internal_shared_httpx.ErrorResponse"];

export class APIError extends Error {
  constructor(
    public status: number,
    public body: ErrorResponse | undefined,
  ) {
    super(body?.error?.message ?? `HTTP ${status}`);
  }
}
