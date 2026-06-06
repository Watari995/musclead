import AxeBuilder from "@axe-core/playwright";
import { test, expect } from "@playwright/test";
import { TARGET_PAGES } from "../helpers/pages";

// critical / serious 以外 (moderate / minor) は初期段階のノイズが多いので落とさない。
// 段階的に閾値を下げる方針。
const FAIL_IMPACTS = new Set(["critical", "serious"]);

for (const target of TARGET_PAGES) {
  test(`a11y (critical+serious): ${target.path}`, async ({ page }) => {
    await page.goto(target.path);
    await (target.waitFor?.(page) ?? page.waitForLoadState("networkidle"));
    if (target.interact) await target.interact(page);

    const results = await new AxeBuilder({ page })
      .withTags(["wcag2a", "wcag2aa", "wcag21aa", "best-practice"])
      .analyze();

    const violations = results.violations
      .filter((v) => v.impact && FAIL_IMPACTS.has(v.impact))
      .map((v) => ({
        id: v.id,
        impact: v.impact,
        help: v.help,
        helpUrl: v.helpUrl,
        nodes: v.nodes.length,
        firstNodeTarget: v.nodes[0]?.target,
      }));

    expect(
      violations,
      `${target.path}: a11y violations (critical/serious)\n${JSON.stringify(violations, null, 2)}`,
    ).toEqual([]);
  });
}
