import { defineConfig } from "@playwright/test";

const baseURL = process.env.E2E_BASE_URL ?? "http://localhost:3000";

// iPhone SE (375×667) 相当のモバイル emulation。 Chromium / WebKit 両エンジンで
// 共有する。 datetime-local 等のネイティブフォームコントロールは Chromium と
// iOS Safari (WebKit) で描画が異なり、 Chromium だけでは Safari 固有のはみ出しを
// 検出できないため、 両エンジンで mobile UI ガードを回す。
const mobileEmulation = {
  viewport: { width: 375, height: 667 },
  deviceScaleFactor: 2,
  isMobile: true,
  hasTouch: true,
  userAgent:
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
  storageState: "e2e/.auth/user.json",
} as const;

export default defineConfig({
  testDir: "./e2e",
  timeout: 30_000,
  expect: { timeout: 5_000 },
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI
    ? [["github"], ["html", { open: "never" }]]
    : [["list"], ["html", { open: "never" }]],
  use: {
    baseURL,
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
  projects: [
    {
      name: "setup",
      testMatch: /auth\.setup\.ts$/,
    },
    {
      name: "mobile-chromium",
      dependencies: ["setup"],
      testMatch: /specs\/.*\.spec\.ts$/,
      use: { browserName: "chromium", ...mobileEmulation },
    },
    {
      name: "mobile-webkit",
      dependencies: ["setup"],
      testMatch: /specs\/.*\.spec\.ts$/,
      use: { browserName: "webkit", ...mobileEmulation },
    },
  ],
});
