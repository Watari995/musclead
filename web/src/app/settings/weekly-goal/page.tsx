"use client";

import { useEffect, useState } from "react";
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
  return (
    <div className="flex flex-col gap-1">
      <label className="text-sm font-medium text-[var(--color-ink)]">{label}</label>
      <div className="flex items-center gap-2">
        <input
          type="number"
          value={value ?? ""}
          placeholder="未設定"
          onChange={(e) =>
            onChange(e.target.value === "" ? null : Number(e.target.value))
          }
          className="w-32 h-9 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] px-3 text-sm text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-[var(--color-ink)]"
        />
        <span className="text-sm text-[var(--color-ink-muted)]">{unit}</span>
      </div>
    </div>
  );
}

export default function WeeklyGoalSettingsPage() {
  const { data, isLoading } = useWeeklyGoalQuery();
  const upsert = useUpsertWeeklyGoalMutation();

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
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <h2 className="text-base font-semibold text-[var(--color-ink)]">週次目標</h2>
        <p className="text-sm text-[var(--color-ink-muted)] mt-1">
          毎週日曜に達成チェックの通知が届きます。設定しない項目は空欄のままにしてください。
        </p>
      </div>

      <div className="space-y-4">
        <NullableIntInput
          label="トレーニング回数"
          unit="回 / 週"
          value={trainingCount}
          onChange={setTrainingCount}
        />
        <NullableIntInput
          label="平均カロリー目標"
          unit="kcal 以内 / 日"
          value={calorieAverage}
          onChange={setCalorieAverage}
        />
        <div className="flex flex-col gap-1">
          <label className="text-sm font-medium text-[var(--color-ink)]">体重変化目標</label>
          <div className="flex items-center gap-2">
            <input
              type="number"
              step="0.1"
              value={weightChangeKg ?? ""}
              placeholder="未設定"
              onChange={(e) =>
                setWeightChangeKg(e.target.value === "" ? null : Number(e.target.value))
              }
              className="w-32 h-9 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] px-3 text-sm text-[var(--color-ink)] focus:outline-none focus:ring-2 focus:ring-[var(--color-ink)]"
            />
            <span className="text-sm text-[var(--color-ink-muted)]">kg（例: -0.5 で減量目標）</span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <Button type="submit" disabled={upsert.isPending}>
          {upsert.isPending ? "保存中…" : "保存"}
        </Button>
        {saved && <span className="text-sm text-green-600">保存しました</span>}
        {upsert.isError && <span className="text-sm text-red-500">保存に失敗しました</span>}
      </div>
    </form>
  );
}
