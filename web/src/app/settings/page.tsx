"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

// /settings はデフォルト項目 (外観) にリダイレクトする。
// サイドバー / 認証ガードは settings/layout.tsx が持つ。
export default function SettingsPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace("/settings/appearance");
  }, [router]);

  return null;
}
