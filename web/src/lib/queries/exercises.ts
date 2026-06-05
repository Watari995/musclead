"use client";

import { useQuery } from "@tanstack/react-query";
import {
  apiClient,
  type ExerciseDTO,
  type ListExercisesResponse,
} from "@/api/client";

export const EXERCISES_QUERY_KEY = ["exercises", "all"] as const;

export function useExercisesQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: EXERCISES_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<ExerciseDTO[]> => {
      const { data, error, response } = await apiClient.GET("/exercises", {
        params: { query: { limit: 100, offset: 0 } },
      });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return (data as ListExercisesResponse).exercises ?? [];
    },
  });
}
