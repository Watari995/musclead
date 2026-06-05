"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type ListTrainingsResponse,
  type RecordTrainingRequest,
  type RecordTrainingResponse,
  type TrainingDTO,
} from "@/shared/api/client";

export const TRAININGS_QUERY_KEY = ["trainings"] as const;
const TRAINING_QUERY_KEY = (id: string) => ["training", id] as const;

export function useTrainingsQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: TRAININGS_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<TrainingDTO[]> => {
      const { data, error, response } = await apiClient.GET("/trainings", {
        params: { query: { limit: 50, offset: 0 } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return (data as ListTrainingsResponse).trainings ?? [];
    },
  });
}

export function useTrainingQuery(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: TRAINING_QUERY_KEY(id),
    enabled: enabled && Boolean(id),
    queryFn: async (): Promise<TrainingDTO> => {
      const { data, error, response } = await apiClient.GET(
        "/trainings/{id}",
        { params: { path: { id } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as TrainingDTO;
    },
  });
}

export function useRecordTrainingMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (
      body: RecordTrainingRequest,
    ): Promise<RecordTrainingResponse> => {
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

export function useUpdateTrainingMutation(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: RecordTrainingRequest) => {
      const { error, response } = await apiClient.PUT("/trainings/{id}", {
        params: { path: { id } },
        body,
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TRAININGS_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: TRAINING_QUERY_KEY(id) });
    },
  });
}

export function useDeleteTrainingMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE("/trainings/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: TRAININGS_QUERY_KEY });
    },
  });
}
