"use client";

import { useState } from "react";
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import {
  useWeightTimeseriesQuery,
  type WeightTimeseriesPeriod,
} from "@/features/weight/api/weights";
import type { Weight } from "@/features/weight/model/weight";
import { Card, ErrorText } from "@/shared/ui";

const PERIOD_OPTIONS: { value: WeightTimeseriesPeriod; label: string }[] = [
  { value: "1week", label: "1週間" },
  { value: "1month", label: "1ヶ月" },
  { value: "3months", label: "3ヶ月" },
  { value: "halfyear", label: "半年" },
  { value: "1year", label: "1年" },
];

type WeightType = "weight" | "body_fat" | "muscle";

const TYPE_OPTIONS: { value: WeightType; label: string; unit: string }[] = [
  { value: "weight", label: "体重", unit: "kg" },
  { value: "body_fat", label: "体脂肪率", unit: "%" },
  { value: "muscle", label: "骨格筋量", unit: "kg" },
];

function getValue(w: Weight, type: WeightType): number | null {
  if (type === "weight") return parseFloat(w.weightKg);
  if (type === "body_fat")
    return w.bodyFatPercentage ? parseFloat(w.bodyFatPercentage) : null;
  if (type === "muscle")
    return w.skeletalMuscleKg ? parseFloat(w.skeletalMuscleKg) : null;
  return null;
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

export function WeightGraph() {
  const [period, setPeriod] = useState<WeightTimeseriesPeriod>("1month");
  const [type, setType] = useState<WeightType>("weight");
  const query = useWeightTimeseriesQuery(period);

  const typeOption = TYPE_OPTIONS.find((o) => o.value === type);
  const points = (query.data ?? [])
    .map((w) => ({
      date: formatDate(w.measuredAt),
      value: getValue(w, type),
    }))
    .filter((p): p is { date: string; value: number } => p.value !== null);

  return (
    <Card className="p-4 sm:p-5">
      <div className="space-y-4">
        {/* type タブ */}
        <div className="flex gap-1 overflow-x-auto">
          {TYPE_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              type="button"
              onClick={() => setType(opt.value)}
              className={`px-3 py-1.5 text-xs rounded-md border whitespace-nowrap ${
                type === opt.value
                  ? "bg-[var(--color-ink)] text-[var(--color-surface)] border-[var(--color-ink)]"
                  : "border-[var(--color-line)] text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)]"
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>

        {/* グラフエリア */}
        <div className="h-64">
          {query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)] py-12 text-center">
              読み込み中…
            </p>
          )}
          {query.isError && (
            <ErrorText>{(query.error as Error).message}</ErrorText>
          )}
          {query.data && points.length === 0 && !query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)] py-12 text-center">
              データがありません
            </p>
          )}
          {points.length > 0 && (
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={points}>
                <CartesianGrid stroke="var(--color-line)" strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} />
                <YAxis
                  tick={{ fontSize: 11 }}
                  domain={["auto", "auto"]}
                  width={40}
                />
                <Tooltip
                  formatter={(value) => [
                    `${String(value)} ${typeOption?.unit ?? ""}`,
                    typeOption?.label ?? "",
                  ]}
                />
                <Line
                  type="monotone"
                  dataKey="value"
                  stroke="var(--color-ink)"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          )}
        </div>

        {/* period タブ */}
        <div className="flex gap-1 overflow-x-auto border-t border-[var(--color-line)] pt-3">
          {PERIOD_OPTIONS.map((opt) => (
            <button
              key={opt.value}
              type="button"
              onClick={() => setPeriod(opt.value)}
              className={`flex-1 px-2 py-1.5 text-xs rounded-md whitespace-nowrap ${
                period === opt.value
                  ? "bg-[var(--color-ink)] text-[var(--color-surface)] font-bold"
                  : "text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)]"
              }`}
            >
              {opt.label}
            </button>
          ))}
        </div>
      </div>
    </Card>
  );
}
