"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  apiClient,
  type ExerciseDTO,
  type UpsertExerciseRequest,
} from "@/api/client";
import { useAccessToken } from "@/lib/access-token";
import { EXERCISES_QUERY_KEY } from "@/lib/queries/exercises";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/components/ui";

export default function EditExercisePage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();
  const [name, setName] = useState("");
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useQuery({
    queryKey: ["exercise", params.id],
    enabled: Boolean(token && params.id),
    queryFn: async (): Promise<ExerciseDTO> => {
      const { data, error, response } = await apiClient.GET("/exercises/{id}", {
        params: { path: { id: params.id } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as ExerciseDTO;
    },
  });

  useEffect(() => {
    if (query.data && !initialized) {
      setName(query.data.name ?? "");
      setInitialized(true);
    }
  }, [query.data, initialized]);

  const mutation = useMutation({
    mutationFn: async (body: UpsertExerciseRequest) => {
      const { error, response } = await apiClient.PUT("/exercises/{id}", {
        params: { path: { id: params.id } },
        body,
      });
      if (error) {
        const code = error.error?.code;
        if (code === "training.exercise_name_already_exists_error") {
          throw new Error("同じ名前の種目が既に登録されています。");
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: EXERCISES_QUERY_KEY });
      queryClient.invalidateQueries({ queryKey: ["exercise", params.id] });
      router.replace("/exercises");
    },
  });

  if (!ready || !token) return null;

  if (query.isLoading || !initialized) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }

  return (
    <div className="space-y-6">
      <SectionTitle>種目を編集</SectionTitle>
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate({ name });
        }}
      >
        <Card className="p-5 space-y-4">
          <Label label="種目名">
            <TextInput
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              maxLength={50}
              disabled={mutation.isPending}
            />
          </Label>
        </Card>
        {mutation.isError && (
          <ErrorText>{(mutation.error as Error).message}</ErrorText>
        )}
        <div className="flex gap-3">
          <Button
            type="button"
            variant="ghost"
            onClick={() => router.back()}
            disabled={mutation.isPending}
          >
            キャンセル
          </Button>
          <Button
            type="submit"
            disabled={mutation.isPending || !name.trim()}
            className="flex-1"
          >
            {mutation.isPending ? "更新中…" : "更新する"}
          </Button>
        </div>
      </form>
    </div>
  );
}
