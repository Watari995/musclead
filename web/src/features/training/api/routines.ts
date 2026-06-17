"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type ListRoutinesResponse,
  type RecordTrainingRequest,
  type RecordTrainingResponse,
  type ReorderRoutinesRequest,
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

// RoutineLimitReachedError は無料プランのルーティン上限 (3件) に達した時。
// 呼び出し側は Pro へのアップグレード導線を出す。
export class RoutineLimitReachedError extends Error {
  constructor() {
    super("ルーティンは無料プランで3件までです。");
    this.name = "RoutineLimitReachedError";
  }
}

export function useRoutinesQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: ROUTINES_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<RoutineDTO[]> => {
      const { data, error, response } = await apiClient.GET("/routines", {
        params: { query: { limit: 100, offset: 0 } },
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
        if (error.error?.code === "training.routine_limit_reached_error") {
          throw new RoutineLimitReachedError();
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

export function useReorderRoutinesMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    // ordered には並び替え後の全件を渡す。 楽観的にキャッシュを更新する。
    mutationFn: async (ordered: RoutineDTO[]) => {
      const body: ReorderRoutinesRequest = {
        routine_ids: ordered.map((r) => r.id ?? ""),
      };
      const { error, response } = await apiClient.POST("/routines/reorder", {
        body,
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onMutate: async (ordered: RoutineDTO[]) => {
      await queryClient.cancelQueries({ queryKey: ROUTINES_QUERY_KEY });
      const previous =
        queryClient.getQueryData<RoutineDTO[]>(ROUTINES_QUERY_KEY);
      queryClient.setQueryData<RoutineDTO[]>(ROUTINES_QUERY_KEY, ordered);
      return { previous };
    },
    onError: (_err, _ordered, context) => {
      if (context?.previous) {
        queryClient.setQueryData(ROUTINES_QUERY_KEY, context.previous);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ROUTINES_QUERY_KEY });
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
      // 即時 refetch すると、 まだ表示中の一覧に新規分が一瞬描画されてから
      // 記録画面へ遷移してしまう(チラつき)。 stale マークのみに留め、
      // 一覧へ戻った時に再取得させる。
      queryClient.invalidateQueries({
        queryKey: TRAININGS_QUERY_KEY,
        refetchType: "none",
      });
    },
  });
}
