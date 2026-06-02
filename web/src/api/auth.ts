import type { components } from "./schema";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type Schemas = components["schemas"];
export type LoginRequest =
  Schemas["github_com_Watari995_musclead_internal_auth_dto.LoginRequest"];
export type AccessTokenResponse =
  Schemas["github_com_Watari995_musclead_internal_auth_dto.AccessTokenResponse"];

export async function loginRequest(
  email: string,
  password: string,
): Promise<Required<AccessTokenResponse>> {
  const body: LoginRequest = { email, password };
  const res = await fetch(`${BASE_URL}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(await readError(res));
  }
  return ensureTokens((await res.json()) as AccessTokenResponse, res.status);
}

export async function refreshRequest(): Promise<Required<AccessTokenResponse> | null> {
  const res = await fetch(`${BASE_URL}/auth/refresh`, {
    method: "POST",
    credentials: "include",
  });
  if (!res.ok) return null;
  return ensureTokens((await res.json()) as AccessTokenResponse, res.status);
}

export async function logoutRequest(): Promise<void> {
  await fetch(`${BASE_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
}

function ensureTokens(
  body: AccessTokenResponse,
  status: number,
): Required<AccessTokenResponse> {
  if (!body.access_token || !body.access_token_expires_at) {
    throw new Error(`malformed auth response (status ${status})`);
  }
  return {
    access_token: body.access_token,
    access_token_expires_at: body.access_token_expires_at,
  };
}

async function readError(res: Response): Promise<string> {
  try {
    const body = (await res.json()) as { error?: { message?: string } };
    return body?.error?.message ?? `HTTP ${res.status}`;
  } catch {
    return `HTTP ${res.status}`;
  }
}
