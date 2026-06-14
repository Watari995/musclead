"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { SectionTitle } from "@/shared/ui";

// 設定はメニュー(ハブ) + 項目ごとのページ構成。
// 各項目は /settings/<item> の独立ページ。 ここはその一覧。
//
// プランは Stripe の本番設定が整うまで動線を非表示にする。
// 再開時は下のコメントアウトを外すだけでよい (ページ自体は /settings/plan に存在)。
const ITEMS = [
  { href: "/settings/appearance", label: "外観", desc: "テーマモードの選択" },
  // { href: "/settings/plan", label: "プラン", desc: "Pro へのアップグレード" },
  { href: "/settings/account", label: "アカウント", desc: "サインアウト" },
];

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

      <ul className="divide-y divide-[var(--color-line)] rounded-md border border-[var(--color-line)]">
        {ITEMS.map((item) => (
          <li key={item.href}>
            <Link
              href={item.href}
              className="flex items-center justify-between gap-3 px-4 py-3 hover:bg-[var(--color-surface-alt)]"
            >
              <span className="min-w-0">
                <span className="block font-medium tracking-tight">
                  {item.label}
                </span>
                <span className="block text-sm text-[var(--color-ink-muted)]">
                  {item.desc}
                </span>
              </span>
              <span aria-hidden className="text-[var(--color-ink-muted)]">
                ›
              </span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
