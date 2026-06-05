"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import {
  apiClient,
  type TrainingDTO,
  type TrainingExerciseDTO,
} from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import { useExercisesQuery } from "@/lib/queries/exercises";
import { formatDateTime, resolveRestSeconds } from "@/lib/training-form";
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

  const query = useQuery({
    queryKey: ["training", params.id],
    enabled: Boolean(token && params.id),
    queryFn: async (): Promise<TrainingDTO> => {
      const { data, error, response } = await apiClient.GET("/trainings/{id}", {
        params: { path: { id: params.id } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as TrainingDTO;
    },
  });

  const exercisesQuery = useExercisesQuery(Boolean(token));
  const exerciseNameByID = new Map<string, string>();
  for (const ex of exercisesQuery.data ?? []) {
    if (ex.id && ex.name) exerciseNameByID.set(ex.id, ex.name);
  }

  if (!ready || !token) return null;

  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
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
              終了: {formatDateTime(training.ended_at)}
            </p>
          )}
        </div>
        <Link href={`/trainings/${training.id}/edit`}>
          <Button variant="ghost">編集</Button>
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
              ex.exercise_id ? exerciseNameByID.get(ex.exercise_id) ?? "(削除済み)" : "-"
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
  return (
    <Card className="p-4 space-y-3">
      <div className="flex items-baseline justify-between">
        <h3 className="text-sm font-bold tracking-tight">{exerciseName}</h3>
        <span className="text-xs text-[var(--color-ink-muted)]">
          {(exercise.sets ?? []).length} セット
          {exercise.rest_seconds != null && ` / 既定休憩 ${exercise.rest_seconds}秒`}
        </span>
      </div>

      {exercise.memo && (
        <p className="text-xs text-[var(--color-ink-muted)]">{exercise.memo}</p>
      )}

      <table className="w-full text-sm">
        <thead>
          <tr className="text-xs text-[var(--color-ink-muted)] text-left">
            <th className="w-12 font-medium pb-1">#</th>
            <th className="font-medium pb-1">重量</th>
            <th className="font-medium pb-1">レップ</th>
            <th className="font-medium pb-1">休憩</th>
            <th className="font-medium pb-1">メモ</th>
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
                  {rest !== null ? `${rest}秒` : "—"}
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
