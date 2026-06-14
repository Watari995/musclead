"use client";

import { useRouter } from "next/navigation";
import { useLogoutMutation } from "@/features/user/api/user";
import { Button } from "@/shared/ui";

// 認証ガード / サイドバーは settings/layout.tsx が持つ。 ここは内容のみ。
export default function AccountSettingsPage() {
  const router = useRouter();
  const logout = useLogoutMutation();

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace("/login"),
    });
  };

  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="text-lg font-bold tracking-tight">アカウント</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          現在のデバイスからサインアウトします。
        </p>
      </header>
      <Button
        type="button"
        variant="ghost"
        onClick={handleLogout}
        disabled={logout.isPending}
      >
        {logout.isPending ? "ログアウト中…" : "ログアウト"}
      </Button>
    </section>
  );
}
