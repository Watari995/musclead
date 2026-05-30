"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useUserId } from "@/lib/auth";

export default function HomePage() {
  const router = useRouter();
  const { userId, ready } = useUserId();

  useEffect(() => {
    if (!ready) return;
    router.replace(userId ? "/meals" : "/register");
  }, [ready, userId, router]);

  return <div className="text-slate-500 text-sm">読み込み中…</div>;
}
