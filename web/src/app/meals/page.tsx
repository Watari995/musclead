"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { useMealsQuery } from "@/features/meal/api/meals";
import { MealRow } from "@/features/meal/ui/MealRow";
import { Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function MealsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const t = useTranslations("meals");
  const tc = useTranslations("common");
  const query = useMealsQuery(Boolean(token));

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>{t("title")}</SectionTitle>
        <Link
          href="/meals/record"
          className="h-9 px-4 rounded-md bg-[var(--color-ink)] text-[var(--color-surface)] text-sm font-medium inline-flex items-center hover:opacity-90 transition-opacity"
        >
          {tc("record")}
        </Link>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {query.data && query.data.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          {t("noMeals")}
        </Card>
      )}
      {query.data && query.data.length > 0 && (
        <ul className="divide-y divide-[var(--color-line)] rough overflow-hidden bg-[var(--color-surface)]">
          {query.data.map((m) => (
            <MealRow key={m.id} meal={m} />
          ))}
        </ul>
      )}
    </div>
  );
}
