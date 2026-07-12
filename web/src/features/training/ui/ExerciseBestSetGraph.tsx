"use client";

import { useState } from "react";
import {
  Bar,
  BarChart,
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
  useExerciseBestSetTimeseriesQuery,
  type BestSetTimeseriesPeriod,
} from "@/features/training/api/exercises";
import { Card, ErrorText } from "@/shared/ui";

const PERIOD_VALUES: BestSetTimeseriesPeriod[] = [
  "1week", "1month", "3months", "halfyear", "1year",
];

const PERIOD_KEY_MAP: Record<BestSetTimeseriesPeriod, string> = {
  "1week": "period1Week",
  "1month": "period1Month",
  "3months": "period3Months",
  "halfyear": "periodHalfYear",
  "1year": "period1Year",
};

function formatDate(iso: string): string {
  const d = new Date(iso);
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

const PERIOD_BUTTON_CLASS = (active: boolean) =>
  `flex-1 px-2 py-1.5 text-xs rounded-md whitespace-nowrap ${
    active
      ? "bg-[var(--color-ink)] text-[var(--color-surface)] font-bold"
      : "text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)]"
  }`;

function EmptyState({ loading, error }: { loading: boolean; error: Error | null }) {
  const tCommon = useTranslations("common");
  if (loading) return <p className="text-sm text-center py-16 text-[var(--color-ink-muted)]">{tCommon("loading")}</p>;
  if (error) return <ErrorText>{error.message}</ErrorText>;
  return <p className="text-sm text-center py-16 text-[var(--color-ink-muted)]">{tCommon("noData")}</p>;
}

export function ExerciseBestSetGraph({ exerciseId }: { exerciseId: string | null }) {
  const t = useTranslations("graph");
  const [period, setPeriod] = useState<BestSetTimeseriesPeriod>("1month");
  const query = useExerciseBestSetTimeseriesQuery(exerciseId, period, Boolean(exerciseId));

  const points = (query.data?.data_points ?? []).map((p) => ({
    date: formatDate(p.performed_at),
    weight: parseFloat(p.weight_kg),
    reps: p.reps,
  }));

  const hasData = points.length > 0;
  const showEmpty = !hasData && !query.isLoading;

  return (
    <div className="space-y-4">
      {/* 重量 LineChart */}
      <Card className="p-4 sm:p-5">
        <p className="text-sm font-semibold text-[var(--color-ink-muted)] mb-3">{t("weightProgress")}</p>
        <div className="h-52">
          {showEmpty && (
            <EmptyState
              loading={query.isLoading}
              error={query.isError ? (query.error as Error) : null}
            />
          )}
          {hasData && (
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={points}>
                <CartesianGrid stroke="var(--color-line)" strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} domain={["auto", "auto"]} width={45} />
                <Tooltip formatter={(v) => [`${String(v)} kg`, t("weightLabel")]} />
                <Line
                  type="monotone"
                  dataKey="weight"
                  stroke="var(--color-ink)"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          )}
        </div>
      </Card>

      {/* レップス BarChart */}
      <Card className="p-4 sm:p-5">
        <p className="text-sm font-semibold text-[var(--color-ink-muted)] mb-3">{t("repsProgress")}</p>
        <div className="h-52">
          {showEmpty && (
            <EmptyState
              loading={query.isLoading}
              error={query.isError ? (query.error as Error) : null}
            />
          )}
          {hasData && (
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={points}>
                <CartesianGrid stroke="var(--color-line)" strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 11 }} />
                <YAxis
                  tick={{ fontSize: 11 }}
                  domain={[0, "auto"]}
                  width={45}
                  allowDecimals={false}
                />
                <Tooltip formatter={(v) => [`${String(v)} reps`, t("repsLabel")]} />
                <Bar dataKey="reps" fill="var(--color-ink)" radius={[3, 3, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          )}
        </div>
      </Card>

      {/* 期間タブ */}
      <div className="flex gap-1 overflow-x-auto border-t border-[var(--color-line)] pt-3">
        {PERIOD_VALUES.map((value) => (
          <button
            key={value}
            type="button"
            onClick={() => setPeriod(value)}
            className={PERIOD_BUTTON_CLASS(period === value)}
          >
            {t(PERIOD_KEY_MAP[value] as "period1Week" | "period1Month" | "period3Months" | "periodHalfYear" | "period1Year")}
          </button>
        ))}
      </div>
    </div>
  );
}
