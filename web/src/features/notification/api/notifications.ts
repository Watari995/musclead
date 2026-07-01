"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAccessToken } from "@/shared/auth/access-token";

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export type NotificationMetadata = Record<string, unknown>;

export type NotificationDTO = {
  id: string;
  notification_type: string;
  metadata: NotificationMetadata;
  is_read: boolean;
  read_at?: string;
  created_at: string;
};

export type GetNotificationsResponse = {
  notifications: NotificationDTO[];
  unread_count: number;
};

export type WeeklyGoalDTO = {
  training_count: number | null;
  calorie_average: number | null;
  weight_change_kg: number | null;
  created_at: string;
  updated_at: string;
};

export type UpsertWeeklyGoalRequest = {
  training_count: number | null;
  calorie_average: number | null;
  weight_change_kg: number | null;
};

async function authFetch(token: string, input: RequestInfo, init?: RequestInit) {
  const res = await fetch(`${BASE_URL}${input}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...init?.headers,
    },
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res;
}

export const NOTIFICATIONS_QUERY_KEY = ["notifications"] as const;
export const WEEKLY_GOAL_QUERY_KEY = ["weekly-goal"] as const;

export function useNotificationsQuery() {
  const { token } = useAccessToken();
  return useQuery({
    queryKey: NOTIFICATIONS_QUERY_KEY,
    enabled: Boolean(token),
    queryFn: async (): Promise<GetNotificationsResponse> => {
      const res = await authFetch(token!, "/notifications");
      return res.json();
    },
  });
}

export function useReadNotificationMutation() {
  const { token } = useAccessToken();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await authFetch(token!, `/notifications/${id}/read`, { method: "PUT" });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: NOTIFICATIONS_QUERY_KEY });
    },
  });
}

export function useWeeklyGoalQuery() {
  const { token } = useAccessToken();
  return useQuery({
    queryKey: WEEKLY_GOAL_QUERY_KEY,
    enabled: Boolean(token),
    queryFn: async (): Promise<WeeklyGoalDTO | null> => {
      const res = await authFetch(token!, "/users/me/weekly-goal");
      const data = await res.json();
      return data ?? null;
    },
  });
}

export function useUpsertWeeklyGoalMutation() {
  const { token } = useAccessToken();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertWeeklyGoalRequest): Promise<WeeklyGoalDTO> => {
      const res = await authFetch(token!, "/users/me/weekly-goal", {
        method: "PUT",
        body: JSON.stringify(body),
      });
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: WEEKLY_GOAL_QUERY_KEY });
    },
  });
}
