"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  ExerciseInUseError,
  useDeleteExerciseMutation,
  useExercisesQuery,
  useReorderExercisesMutation,
} from "@/features/training/api/exercises";
import type { Exercise } from "@/features/training/model/exercise";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function ExercisesPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useExercisesQuery(Boolean(token));
  const del = useDeleteExerciseMutation();
  const reorder = useReorderExercisesMutation();

  if (!ready || !token) return null;

  const exercises = query.data ?? [];

  const move = (index: number, direction: -1 | 1) => {
    const to = index + direction;
    if (to < 0 || to >= exercises.length) return;
    const next = [...exercises];
    [next[index], next[to]] = [next[to], next[index]];
    reorder.mutate(next);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>種目マスタ</SectionTitle>
        <Link href="/exercises/new">
          <Button>+ 新しい種目</Button>
        </Link>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {del.isError && (
        <ErrorText>
          {del.error instanceof ExerciseInUseError
            ? del.error.message
            : (del.error as Error).message}
        </ErrorText>
      )}
      {reorder.isError && (
        <ErrorText>並び替えに失敗しました。 時間をおいて再度お試しください。</ErrorText>
      )}

      {query.data && exercises.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          まだ種目が登録されていません。 「+ 新しい種目」 から作成してください。
        </Card>
      )}

      {exercises.length > 0 && (
        <ul className="space-y-2">
          {exercises.map((ex, index) => (
            <ExerciseRow
              key={ex.id}
              exercise={ex}
              onMoveUp={index > 0 ? () => move(index, -1) : undefined}
              onMoveDown={
                index < exercises.length - 1 ? () => move(index, 1) : undefined
              }
              reordering={reorder.isPending}
              onDelete={() => {
                if (confirm(`「${ex.name}」 を削除しますか?`)) {
                  del.mutate(ex.id);
                }
              }}
              deleting={del.isPending}
            />
          ))}
        </ul>
      )}
    </div>
  );
}

function ExerciseRow({
  exercise,
  onMoveUp,
  onMoveDown,
  reordering,
  onDelete,
  deleting,
}: {
  exercise: Exercise;
  onMoveUp?: () => void;
  onMoveDown?: () => void;
  reordering: boolean;
  onDelete: () => void;
  deleting: boolean;
}) {
  return (
    <li className="bg-[var(--color-surface)] border border-[var(--color-line)] rounded-lg p-4 flex items-center justify-between gap-2">
      <div className="flex items-center gap-1 shrink-0">
        <button
          type="button"
          onClick={onMoveUp}
          disabled={reordering || !onMoveUp}
          className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-8"
          aria-label="上へ"
        >
          ↑
        </button>
        <button
          type="button"
          onClick={onMoveDown}
          disabled={reordering || !onMoveDown}
          className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-30 px-2 h-8"
          aria-label="下へ"
        >
          ↓
        </button>
      </div>
      <Link
        href={`/exercises/${exercise.id}/edit`}
        className="flex-1 min-w-0 hover:opacity-70 transition-opacity"
      >
        <p className="text-sm font-bold tracking-tight">{exercise.name}</p>
      </Link>
      <div className="flex gap-1 shrink-0">
        <Link
          href={`/exercises/${exercise.id}/edit`}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] px-2 h-8 inline-flex items-center"
        >
          編集
        </Link>
        <button
          type="button"
          onClick={onDelete}
          disabled={deleting}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2 h-8"
        >
          削除
        </button>
      </div>
    </li>
  );
}
