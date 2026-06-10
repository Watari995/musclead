import type { Page } from "@playwright/test";

export type TargetPage = {
  path: string;
  /**
   * インタラクション (ポップオーバーを開くなど) を加えた状態でも崩れていないか
   * を検査したいときに指定する。
   */
  interact?: (page: Page) => Promise<void>;
  /**
   * axe-core / overflow チェックの基準になる「ページが落ち着いた」状態を待つ。
   * デフォルトでは networkidle を待つ。
   */
  waitFor?: (page: Page) => Promise<void>;
};

export const TARGET_PAGES: readonly TargetPage[] = [
  { path: "/login" },
  { path: "/register" },
  { path: "/meals" },
  { path: "/weights" },
  {
    path: "/trainings",
    interact: async (page) => {
      const btn = page.getByRole("button", { name: /ルーティンから始める/ });
      if ((await btn.count()) > 0 && (await btn.isEnabled())) {
        await btn.click();
        // Popover 内容 (Popover.tsx で role="dialog") の出現を確実に待つ
        await page.getByRole("dialog").waitFor({ state: "visible" });
      }
    },
  },
  { path: "/routines" },
  { path: "/exercises" },
  { path: "/profile" },
  { path: "/settings" },
  { path: "/trainings/new" },
  { path: "/routines/new" },
  { path: "/exercises/new" },
];
