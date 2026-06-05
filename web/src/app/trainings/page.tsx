"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import type { TrainingDTO } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import { useExercisesQuery } from "@/features/training/api/exercises";
import {
  useDeleteTrainingMutation,
  useTrainingsQuery,
} from "@/features/training/api/trainings";
import { formatDateTime } from "@/features/training/model/training-draft";
import {
  Button,
  Card,
  ErrorText,
  SectionTitle,
  useConfirm,
} from "@/shared/ui";

export default function TrainingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const confirm = useConfirm();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useTrainingsQuery(Boolean(token));
  const exercisesQuery = useExercisesQuery(Boolean(token));
  const exerciseNameByID = new Map<string, string>();
  for (const ex of exercisesQuery.data ?? []) {
    exerciseNameByID.set(ex.id, ex.name);
  }

  const del = useDeleteTrainingMutation();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>トレーニング履歴</SectionTitle>
        <Link href="/trainings/new">
          <Button>+ 新規記録</Button>
        </Link>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}

      {query.data && query.data.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          まだトレーニングが記録されていません。
        </Card>
      )}

      {query.data && query.data.length > 0 && (
        <ul className="space-y-3">
          {query.data.map((t) => (
            <TrainingCard
              key={t.id}
              training={t}
              exerciseNameByID={exerciseNameByID}
              onDelete={async () => {
                const ok = await confirm({
                  title: "トレーニングを削除しますか?",
                  description: `${formatDateTime(t.started_at)} のトレーニング記録を削除します。 この操作は取り消せません。`,
                  confirmLabel: "削除する",
                  destructive: true,
                });
                if (ok) del.mutate(t.id ?? "");
              }}
              deleting={del.isPending}
            />
          ))}
        </ul>
      )}
    </div>
  );
}

function TrainingCard({
  training,
  exerciseNameByID,
  onDelete,
  deleting,
}: {
  training: TrainingDTO;
  exerciseNameByID: Map<string, string>;
  onDelete: () => void;
  deleting: boolean;
}) {
  const totalSets = (training.exercises ?? []).reduce(
    (n, ex) => n + (ex.sets?.length ?? 0),
    0,
  );
  const exerciseNames = (training.exercises ?? [])
    .map((ex) => (ex.exercise_id ? exerciseNameByID.get(ex.exercise_id) : undefined))
    .filter((n): n is string => Boolean(n))
    .slice(0, 3);

  return (
    <li className="bg-white border border-[var(--color-line)] rounded-lg p-4 flex items-start justify-between gap-4">
      <Link
        href={`/trainings/${training.id}`}
        className="flex-1 min-w-0 space-y-1 hover:opacity-70 transition-opacity"
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {formatDateTime(training.started_at)}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {(training.exercises ?? []).length} 種目 / {totalSets} セット
          </span>
        </div>
        {exerciseNames.length > 0 && (
          <p className="text-sm text-[var(--color-ink)] line-clamp-1">
            {exerciseNames.join(" / ")}
            {(training.exercises ?? []).length > 3 && " …"}
          </p>
        )}
        {training.memo && (
          <p className="text-xs text-[var(--color-ink-muted)] line-clamp-1">
            {training.memo}
          </p>
        )}
      </Link>
      <div className="flex flex-col gap-1 shrink-0">
        <Link
          href={`/trainings/${training.id}/edit`}
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
