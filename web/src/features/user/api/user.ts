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
