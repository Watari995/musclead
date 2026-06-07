"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useLogoutMutation } from "@/features/user/api/user";
import { ThemePicker } from "@/features/user/ui/ThemePicker";
import { useAccessToken } from "@/shared/auth/access-token";
import { Button, SectionTitle } from "@/shared/ui";

export default function SettingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>設定</SectionTitle>

      <div className="grid grid-cols-1 md:grid-cols-[200px_1fr] gap-6 md:gap-8">
        <SettingsSidebar />
        <div className="space-y-10 min-w-0">
          <AppearanceSection />
          <AccountSection />
        </div>
      </div>
    </div>
  );
}

// 将来 通知 など増えた時にも 同 sidebar に並べる前提。
// mobile ではタブ風の横スクロール、 desktop は縦並びの GitHub 風。
function SettingsSidebar() {
  return (
    <nav className="md:sticky md:top-20 self-start">
      <ul className="flex md:flex-col gap-1 overflow-x-auto md:overflow-visible -mx-4 px-4 md:mx-0 md:px-0">
        <li>
          <a
            href="#appearance"
            className="inline-block whitespace-nowrap px-3 py-2 rounded-md hover:bg-[var(--color-surface-alt)] text-[var(--color-ink)] text-sm font-medium border-l-2 border-transparent md:border-l-2"
          >
            外観
          </a>
        </li>
        <li>
          <a
            href="#account"
            className="inline-block whitespace-nowrap px-3 py-2 rounded-md hover:bg-[var(--color-surface-alt)] text-[var(--color-ink)] text-sm font-medium border-l-2 border-transparent md:border-l-2"
          >
            アカウント
          </a>
        </li>
      </ul>
    </nav>
  );
}

function AppearanceSection() {
  return (
    <section
      id="appearance"
      aria-labelledby="appearance-title"
      className="space-y-4"
    >
      <header className="space-y-1">
        <h2
          id="appearance-title"
          className="text-lg font-bold tracking-tight"
        >
          外観
        </h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          テーマモードを選択してください。 選んだ瞬間に保存されます。
        </p>
      </header>
      <ThemePicker />
    </section>
  );
}

function AccountSection() {
  const router = useRouter();
  const logout = useLogoutMutation();
  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace("/login"),
    });
  };
  return (
    <section
      id="account"
      aria-labelledby="account-title"
      className="space-y-4"
    >
      <header className="space-y-1">
        <h2 id="account-title" className="text-lg font-bold tracking-tight">
          アカウント
        </h2>
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
