"use client";

import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useCreateRoutineMutation } from "@/features/training/api/routines";
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

  return (
    <div className="space-y-6">
      <SectionTitle>新しいルーティン</SectionTitle>
      <RoutineForm
        initial={initial}
        submitLabel="作成する"
        submittingLabel="作成中…"
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError ? (mutation.error as Error).message : null
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
