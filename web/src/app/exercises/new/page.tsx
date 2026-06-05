"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useCreateExerciseMutation } from "@/features/training/api/exercises";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/shared/ui";

export default function NewExercisePage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const [name, setName] = useState("");

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const mutation = useCreateExerciseMutation();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>新しい種目</SectionTitle>
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate(
            { name },
            { onSuccess: () => router.replace("/exercises") },
          );
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
