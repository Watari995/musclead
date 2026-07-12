"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useExerciseQuery,
  useUpdateExerciseMutation,
} from "@/features/training/api/exercises";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/shared/ui";

export default function EditExercisePage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const tc = useTranslations("common");

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useExerciseQuery(params.id, Boolean(token));

  if (!ready || !token) return null;
  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  return (
    <EditExerciseForm
      id={params.id}
      initialName={query.data.name}
      onDone={() => router.replace("/exercises")}
      onCancel={() => router.back()}
    />
  );
}

function EditExerciseForm({
  id,
  initialName,
  onDone,
  onCancel,
}: {
  id: string;
  initialName: string;
  onDone: () => void;
  onCancel: () => void;
}) {
  const t = useTranslations("exercises");
  const tc = useTranslations("common");
  const [name, setName] = useState(initialName);
  const mutation = useUpdateExerciseMutation(id);

  return (
    <div className="space-y-6">
      <SectionTitle>{t("editExercise")}</SectionTitle>
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate({ name }, { onSuccess: onDone });
        }}
      >
        <Card className="p-5 space-y-4">
          <Label label={t("exerciseName")}>
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
            onClick={onCancel}
            disabled={mutation.isPending}
          >
            {tc("cancel")}
          </Button>
          <Button
            type="submit"
            disabled={mutation.isPending || !name.trim()}
            className="flex-1"
          >
            {mutation.isPending ? tc("updating") : tc("update")}
          </Button>
        </div>
      </form>
    </div>
  );
}
