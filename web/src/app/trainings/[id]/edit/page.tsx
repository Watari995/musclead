"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useTrainingQuery,
  useUpdateTrainingMutation,
} from "@/features/training/api/trainings";
import {
  type TrainingDraft,
  fromTrainingDTO,
} from "@/features/training/model/training-draft";
import { type TrainingDTO } from "@/shared/api/client";
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

  if (!ready || !token) return null;
  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  // データ確定後に form をマウントし、 draft を一度だけ遅延初期化する
  // (effect での setState を避ける)。
  return <EditTrainingForm id={params.id} dto={query.data} />;
}

function EditTrainingForm({ id, dto }: { id: string; dto: TrainingDTO }) {
  const router = useRouter();
  const [draft, setDraft] = useState<TrainingDraft>(() => fromTrainingDTO(dto));
  const mutation = useUpdateTrainingMutation(id);

  return (
    <div className="space-y-6">
      <SectionTitle>トレーニングを編集</SectionTitle>
      <TrainingForm
        value={draft}
        onChange={setDraft}
        submitLabel="保存する"
        submittingLabel="保存中…"
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError ? (mutation.error as Error).message : null
        }
        onSubmit={(payload) =>
          mutation.mutate(payload, {
            onSuccess: () => router.replace(`/trainings/${id}`),
          })
        }
        onCancel={() => router.back()}
      />
    </div>
  );
}
