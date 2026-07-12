"use client";

import { useMemo } from "react";
import { useTranslations } from "next-intl";
import { DayPicker } from "react-day-picker";
import { Popover } from "@/shared/ui";

type Props = {
  /** ISO short date "YYYY-MM-DD" or "" (空欄) */
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  required?: boolean;
};

export function BirthdayInput({
  value,
  onChange,
  disabled,
  required,
}: Props) {
  const t = useTranslations("birthdayInput");
  const selected = useMemo(() => parseISODate(value), [value]);

  const today = useMemo(() => new Date(), []);
  const startMonth = useMemo(() => new Date(1900, 0), []);
  const endMonth = useMemo(
    () => new Date(today.getFullYear(), 11),
    [today],
  );
  const defaultMonth = useMemo(
    () => selected ?? new Date(today.getFullYear() - 30, 0),
    [selected, today],
  );

  return (
    <>
      <Popover
        trigger={(p) => (
          <button
            type="button"
            onClick={p.onClick}
            disabled={disabled}
            aria-expanded={p["aria-expanded"]}
            aria-controls={p["aria-controls"]}
            aria-haspopup="dialog"
            className="block w-full h-11 px-3 rough bg-[var(--color-surface)] text-left text-[var(--color-ink)] focus:outline-none focus:[--rough-color:var(--color-accent)] transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-between"
          >
            <span
              className={
                selected
                  ? "text-[var(--color-ink)]"
                  : "text-[var(--color-ink-muted)]"
              }
            >
              {selected ? t("dateFormat", { year: selected.getFullYear(), month: selected.getMonth() + 1, day: selected.getDate() }) : t("placeholder")}
            </span>
            <CalendarIcon />
          </button>
        )}
      >
        <div className="p-2">
          <DayPicker
            mode="single"
            selected={selected}
            onSelect={(d) => onChange(d ? toISOShort(d) : "")}
            captionLayout="dropdown"
            startMonth={startMonth}
            endMonth={endMonth}
            defaultMonth={defaultMonth}
            disabled={{ after: today }}
          />
        </div>
      </Popover>
      {/* form submit validation 用の hidden input */}
      {required && (
        <input
          tabIndex={-1}
          aria-hidden
          required
          value={value}
          onChange={() => {}}
          className="sr-only"
        />
      )}
    </>
  );
}

function parseISODate(iso: string): Date | undefined {
  if (!iso) return undefined;
  // 'YYYY-MM-DD' を Local timezone の Date に展開 (UTC 解釈で日付がズレる事象を回避)
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(iso);
  if (!m) return undefined;
  return new Date(Number(m[1]), Number(m[2]) - 1, Number(m[3]));
}

function toISOShort(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}


function CalendarIcon() {
  return (
    <svg
      width="18"
      height="18"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      className="text-[var(--color-ink-muted)] shrink-0"
      aria-hidden
    >
      <rect x="3" y="4" width="18" height="18" rx="2" />
      <path d="M16 2v4M8 2v4M3 10h18" />
    </svg>
  );
}
