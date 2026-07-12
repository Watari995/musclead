"use client";

import Link from "next/link";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { type UpsertRoutineRequest } from "@/shared/api/client";
import {
  addExercise,
  moveExercise,
  removeExercise,
  setExerciseID,
  setName,
  toUpsertRequest,
  type RoutineDraft,
} from "@/features/training/model/routine-draft";
import { useExercisesQuery } from "@/features/training/api/exercises";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

type Props = {
  initial: RoutineDraft;
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (payload: UpsertRoutineRequest) => void | Promise<void>;
  submitting?: boolean;
  errorMessage?: string | null;
  onCancel?: () => void;
};

export function RoutineForm({
  initial,
  submitLabel,
  submittingLabel,
  onSubmit,
  submitting = false,
  errorMessage,
  onCancel,
}: Props) {
  const t = useTranslations("routines");
  const tCommon = useTranslations("common");
  const tTrainings = useTranslations("trainings");
  const [draft, setDraft] = useState<RoutineDraft>(initial);

  const exercisesQuery = useExercisesQuery();
  const exercises = exercisesQuery.data ?? [];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    void onSubmit(toUpsertRequest(draft));
  };

  return (
    <form className="space-y-6" onSubmit={handleSubmit}>
      <Card className="p-5 space-y-4">
        <Label label={t("routineName")}>
          <TextInput
            value={draft.name}
            onChange={(e) =>
              setDraft((d) => setName(d, e.target.value))
            }
            placeholder="PPL Day1 (Push)"
            required
            maxLength={50}
            disabled={submitting}
          />
        </Label>
      </Card>

      <div className="space-y-3">
        {draft.exercises.map((exercise, index) => {
          const last = draft.exercises.length - 1;
          return (
            <Card key={exercise.key} className="p-4 space-y-3">
              <div className="flex items-start gap-2">
                <div className="flex-1">
                  <Label label={tTrainings("exercise", { index: index + 1 })}>
                    <select
                      value={exercise.exerciseID}
                      onChange={(e) =>
                        setDraft((d) =>
                          setExerciseID(d, index, e.target.value),
                        )
                      }
                      disabled={submitting}
                      required
                      aria-label={tTrainings("exerciseSelect", { index: index + 1 })}
                      className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] text-[var(--color-ink)] focus:outline-none focus:border-[var(--color-ink)] transition-colors"
                    >
                      <option value="" disabled>
                        {t("selectExercise")}
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
                      {t("noExercisesCreate")}{" "}
                      <Link
                        href="/exercises/new"
                        className="underline hover:opacity-70"
                      >
                        {t("createFirst")}
                      </Link>
                    </p>
                  )}
                </div>
                <div className="flex items-end gap-1 pb-px">
                  <button
                    type="button"
                    onClick={() =>
                      setDraft((d) => moveExercise(d, index, index - 1))
                    }
                    disabled={submitting || index === 0}
                    className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-11"
                    aria-label={tCommon("up")}
                  >
                    ↑
                  </button>
                  <button
                    type="button"
                    onClick={() =>
                      setDraft((d) => moveExercise(d, index, index + 1))
                    }
                    disabled={submitting || index === last}
                    className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-11"
                    aria-label={tCommon("down")}
                  >
                    ↓
                  </button>
                  <button
                    type="button"
                    onClick={() => setDraft((d) => removeExercise(d, index))}
                    disabled={submitting || draft.exercises.length === 1}
                    className="text-xs text-[var(--color-accent)] disabled:opacity-50 px-2 h-11"
                    aria-label={tTrainings("exerciseDelete", { index: index + 1 })}
                  >
                    {tCommon("delete")}
                  </button>
                </div>
              </div>
            </Card>
          );
        })}
        <Button
          type="button"
          variant="ghost"
          onClick={() => setDraft((d) => addExercise(d))}
          disabled={submitting}
          fullWidth
        >
          {t("addExercise")}
        </Button>
      </div>

      {errorMessage && <ErrorText>{errorMessage}</ErrorText>}

      <div className="flex gap-3 pt-2">
        {onCancel && (
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={submitting}
          >
            {tCommon("cancel")}
          </Button>
        )}
        <Button
          type="submit"
          disabled={
            submitting ||
            draft.exercises.length === 0 ||
            !draft.name.trim() ||
            draft.exercises.some((e) => !e.exerciseID)
          }
          className="flex-1"
        >
          {submitting ? submittingLabel : submitLabel}
        </Button>
      </div>
    </form>
  );
}
