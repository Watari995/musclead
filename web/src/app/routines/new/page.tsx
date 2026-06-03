"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useMemo } from "react";
import {
  apiClient,
  type UpsertRoutineRequest,
  type UpsertRoutineResponse,
} from "@/api/client";
import { useAccessToken } from "@/lib/access-token";
import { createInitialRoutine } from "@/lib/routine-form";
import { SectionTitle } from "@/components/ui";
import { RoutineForm } from "@/components/routine/RoutineForm";

export default function NewRoutinePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const initial = useMemo(() => createInitialRoutine(), []);

  const mutation = useMutation({
    mutationFn: async (body: UpsertRoutineRequest) => {
      const { data, error, response } = await apiClient.POST("/routines", {
        body,
      });
      if (error) {
        const code = error.error?.code;
        if (code === "training.routine_name_already_exists_error") {
          throw new Error("同じ名前のルーティンが既に登録されています。");
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as UpsertRoutineResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["routines"] });
      router.replace("/routines");
    },
  });

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
        onSubmit={(payload) => mutation.mutate(payload)}
        onCancel={() => router.back()}
      />
    </div>
  );
}
