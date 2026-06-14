"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { ThemePicker } from "@/features/user/ui/ThemePicker";
import { useAccessToken } from "@/shared/auth/access-token";
import { SettingsDetailHeader } from "../_components/SettingsDetailHeader";

export default function AppearanceSettingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SettingsDetailHeader title="外観" />
      <p className="text-sm text-[var(--color-ink-muted)]">
        テーマモードを選択してください。 選んだ瞬間に保存されます。
      </p>
      <ThemePicker />
    </div>
  );
}
