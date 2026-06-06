import type { APIRequestContext } from "@playwright/test";

export type TestUser = {
  email: string;
  password: string;
  name: string;
};

export const API_BASE =
  process.env.E2E_API_BASE_URL ?? "http://localhost:8080";

export async function registerTestUser(
  request: APIRequestContext,
): Promise<TestUser> {
  const suffix = `${Date.now()}-${Math.floor(Math.random() * 1_000_000)}`;
  const user: TestUser = {
    name: "E2E User",
    email: `e2e+${suffix}@example.com`,
    password: "SuperSecret123!",
  };
  const res = await request.post(`${API_BASE}/users`, {
    data: {
      name: user.name,
      email: user.email,
      password: user.password,
      birthday: "1990-01-01",
    },
    failOnStatusCode: false,
  });
  if (!res.ok()) {
    const body = await res.text();
    throw new Error(`register failed: ${res.status()} ${body}`);
  }
  return user;
}

export async function loginForAccessToken(
  request: APIRequestContext,
  user: TestUser,
): Promise<string> {
  const res = await request.post(`${API_BASE}/auth/login`, {
    data: { email: user.email, password: user.password },
    failOnStatusCode: false,
  });
  if (!res.ok()) {
    throw new Error(`login failed: ${res.status()} ${await res.text()}`);
  }
  const json = (await res.json()) as { access_token: string };
  return json.access_token;
}

/**
 * Popover や 「ルーティンから始める」 ボタンなど "データがあるとき出る UI" を
 * テスト対象に含めるための最小シード。 種目 1 件 + ルーティン 1 件を作る。
 */
export async function seedMinimumTrainingData(
  request: APIRequestContext,
  accessToken: string,
): Promise<void> {
  const auth = { Authorization: `Bearer ${accessToken}` };

  const exRes = await request.post(`${API_BASE}/exercises`, {
    headers: auth,
    data: { name: "E2E ベンチプレス" },
    failOnStatusCode: false,
  });
  if (!exRes.ok()) {
    throw new Error(
      `create exercise failed: ${exRes.status()} ${await exRes.text()}`,
    );
  }
  const { id: exerciseID } = (await exRes.json()) as { id: string };

  const rtRes = await request.post(`${API_BASE}/routines`, {
    headers: auth,
    data: {
      name: "E2E Day A",
      exercises: [{ exercise_id: exerciseID, display_order: 0 }],
    },
    failOnStatusCode: false,
  });
  if (!rtRes.ok()) {
    throw new Error(
      `create routine failed: ${rtRes.status()} ${await rtRes.text()}`,
    );
  }
}
