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

type Schemas = components["schemas"];
export type UserDTO =
  Schemas["userdto.UserDTO"];
export type MealDTO =
  Schemas["mealdto.MealDTO"];
export type RegisterRequest =
  Schemas["userdto.RegisterRequest"];
export type RegisterResponse =
  Schemas["userdto.RegisterResponse"];
export type RecordMealRequest =
  Schemas["mealdto.RecordMealRequest"];
export type ListMealsResponse =
  Schemas["mealdto.ListMealsResponse"];

export type TrainingDTO =
  Schemas["trainingdto.TrainingDTO"];
export type TrainingExerciseDTO =
  Schemas["trainingdto.TrainingExerciseDTO"];
export type TrainingSetDTO =
  Schemas["trainingdto.TrainingSetDTO"];
export type RecordTrainingRequest =
  Schemas["trainingdto.RecordTrainingRequest"];
export type RecordTrainingExerciseRequest =
  Schemas["trainingdto.RecordTrainingExerciseRequest"];
export type RecordTrainingSetRequest =
  Schemas["trainingdto.RecordTrainingSetRequest"];
export type RecordTrainingResponse =
  Schemas["trainingdto.RecordTrainingResponse"];
export type ListTrainingsResponse =
  Schemas["trainingdto.ListTrainingsResponse"];

export type ExerciseDTO =
  Schemas["trainingdto.ExerciseDTO"];
export type ListExercisesResponse =
  Schemas["trainingdto.ListExercisesResponse"];
export type UpsertExerciseRequest =
  Schemas["trainingdto.UpsertExerciseRequest"];
export type UpsertExerciseResponse =
  Schemas["trainingdto.UpsertExerciseResponse"];

export type RoutineDTO =
  Schemas["trainingdto.RoutineDTO"];
export type RoutineExerciseDTO =
  Schemas["trainingdto.RoutineExerciseDTO"];
export type ListRoutinesResponse =
  Schemas["trainingdto.ListRoutinesResponse"];
export type UpsertRoutineRequest =
  Schemas["trainingdto.UpsertRoutineRequest"];
export type UpsertRoutineExerciseRequest =
  Schemas["trainingdto.UpsertRoutineExerciseRequest"];
export type UpsertRoutineResponse =
  Schemas["trainingdto.UpsertRoutineResponse"];

export type ErrorResponse =
  Schemas["httpx.ErrorResponse"];

export class APIError extends Error {
  constructor(
    public status: number,
    public body: ErrorResponse | undefined,
  ) {
    super(body?.error?.message ?? `HTTP ${status}`);
  }
}
