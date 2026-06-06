import { test, expect, type Page } from "@playwright/test";
import { TARGET_PAGES } from "../helpers/pages";

type Offender = {
  tag: string;
  cls: string;
  id: string;
  left: number;
  right: number;
  width: number;
};

// 1px はサブピクセル丸めで起きうるノイズなので許容する
const SUBPIXEL_TOLERANCE = 1;

async function assertNoHorizontalScroll(page: Page, label: string) {
  const overflow = await page.evaluate(() => {
    const root = document.documentElement;
    return root.scrollWidth - root.clientWidth;
  });
  expect(
    overflow,
    `${label}: 水平スクロールが発生しています (${overflow}px はみ出し)`,
  ).toBeLessThanOrEqual(SUBPIXEL_TOLERANCE);
}

async function assertNoElementOutsideViewport(page: Page, label: string) {
  const offenders: Offender[] = await page.evaluate((tol) => {
    const vw = document.documentElement.clientWidth;
    const acc: Offender[] = [];
    const all = document.querySelectorAll<HTMLElement>("body *");
    for (const el of all) {
      const style = window.getComputedStyle(el);
      if (style.display === "none" || style.visibility === "hidden") continue;
      // 画面外スライドインなど、 意図的に viewport の外に置かれる UI は除外
      if (style.position === "fixed") continue;
      const r = el.getBoundingClientRect();
      if (r.width <= 0 || r.height <= 0) continue;
      if (r.left < -tol || r.right > vw + tol) {
        acc.push({
          tag: el.tagName.toLowerCase(),
          cls: (el.className || "").toString().slice(0, 120),
          id: el.id ?? "",
          left: Math.round(r.left),
          right: Math.round(r.right),
          width: Math.round(r.width),
        });
      }
    }
    return acc;
  }, SUBPIXEL_TOLERANCE);

  expect(
    offenders,
    `${label}: viewport 外にはみ出した要素があります\n${JSON.stringify(offenders, null, 2)}`,
  ).toEqual([]);
}

for (const target of TARGET_PAGES) {
  test(`mobile overflow: ${target.path}`, async ({ page }) => {
    await page.goto(target.path);
    await (target.waitFor?.(page) ?? page.waitForLoadState("networkidle"));
    if (target.interact) await target.interact(page);

    await assertNoHorizontalScroll(page, target.path);
    await assertNoElementOutsideViewport(page, target.path);
  });
}
