"use client";

import { useTranslations } from "next-intl";
import type { BestSetDTO } from "@/shared/api/client";

type Props = {
  bestSet: BestSetDTO | null;
  loading: boolean;
};

// "100.00" → "100" / "97.50" → "97.5" のように末尾ゼロを落とす。
function formatWeight(weightKg: string): string {
  const n = Number(weightKg);
  return Number.isFinite(n) ? String(n) : weightKg;
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${d.getFullYear()}/${m}/${day}`;
}

export function BestSetBadge({ bestSet, loading }: Props) {
  const t = useTranslations("bestSet");

  // 取得中はチラつき防止で何も出さない。
  if (loading) return null;

  if (!bestSet) {
    return (
      <p className="mt-1 text-xs text-[var(--color-ink-muted)]">
        {t("noRecord")}
      </p>
    );
  }

  const weight = formatWeight(bestSet.weight_kg ?? "");
  const date = bestSet.performed_at ? formatDate(bestSet.performed_at) : "";

  return (
    <p className="mt-1 text-xs text-[var(--color-ink-muted)]">
      <span className="text-[var(--color-accent)]">★</span>{" "}
      <span className="font-medium text-[var(--color-ink)]">{t("bestRecord")}</span>{" "}
      {weight}kg × {bestSet.reps ?? 0}回
      {date && <span className="ml-1 opacity-70">({date})</span>}
    </p>
  );
}
