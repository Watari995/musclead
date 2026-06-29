"use client";

import { ThemePicker } from "@/features/user/ui/ThemePicker";
import { CalendarColorPicker } from "@/features/user/ui/CalendarColorPicker";

export default function AppearanceSettingsPage() {
  return (
    <div className="space-y-8">
      <section className="space-y-4">
        <header className="space-y-1">
          <h2 className="text-lg font-bold tracking-tight">外観</h2>
          <p className="text-sm text-[var(--color-ink-muted)]">
            テーマモードを選択してください。 選んだ瞬間に保存されます。
          </p>
        </header>
        <ThemePicker />
      </section>
      <section className="space-y-4">
        <header className="space-y-1">
          <h2 className="text-lg font-bold tracking-tight">カレンダーの色</h2>
          <p className="text-sm text-[var(--color-ink-muted)]">
            カレンダーに表示する点の色を設定します。選んだ瞬間に保存されます。
          </p>
        </header>
        <CalendarColorPicker />
      </section>
    </div>
  );
}
