"use client";

import Link from "next/link";
import type { Exercise } from "@/features/training/model/exercise";
import type { ExerciseDraft, SetDraft } from "@/features/training/model/training-draft";
import { Button, Card, Label, TextInput } from "@/shared/ui";
import { SetField } from "./SetField";

type Props = {
  exercise: ExerciseDraft;
  index: number;
  exercises: Exercise[];
  onChange: (patch: Partial<Omit<ExerciseDraft, "key" | "sets" | "displayOrder">>) => void;
  onRemove: () => void;
  onMoveUp?: () => void;
  onMoveDown?: () => void;
  onAddSet: () => void;
  onChangeSet: (
    setIndex: number,
    patch: Partial<Omit<SetDraft, "key" | "setNumber">>,
  ) => void;
  onRemoveSet: (setIndex: number) => void;
  disabled?: boolean;
};

export function ExerciseField({
  exercise,
  index,
  exercises,
  onChange,
  onRemove,
  onMoveUp,
  onMoveDown,
  onAddSet,
  onChangeSet,
  onRemoveSet,
  disabled,
}: Props) {
  return (
    <Card className="p-4 space-y-3">
      {/* header: 種目選択 + 並び替え/削除 */}
      <div className="flex items-start gap-2">
        <div className="flex-1">
          <Label label={`種目 ${index + 1}`}>
            <select
              value={exercise.exerciseID}
              onChange={(e) => onChange({ exerciseID: e.target.value })}
              disabled={disabled}
              required
              aria-label={`種目${index + 1}を選択`}
              className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-white text-[var(--color-ink)] focus:outline-none focus:border-[var(--color-ink)] transition-colors"
            >
              <option value="" disabled>
                種目を選択…
              </option>
              {exercises.map((ex) => (
                <option key={ex.id} value={ex.id}>
                  {ex.name}
                </option>
              ))}
            </select>
          </Label>
          {exercises.length === 0 && (
            <p className="text-xs text-[var(--color-ink-muted)] mt-1">
              まだ種目が登録されていません。{" "}
              <Link
                href="/exercises/new"
                className="underline hover:opacity-70"
              >
                先に作成
              </Link>
            </p>
          )}
        </div>
        <div className="flex items-end gap-1 pb-px">
          <button
            type="button"
            onClick={onMoveUp}
            disabled={disabled || !onMoveUp}
            className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-11"
            aria-label="上へ"
          >
            ↑
          </button>
          <button
            type="button"
            onClick={onMoveDown}
            disabled={disabled || !onMoveDown}
            className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-11"
            aria-label="下へ"
          >
            ↓
          </button>
          <button
            type="button"
            onClick={onRemove}
            disabled={disabled}
            className="text-xs text-[var(--color-accent)] disabled:opacity-50 px-2 h-11"
            aria-label={`種目${index + 1}を削除`}
          >
            削除
          </button>
        </div>
      </div>

      {/* exercise メタ: 既定の休憩 + メモ */}
      <div className="grid grid-cols-2 gap-3">
        <Label label="既定の休憩(秒)">
          <TextInput
            type="number"
            min={0}
            value={exercise.restSeconds ?? ""}
            placeholder="例: 90"
            onChange={(e) =>
              onChange({
                restSeconds:
                  e.target.value === "" ? null : Number(e.target.value),
              })
            }
            disabled={disabled}
          />
        </Label>
        <Label label="メモ">
          <TextInput
            value={exercise.memo}
            onChange={(e) => onChange({ memo: e.target.value })}
            placeholder="任意"
            disabled={disabled}
          />
        </Label>
      </div>

      {/* sets */}
      <div className="space-y-2">
        <div className="grid grid-cols-[auto_minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)_auto] gap-2 text-xs font-medium text-[var(--color-ink-muted)]">
          <span className="w-8" />
          <span>重量(kg)</span>
          <span>レップ</span>
          <span>休憩(秒、 空欄で既定)</span>
          <span />
        </div>
        {exercise.sets.map((set, setIndex) => (
          <SetField
            key={set.key}
            set={set}
            exerciseDefaultRest={exercise.restSeconds}
            onChange={(patch) => onChangeSet(setIndex, patch)}
            onRemove={() => onRemoveSet(setIndex)}
            disabled={disabled}
          />
        ))}
        <Button
          type="button"
          variant="ghost"
          onClick={onAddSet}
          disabled={disabled}
          className="w-full"
        >
          + セットを追加
        </Button>
      </div>
    </Card>
  );
}
