"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import {
  apiClient,
  type ExerciseDTO,
  type ListExercisesResponse,
} from "@/api/client";
import { useAccessToken } from "@/lib/access-token";
import { Button, Card, ErrorText, SectionTitle } from "@/components/ui";

const EXERCISES_QUERY_KEY = ["exercises", "all"] as const;

export default function ExercisesPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useQuery({
    queryKey: EXERCISES_QUERY_KEY,
    enabled: Boolean(token),
    queryFn: async (): Promise<ListExercisesResponse> => {
      const { data, error, response } = await apiClient.GET("/exercises", {
        params: { query: { limit: 100, offset: 0 } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as ListExercisesResponse;
    },
  });

  const del = useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE("/exercises/{id}", {
        params: { path: { id } },
      });
      if (error) {
        const code = error.error?.code;
        if (code === "training.exercise_in_use_error") {
          throw new Error(
            "この種目はトレーニング履歴で使われているため削除できません。",
          );
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY }),
  });

  if (!ready || !token) return null;

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
      {del.isError && <ErrorText>{(del.error as Error).message}</ErrorText>}

      {query.data && (query.data.exercises ?? []).length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          まだ種目が登録されていません。 「+ 新しい種目」 から作成してください。
        </Card>
      )}

      {query.data && (query.data.exercises ?? []).length > 0 && (
        <ul className="space-y-2">
          {(query.data.exercises ?? []).map((ex) => (
            <ExerciseRow
              key={ex.id}
              exercise={ex}
              onDelete={() => {
                if (confirm(`「${ex.name}」を削除しますか?`)) {
                  del.mutate(ex.id ?? "");
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
  onDelete,
  deleting,
}: {
  exercise: ExerciseDTO;
  onDelete: () => void;
  deleting: boolean;
}) {
  return (
    <li className="bg-white border border-[var(--color-line)] rounded-lg p-4 flex items-center justify-between gap-4">
      <Link
        href={`/exercises/${exercise.id}/edit`}
        className="flex-1 min-w-0 hover:opacity-70 transition-opacity"
      >
        <p className="text-sm font-bold tracking-tight">{exercise.name}</p>
        <p className="text-xs text-[var(--color-ink-muted)]">
          登録: {new Date(exercise.created_at ?? "").toLocaleDateString("ja-JP")}
        </p>
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
