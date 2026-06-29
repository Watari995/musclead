"use client";

import { useQuery } from "@tanstack/react-query";
import {
  apiClient,
  type GetMonthlySummaryResponse,
  type GetDailySummaryResponse,
} from "@/shared/api/client";

export const MONTHLY_SUMMARY_QUERY_KEY = (year: number, month: number) =>
  ["calendar", "monthly", year, month] as const;

export const DAILY_SUMMARY_QUERY_KEY = (date: string) =>
  ["calendar", "daily", date] as const;

export function useMonthlyCalendarQuery(year: number, month: number, enabled = true) {
  return useQuery({
    queryKey: MONTHLY_SUMMARY_QUERY_KEY(year, month),
    enabled,
    queryFn: async (): Promise<GetMonthlySummaryResponse> => {
      const { data, error, response } = await apiClient.GET(
        "/calendar/monthly-summary",
        { params: { query: { year, month } } },
      );
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as GetMonthlySummaryResponse;
    },
  });
}

export function useDailyCalendarQuery(date: string, enabled = true) {
  return useQuery({
    queryKey: DAILY_SUMMARY_QUERY_KEY(date),
    enabled: enabled && Boolean(date),
    queryFn: async (): Promise<GetDailySummaryResponse> => {
      const { data, error, response } = await apiClient.GET(
        "/calendar/daily-summary",
        { params: { query: { date } } },
      );
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as GetDailySummaryResponse;
    },
  });
}
