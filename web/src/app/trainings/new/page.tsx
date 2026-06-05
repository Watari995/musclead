"use client";

import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useRecordTrainingMutation } from "@/features/training/api/trainings";
import { createInitialTraining } from "@/features/training/model/training-draft";
import { SectionTitle } from "@/shared/ui";
import { TrainingForm } from "@/features/training/ui/TrainingForm";

export default function NewTrainingPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const initial = useMemo(() => createInitialTraining(), []);
  const mutation = useRecordTrainingMutation();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>トレーニングを記録</SectionTitle>
      <TrainingForm
        initial={initial}
        submitLabel="記録する"
        submittingLabel="記録中…"
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError ? (mutation.error as Error).message : null
        }
        onSubmit={(payload) =>
          mutation.mutate(payload, {
            onSuccess: () => router.replace("/trainings"),
          })
        }
        onCancel={() => router.back()}
      />
    </div>
  );
}
