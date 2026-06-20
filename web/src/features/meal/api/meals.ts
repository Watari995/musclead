"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type RecordMealRequest,
  type UpdateMealRequest,
} from "@/shared/api/client";
import { getAccessToken } from "@/shared/auth/access-token";
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

// 食事写真アップロード:
//   1) POST /meals/photos/presigned-url で {url, path} 取得
//   2) その url に PUT で blob を S3 へ直接アップロード
//   3) 戻り値 path を呼び出し側が POST /meals の photos[].image_path に使う
export function useUploadMealPhotoMutation() {
  return useMutation({
    mutationFn: async ({
      file,
    }: {
      file: File;
    }): Promise<{ path: string }> => {
      const token = getAccessToken();
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

      const presignedRes = await fetch(
        `${baseUrl}/meals/photos/presigned-url`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({ content_type: file.type }),
        },
      );
      if (!presignedRes.ok) {
        throw new Error(
          `failed to get presigned URL (HTTP ${presignedRes.status})`,
        );
      }
      const { url, path } = (await presignedRes.json()) as {
        url: string;
        path: string;
      };

      const putRes = await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": file.type },
        body: file,
      });
      if (!putRes.ok) {
        throw new Error(`failed to upload to S3 (HTTP ${putRes.status})`);
      }

      return { path };
    },
  });
}

export function useFindMealQuery(id: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ["meals", id],
    enabled,
    queryFn: async (): Promise<Meal> => {
      const { data, error, response } = await apiClient.GET("/meals/{id}", {
        params: { path: { id } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return toMeal(data!);
    },
  });
}

export function useUpdateMealMutation(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpdateMealRequest) => {
      const { data, error, response } = await apiClient.PUT("/meals/{id}", {
        params: { path: { id } },
        body,
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEALS_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: ["meals", id] });
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
