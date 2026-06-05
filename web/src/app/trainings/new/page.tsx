"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import {
  apiClient,
  type RecordTrainingRequest,
  type RecordTrainingResponse,
} from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import { createInitialTraining } from "@/lib/training-form";
import { SectionTitle } from "@/shared/ui";
import { TrainingForm } from "@/components/training/TrainingForm";

export default function NewTrainingPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const initial = useMemo(() => createInitialTraining(), []);

  const mutation = useMutation({
    mutationFn: async (body: RecordTrainingRequest) => {
      const { data, error, response } = await apiClient.POST("/trainings", {
        body,
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as RecordTrainingResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["trainings"] });
      router.replace("/trainings");
    },
  });

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
        onSubmit={(payload) => mutation.mutate(payload)}
        onCancel={() => router.back()}
      />
    </div>
  );
}
