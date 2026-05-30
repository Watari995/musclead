import createClient, { type Middleware } from "openapi-fetch";
import { USER_ID_STORAGE_KEY } from "@/lib/auth";
import type { components, paths } from "./schema";

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    if (typeof window === "undefined") return request;
    const userId = window.localStorage.getItem(USER_ID_STORAGE_KEY);
    if (userId) request.headers.set("X-User-ID", userId);
    return request;
  },
};

export const apiClient = createClient<paths>({ baseUrl: BASE_URL });
apiClient.use(authMiddleware);

type Schemas = components["schemas"];
export type UserDTO = Schemas["github_com_Watari995_musclead_internal_user_dto.UserDTO"];
export type MealDTO = Schemas["github_com_Watari995_musclead_internal_meal_dto.MealDTO"];
export type RegisterRequest = Schemas["internal_user_internal_handler.RegisterRequest"];
export type RegisterResponse = Schemas["internal_user_internal_handler.RegisterResponse"];
export type RecordMealRequest = Schemas["internal_meal_internal_handler.RecordMealRequest"];
export type ListMealsResponse = Schemas["internal_meal_internal_handler.ListMealsResponse"];
export type ErrorResponse = Schemas["github_com_Watari995_musclead_internal_shared_httpx.ErrorResponse"];

export class APIError extends Error {
  constructor(public status: number, public body: ErrorResponse | undefined) {
    super(body?.error?.message ?? `HTTP ${status}`);
  }
}
