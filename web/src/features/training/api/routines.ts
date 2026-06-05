"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type ListRoutinesResponse,
  type RecordTrainingRequest,
  type RecordTrainingResponse,
  type RoutineDTO,
  type UpsertRoutineRequest,
  type UpsertRoutineResponse,
} from "@/shared/api/client";
import { TRAININGS_QUERY_KEY } from "./trainings";

export const ROUTINES_QUERY_KEY = ["routines"] as const;
const ROUTINE_QUERY_KEY = (id: string) => ["routine", id] as const;

export class RoutineNameTakenError extends Error {
  constructor() {
    super("同じ名前のルーティンが既に登録されています。");
    this.name = "RoutineNameTakenError";
  }
}

export function useRoutinesQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: ROUTINES_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<RoutineDTO[]> => {
      const { data, error, response } = await apiClient.GET("/routines", {
        params: { query: { limit: 50, offset: 0 } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return (data as ListRoutinesResponse).routines ?? [];
    },
  });
}

export function useRoutineQuery(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ROUTINE_QUERY_KEY(id),
    enabled: enabled && Boolean(id),
    queryFn: async (): Promise<RoutineDTO> => {
      const { data, error, response } = await apiClient.GET(
        "/routines/{id}",
        { params: { path: { id } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as RoutineDTO;
    },
  });
}

export function useCreateRoutineMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertRoutineRequest) => {
      const { data, error, response } = await apiClient.POST("/routines", {
        body,
      });
      if (error) {
        if (error.error?.code === "training.routine_name_already_exists_error") {
          throw new RoutineNameTakenError();
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as UpsertRoutineResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTINES_QUERY_KEY });
    },
  });
}

export function useUpdateRoutineMutation(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertRoutineRequest) => {
      const { error, response } = await apiClient.PUT("/routines/{id}", {
        params: { path: { id } },
        body,
      });
      if (error) {
        if (error.error?.code === "training.routine_name_already_exists_error") {
          throw new RoutineNameTakenError();
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTINES_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: ROUTINE_QUERY_KEY(id) });
    },
  });
}

export function useDeleteRoutineMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE("/routines/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ROUTINES_QUERY_KEY });
    },
  });
}

// Routine から空 Training を派生させる(その日に決めたセット数値で記録するため)。
// ADR 0006 §4 copy-on-use。
export function useStartTrainingFromRoutineMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (routine: RoutineDTO): Promise<RecordTrainingResponse> => {
      const body: RecordTrainingRequest = {
        started_at: new Date().toISOString(),
        exercises: (routine.routine_exercises ?? []).map((re) => ({
          exercise_id: re.exercise_id ?? "",
          display_order: re.display_order ?? 1,
          sets: [{ set_number: 1, weight_kg: "0", reps: 0 }],
        })),
      };
      const { data, error, response } = await apiClient.POST("/trainings", {
        body,
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as RecordTrainingResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TRAININGS_QUERY_KEY });
    },
  });
}
