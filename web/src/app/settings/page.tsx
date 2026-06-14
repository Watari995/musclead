"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useSubscribeMutation } from "@/features/purchase/api/purchase";
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
          <PlanSection />
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
            href="#plan"
            className="inline-block whitespace-nowrap px-3 py-2 rounded-md hover:bg-[var(--color-surface-alt)] text-[var(--color-ink)] text-sm font-medium border-l-2 border-transparent md:border-l-2"
          >
            プラン
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

// PlanSection は Pro 申込み導線。
// 「Pro にアップグレード」 → POST /purchase/subscribe → 返却 checkout_url へ遷移 (Stripe Checkout)。
// Checkout からの戻りは success_url / cancel_url に付与された ?purchase=success|cancel で判定し、
// 簡易メッセージを表示する (専用ページは作らない)。
function PlanSection() {
  const subscribe = useSubscribeMutation();
  // Stripe Checkout からの戻りクエリ (?purchase=success|cancel) を読む。
  // useSearchParams は Suspense 境界が要るため、 既存流儀に倣い window から直接読む。
  // lazy 初期化で読み (effect 内 setState を避ける)、 URL クリーンアップだけ effect で行う。
  const [returnStatus] = useState<"success" | "cancel" | null>(() => {
    if (typeof window === "undefined") return null;
    const status = new URLSearchParams(window.location.search).get("purchase");
    return status === "success" || status === "cancel" ? status : null;
  });

  useEffect(() => {
    if (returnStatus === null || typeof window === "undefined") return;
    // クエリを URL から消してリロード時の再表示を防ぐ
    const url = new URL(window.location.href);
    url.searchParams.delete("purchase");
    window.history.replaceState({}, "", url.toString());
  }, [returnStatus]);

  const handleUpgrade = () => {
    subscribe.mutate(
      { plan: "pro" },
      {
        onSuccess: (data) => {
          if (data.checkout_url) window.location.href = data.checkout_url;
        },
      },
    );
  };

  return (
    <section id="plan" aria-labelledby="plan-title" className="space-y-4">
      <header className="space-y-1">
        <h2 id="plan-title" className="text-lg font-bold tracking-tight">
          プラン
        </h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          Pro にアップグレードすると、 すべての機能が利用できます。
        </p>
      </header>

      {returnStatus === "success" && (
        <p className="text-sm text-[var(--color-ink)] rounded-md border border-[var(--color-line)] bg-[var(--color-surface-alt)] px-3 py-2">
          お申し込みありがとうございます。 Pro が有効になりました。
        </p>
      )}
      {returnStatus === "cancel" && (
        <p className="text-sm text-[var(--color-ink-muted)] rounded-md border border-[var(--color-line)] px-3 py-2">
          お申し込みはキャンセルされました。
        </p>
      )}

      <div className="rounded-md border border-[var(--color-line)] p-4 space-y-3">
        <div className="flex items-baseline justify-between gap-2">
          <span className="font-bold tracking-tight">Pro</span>
          <span className="text-sm text-[var(--color-ink-muted)]">
            ¥480 / 月
          </span>
        </div>
        <Button
          type="button"
          onClick={handleUpgrade}
          disabled={subscribe.isPending}
        >
          {subscribe.isPending ? "リダイレクト中…" : "Pro にアップグレード"}
        </Button>
        {subscribe.isError && (
          <p className="text-sm text-[var(--color-accent)]">
            {subscribe.error instanceof Error
              ? subscribe.error.message
              : "エラーが発生しました。 時間をおいて再度お試しください。"}
          </p>
        )}
      </div>
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
