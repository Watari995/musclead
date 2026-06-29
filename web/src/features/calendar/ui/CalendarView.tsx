"use client";

import type { MonthlySummaryDayDTO } from "@/shared/api/client";
import type { PreferencesDTO } from "@/shared/api/client";

type Props = {
  year: number;
  month: number;
  days: MonthlySummaryDayDTO[];
  preferences: PreferencesDTO | undefined;
  selectedDate: string | null;
  onSelectDate: (date: string) => void;
  onPrevMonth: () => void;
  onNextMonth: () => void;
};

const WEEKDAYS = ["日", "月", "火", "水", "木", "金", "土"];

export function CalendarView({
  year,
  month,
  days,
  preferences,
  selectedDate,
  onSelectDate,
  onPrevMonth,
  onNextMonth,
}: Props) {
  const dayMap = new Map(days.map((d) => [d.date ?? "", d]));
  const firstDay = new Date(year, month - 1, 1).getDay();
  const daysInMonth = new Date(year, month, 0).getDate();
  const todayStr = new Date().toISOString().slice(0, 10);

  const cells: (number | null)[] = [
    ...Array<null>(firstDay).fill(null),
    ...Array.from({ length: daysInMonth }, (_, i) => i + 1),
  ];

  const trainingColor = preferences?.training_color ?? "#4A90E2";
  const mealColor = preferences?.meal_color ?? "#7ED321";
  const weightColor = preferences?.weight_color ?? "#FF6B6B";

  return (
    <div>
      {/* header */}
      <div className="flex items-center justify-between mb-4">
        <button
          type="button"
          onClick={onPrevMonth}
          className="p-2 rounded-md hover:bg-[var(--color-surface-alt)] transition-colors"
          aria-label="前の月"
        >
          ‹
        </button>
        <span className="text-base font-bold">
          {year}年{month}月
        </span>
        <button
          type="button"
          onClick={onNextMonth}
          className="p-2 rounded-md hover:bg-[var(--color-surface-alt)] transition-colors"
          aria-label="次の月"
        >
          ›
        </button>
      </div>

      {/* weekday labels */}
      <div className="grid grid-cols-7 mb-1">
        {WEEKDAYS.map((d, i) => (
          <div
            key={d}
            className={`text-center text-xs font-semibold py-1 ${i === 0 ? "text-[var(--color-accent)]" : i === 6 ? "text-blue-500" : "text-[var(--color-ink-muted)]"}`}
          >
            {d}
          </div>
        ))}
      </div>

      {/* grid */}
      <div className="grid grid-cols-7 gap-y-1">
        {cells.map((day, idx) => {
          if (day === null) return <div key={`empty-${idx}`} />;
          const dateStr = `${year}-${String(month).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
          const info = dayMap.get(dateStr);
          const isToday = dateStr === todayStr;
          const isSelected = dateStr === selectedDate;
          const dow = (firstDay + day - 1) % 7;

          return (
            <button
              key={dateStr}
              type="button"
              onClick={() => onSelectDate(dateStr)}
              className={`flex flex-col items-center py-1 rounded-lg transition-colors ${
                isSelected
                  ? "bg-[var(--color-ink)] text-[var(--color-surface)]"
                  : "hover:bg-[var(--color-surface-alt)]"
              }`}
            >
              <span
                className={`text-sm font-medium leading-6 w-7 h-7 flex items-center justify-center rounded-full ${
                  isSelected
                    ? "text-[var(--color-surface)]"
                    : isToday
                      ? "border border-[var(--color-ink)]"
                      : dow === 0
                        ? "text-[var(--color-accent)]"
                        : dow === 6
                          ? "text-blue-500"
                          : ""
                }`}
              >
                {day}
              </span>
              <div className="flex gap-0.5 h-2 mt-0.5">
                {info?.has_training && (
                  <span
                    className="w-1.5 h-1.5 rounded-full"
                    style={{ background: trainingColor }}
                  />
                )}
                {info?.has_meal && (
                  <span
                    className="w-1.5 h-1.5 rounded-full"
                    style={{ background: mealColor }}
                  />
                )}
                {info?.has_weight && (
                  <span
                    className="w-1.5 h-1.5 rounded-full"
                    style={{ background: weightColor }}
                  />
                )}
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
