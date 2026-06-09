"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { useWeightsQuery } from "@/features/weight/api/weights";
import { RecordWeightForm } from "@/features/weight/ui/RecordWeightForm";
import { WeightGraph } from "@/features/weight/ui/WeightGraph";
import { WeightRow } from "@/features/weight/ui/WeightRow";
import { Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function WeightsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useWeightsQuery(Boolean(token));

  if (!ready || !token) return null;

  return (
    <div className="space-y-8">
      <section>
        <SectionTitle>グラフ</SectionTitle>
        <WeightGraph />
      </section>

      <div className="grid gap-8 lg:grid-cols-[1fr_360px]">
        <section>
          <SectionTitle>体重一覧</SectionTitle>
          {query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
          )}
          {query.isError && (
            <ErrorText>{(query.error as Error).message}</ErrorText>
          )}
          {query.data && query.data.length === 0 && (
            <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
              まだ体重が記録されていません。
            </Card>
          )}
          {query.data && query.data.length > 0 && (
            <ul className="divide-y divide-[var(--color-line)] border border-[var(--color-line)] rounded-lg overflow-hidden bg-[var(--color-surface)]">
              {query.data.map((w) => (
                <WeightRow key={w.id} weight={w} />
              ))}
            </ul>
          )}
        </section>
        <aside>
          <SectionTitle>体重を記録</SectionTitle>
          <RecordWeightForm />
        </aside>
      </div>
    </div>
  );
}
