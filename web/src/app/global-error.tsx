"use client";

import * as Sentry from "@sentry/nextjs";
import { useEffect } from "react";
import { useTranslations } from "next-intl";

export default function GlobalError({
  error,
  unstable_retry,
}: {
  error: Error & { digest?: string };
  unstable_retry: () => void;
}) {
  const t = useTranslations("globalError");

  useEffect(() => {
    Sentry.captureException(error);
  }, [error]);

  return (
    <html lang="ja">
      <body
        style={{
          fontFamily:
            "system-ui, -apple-system, 'Hiragino Sans', sans-serif",
          margin: 0,
          padding: "2rem",
          backgroundColor: "#fff",
          color: "#111",
        }}
      >
        <div style={{ maxWidth: "32rem", margin: "0 auto" }}>
          <h1 style={{ fontSize: "1.25rem", fontWeight: 700 }}>
            {t("title")}
          </h1>
          <p style={{ fontSize: "0.875rem", marginTop: "1rem" }}>
            {t("message")}
          </p>
          <p
            style={{
              fontSize: "0.75rem",
              color: "#c00",
              marginTop: "1rem",
              wordBreak: "break-word",
            }}
          >
            {error.message}
          </p>
          <button
            type="button"
            onClick={() => unstable_retry()}
            style={{
              marginTop: "1.5rem",
              padding: "0.5rem 1.25rem",
              fontSize: "0.875rem",
              borderRadius: "0.375rem",
              backgroundColor: "#111",
              color: "#fff",
              border: "none",
              cursor: "pointer",
            }}
          >
            {t("retry")}
          </button>
        </div>
      </body>
    </html>
  );
}
