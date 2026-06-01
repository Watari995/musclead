"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/lib/access-token";

export default function HomePage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (!ready) return;
    router.replace(token ? "/meals" : "/login");
  }, [ready, token, router]);

  return <div className="text-[var(--color-ink-muted)] text-sm">読み込み中…</div>;
}
