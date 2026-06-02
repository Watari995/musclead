"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import {
  apiClient,
  type RecordTrainingRequest,
  type TrainingDTO,
} from "@/api/client";
import { useAccessToken } from "@/lib/access-token";
import { fromTrainingDTO } from "@/lib/training-form";
import { ErrorText, SectionTitle } from "@/components/ui";
import { TrainingForm } from "@/components/training/TrainingForm";

export default function EditTrainingPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useQuery({
    queryKey: ["training", params.id],
    enabled: Boolean(token && params.id),
    queryFn: async (): Promise<TrainingDTO> => {
      const { data, error, response } = await apiClient.GET("/trainings/{id}", {
        params: { path: { id: params.id } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as TrainingDTO;
    },
  });

  const initial = useMemo(
    () => (query.data ? fromTrainingDTO(query.data) : null),
    [query.data],
  );

  const mutation = useMutation({
    mutationFn: async (body: RecordTrainingRequest) => {
      const { error, response } = await apiClient.PUT("/trainings/{id}", {
        params: { path: { id: params.id } },
        body,
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["trainings"] });
      queryClient.invalidateQueries({ queryKey: ["training", params.id] });
      router.replace(`/trainings/${params.id}`);
    },
  });

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
        onSubmit={(payload) => mutation.mutate(payload)}
        onCancel={() => router.back()}
      />
    </div>
  );
}
