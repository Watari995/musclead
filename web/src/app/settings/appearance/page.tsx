"use client";

import { ThemePicker } from "@/features/user/ui/ThemePicker";

// 認証ガード / サイドバーは settings/layout.tsx が持つ。 ここは内容のみ。
export default function AppearanceSettingsPage() {
  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="text-lg font-bold tracking-tight">外観</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          テーマモードを選択してください。 選んだ瞬間に保存されます。
        </p>
      </header>
      <ThemePicker />
    </section>
  );
}
