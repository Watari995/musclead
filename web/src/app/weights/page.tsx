"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
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

  const t = useTranslations("weights");
  const tc = useTranslations("common");
  const query = useWeightsQuery(Boolean(token));

  if (!ready || !token) return null;

  return (
    <div className="space-y-8">
      <section>
        <SectionTitle>{t("graph")}</SectionTitle>
        <WeightGraph />
      </section>

      <div className="grid gap-8 lg:grid-cols-[1fr_360px]">
        <section>
          <SectionTitle>{t("list")}</SectionTitle>
          {query.isLoading && (
            <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
          )}
          {query.isError && (
            <ErrorText>{(query.error as Error).message}</ErrorText>
          )}
          {query.data && query.data.length === 0 && (
            <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
              {t("noWeights")}
            </Card>
          )}
          {query.data && query.data.length > 0 && (
            <ul className="divide-y divide-[var(--color-line)] rough overflow-hidden bg-[var(--color-surface)]">
              {query.data.map((w) => (
                <WeightRow key={w.id} weight={w} />
              ))}
            </ul>
          )}
        </section>
        <aside>
          <SectionTitle>{t("record")}</SectionTitle>
          <RecordWeightForm />
        </aside>
      </div>
    </div>
  );
}
