"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, type ReactNode } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { SectionTitle } from "@/shared/ui";

// 設定は GitHub 風: 左サイドバーは全ページ共通で残り、 選んだ項目の内容が右に出る。
// 各項目は /settings/<item> の独立ルート。 サイドバーの active は現在の pathname で判定。
//
// プランは Stripe の本番設定が整うまで動線を非表示にする。
// 再開時は下のコメントアウトを外すだけ (ページ自体は /settings/plan に存在)。
const NAV = [
  { href: "/settings/appearance", label: "外観" },
  // { href: "/settings/plan", label: "プラン" },
  { href: "/settings/account", label: "アカウント" },
];

export default function SettingsLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>設定</SectionTitle>

      <div className="grid grid-cols-1 md:grid-cols-[200px_1fr] gap-6 md:gap-8">
        <nav className="md:sticky md:top-20 self-start">
          <ul className="flex md:flex-col gap-1 overflow-x-auto md:overflow-visible -mx-4 px-4 md:mx-0 md:px-0">
            {NAV.map((item) => {
              const active = pathname === item.href;
              return (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    aria-current={active ? "page" : undefined}
                    className={`inline-block whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium border-l-2 ${
                      active
                        ? "bg-[var(--color-surface-alt)] border-[var(--color-ink)] text-[var(--color-ink)]"
                        : "border-transparent text-[var(--color-ink-muted)] hover:bg-[var(--color-surface-alt)] hover:text-[var(--color-ink)]"
                    }`}
                  >
                    {item.label}
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>

        <div className="space-y-6 min-w-0">{children}</div>
      </div>
    </div>
  );
}
