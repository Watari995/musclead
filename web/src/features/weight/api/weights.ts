"use client";

import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { getAccessToken } from "@/shared/auth/access-token";
import { toWeight, type Weight, type WeightDTO } from "../model/weight";

export const WEIGHTS_QUERY_KEY = ["weights"] as const;

export type UpsertWeightRequest = {
  weight_kg: string;
  body_fat_percentage?: string;
  skeletal_muscle_kg?: string;
  measured_at: string;
};

// 後方互換のため Record 用に同じ型を別名でも公開
export type RecordWeightRequest = UpsertWeightRequest;

export type UpsertWeightResponse = {
  weight_id: string;
};
export type RecordWeightResponse = UpsertWeightResponse;

function baseUrl(): string {
  return process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
}

function authHeaders(): HeadersInit {
  const token = getAccessToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export function useWeightsQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: WEIGHTS_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<Weight[]> => {
      const res = await fetch(
        `${baseUrl()}/weights?limit=50&offset=0`,
        {
          credentials: "include",
          headers: { ...authHeaders() },
        },
      );
      if (!res.ok) {
        throw new Error(`failed to fetch weights (HTTP ${res.status})`);
      }
      const data = (await res.json()) as { weights?: WeightDTO[] };
      return (data.weights ?? []).map(toWeight);
    },
  });
}

export function useRecordWeightMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (
      body: UpsertWeightRequest,
    ): Promise<UpsertWeightResponse> => {
      const res = await fetch(`${baseUrl()}/weights`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json", ...authHeaders() },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        throw new Error(`failed to record weight (HTTP ${res.status})`);
      }
      return (await res.json()) as UpsertWeightResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WEIGHTS_QUERY_KEY });
    },
  });
}

export function useUpdateWeightMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: {
      id: string;
      body: UpsertWeightRequest;
    }): Promise<UpsertWeightResponse> => {
      const res = await fetch(`${baseUrl()}/weights/${input.id}`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json", ...authHeaders() },
        body: JSON.stringify(input.body),
      });
      if (!res.ok) {
        throw new Error(`failed to update weight (HTTP ${res.status})`);
      }
      return (await res.json()) as UpsertWeightResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WEIGHTS_QUERY_KEY });
    },
  });
}

export function useDeleteWeightMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      const res = await fetch(`${baseUrl()}/weights/${id}`, {
        method: "DELETE",
        credentials: "include",
        headers: { ...authHeaders() },
      });
      if (!res.ok) {
        throw new Error(`failed to delete weight (HTTP ${res.status})`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WEIGHTS_QUERY_KEY });
    },
  });
}
