"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useRoutineQuery,
  useUpdateRoutineMutation,
} from "@/features/training/api/routines";
import {
  fromRoutineDTO,
  type RoutineDraft,
} from "@/features/training/model/routine-draft";
import { ErrorText, SectionTitle } from "@/shared/ui";
import { RoutineForm } from "@/features/training/ui/RoutineForm";

export default function EditRoutinePage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const [initial, setInitial] = useState<RoutineDraft | null>(null);

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useRoutineQuery(params.id, Boolean(token));

  // 初期データが取れたら 1 度だけ initial draft を確定。
  // フォーム編集中に query が再フェッチされても初期値を上書きしないため。
  useEffect(() => {
    if (query.data && !initial) {
      setInitial(fromRoutineDTO(query.data));
    }
  }, [query.data, initial]);

  const mutation = useUpdateRoutineMutation(params.id);

  if (!ready || !token) return null;

  if (query.isLoading || !initial) {
    return (
      <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
    );
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }

  return (
    <div className="space-y-6">
      <SectionTitle>ルーティンを編集</SectionTitle>
      <RoutineForm
        initial={initial}
        submitLabel="更新する"
        submittingLabel="更新中…"
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
