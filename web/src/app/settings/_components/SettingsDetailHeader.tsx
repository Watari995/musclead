"use client";

import Link from "next/link";
import { SectionTitle } from "@/shared/ui";

// 設定サブページ共通のヘッダー。 「← 設定」 戻りリンク + タイトル。
// _components はアンダースコア始まりのため route にはならない (App Router)。
export function SettingsDetailHeader({ title }: { title: string }) {
  return (
    <div className="space-y-2">
      <Link
        href="/settings"
        className="inline-block text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
      >
        ← 設定
      </Link>
      <SectionTitle>{title}</SectionTitle>
    </div>
  );
}
