"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useMealsQuery } from "@/features/meal/api/meals";
import { MealRow } from "@/features/meal/ui/MealRow";
import { RecordMealForm } from "@/features/meal/ui/RecordMealForm";
import { Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function MealsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useMealsQuery(Boolean(token));

  if (!ready || !token) return null;

  return (
    <div className="grid gap-8 lg:grid-cols-[1fr_360px]">
      <section>
        <SectionTitle>食事一覧</SectionTitle>
        {query.isLoading && (
          <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
        )}
        {query.isError && (
          <ErrorText>{(query.error as Error).message}</ErrorText>
        )}
        {query.data && query.data.length === 0 && (
          <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
            まだ食事が記録されていません。
          </Card>
        )}
        {query.data && query.data.length > 0 && (
          <ul className="divide-y divide-[var(--color-line)] border border-[var(--color-line)] rounded-lg overflow-hidden bg-white">
            {query.data.map((m) => (
              <MealRow key={m.id} meal={m} />
            ))}
          </ul>
        )}
      </section>
      <aside>
        <SectionTitle>食事を記録</SectionTitle>
        <RecordMealForm />
      </aside>
    </div>
  );
}
