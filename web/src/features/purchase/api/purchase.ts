"use client";

import { useMutation, useQuery } from "@tanstack/react-query";
import {
  apiClient,
  type CreatePortalSessionResponse,
  type GetSubscriptionResponse,
  type SubscribeRequest,
  type SubscribeResponse,
} from "@/shared/api/client";

export const SUBSCRIPTION_QUERY_KEY = ["subscription"] as const;

// useSubscriptionQuery は現在のサブスク状態 (is_pro / plan / expires_at) を取得する。
// プラン画面で「申込み」か「現在Pro+管理」かを出し分けるために使う。
export function useSubscriptionQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: SUBSCRIPTION_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<GetSubscriptionResponse> => {
      const { data, error, response } = await apiClient.GET(
        "/purchase/subscription",
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as GetSubscriptionResponse;
    },
  });
}

// useSubscribeMutation は Pro 申込みを開始する。
// POST /purchase/subscribe で Stripe Checkout Session を作成し、 checkout_url を返す。
// 呼び出し側は返却された checkout_url へ window.location で遷移する。
export function useSubscribeMutation() {
  return useMutation({
    mutationFn: async (body: SubscribeRequest): Promise<SubscribeResponse> => {
      const { data, error, response } = await apiClient.POST(
        "/purchase/subscribe",
        { body },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as SubscribeResponse;
    },
  });
}

// usePortalSessionMutation は Stripe Customer Portal の URL を払い出す。
// 「お支払い・解約の管理」 押下時に呼び、 返却 URL へ window.location で遷移する。
export function usePortalSessionMutation() {
  return useMutation({
    mutationFn: async (): Promise<CreatePortalSessionResponse> => {
      const { data, error, response } = await apiClient.POST(
        "/purchase/portal-session",
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as CreatePortalSessionResponse;
    },
  });
}
