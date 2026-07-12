"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/shared/ui";
import {
  useWeeklyGoalQuery,
  useUpsertWeeklyGoalMutation,
} from "@/features/notification/api/notifications";

function NullableIntInput({
  label,
  unit,
  value,
  onChange,
}: {
  label: string;
  unit: string;
  value: number | null;
  onChange: (v: number | null) => void;
}) {
  const tc = useTranslations("common");
  return (
    <div className="flex flex-col gap-1">
      <label className="text-sm font-medium text-[var(--color-ink)]">{label}</label>
      <div className="flex items-center gap-2">
        <input
          type="number"
          value={value ?? ""}
          placeholder={tc("notSet")}
          onChange={(e) =>
            onChange(e.target.value === "" ? null : Number(e.target.value))
          }
          className="w-32 h-9 rough bg-[var(--color-surface)] px-3 text-sm text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-[var(--color-ink)]"
        />
        <span className="text-sm text-[var(--color-ink-muted)]">{unit}</span>
      </div>
    </div>
  );
}

export default function WeeklyGoalSettingsPage() {
  const { data, isLoading } = useWeeklyGoalQuery();
  const upsert = useUpsertWeeklyGoalMutation();
  const t = useTranslations("weeklyGoal");
  const tc = useTranslations("common");

  const [trainingCount, setTrainingCount] = useState<number | null>(null);
  const [calorieAverage, setCalorieAverage] = useState<number | null>(null);
  const [weightChangeKg, setWeightChangeKg] = useState<number | null>(null);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (data) {
      setTrainingCount(data.training_count);
      setCalorieAverage(data.calorie_average);
      setWeightChangeKg(data.weight_change_kg);
    }
  }, [data]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    await upsert.mutateAsync({ training_count: trainingCount, calorie_average: calorieAverage, weight_change_kg: weightChangeKg });
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  }

  if (isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>;
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <h2 className="text-base font-semibold text-[var(--color-ink)]">{t("title")}</h2>
        <p className="text-sm text-[var(--color-ink-muted)] mt-1">
          {t("desc")}
        </p>
      </div>

      <div className="space-y-4">
        <NullableIntInput
          label={t("trainingCount")}
          unit={t("trainingCountUnit")}
          value={trainingCount}
          onChange={setTrainingCount}
        />
        <NullableIntInput
          label={t("calorieAverage")}
          unit={t("calorieAverageUnit")}
          value={calorieAverage}
          onChange={setCalorieAverage}
        />
        <div className="flex flex-col gap-1">
          <label className="text-sm font-medium text-[var(--color-ink)]">{t("weightChange")}</label>
          <div className="flex items-center gap-2">
            <input
              type="number"
              step="0.1"
              value={weightChangeKg ?? ""}
              placeholder={tc("notSet")}
              onChange={(e) =>
                setWeightChangeKg(e.target.value === "" ? null : Number(e.target.value))
              }
              className="w-32 h-9 rough bg-[var(--color-surface)] px-3 text-sm text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-[var(--color-ink)]"
            />
            <span className="text-sm text-[var(--color-ink-muted)]">{t("weightChangeUnit")}</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <Button type="submit" disabled={upsert.isPending}>
          {upsert.isPending ? tc("saving") : tc("save")}
        </Button>
        {saved && <span className="text-sm text-green-600">{tc("saved")}</span>}
        {upsert.isError && <span className="text-sm text-red-500">{t("saveFailed")}</span>}
      </div>
    </form>
  );
}
