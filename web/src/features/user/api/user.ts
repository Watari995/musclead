"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  loginRequest,
  logoutRequest,
  type AccessTokenResponse,
} from "@/shared/api/auth";
import {
  apiClient,
  type RegisterRequest,
  type UserDTO,
} from "@/shared/api/client";
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/shared/auth/access-token";

export const ME_QUERY_KEY = ["me"] as const;

export function useMeQuery(enabled: boolean) {
  return useQuery({
    queryKey: ME_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<UserDTO> => {
      const { data, error, response } = await apiClient.GET("/users/me");
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as UserDTO;
    },
  });
}

type LoginInput = { email: string; password: string };

export function useLoginMutation() {
  return useMutation({
    mutationFn: async (
      input: LoginInput,
    ): Promise<Required<AccessTokenResponse>> => {
      const tokens = await loginRequest(input.email, input.password);
      setAccessToken(tokens.access_token);
      return tokens;
    },
  });
}

export function useRegisterMutation() {
  return useMutation({
    mutationFn: async (
      body: RegisterRequest,
    ): Promise<Required<AccessTokenResponse>> => {
      const { error, response } = await apiClient.POST("/users", { body });
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      const tokens = await loginRequest(body.email ?? "", body.password ?? "");
      setAccessToken(tokens.access_token);
      return tokens;
    },
  });
}

export function useLogoutMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      await logoutRequest();
      clearAccessToken();
      queryClient.removeQueries({ queryKey: ME_QUERY_KEY });
    },
  });
}

// PATCH /users/me: 部分更新(name と birthday のみ)
// schema.ts に PATCH 未反映のため direct fetch で実装
//   - undefined キー → JSON.stringify が省略 → サーバーは「未送信(更新しない)」 と判定
//   - null            → サーバーは「明示的にクリア」 と判定(birthday のみ許可)
//   - 値              → 更新
export type UpdateUserBody = {
  name?: string;
  birthday?: string | null;
};

export function useUpdateUserMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpdateUserBody): Promise<void> => {
      const token = getAccessToken();
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
      const res = await fetch(`${baseUrl}/users/me`, {
        method: "PATCH",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        let message = `HTTP ${res.status}`;
        try {
          const json = (await res.json()) as { error?: { message?: string } };
          if (json.error?.message) message = json.error.message;
        } catch {
          // fall through
        }
        throw new Error(message);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ME_QUERY_KEY });
    },
  });
}
