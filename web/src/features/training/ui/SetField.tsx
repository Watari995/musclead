"use client";

import type { SetDraft } from "@/features/training/model/training-draft";
import { TextInput } from "@/shared/ui";

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
  return (
    <>
      {/* モバイル */}
      <div className="sm:hidden border border-[var(--color-line)] rounded-md p-3 space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-xs font-bold tracking-tight text-[var(--color-ink-muted)]">
            セット #{set.setNumber}
          </span>
          <button
            type="button"
            onClick={onRemove}
            disabled={disabled}
            className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2"
            aria-label={`セット${set.setNumber}を削除`}
          >
            削除
          </button>
        </div>
        <div className="grid grid-cols-3 gap-2">
          <MobileField label="重量(kg)">
            <TextInput
              type="text"
              inputMode="decimal"
              placeholder="kg"
              value={set.weightKg}
              onChange={(e) => onChange({ weightKg: e.target.value })}
              aria-label={`セット${set.setNumber}の重量`}
              disabled={disabled}
            />
          </MobileField>
          <MobileField label="レップ">
            <TextInput
              type="number"
              min={0}
              placeholder="回"
              value={set.reps || ""}
              onChange={(e) => onChange({ reps: Number(e.target.value) })}
              aria-label={`セット${set.setNumber}のレップ`}
              disabled={disabled}
            />
          </MobileField>
          <MobileField label="休憩(秒)">
            <TextInput
              type="number"
              min={0}
              placeholder={
                exerciseDefaultRest !== null ? `${exerciseDefaultRest}` : "秒"
              }
              value={set.restSeconds ?? ""}
              onChange={(e) =>
                onChange({
                  restSeconds: e.target.value === "" ? null : Number(e.target.value),
                })
              }
              aria-label={`セット${set.setNumber}の休憩秒数`}
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
          aria-label={`セット${set.setNumber}の重量`}
          disabled={disabled}
        />
        <TextInput
          type="number"
          min={0}
          placeholder="回"
          value={set.reps || ""}
          onChange={(e) => onChange({ reps: Number(e.target.value) })}
          aria-label={`セット${set.setNumber}のレップ`}
          disabled={disabled}
        />
        <TextInput
          type="number"
          min={0}
          placeholder={
            exerciseDefaultRest !== null
              ? `${exerciseDefaultRest}秒(既定)`
              : "休憩 秒"
          }
          value={set.restSeconds ?? ""}
          onChange={(e) =>
            onChange({
              restSeconds: e.target.value === "" ? null : Number(e.target.value),
            })
          }
          aria-label={`セット${set.setNumber}の休憩秒数`}
          disabled={disabled}
        />
        <button
          type="button"
          onClick={onRemove}
          disabled={disabled}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2"
          aria-label={`セット${set.setNumber}を削除`}
        >
          削除
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
