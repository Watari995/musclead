"use client";

import type { LastSessionSetsByExerciseDTO } from "@/shared/api/client";

type Props = {
  lastSession: LastSessionSetsByExerciseDTO | null;
  loading: boolean;
};

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

export function LastSessionBadge({ lastSession, loading }: Props) {
  if (loading) return null;
  if (!lastSession?.sets?.length) return null;

  const date = lastSession.performed_at ? formatDate(lastSession.performed_at) : "";

  return (
    <div className="mt-1 text-xs text-[var(--color-ink-muted)]">
      <span className="font-medium">
        📅 前回{date && <span className="ml-1 opacity-70">({date})</span>}
      </span>
      <span className="ml-2 inline-flex flex-wrap gap-x-3">
        {lastSession.sets.map((s) => (
          <span key={s.set_number}>
            {s.set_number}.{" "}
            <span className="text-[var(--color-ink)]">
              {formatWeight(s.weight_kg ?? "")}kg × {s.reps}回
            </span>
          </span>
        ))}
      </span>
    </div>
  );
}
