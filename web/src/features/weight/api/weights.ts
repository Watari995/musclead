"use client";

import { useMutation } from "@tanstack/react-query";
import { getAccessToken } from "@/shared/auth/access-token";

export type RecordWeightRequest = {
  weight_kg: string;
  body_fat_percentage?: string;
  skeletal_muscle_kg?: string;
  measured_at: string;
};

export type RecordWeightResponse = {
  weight_id: string;
};

export function useRecordWeightMutation() {
  return useMutation({
    mutationFn: async (
      body: RecordWeightRequest,
    ): Promise<RecordWeightResponse> => {
      const token = getAccessToken();
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

      const res = await fetch(`${baseUrl}/weights`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        throw new Error(`failed to record weight (HTTP ${res.status})`);
      }
      return (await res.json()) as RecordWeightResponse;
    },
  });
}
