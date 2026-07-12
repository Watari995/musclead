"use client";

import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function Error({
  error,
  unstable_retry,
}: {
  error: Error & { digest?: string };
  unstable_retry: () => void;
}) {
  const t = useTranslations("error");
  const tc = useTranslations("common");

  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="max-w-md mx-auto space-y-4">
      <SectionTitle>{t("title")}</SectionTitle>
      <Card className="p-6 space-y-4">
        <p className="text-sm text-[var(--color-ink)]">{t("message")}</p>
        <ErrorText>{error.message}</ErrorText>
        <div className="flex gap-3">
          <Button type="button" variant="ghost" onClick={() => unstable_retry()}>
            {tc("retry")}
          </Button>
          <Button
            type="button"
            onClick={() => window.location.reload()}
            className="flex-1"
          >
            {tc("reloadPage")}
          </Button>
        </div>
      </Card>
    </div>
  );
}
