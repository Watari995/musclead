"use client";

import { useState } from "react";
import {
  type TrainingDraft,
  addExercise,
  addSet,
  moveExercise,
  removeExercise,
  removeSet,
  toRecordRequest,
  updateExercise,
  updateSet,
  updateTraining,
} from "@/features/training/model/training-draft";
import { type RecordTrainingRequest } from "@/shared/api/client";
import { useExercisesQuery } from "@/features/training/api/exercises";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";
import { ExerciseField } from "./ExerciseField";

type Props = {
  initial: TrainingDraft;
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (payload: RecordTrainingRequest) => void | Promise<void>;
  submitting?: boolean;
  errorMessage?: string | null;
  onCancel?: () => void;
};

export function TrainingForm({
  initial,
  submitLabel,
  submittingLabel,
  onSubmit,
  submitting = false,
  errorMessage,
  onCancel,
}: Props) {
  const [draft, setDraft] = useState<TrainingDraft>(initial);

  const exercisesQuery = useExercisesQuery();
  const exercises = exercisesQuery.data ?? [];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    void onSubmit(toRecordRequest(draft));
  };

  return (
    <form className="space-y-6" onSubmit={handleSubmit}>
      {/* メタ情報 */}
      <Card className="p-4 sm:p-5 space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <Label label="開始時刻">
            <TextInput
              type="datetime-local"
              value={draft.startedAt}
              onChange={(e) =>
                setDraft((d) => updateTraining(d, { startedAt: e.target.value }))
              }
              required
              disabled={submitting}
            />
          </Label>
          <Label label="終了時刻(任意)">
            <TextInput
              type="datetime-local"
              value={draft.endedAt}
              onChange={(e) =>
                setDraft((d) => updateTraining(d, { endedAt: e.target.value }))
              }
              disabled={submitting}
            />
          </Label>
        </div>
        <Label label="メモ">
          <textarea
            value={draft.memo}
            onChange={(e) =>
              setDraft((d) => updateTraining(d, { memo: e.target.value }))
            }
            rows={2}
            className="block w-full px-3 py-2 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] focus:outline-none focus:border-[var(--color-ink)]"
            placeholder="セッション全体のメモ(任意)"
            disabled={submitting}
          />
        </Label>
      </Card>

      {/* 種目リスト */}
      <div className="space-y-4">
        {draft.exercises.map((exercise, index) => {
          const last = draft.exercises.length - 1;
          return (
            <ExerciseField
              key={exercise.key}
              exercise={exercise}
              index={index}
              exercises={exercises}
              onChange={(patch) =>
                setDraft((d) => updateExercise(d, index, patch))
              }
              onRemove={() => setDraft((d) => removeExercise(d, index))}
              onMoveUp={
                index === 0
                  ? undefined
                  : () => setDraft((d) => moveExercise(d, index, index - 1))
              }
              onMoveDown={
                index === last
                  ? undefined
                  : () => setDraft((d) => moveExercise(d, index, index + 1))
              }
              onAddSet={() => setDraft((d) => addSet(d, index))}
              onChangeSet={(setIndex, patch) =>
                setDraft((d) => updateSet(d, index, setIndex, patch))
              }
              onRemoveSet={(setIndex) =>
                setDraft((d) => removeSet(d, index, setIndex))
              }
              disabled={submitting}
            />
          );
        })}
        <Button
          type="button"
          variant="ghost"
          onClick={() => setDraft((d) => addExercise(d))}
          disabled={submitting}
          fullWidth
        >
          + 種目を追加
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
            キャンセル
          </Button>
        )}
        <Button
          type="submit"
          disabled={submitting || draft.exercises.length === 0}
          className="flex-1"
        >
          {submitting ? submittingLabel : submitLabel}
        </Button>
      </div>
    </form>
  );
}
