"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  loginRequest,
  logoutRequest,
  type AccessTokenResponse,
} from "@/shared/api/auth";
import {
  apiClient,
  type MeResponse,
  type PreferencesDTO,
  type RegisterRequest,
  type UserDTO,
} from "@/shared/api/client";
import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/shared/auth/access-token";

export const ME_QUERY_KEY = ["me"] as const;

async function fetchMe(): Promise<MeResponse> {
  const { data, error, response } = await apiClient.GET("/users/me");
  if (error) {
    throw new Error(error.error?.message ?? `HTTP ${response.status}`);
  }
  return data as MeResponse;
}

// 既存 caller との互換のため UserDTO を返す。 内部は MeResponse を fetch し
// select で user 部分だけ取り出す。 PreferencesDTO は usePreferencesQuery が同じ
// queryKey + queryFn で別の select を持つ形で共存する。
export function useMeQuery(enabled: boolean) {
  return useQuery({
    queryKey: ME_QUERY_KEY,
    enabled,
    queryFn: fetchMe,
    select: (data) => data.user as UserDTO,
  });
}

export function usePreferencesQuery(enabled: boolean) {
  return useQuery({
    queryKey: ME_QUERY_KEY,
    enabled,
    queryFn: fetchMe,
    select: (data) => data.preferences as PreferencesDTO,
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

// PATCH /users/me: 部分更新(name, birthday, profile_image_path)
// schema.ts に PATCH 未反映のため direct fetch で実装
//   - undefined キー → JSON.stringify が省略 → サーバーは「未送信(更新しない)」 と判定
//   - null            → サーバーは「明示的にクリア」 と判定
//                        (birthday: クリア、 profile_image_path: default 復帰)
//   - 値              → 更新
export type UpdateUserBody = {
  name?: string;
  birthday?: string | null;
  profile_image_path?: string | null;
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

// PATCH /users/me/preferences: theme を更新。
// 受け取った theme は light / dark / system のいずれか。
// 楽観的 UI(next-themes 即時反映)とは別レイヤーで永続化を担う。
export type Theme = "light" | "dark" | "system";

export type UpdatePreferencesBody = {
  theme?: Theme;
};

export function useUpdatePreferencesMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpdatePreferencesBody): Promise<void> => {
      const token = getAccessToken();
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
      const res = await fetch(`${baseUrl}/users/me/preferences`, {
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

// プロフィール画像アップロード:
//   1) POST /users/me/profile-image/presigned-url で {url, path} を取得
//   2) その url に PUT で画像 blob を直接 S3 アップロード
//   3) 戻り値の path を呼び出し側が PATCH /users/me に渡す
export type UploadProfileImageInput = { blob: Blob; contentType: string };

export function useUploadProfileImageMutation() {
  return useMutation({
    mutationFn: async ({
      blob,
      contentType,
    }: UploadProfileImageInput): Promise<{ path: string }> => {
      const token = getAccessToken();
      const baseUrl =
        process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

      // 1) presigned URL を取得
      const presignedRes = await fetch(
        `${baseUrl}/users/me/profile-image/presigned-url`,
        {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({ content_type: contentType }),
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

      // 2) S3 に直接 PUT(認証不要、 署名済 URL)
      const putRes = await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": contentType },
        body: blob,
      });
      if (!putRes.ok) {
        throw new Error(`failed to upload to S3 (HTTP ${putRes.status})`);
      }

      return { path };
    },
  });
}
