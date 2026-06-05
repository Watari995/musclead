"use client";

import { useEffect } from "react";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function Error({
  error,
  unstable_retry,
}: {
  error: Error & { digest?: string };
  unstable_retry: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="max-w-md mx-auto space-y-4">
      <SectionTitle>エラーが発生しました</SectionTitle>
      <Card className="p-6 space-y-4">
        <p className="text-sm text-[var(--color-ink)]">
          画面の表示中に問題が発生しました。
          一時的な不具合の可能性があるため、 「再試行」 を押すか、
          ページを再読み込みしてください。
        </p>
        <ErrorText>{error.message}</ErrorText>
        <div className="flex gap-3">
          <Button type="button" variant="ghost" onClick={() => unstable_retry()}>
            再試行
          </Button>
          <Button
            type="button"
            onClick={() => window.location.reload()}
            className="flex-1"
          >
            ページを再読み込み
          </Button>
        </div>
      </Card>
    </div>
  );
}
