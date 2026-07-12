"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import {
  usePortalSessionMutation,
  useSubscribeMutation,
  useSubscriptionQuery,
} from "@/features/purchase/api/purchase";
import { Button } from "@/shared/ui";

// プランページ。 認証ガード / サイドバーは settings/layout.tsx が持つ。
// 現在 layout のサイドバー動線はコメントアウトで非表示 (Stripe 本番設定が test のため)。
// ページ自体は残しており、 サイドバーのリンクを復活させれば再開できる。
//
// 表示は GET /purchase/subscription の is_pro で出し分け:
//   - free → 「Pro にアップグレード」 (POST /subscribe → Checkout)
//   - pro  → 「現在 Pro (期限)」 + 「お支払い・解約の管理」 (POST /portal-session → Customer Portal)
export default function PlanSettingsPage() {
  const subscriptionQuery = useSubscriptionQuery();
  const subscribe = useSubscribeMutation();
  const portal = usePortalSessionMutation();
  const t = useTranslations("plan");
  const tc = useTranslations("common");

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

  const handleManage = () => {
    portal.mutate(undefined, {
      onSuccess: (data) => {
        if (data.portal_url) window.location.href = data.portal_url;
      },
    });
  };

  const isPro = subscriptionQuery.data?.is_pro ?? false;
  const expiresAt = subscriptionQuery.data?.expires_at;

  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="text-lg font-bold tracking-tight">{t("title")}</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          {t("desc")}
        </p>
      </header>

      {returnStatus === "success" && (
        <p className="text-sm text-[var(--color-ink)] rough bg-[var(--color-surface-alt)] px-3 py-2">
          {t("subscribed")}
        </p>
      )}
      {returnStatus === "cancel" && (
        <p className="text-sm text-[var(--color-ink-muted)] rough px-3 py-2">
          {t("cancelled")}
        </p>
      )}

      {subscriptionQuery.isPending ? (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      ) : isPro ? (
        <div className="rough p-4 space-y-3">
          <div className="flex items-baseline justify-between gap-2">
            <span className="font-bold tracking-tight">{t("currentPlan")}</span>
            <span className="text-sm text-[var(--color-ink-muted)]">
              {t("pricePerMonth")}
            </span>
          </div>
          {expiresAt && (
            <p className="text-sm text-[var(--color-ink-muted)]">
              {t("validUntil", { date: new Date(expiresAt).toLocaleDateString("ja-JP") })}
            </p>
          )}
          <Button
            type="button"
            variant="ghost"
            onClick={handleManage}
            disabled={portal.isPending}
          >
            {portal.isPending ? t("redirecting") : t("managePlan")}
          </Button>
          {portal.isError && (
            <p className="text-sm text-[var(--color-accent)]">
              {portal.error instanceof Error
                ? portal.error.message
                : tc("errorOccurred")}
            </p>
          )}
        </div>
      ) : (
        <div className="rough p-4 space-y-3">
          <div className="flex items-baseline justify-between gap-2">
            <span className="font-bold tracking-tight">Pro</span>
            <span className="text-sm text-[var(--color-ink-muted)]">
              {t("pricePerMonth")}
            </span>
          </div>
          <Button
            type="button"
            onClick={handleUpgrade}
            disabled={subscribe.isPending}
          >
            {subscribe.isPending ? t("redirecting") : t("upgradeToPro")}
          </Button>
          {subscribe.isError && (
            <p className="text-sm text-[var(--color-accent)]">
              {subscribe.error instanceof Error
                ? subscribe.error.message
                : tc("errorOccurred")}
            </p>
          )}
        </div>
      )}
    </section>
  );
}
