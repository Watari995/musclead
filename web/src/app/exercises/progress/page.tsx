"use client";

import { Suspense, useEffect } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useAccessToken } from "@/shared/auth/access-token";
import { useExercisesQuery } from "@/features/training/api/exercises";
import { ExerciseBestSetGraph } from "@/features/training/ui/ExerciseBestSetGraph";
import { SectionTitle } from "@/shared/ui";

function ExerciseProgressContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const exercisesQuery = useExercisesQuery(Boolean(token));
  const exercises = exercisesQuery.data ?? [];

  const exerciseIdParam = searchParams.get("exercise_id");
  const exerciseId = exerciseIdParam ?? (exercises[0]?.id ?? null);

  const selectedExercise = exercises.find((e) => e.id === exerciseId);

  const handleChange = (id: string) => {
    const params = new URLSearchParams();
    params.set("exercise_id", id);
    router.replace(`/exercises/progress?${params.toString()}`);
  };

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <SectionTitle>記録グラフ</SectionTitle>
        <Link
          href="/exercises"
          className="text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
        >
          ← 種目一覧
        </Link>
      </div>

      <select
        value={exerciseId ?? ""}
        onChange={(e) => handleChange(e.target.value)}
        disabled={exercisesQuery.isLoading || exercises.length === 0}
        className="w-full sm:w-72 border border-[var(--color-line)] rounded-md px-3 py-2 text-sm bg-[var(--color-surface)] text-[var(--color-ink)] focus:outline-none focus:ring-1 focus:ring-[var(--color-ink)]"
      >
        {exercises.length === 0 && <option value="">種目なし</option>}
        {exercises.map((ex) => (
          <option key={ex.id} value={ex.id}>
            {ex.name}
          </option>
        ))}
      </select>

      {selectedExercise && (
        <p className="text-lg font-bold tracking-tight">{selectedExercise.name}</p>
      )}

      {exerciseId ? (
        <ExerciseBestSetGraph exerciseId={exerciseId} />
      ) : (
        <p className="text-sm text-[var(--color-ink-muted)]">種目を選択してください。</p>
      )}
    </div>
  );
}

export default function ExerciseProgressPage() {
  return (
    <Suspense>
      <ExerciseProgressContent />
    </Suspense>
  );
}
