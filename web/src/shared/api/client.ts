import createClient, { type Middleware } from "openapi-fetch";
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/shared/auth/access-token";
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

// swag v1.16.5+ は schema 名に fully qualified Go path を出すようになったため、
// 各 package ごとに短い alias を定義しておく(呼び出し側のコードを変えなくていいように)
type Schemas = components["schemas"];
type UserDto<K extends string> =
  `github_com_Watari995_musclead_internal_user_dto.${K}` extends infer P
    ? P extends keyof Schemas
      ? Schemas[P]
      : never
    : never;
type MealDto<K extends string> =
  `github_com_Watari995_musclead_internal_meal_dto.${K}` extends infer P
    ? P extends keyof Schemas
      ? Schemas[P]
      : never
    : never;
type TrainingDto<K extends string> =
  `github_com_Watari995_musclead_internal_training_dto.${K}` extends infer P
    ? P extends keyof Schemas
      ? Schemas[P]
      : never
    : never;
type Httpx<K extends string> =
  `github_com_Watari995_musclead_internal_shared_httpx.${K}` extends infer P
    ? P extends keyof Schemas
      ? Schemas[P]
      : never
    : never;

export type UserDTO = UserDto<"UserDTO">;
export type PreferencesDTO = UserDto<"PreferencesDTO">;
export type MeResponse = UserDto<"MeResponse">;
export type RegisterRequest = UserDto<"RegisterRequest">;
export type RegisterResponse = UserDto<"RegisterResponse">;
export type UpdatePreferencesRequest = UserDto<"UpdatePreferencesRequest">;
export type UpdatePreferencesResponse = UserDto<"UpdatePreferencesResponse">;

export type MealDTO = MealDto<"MealDTO">;
export type RecordMealRequest = MealDto<"RecordMealRequest">;
export type ListMealsResponse = MealDto<"ListMealsResponse">;

export type TrainingDTO = TrainingDto<"TrainingDTO">;
export type TrainingExerciseDTO = TrainingDto<"TrainingExerciseDTO">;
export type TrainingSetDTO = TrainingDto<"TrainingSetDTO">;
export type RecordTrainingRequest = TrainingDto<"RecordTrainingRequest">;
export type RecordTrainingExerciseRequest = TrainingDto<"RecordTrainingExerciseRequest">;
export type RecordTrainingSetRequest = TrainingDto<"RecordTrainingSetRequest">;
export type RecordTrainingResponse = TrainingDto<"RecordTrainingResponse">;
export type ListTrainingsResponse = TrainingDto<"ListTrainingsResponse">;

export type ExerciseDTO = TrainingDto<"ExerciseDTO">;
export type ListExercisesResponse = TrainingDto<"ListExercisesResponse">;
export type UpsertExerciseRequest = TrainingDto<"UpsertExerciseRequest">;
export type UpsertExerciseResponse = TrainingDto<"UpsertExerciseResponse">;
export type ReorderExercisesRequest = TrainingDto<"ReorderExercisesRequest">;

export type RoutineDTO = TrainingDto<"RoutineDTO">;
export type RoutineExerciseDTO = TrainingDto<"RoutineExerciseDTO">;
export type ListRoutinesResponse = TrainingDto<"ListRoutinesResponse">;
export type UpsertRoutineRequest = TrainingDto<"UpsertRoutineRequest">;
export type UpsertRoutineExerciseRequest = TrainingDto<"UpsertRoutineExerciseRequest">;
export type UpsertRoutineResponse = TrainingDto<"UpsertRoutineResponse">;

export type ErrorResponse = Httpx<"ErrorResponse">;

export class APIError extends Error {
  constructor(
    public status: number,
    public body: ErrorResponse | undefined,
  ) {
    super(body?.error?.message ?? `HTTP ${status}`);
  }
}
