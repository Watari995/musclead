"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useSubscribeMutation } from "@/features/purchase/api/purchase";
import { useAccessToken } from "@/shared/auth/access-token";
import { Button } from "@/shared/ui";
import { SettingsDetailHeader } from "../_components/SettingsDetailHeader";

// プランページ (Pro 申込み)。
// 現在 /settings ハブの動線はコメントアウトで非表示 (Stripe 本番設定が整うまで)。
// ページ自体は残しており、 ハブのリンクを復活させれば再開できる。
export default function PlanSettingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const subscribe = useSubscribeMutation();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

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
    const url = new URL(window.location.href);
    url.searchParams.delete("purchase");
    window.history.replaceState({}, "", url.toString());
  }, [returnStatus]);

  if (!ready || !token) return null;

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
    <div className="space-y-6">
      <SettingsDetailHeader title="プラン" />
      <p className="text-sm text-[var(--color-ink-muted)]">
        Pro にアップグレードすると、 すべての機能が利用できます。
      </p>

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
    </div>
  );
}
