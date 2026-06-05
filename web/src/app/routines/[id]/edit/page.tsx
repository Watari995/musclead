"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import {
  apiClient,
  type RoutineDTO,
  type UpsertRoutineRequest,
} from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  createInitialRoutine,
  fromRoutineDTO,
  type RoutineDraft,
} from "@/lib/routine-form";
import { ErrorText, SectionTitle } from "@/shared/ui";
import { RoutineForm } from "@/components/routine/RoutineForm";

export default function EditRoutinePage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();
  const [initial, setInitial] = useState<RoutineDraft | null>(null);

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useQuery({
    queryKey: ["routine", params.id],
    enabled: Boolean(token && params.id),
    queryFn: async (): Promise<RoutineDTO> => {
      const { data, error, response } = await apiClient.GET(
        "/routines/{id}",
        { params: { path: { id: params.id } } },
      );
      if (error)
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as RoutineDTO;
    },
  });

  // 初期データが取れたら 1 度だけ initial draft を確定
  useEffect(() => {
    if (query.data && !initial) {
      setInitial(fromRoutineDTO(query.data));
    }
  }, [query.data, initial]);

  const fallbackInitial = useMemo(() => createInitialRoutine(), []);

  const mutation = useMutation({
    mutationFn: async (body: UpsertRoutineRequest) => {
      const { error, response } = await apiClient.PUT("/routines/{id}", {
        params: { path: { id: params.id } },
        body,
      });
      if (error) {
        const code = error.error?.code;
        if (code === "training.routine_name_already_exists_error") {
          throw new Error("同じ名前のルーティンが既に登録されています。");
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["routines"] });
      queryClient.invalidateQueries({ queryKey: ["routine", params.id] });
      router.replace("/routines");
    },
  });

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
        initial={initial ?? fallbackInitial}
        submitLabel="更新する"
        submittingLabel="更新中…"
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
