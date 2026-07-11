"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import type { Exercise } from "@/features/training/model/exercise";
import type { ExerciseDraft, SetDraft } from "@/features/training/model/training-draft";
import type { BestSetDTO, LastSessionSetsByExerciseDTO } from "@/shared/api/client";
import { Button, Card, Label, NumberField, TextInput } from "@/shared/ui";
import { BestSetBadge } from "./BestSetBadge";
import { LastSessionBadge } from "./LastSessionBadge";
import { SetField } from "./SetField";

type Props = {
  exercise: ExerciseDraft;
  index: number;
  exercises: Exercise[];
  bestSet: BestSetDTO | null;
  bestSetLoading: boolean;
  lastSession: LastSessionSetsByExerciseDTO | null;
  lastSessionLoading: boolean;
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
  bestSet,
  bestSetLoading,
  lastSession,
  lastSessionLoading,
  onChange,
  onRemove,
  onMoveUp,
  onMoveDown,
  onAddSet,
  onChangeSet,
  onRemoveSet,
  disabled,
}: Props) {
  const t = useTranslations("trainings");
  const tCommon = useTranslations("common");
  const tRoutines = useTranslations("routines");

  return (
    <Card className="p-3 sm:p-4 space-y-3">
      {/* header: 種目選択 + 並び替え/削除 */}
      <div className="flex items-start gap-1 sm:gap-2">
        <div className="flex-1">
          <Label label={t("exercise", { index: index + 1 })}>
            <select
              value={exercise.exerciseID}
              onChange={(e) => onChange({ exerciseID: e.target.value })}
              disabled={disabled}
              required
              aria-label={t("exerciseSelect", { index: index + 1 })}
              className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] text-[var(--color-ink)] focus:outline-none focus:border-[var(--color-ink)] transition-colors"
            >
              <option value="" disabled>
                {tRoutines("selectExercise")}
              </option>
              {exercises.map((ex) => (
                <option key={ex.id} value={ex.id}>
                  {ex.name}
                </option>
              ))}
            </select>
          </Label>
          {exercise.exerciseID && (
            <>
              <BestSetBadge bestSet={bestSet} loading={bestSetLoading} />
              <LastSessionBadge lastSession={lastSession} loading={lastSessionLoading} />
            </>
          )}
          {exercises.length === 0 && (
            <p className="text-xs text-[var(--color-ink-muted)] mt-1">
              {tRoutines("noExercisesCreate")}{" "}
              <Link
                href="/exercises/new"
                className="underline hover:opacity-70"
              >
                {tRoutines("createFirst")}
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
            aria-label={tCommon("up")}
          >
            ↑
          </button>
          <button
            type="button"
            onClick={onMoveDown}
            disabled={disabled || !onMoveDown}
            className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-11"
            aria-label={tCommon("down")}
          >
            ↓
          </button>
          <button
            type="button"
            onClick={onRemove}
            disabled={disabled}
            className="text-xs text-[var(--color-accent)] disabled:opacity-50 px-2 h-11"
            aria-label={t("exerciseDelete", { index: index + 1 })}
          >
            {tCommon("delete")}
          </button>
        </div>
      </div>

      {/* exercise メタ: 既定の休憩 + メモ */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <Label label={t("restExerciseDefault")}>
          <NumberField
            min={0}
            value={exercise.restSeconds ?? undefined}
            placeholder="90"
            onChange={(v) => onChange({ restSeconds: v ?? null })}
            disabled={disabled}
          />
        </Label>
        <Label label={tCommon("memo")}>
          <TextInput
            value={exercise.memo}
            onChange={(e) => onChange({ memo: e.target.value })}
            placeholder={t("memoOptional")}
            disabled={disabled}
          />
        </Label>
      </div>

      {/* sets */}
      <div className="space-y-2">
        <div className="hidden sm:grid grid-cols-[auto_minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)_auto] gap-2 text-xs font-medium text-[var(--color-ink-muted)]">
          <span className="w-8" />
          <span>{t("weightKg")}</span>
          <span>{t("reps")}</span>
          <span>{t("restSecondsHeader")}</span>
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
          {t("addSet")}
        </Button>
      </div>
    </Card>
  );
}
