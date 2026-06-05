"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type ListExercisesResponse,
  type UpsertExerciseRequest,
  type UpsertExerciseResponse,
} from "@/shared/api/client";
import { toExercise, type Exercise } from "../model/exercise";

export const EXERCISES_QUERY_KEY = ["exercises", "all"] as const;
const EXERCISE_QUERY_KEY = (id: string) => ["exercise", id] as const;

export class ExerciseNameTakenError extends Error {
  constructor() {
    super("同じ名前の種目が既に登録されています。");
    this.name = "ExerciseNameTakenError";
  }
}

export class ExerciseInUseError extends Error {
  constructor() {
    super("この種目はトレーニング履歴で使われているため削除できません。");
    this.name = "ExerciseInUseError";
  }
}

export function useExercisesQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: EXERCISES_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<Exercise[]> => {
      const { data, error, response } = await apiClient.GET("/exercises", {
        params: { query: { limit: 100, offset: 0 } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      const payload = data as ListExercisesResponse;
      return (payload.exercises ?? []).map(toExercise);
    },
  });
}

export function useExerciseQuery(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: EXERCISE_QUERY_KEY(id),
    enabled: enabled && Boolean(id),
    queryFn: async (): Promise<Exercise> => {
      const { data, error, response } = await apiClient.GET(
        "/exercises/{id}",
        { params: { path: { id } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return toExercise(data);
    },
  });
}

export function useCreateExerciseMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertExerciseRequest) => {
      const { data, error, response } = await apiClient.POST("/exercises", {
        body,
      });
      if (error) {
        if (
          error.error?.code === "training.exercise_name_already_exists_error"
        ) {
          throw new ExerciseNameTakenError();
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as UpsertExerciseResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY });
    },
  });
}

export function useUpdateExerciseMutation(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertExerciseRequest) => {
      const { error, response } = await apiClient.PUT("/exercises/{id}", {
        params: { path: { id } },
        body,
      });
      if (error) {
        if (
          error.error?.code === "training.exercise_name_already_exists_error"
        ) {
          throw new ExerciseNameTakenError();
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: EXERCISE_QUERY_KEY(id) });
    },
  });
}

export function useDeleteExerciseMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE("/exercises/{id}", {
        params: { path: { id } },
      });
      if (error) {
        if (error.error?.code === "training.exercise_used_in_training_error") {
          throw new ExerciseInUseError();
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY });
    },
  });
}
