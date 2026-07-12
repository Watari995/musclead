"use client";

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
import {
  useBestSetsQuery,
  useExercisesQuery,
  useLastSessionSetsQuery,
} from "@/features/training/api/exercises";
import { useTranslations } from "next-intl";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";
import { ExerciseField } from "./ExerciseField";

type Props = {
  /** 親が保持する draft(controlled)。 親側で下書きの自動退避などを担う。 */
  value: TrainingDraft;
  onChange: (draft: TrainingDraft) => void;
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (payload: RecordTrainingRequest) => void | Promise<void>;
  submitting?: boolean;
  errorMessage?: string | null;
  onCancel?: () => void;
};

export function TrainingForm({
  value: draft,
  onChange,
  submitLabel,
  submittingLabel,
  onSubmit,
  submitting = false,
  errorMessage,
  onCancel,
}: Props) {
  const t = useTranslations("trainings");
  const tCommon = useTranslations("common");

  // 既存の `setDraft((d) => ...)` 呼び出しをそのまま使えるよう、
  // 関数型アップデートを controlled な onChange に橋渡しするアダプタ。
  const setDraft = (
    updater: TrainingDraft | ((prev: TrainingDraft) => TrainingDraft),
  ) =>
    onChange(typeof updater === "function" ? updater(draft) : updater);

  const exercisesQuery = useExercisesQuery();
  const exercises = exercisesQuery.data ?? [];

  // 選択中の全種目の最高記録を 1 リクエストでまとめて取得(N+1 回避)。
  const selectedExerciseIDs = draft.exercises
    .map((e) => e.exerciseID)
    .filter(Boolean);
  const bestSetsQuery = useBestSetsQuery(selectedExerciseIDs);
  const bestSets = bestSetsQuery.data;

  const lastSessionSetsQuery = useLastSessionSetsQuery(selectedExerciseIDs);
  const lastSessionSets = lastSessionSetsQuery.data;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    void onSubmit(toRecordRequest(draft));
  };

  return (
    <form className="space-y-6" onSubmit={handleSubmit}>
      {/* メタ情報 */}
      <Card className="p-4 sm:p-5 space-y-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <Label label={t("startTime")}>
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
          <Label label={t("endTime")}>
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
        <Label label={tCommon("memo")}>
          <textarea
            value={draft.memo}
            onChange={(e) =>
              setDraft((d) => updateTraining(d, { memo: e.target.value }))
            }
            rows={2}
            className="block w-full px-3 py-2 rough bg-[var(--color-surface)] focus:outline-none focus:[--rough-color:var(--color-accent)]"
            placeholder={t("sessionMemo")}
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
              bestSet={bestSets?.get(exercise.exerciseID) ?? null}
              bestSetLoading={bestSetsQuery.isLoading}
              lastSession={lastSessionSets?.get(exercise.exerciseID) ?? null}
              lastSessionLoading={lastSessionSetsQuery.isLoading}
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
          disabled={submitting || draft.exercises.length === 0}
          className="flex-1"
        >
          {submitting ? submittingLabel : submitLabel}
        </Button>
      </div>
    </form>
  );
}
