"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  RoutineLimitReachedError,
  useCreateRoutineMutation,
} from "@/features/training/api/routines";
import { createInitialRoutine } from "@/features/training/model/routine-draft";
import { SectionTitle } from "@/shared/ui";
import { RoutineForm } from "@/features/training/ui/RoutineForm";

export default function NewRoutinePage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const initial = useMemo(() => createInitialRoutine(), []);
  const mutation = useCreateRoutineMutation();

  if (!ready || !token) return null;

  // 無料プランの上限到達は通常エラーではなく Pro 誘導として表示する。
  const isLimitReached = mutation.error instanceof RoutineLimitReachedError;

  return (
    <div className="space-y-6">
      <SectionTitle>新しいルーティン</SectionTitle>

      {isLimitReached && (
        <div className="rounded-md border border-[var(--color-line)] bg-[var(--color-surface-alt)] p-4 space-y-2">
          <p className="text-sm text-[var(--color-ink)]">
            ルーティンは無料プランで3件までです。 Pro
            にアップグレードすると無制限に作成できます。
          </p>
          <Link
            href="/settings/plan"
            className="inline-block text-sm font-medium text-[var(--color-ink)] underline"
          >
            プランを見る →
          </Link>
        </div>
      )}

      <RoutineForm
        initial={initial}
        submitLabel="作成する"
        submittingLabel="作成中…"
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError && !isLimitReached
            ? (mutation.error as Error).message
            : null
        }
        onSubmit={(payload) =>
          mutation.mutate(payload, {
            onSuccess: () => router.replace("/routines"),
          })
        }
        onCancel={() => router.back()}
      />
    </div>
  );
}
