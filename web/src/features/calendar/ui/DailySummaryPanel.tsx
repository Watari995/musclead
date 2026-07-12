"use client";

import { useTranslations } from "next-intl";
import type { GetDailySummaryResponse } from "@/shared/api/client";

type Props = {
  date: string;
  data: GetDailySummaryResponse | undefined;
  isLoading: boolean;
};

export function DailySummaryPanel({ date, data, isLoading }: Props) {
  const t = useTranslations("calendar");
  const tCommon = useTranslations("common");

  const label = new Date(date + "T00:00:00").toLocaleDateString("ja-JP", {
    month: "long",
    day: "numeric",
    weekday: "short",
  });

  return (
    <div className="mt-4 border-t border-[var(--color-line)] pt-4">
      <h3 className="text-sm font-bold mb-3 text-[var(--color-ink)]">{label}</h3>
      {isLoading ? (
        <p className="text-xs text-[var(--color-ink-muted)]">{tCommon("loading")}</p>
      ) : !data ? null : (
        <div className="space-y-3">
          {(data.trainings ?? []).length > 0 && (
            <Section title={t("training")}>
              {data.trainings!.map((tr) => (
                <Row key={tr.training_id}>
                  <span>{t("exerciseCount", { count: tr.exercise_count ?? 0 })} / {t("setsCount", { count: tr.set_count ?? 0 })}</span>
                  <span className="text-[var(--color-ink-muted)] text-xs">
                    {new Date(tr.started_at!).toLocaleTimeString("ja-JP", { hour: "2-digit", minute: "2-digit" })}
                  </span>
                </Row>
              ))}
            </Section>
          )}
          {(data.meals ?? []).length > 0 && (
            <Section title={t("meal")}>
              <Row>
                <span>{t("total")}</span>
                <span className="text-[var(--color-ink-muted)] text-xs">{data.total_calories ?? 0}kcal</span>
              </Row>
            </Section>
          )}
          {(data.weights ?? []).length > 0 && (
            <Section title={t("weight")}>
              {data.weights!.map((w) => (
                <Row key={w.weight_id}>
                  <span>{w.weight_kg}kg</span>
                  {w.body_fat_percentage && (
                    <span className="text-[var(--color-ink-muted)] text-xs">{t("bodyFat", { value: w.body_fat_percentage })}</span>
                  )}
                </Row>
              ))}
            </Section>
          )}
          {(data.trainings ?? []).length === 0 &&
           (data.meals ?? []).length === 0 &&
           (data.weights ?? []).length === 0 && (
            <p className="text-xs text-[var(--color-ink-muted)]">{t("noRecords")}</p>
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
