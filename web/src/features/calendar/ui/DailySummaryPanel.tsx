import type { GetDailySummaryResponse } from "@/shared/api/client";

type Props = {
  date: string;
  data: GetDailySummaryResponse | undefined;
  isLoading: boolean;
};

export function DailySummaryPanel({ date, data, isLoading }: Props) {
  const label = new Date(date + "T00:00:00").toLocaleDateString("ja-JP", {
    month: "long",
    day: "numeric",
    weekday: "short",
  });

  return (
    <div className="mt-4 border-t border-[var(--color-line)] pt-4">
      <h3 className="text-sm font-bold mb-3 text-[var(--color-ink)]">{label}</h3>
      {isLoading ? (
        <p className="text-xs text-[var(--color-ink-muted)]">読み込み中…</p>
      ) : !data ? null : (
        <div className="space-y-3">
          {(data.trainings ?? []).length > 0 && (
            <Section title="トレーニング">
              {data.trainings!.map((t) => (
                <Row key={t.training_id}>
                  <span>{t.exercise_count}種目 / {t.set_count}セット</span>
                  <span className="text-[var(--color-ink-muted)] text-xs">
                    {new Date(t.started_at!).toLocaleTimeString("ja-JP", { hour: "2-digit", minute: "2-digit" })}
                  </span>
                </Row>
              ))}
            </Section>
          )}
          {(data.meals ?? []).length > 0 && (
            <Section title="食事">
              {data.meals!.map((m) => (
                <Row key={m.meal_id}>
                  <span>{m.meal_type}</span>
                  <span className="text-[var(--color-ink-muted)] text-xs">{m.calories}kcal</span>
                </Row>
              ))}
            </Section>
          )}
          {(data.weights ?? []).length > 0 && (
            <Section title="体重">
              {data.weights!.map((w) => (
                <Row key={w.weight_id}>
                  <span>{w.weight_kg}kg</span>
                  {w.body_fat_percentage && (
                    <span className="text-[var(--color-ink-muted)] text-xs">体脂肪 {w.body_fat_percentage}%</span>
                  )}
                </Row>
              ))}
            </Section>
          )}
          {(data.trainings ?? []).length === 0 &&
           (data.meals ?? []).length === 0 &&
           (data.weights ?? []).length === 0 && (
            <p className="text-xs text-[var(--color-ink-muted)]">この日の記録はありません</p>
          )}
        </div>
      )}
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <p className="text-xs font-semibold text-[var(--color-ink-muted)] mb-1">{title}</p>
      <div className="space-y-1">{children}</div>
    </div>
  );
}

function Row({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center justify-between text-sm px-2 py-1 rounded bg-[var(--color-surface-alt)]">
      {children}
    </div>
  );
}
