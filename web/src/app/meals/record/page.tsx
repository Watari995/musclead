"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { RecordMealForm } from "@/features/meal/ui/RecordMealForm";
import { SectionTitle } from "@/shared/ui";

export default function MealRecordPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const t = useTranslations("meals");
  const tc = useTranslations("common");

  if (!ready || !token) return null;

  return (
    <div className="max-w-lg mx-auto space-y-4">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => router.back()}
          className="text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] transition-colors"
        >
          {tc("back")}
        </button>
        <SectionTitle>{t("recordMeal")}</SectionTitle>
      </div>
      <RecordMealForm onSuccess={() => router.push("/meals")} />
    </div>
  );
}
