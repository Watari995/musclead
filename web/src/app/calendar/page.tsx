"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { useMonthlyCalendarQuery, useDailyCalendarQuery } from "@/features/calendar/api/calendar";
import { usePreferencesQuery } from "@/features/user/api/user";
import { CalendarView } from "@/features/calendar/ui/CalendarView";
import { DailySummaryPanel } from "@/features/calendar/ui/DailySummaryPanel";
import { Card } from "@/shared/ui";

export default function CalendarPage() {
  const { token, ready } = useAccessToken();
  const authenticated = ready && Boolean(token);

  const now = new Date();
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [selectedDate, setSelectedDate] = useState<string | null>(
    now.toISOString().slice(0, 10),
  );

  const { data: monthlySummary, isLoading: monthlyLoading } =
    useMonthlyCalendarQuery(year, month, authenticated);
  const { data: preferences } = usePreferencesQuery(authenticated);
  const { data: dailySummary, isLoading: dailyLoading } = useDailyCalendarQuery(
    selectedDate ?? "",
    authenticated && Boolean(selectedDate),
  );

  const handlePrevMonth = () => {
    if (month === 1) {
      setYear((y) => y - 1);
      setMonth(12);
    } else {
      setMonth((m) => m - 1);
    }
    setSelectedDate(null);
  };

  const handleNextMonth = () => {
    if (month === 12) {
      setYear((y) => y + 1);
      setMonth(1);
    } else {
      setMonth((m) => m + 1);
    }
    setSelectedDate(null);
  };

  if (!ready) {
    return <p className="text-[var(--color-ink-muted)] text-sm">読み込み中…</p>;
  }

  return (
    <div className="max-w-lg mx-auto">
      <h1 className="text-2xl font-bold mb-4">カレンダー</h1>
      <Card className="p-4">
        {monthlyLoading ? (
          <p className="text-[var(--color-ink-muted)] text-sm text-center py-8">読み込み中…</p>
        ) : (
          <CalendarView
            year={year}
            month={month}
            days={monthlySummary?.days ?? []}
            preferences={preferences}
            selectedDate={selectedDate}
            onSelectDate={setSelectedDate}
            onPrevMonth={handlePrevMonth}
            onNextMonth={handleNextMonth}
          />
        )}
        {selectedDate && (
          <DailySummaryPanel
            date={selectedDate}
            data={dailySummary}
            isLoading={dailyLoading}
          />
        )}
      </Card>
      <div className="mt-4">
        <div className="flex gap-4 text-xs text-[var(--color-ink-muted)]">
          <span className="flex items-center gap-1">
            <span
              className="inline-block w-2 h-2 rounded-full"
              style={{ background: preferences?.training_color ?? "#4A90E2" }}
            />
            トレーニング
          </span>
          <span className="flex items-center gap-1">
            <span
              className="inline-block w-2 h-2 rounded-full"
              style={{ background: preferences?.meal_color ?? "#7ED321" }}
            />
            食事
          </span>
          <span className="flex items-center gap-1">
            <span
              className="inline-block w-2 h-2 rounded-full"
              style={{ background: preferences?.weight_color ?? "#FF6B6B" }}
            />
            体重
          </span>
        </div>
      </div>
    </div>
  );
}
