"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useTrainingQuery,
  useUpdateTrainingMutation,
} from "@/features/training/api/trainings";
import { fromTrainingDTO } from "@/features/training/model/training-draft";
import { ErrorText, SectionTitle } from "@/shared/ui";
import { TrainingForm } from "@/features/training/ui/TrainingForm";

export default function EditTrainingPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useTrainingQuery(params.id, Boolean(token));

  const initial = useMemo(
    () => (query.data ? fromTrainingDTO(query.data) : null),
    [query.data],
  );

  const mutation = useUpdateTrainingMutation(params.id);

  if (!ready || !token) return null;
  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!initial) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>トレーニングを編集</SectionTitle>
      <TrainingForm
        initial={initial}
        submitLabel="保存する"
        submittingLabel="保存中…"
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError ? (mutation.error as Error).message : null
        }
        onSubmit={(payload) =>
          mutation.mutate(payload, {
            onSuccess: () => router.replace(`/trainings/${params.id}`),
          })
        }
        onCancel={() => router.back()}
      />
    </div>
  );
}
