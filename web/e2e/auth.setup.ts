import { test as setup } from "@playwright/test";
import {
  loginForAccessToken,
  registerTestUser,
  seedMinimumTrainingData,
} from "./helpers/auth";

const STORAGE_STATE_PATH = "e2e/.auth/user.json";

setup("register e2e user and persist auth state", async ({ page, request }) => {
  const user = await registerTestUser(request);

  // Popover や 「ルーティンから始める」 など "データがあるときに出る UI" も
  // 検査対象に含めるため、 最低限のシードを入れる
  const accessToken = await loginForAccessToken(request, user);
  await seedMinimumTrainingData(request, accessToken);

  await page.goto("/login");
  await page.getByLabel("メールアドレス").fill(user.email);
  await page.getByLabel("パスワード").fill(user.password);
  await page.getByRole("button", { name: "ログイン" }).click();
  await page.waitForURL("**/meals");

  await page.context().storageState({ path: STORAGE_STATE_PATH });
});
