"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
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

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useRoutineQuery(params.id, Boolean(token));

  if (!ready || !token) return null;
  if (query.isLoading) {
    return (
      <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
    );
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  return (
    <EditRoutineFormShell
      id={params.id}
      initial={fromRoutineDTO(query.data)}
      onDone={() => router.replace("/routines")}
      onCancel={() => router.back()}
    />
  );
}

// RoutineForm itself initializes its useState with `initial` only once on
// mount, so background refetches do not clobber in-progress edits as long as
// this shell is not remounted.
function EditRoutineFormShell({
  id,
  initial,
  onDone,
  onCancel,
}: {
  id: string;
  initial: RoutineDraft;
  onDone: () => void;
  onCancel: () => void;
}) {
  const mutation = useUpdateRoutineMutation(id);

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
          mutation.mutate(payload, { onSuccess: onDone })
        }
        onCancel={onCancel}
      />
    </div>
  );
}
