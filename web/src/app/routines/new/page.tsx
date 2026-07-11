"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useTranslations } from "next-intl";
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

  const t = useTranslations("routines");
  const tc = useTranslations("common");
  const initial = useMemo(() => createInitialRoutine(), []);
  const mutation = useCreateRoutineMutation();

  if (!ready || !token) return null;

  // 無料プランの上限到達は通常エラーではなく Pro 誘導として表示する。
  const isLimitReached = mutation.error instanceof RoutineLimitReachedError;

  return (
    <div className="space-y-6">
      <SectionTitle>{t("newRoutinePage")}</SectionTitle>

      {isLimitReached && (
        <div className="rounded-md border border-[var(--color-line)] bg-[var(--color-surface-alt)] p-4 space-y-2">
          <p className="text-sm text-[var(--color-ink)]">
            {t("proLimitReached")}
          </p>
          <Link
            href="/settings/plan"
            className="inline-block text-sm font-medium text-[var(--color-ink)] underline"
          >
            {t("viewPlan")}
          </Link>
        </div>
      )}

      <RoutineForm
        initial={initial}
        submitLabel={tc("create")}
        submittingLabel={tc("creating")}
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
