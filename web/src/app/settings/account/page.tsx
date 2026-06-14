"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useLogoutMutation } from "@/features/user/api/user";
import { useAccessToken } from "@/shared/auth/access-token";
import { Button } from "@/shared/ui";
import { SettingsDetailHeader } from "../_components/SettingsDetailHeader";

export default function AccountSettingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const logout = useLogoutMutation();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  if (!ready || !token) return null;

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace("/login"),
    });
  };

  return (
    <div className="space-y-6">
      <SettingsDetailHeader title="アカウント" />
      <p className="text-sm text-[var(--color-ink-muted)]">
        現在のデバイスからサインアウトします。
      </p>
      <Button
        type="button"
        variant="ghost"
        onClick={handleLogout}
        disabled={logout.isPending}
      >
        {logout.isPending ? "ログアウト中…" : "ログアウト"}
      </Button>
    </div>
  );
}
