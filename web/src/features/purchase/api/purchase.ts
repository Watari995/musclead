"use client";

import { useMutation } from "@tanstack/react-query";
import {
  apiClient,
  type SubscribeRequest,
  type SubscribeResponse,
} from "@/shared/api/client";

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
