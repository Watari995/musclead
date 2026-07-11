"use client";

import { useTranslations } from "next-intl";

export default function Loading() {
  const t = useTranslations("common");
  return (
    <p className="text-sm text-[var(--color-ink-muted)]">{t("loading")}</p>
  );
}
