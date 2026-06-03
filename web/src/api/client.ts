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
  Schemas["github_com_Watari995_musclead_internal_user_dto.RegisterRequest"];
export type RegisterResponse =
  Schemas["github_com_Watari995_musclead_internal_user_dto.RegisterResponse"];
export type RecordMealRequest =
  Schemas["github_com_Watari995_musclead_internal_meal_dto.RecordMealRequest"];
export type ListMealsResponse =
  Schemas["github_com_Watari995_musclead_internal_meal_dto.ListMealsResponse"];

export type TrainingDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.TrainingDTO"];
export type TrainingExerciseDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.TrainingExerciseDTO"];
export type TrainingSetDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.TrainingSetDTO"];
export type RecordTrainingRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RecordTrainingRequest"];
export type RecordTrainingExerciseRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RecordTrainingExerciseRequest"];
export type RecordTrainingSetRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RecordTrainingSetRequest"];
export type RecordTrainingResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RecordTrainingResponse"];
export type ListTrainingsResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.ListTrainingsResponse"];

export type ExerciseDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.ExerciseDTO"];
export type ListExercisesResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.ListExercisesResponse"];
export type UpsertExerciseRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.UpsertExerciseRequest"];
export type UpsertExerciseResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.UpsertExerciseResponse"];

export type RoutineDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RoutineDTO"];
export type RoutineExerciseDTO =
  Schemas["github_com_Watari995_musclead_internal_training_dto.RoutineExerciseDTO"];
export type ListRoutinesResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.ListRoutinesResponse"];
export type UpsertRoutineRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.UpsertRoutineRequest"];
export type UpsertRoutineExerciseRequest =
  Schemas["github_com_Watari995_musclead_internal_training_dto.UpsertRoutineExerciseRequest"];
export type UpsertRoutineResponse =
  Schemas["github_com_Watari995_musclead_internal_training_dto.UpsertRoutineResponse"];

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
