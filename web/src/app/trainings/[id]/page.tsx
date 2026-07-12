"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import type { TrainingExerciseDTO } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import { useExercisesQuery } from "@/features/training/api/exercises";
import { useTrainingQuery } from "@/features/training/api/trainings";
import { formatDateTime, resolveRestSeconds } from "@/features/training/model/training-draft";
import {
  Button,
  Card,
  ErrorText,
  SectionTitle,
} from "@/shared/ui";

export default function TrainingDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useTrainingQuery(params.id, Boolean(token));

  const exercisesQuery = useExercisesQuery(Boolean(token));
  const exerciseNameByID = new Map<string, string>();
  for (const ex of exercisesQuery.data ?? []) {
    exerciseNameByID.set(ex.id, ex.name);
  }

  const t = useTranslations("trainings");
  const tc = useTranslations("common");

  if (!ready || !token) return null;

  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  const training = query.data;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <SectionTitle>
            {formatDateTime(training.started_at)}
          </SectionTitle>
          {training.ended_at && (
            <p className="text-xs text-[var(--color-ink-muted)] -mt-2 mb-3">
              {t("ended")} {formatDateTime(training.ended_at)}
            </p>
          )}
        </div>
        <Link href={`/trainings/${training.id}/edit`}>
          <Button variant="ghost">{tc("edit")}</Button>
        </Link>
      </div>

      {training.memo && (
        <Card className="p-4 text-sm text-[var(--color-ink)]">
          {training.memo}
        </Card>
      )}

      <div className="space-y-3">
        {(training.exercises ?? []).map((ex) => (
          <ExerciseSummary
            key={ex.id}
            exercise={ex}
            exerciseName={
              ex.exercise_id ? exerciseNameByID.get(ex.exercise_id) ?? t("deleted") : "-"
            }
          />
        ))}
      </div>
    </div>
  );
}

function ExerciseSummary({
  exercise,
  exerciseName,
}: {
  exercise: TrainingExerciseDTO;
  exerciseName: string;
}) {
  const t = useTranslations("trainings");
  const tc = useTranslations("common");
  return (
    <Card className="p-4 space-y-3">
      <div className="flex items-baseline justify-between">
        <h3 className="text-sm font-bold tracking-tight">{exerciseName}</h3>
        <span className="text-xs text-[var(--color-ink-muted)]">
          {(exercise.sets ?? []).length} {t("sets")}
          {exercise.rest_seconds != null && ` / ${t("defaultRest", { seconds: exercise.rest_seconds })}`}
        </span>
      </div>

      {exercise.memo && (
        <p className="text-xs text-[var(--color-ink-muted)]">{exercise.memo}</p>
      )}

      <table className="w-full text-sm">
        <thead>
          <tr className="text-xs text-[var(--color-ink-muted)] text-left">
            <th className="w-12 font-medium pb-1">#</th>
            <th className="font-medium pb-1">{t("weight")}</th>
            <th className="font-medium pb-1">{t("reps")}</th>
            <th className="font-medium pb-1">{t("rest")}</th>
            <th className="font-medium pb-1">{tc("memo")}</th>
          </tr>
        </thead>
        <tbody>
          {(exercise.sets ?? []).map((set) => {
            const rest = resolveRestSeconds(
              set.rest_seconds,
              exercise.rest_seconds,
            );
            return (
              <tr
                key={set.id}
                className="border-t border-[var(--color-line)]"
              >
                <td className="py-2 text-[var(--color-ink-muted)]">
                  {set.set_number}
                </td>
                <td className="py-2 font-medium">{set.weight_kg} kg</td>
                <td className="py-2">{set.reps}</td>
                <td className="py-2 text-[var(--color-ink-muted)]">
                  {rest !== null ? `${rest}${t("secUnit")}` : "—"}
                </td>
                <td className="py-2 text-[var(--color-ink-muted)] text-xs">
                  {set.memo ?? ""}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </Card>
  );
}
