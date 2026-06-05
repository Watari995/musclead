"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import type { RoutineDTO } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useDeleteRoutineMutation,
  useRoutinesQuery,
} from "@/features/training/api/routines";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function RoutinesPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useRoutinesQuery(Boolean(token));
  const del = useDeleteRoutineMutation();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>ルーティン</SectionTitle>
        <Link href="/routines/new">
          <Button>+ 新しいルーティン</Button>
        </Link>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {del.isError && <ErrorText>{(del.error as Error).message}</ErrorText>}

      {query.data && query.data.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          まだルーティンが登録されていません。 「+ 新しいルーティン」 から作成してください。
        </Card>
      )}

      {query.data && query.data.length > 0 && (
        <ul className="space-y-3">
          {query.data.map((r) => (
            <RoutineCard
              key={r.id}
              routine={r}
              onDelete={() => {
                if (confirm(`「${r.name}」 を削除しますか?`)) {
                  del.mutate(r.id ?? "");
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

function RoutineCard({
  routine,
  onDelete,
  deleting,
}: {
  routine: RoutineDTO;
  onDelete: () => void;
  deleting: boolean;
}) {
  const exerciseNames = (routine.routine_exercises ?? [])
    .map((e) => e.exercise_name)
    .filter(Boolean)
    .slice(0, 3);

  return (
    <li className="bg-white border border-[var(--color-line)] rounded-lg p-4 flex items-start justify-between gap-4">
      <Link
        href={`/routines/${routine.id}`}
        className="flex-1 min-w-0 space-y-1 hover:opacity-70 transition-opacity"
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {routine.name}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {(routine.routine_exercises ?? []).length} 種目
          </span>
        </div>
        {exerciseNames.length > 0 && (
          <p className="text-sm text-[var(--color-ink)] line-clamp-1">
            {exerciseNames.join(" / ")}
            {(routine.routine_exercises ?? []).length > 3 && " …"}
          </p>
        )}
      </Link>
      <div className="flex flex-col gap-1 shrink-0">
        <Link
          href={`/routines/${routine.id}/edit`}
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
