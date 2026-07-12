"use client";

import { useTranslations } from "next-intl";
import type { SetDraft } from "@/features/training/model/training-draft";
import { NumberField, TextInput } from "@/shared/ui";

type Props = {
  set: SetDraft;
  exerciseDefaultRest: number | null;
  onChange: (patch: Partial<Omit<SetDraft, "key" | "setNumber">>) => void;
  onRemove: () => void;
  disabled?: boolean;
};

// モバイル: カード状にスタック (#番号 + 削除を上段、 重量/レップ/休憩を 3 カラム)
// sm 以上: 1 行に並べる従来レイアウト
export function SetField({
  set,
  exerciseDefaultRest,
  onChange,
  onRemove,
  disabled,
}: Props) {
  const t = useTranslations("trainings");
  const tCommon = useTranslations("common");

  return (
    <>
      {/* モバイル */}
      <div className="sm:hidden rough p-3 space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-xs font-bold tracking-tight text-[var(--color-ink-muted)]">
            {t("setNumber", { num: set.setNumber })}
          </span>
          <button
            type="button"
            onClick={onRemove}
            disabled={disabled}
            className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2"
            aria-label={t("setDelete", { num: set.setNumber })}
          >
            {tCommon("delete")}
          </button>
        </div>
        <div className="grid grid-cols-3 gap-2">
          <MobileField label={t("weightKg")}>
            <TextInput
              type="text"
              inputMode="decimal"
              placeholder="kg"
              value={set.weightKg}
              onChange={(e) => onChange({ weightKg: e.target.value })}
              aria-label={t("setWeightAria", { num: set.setNumber })}
              disabled={disabled}
            />
          </MobileField>
          <MobileField label={t("reps")}>
            <NumberField
              min={0}
              placeholder={t("repsUnit")}
              value={set.reps || undefined}
              onChange={(v) => onChange({ reps: v ?? 0 })}
              aria-label={t("setRepsAria", { num: set.setNumber })}
              disabled={disabled}
            />
          </MobileField>
          <MobileField label={t("restSecondsMobile")}>
            <NumberField
              min={0}
              placeholder={
                exerciseDefaultRest !== null ? `${exerciseDefaultRest}` : t("secPlaceholder")
              }
              value={set.restSeconds ?? undefined}
              onChange={(v) => onChange({ restSeconds: v ?? null })}
              aria-label={t("setRestAria", { num: set.setNumber })}
              disabled={disabled}
            />
          </MobileField>
        </div>
      </div>

      {/* sm 以上 */}
      <div className="hidden sm:grid grid-cols-[auto_minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)_auto] gap-2 items-center">
        <span className="text-xs font-bold tracking-tight text-[var(--color-ink-muted)] w-8 shrink-0">
          #{set.setNumber}
        </span>
        <TextInput
          type="text"
          inputMode="decimal"
          placeholder="kg"
          value={set.weightKg}
          onChange={(e) => onChange({ weightKg: e.target.value })}
          aria-label={t("setWeightAria", { num: set.setNumber })}
          disabled={disabled}
        />
        <NumberField
          min={0}
          placeholder={t("repsUnit")}
          value={set.reps || undefined}
          onChange={(v) => onChange({ reps: v ?? 0 })}
          aria-label={t("setRepsAria", { num: set.setNumber })}
          disabled={disabled}
        />
        <NumberField
          min={0}
          placeholder={
            exerciseDefaultRest !== null
              ? t("defaultRestPlaceholder", { seconds: exerciseDefaultRest })
              : t("defaultRestSec")
          }
          value={set.restSeconds ?? undefined}
          onChange={(v) => onChange({ restSeconds: v ?? null })}
          aria-label={t("setRestAria", { num: set.setNumber })}
          disabled={disabled}
        />
        <button
          type="button"
          onClick={onRemove}
          disabled={disabled}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2"
          aria-label={t("setDelete", { num: set.setNumber })}
        >
          {tCommon("delete")}
        </button>
      </div>
    </>
  );
}

function MobileField({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <label className="block">
      <span className="block text-[10px] text-[var(--color-ink-muted)] mb-1">
        {label}
      </span>
      {children}
    </label>
  );
}
