"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  apiClient,
  type UpsertExerciseRequest,
  type UpsertExerciseResponse,
} from "@/api/client";
import { useAccessToken } from "@/lib/access-token";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/components/ui";

export default function NewExercisePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();
  const [name, setName] = useState("");

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const mutation = useMutation({
    mutationFn: async (body: UpsertExerciseRequest) => {
      const { data, error, response } = await apiClient.POST("/exercises", {
        body,
      });
      if (error) {
        const code = error.error?.code;
        if (code === "training.exercise_name_already_exists_error") {
          throw new Error("同じ名前の種目が既に登録されています。");
        }
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data as UpsertExerciseResponse;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["exercises", "all"] });
      router.replace("/exercises");
    },
  });

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>新しい種目</SectionTitle>
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
              placeholder="例: ベンチプレス"
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
            {mutation.isPending ? "作成中…" : "作成する"}
          </Button>
        </div>
      </form>
    </div>
  );
}
