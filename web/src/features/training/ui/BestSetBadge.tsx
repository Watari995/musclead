"use client";

import { useExerciseBestSetQuery } from "../api/exercises";

type Props = {
  exerciseID: string;
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

export function BestSetBadge({ exerciseID }: Props) {
  const { data, isLoading } = useExerciseBestSetQuery(exerciseID);

  // 種目未選択 / 取得中はチラつき防止で何も出さない。
  if (!exerciseID || isLoading) return null;

  if (!data) {
    return (
      <p className="mt-1 text-xs text-[var(--color-ink-muted)]">
        まだ記録がありません
      </p>
    );
  }

  const weight = formatWeight(data.weight_kg ?? "");
  const date = data.performed_at ? formatDate(data.performed_at) : "";

  return (
    <p className="mt-1 text-xs text-[var(--color-ink-muted)]">
      <span className="text-[var(--color-accent)]">★</span>{" "}
      <span className="font-medium text-[var(--color-ink)]">最高記録</span>{" "}
      {weight}kg × {data.reps ?? 0}回
      {date && <span className="ml-1 opacity-70">({date})</span>}
    </p>
  );
}
