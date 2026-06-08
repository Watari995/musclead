"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { RecordWeightForm } from "@/features/weight/ui/RecordWeightForm";
import { SectionTitle } from "@/shared/ui";

export default function WeightsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  if (!ready || !token) return null;

  return (
    <div className="max-w-md mx-auto">
      <SectionTitle>体重を記録</SectionTitle>
      <RecordWeightForm />
    </div>
  );
}
