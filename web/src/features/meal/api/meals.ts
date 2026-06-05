"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type RecordMealRequest,
} from "@/shared/api/client";
import { toMeal, type Meal } from "../model/meal";

export const MEALS_QUERY_KEY = ["meals"] as const;

export function useMealsQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: MEALS_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<Meal[]> => {
      const { data, error, response } = await apiClient.GET("/meals", {
        params: { query: { limit: 50, offset: 0 } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return (data?.meals ?? []).map(toMeal);
    },
  });
}

export function useRecordMealMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: RecordMealRequest) => {
      const { data, error, response } = await apiClient.POST("/meals", {
        body,
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEALS_QUERY_KEY });
    },
  });
}

export function useDeleteMealMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE("/meals/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEALS_QUERY_KEY });
    },
  });
}
