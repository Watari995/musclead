"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useRoutineQuery,
  useStartTrainingFromRoutineMutation,
} from "@/features/training/api/routines";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function RoutineDetailPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useRoutineQuery(params.id, Boolean(token));
  const startTraining = useStartTrainingFromRoutineMutation();

  if (!ready || !token) return null;

  if (query.isLoading) {
    return (
      <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
    );
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  const routine = query.data;
  const exercises = routine.routine_exercises ?? [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-4">
        <SectionTitle>{routine.name}</SectionTitle>
        <Link href={`/routines/${routine.id}/edit`}>
          <Button variant="ghost">編集</Button>
        </Link>
      </div>

      <Card className="p-4">
        <p className="text-xs text-[var(--color-ink-muted)] mb-3">
          {exercises.length} 種目
        </p>
        {exercises.length === 0 ? (
          <p className="text-sm text-[var(--color-ink-muted)]">
            種目が登録されていません。
          </p>
        ) : (
          <ol className="space-y-2">
            {exercises.map((re, idx) => (
              <li
                key={re.id}
                className="flex items-center justify-between text-sm border-b border-[var(--color-line)] last:border-b-0 pb-2 last:pb-0"
              >
                <span>
                  <span className="text-[var(--color-ink-muted)] mr-2">
                    {idx + 1}.
                  </span>
                  {re.exercise_name}
                </span>
              </li>
            ))}
          </ol>
        )}
      </Card>

      {startTraining.isError && (
        <ErrorText>{(startTraining.error as Error).message}</ErrorText>
      )}

      <Button
        type="button"
        onClick={() =>
          startTraining.mutate(routine, {
            onSuccess: (data) => {
              const id = data.training_id ?? "";
              if (id) router.push(`/trainings/${id}/edit`);
            },
          })
        }
        disabled={exercises.length === 0 || startTraining.isPending}
        fullWidth
      >
        {startTraining.isPending
          ? "Training を作成中…"
          : "このルーティンで Training を開始"}
      </Button>
    </div>
  );
}
