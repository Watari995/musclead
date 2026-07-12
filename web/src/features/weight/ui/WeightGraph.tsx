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
import { useTranslations } from "next-intl";
import {
  useWeightTimeseriesQuery,
  type WeightTimeseriesPeriod,
} from "@/features/weight/api/weights";
import type { Weight } from "@/features/weight/model/weight";
import { Card, ErrorText } from "@/shared/ui";

const PERIOD_VALUES: WeightTimeseriesPeriod[] = [
  "1week", "1month", "3months", "halfyear", "1year",
];

const PERIOD_KEY_MAP: Record<WeightTimeseriesPeriod, string> = {
  "1week": "period1Week",
  "1month": "period1Month",
  "3months": "period3Months",
  "halfyear": "periodHalfYear",
  "1year": "period1Year",
};

type WeightType = "weight" | "body_fat" | "muscle";

const WEIGHT_TYPE_VALUES: WeightType[] = ["weight", "body_fat", "muscle"];

function getValue(w: Weight, type: WeightType): number | null {
  if (type === "weight") return parseFloat(w.weightKg);
  if (type === "body_fat")
    return w.bodyFatPercentage ? parseFloat(w.bodyFatPercentage) : null;
  if (type === "muscle")
    return w.skeletalMuscleKg ? parseFloat(w.skeletalMuscleKg) : null;
  return null;
}

function getUnit(type: WeightType): string {
  if (type === "body_fat") return "%";
  return "kg";
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

export function WeightGraph() {
  const t = useTranslations("weights");
  const tGraph = useTranslations("graph");
  const tCommon = useTranslations("common");
  const [period, setPeriod] = useState<WeightTimeseriesPeriod>("1month");
  const [type, setType] = useState<WeightType>("weight");
  const query = useWeightTimeseriesQuery(period);

  const typeLabels: Record<WeightType, string> = {
    weight: t("typeWeight"),
    body_fat: t("typeBodyFat"),
    muscle: t("typeMuscle"),
  };

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
          {WEIGHT_TYPE_VALUES.map((opt) => (
            <button
              key={opt}
              type="button"
              onClick={() => setType(opt)}
              className={`rough px-3 py-1.5 text-xs whitespace-nowrap ${
                type === opt
                  ? "bg-[var(--color-ink)] text-[var(--color-surface)]"
                  : "text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)]"
              }`}
            >
              {typeLabels[opt]}
            </button>
          ))}
        </div>

        {/* グラフエリア */}
        <div className="h-64">
          {query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)] py-12 text-center">
              {tCommon("loading")}
            </p>
          )}
          {query.isError && (
            <ErrorText>{(query.error as Error).message}</ErrorText>
          )}
          {query.data && points.length === 0 && !query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)] py-12 text-center">
              {tCommon("noData")}
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
                    `${String(value)} ${getUnit(type)}`,
                    typeLabels[type],
                  ]}
                />
                <Line
                  type="monotone"
                  dataKey="value"
                  stroke="var(--color-accent)"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          )}
        </div>

        {/* period タブ */}
        <div className="flex gap-1 overflow-x-auto border-t border-[var(--color-line)] pt-3">
          {PERIOD_VALUES.map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => setPeriod(value)}
              className={`flex-1 px-2 py-1.5 text-xs rounded-md whitespace-nowrap ${
                period === value
                  ? "bg-[var(--color-ink)] text-[var(--color-surface)] font-bold"
                  : "text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)]"
              }`}
            >
              {tGraph(PERIOD_KEY_MAP[value] as "period1Week" | "period1Month" | "period3Months" | "periodHalfYear" | "period1Year")}
            </button>
          ))}
        </div>
      </div>
    </Card>
  );
}
