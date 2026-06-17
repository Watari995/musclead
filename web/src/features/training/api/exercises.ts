"use client";

import {
  keepPreviousData,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import {
  apiClient,
  type BestSetDTO,
  type ListBestSetsResponse,
  type ListExercisesResponse,
  type ReorderExercisesRequest,
  type UpsertExerciseRequest,
  type UpsertExerciseResponse,
} from "@/shared/api/client";
import { toExercise, type Exercise } from "../model/exercise";

export const EXERCISES_QUERY_KEY = ["exercises", "all"] as const;
const EXERCISE_QUERY_KEY = (id: string) => ["exercise", id] as const;
const BEST_SETS_QUERY_KEY = (ids: string[]) =>
  ["exercises", "best-sets", ids] as const;

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

// 複数種目の最高記録(最重量セット)を 1 リクエストでまとめて取得する。
// 記録のある種目だけ返ってくるので、exercise_id をキーにした Map にして返す。
export function useBestSetsQuery(exerciseIDs: string[]) {
  const ids = Array.from(new Set(exerciseIDs.filter(Boolean))).sort();
  return useQuery({
    queryKey: BEST_SETS_QUERY_KEY(ids),
    enabled: ids.length > 0,
    placeholderData: keepPreviousData, // 種目追加時に既存バッジを消さない
    staleTime: 60_000,
    queryFn: async (): Promise<Map<string, BestSetDTO>> => {
      const { data, error, response } = await apiClient.GET(
        "/exercises/best-sets",
        { params: { query: { exercise_ids: ids } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      const list = (data as ListBestSetsResponse).best_sets ?? [];
      const map = new Map<string, BestSetDTO>();
      for (const b of list) {
        if (b.exercise_id) map.set(b.exercise_id, b);
      }
      return map;
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

export function useReorderExercisesMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    // ordered には並び替え後の全件を渡す。 楽観的にキャッシュを更新する。
    mutationFn: async (ordered: Exercise[]) => {
      const body: ReorderExercisesRequest = {
        exercise_ids: ordered.map((e) => e.id),
      };
      const { error, response } = await apiClient.POST("/exercises/reorder", {
        body,
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onMutate: async (ordered: Exercise[]) => {
      await queryClient.cancelQueries({ queryKey: EXERCISES_QUERY_KEY });
      const previous = queryClient.getQueryData<Exercise[]>(EXERCISES_QUERY_KEY);
      queryClient.setQueryData<Exercise[]>(
        EXERCISES_QUERY_KEY,
        ordered.map((e, i) => ({ ...e, displayOrder: i })),
      );
      return { previous };
    },
    onError: (_err, _ordered, context) => {
      if (context?.previous) {
        queryClient.setQueryData(EXERCISES_QUERY_KEY, context.previous);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY });
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
