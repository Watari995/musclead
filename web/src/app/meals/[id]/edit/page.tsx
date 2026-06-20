"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useFindMealQuery } from "@/features/meal/api/meals";
import { EditMealForm } from "@/features/meal/ui/EditMealForm";
import { ErrorText, SectionTitle } from "@/shared/ui";
import type { Meal } from "@/features/meal/model/meal";

export default function MealEditPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useFindMealQuery(params.id, Boolean(token));

  if (!ready || !token) return null;
  if (query.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }
  if (query.isError) {
    return <ErrorText>{(query.error as Error).message}</ErrorText>;
  }
  if (!query.data) return null;

  return <MealEditContent meal={query.data} />;
}

function MealEditContent({ meal }: { meal: Meal }) {
  const router = useRouter();

  return (
    <div className="max-w-lg mx-auto space-y-4">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => router.back()}
          className="text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] transition-colors"
        >
          ← 戻る
        </button>
        <SectionTitle>食事を編集</SectionTitle>
      </div>
      <EditMealForm
        meal={meal}
        onSuccess={() => router.push("/meals")}
        onCancel={() => router.back()}
      />
    </div>
  );
}
