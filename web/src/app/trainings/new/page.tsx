"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { useRecordTrainingMutation } from "@/features/training/api/trainings";
import { useNewTrainingDraft } from "@/features/training/model/useNewTrainingDraft";
import { Button, Card, SectionTitle } from "@/shared/ui";
import { TrainingForm } from "@/features/training/ui/TrainingForm";

export default function NewTrainingPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const t = useTranslations("trainings");
  const tc = useTranslations("common");
  const { draft, setDraft, restorable, restore, discard, clear } =
    useNewTrainingDraft();
  const mutation = useRecordTrainingMutation();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>{t("recordTraining")}</SectionTitle>

      {restorable && (
        <Card className="p-4 sm:p-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-sm text-[var(--color-ink)]">
            {t("restoreData")}
          </p>
          <div className="flex shrink-0 gap-2">
            <Button type="button" variant="ghost" onClick={discard}>
              {t("discard")}
            </Button>
            <Button type="button" onClick={restore}>
              {t("restore")}
            </Button>
          </div>
        </Card>
      )}

      <TrainingForm
        value={draft}
        onChange={setDraft}
        submitLabel={tc("record")}
        submittingLabel={tc("recording")}
        submitting={mutation.isPending}
        errorMessage={
          mutation.isError ? (mutation.error as Error).message : null
        }
        onSubmit={(payload) =>
          mutation.mutate(payload, {
            onSuccess: () => {
              clear();
              router.replace("/trainings");
            },
          })
        }
        onCancel={() => router.back()}
      />
    </div>
  );
}
