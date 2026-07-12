"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";

export default function HomePage() {
  const router = useRouter();
  const tc = useTranslations("common");
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (!ready) return;
    router.replace(token ? "/calendar" : "/login");
  }, [ready, token, router]);

  return <div className="text-[var(--color-ink-muted)] text-sm">{tc("loading")}</div>;
}
