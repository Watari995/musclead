"use client";

import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { ThemePicker } from "@/features/user/ui/ThemePicker";
import { CalendarColorPicker } from "@/features/user/ui/CalendarColorPicker";

export default function AppearanceSettingsPage() {
  const t = useTranslations("appearance");
  const router = useRouter();

  const setLocale = (locale: string) => {
    document.cookie = `NEXT_LOCALE=${locale}; max-age=${60 * 60 * 24 * 365}; path=/`;
    router.refresh();
  };

  return (
    <div className="space-y-8">
      <section className="space-y-4">
        <header className="space-y-1">
          <h2 className="text-lg font-bold tracking-tight">{t("title")}</h2>
          <p className="text-sm text-[var(--color-ink-muted)]">
            {t("themeDesc")}
          </p>
        </header>
        <ThemePicker />
      </section>
      <section className="space-y-4">
        <header className="space-y-1">
          <h2 className="text-lg font-bold tracking-tight">{t("calendarColor")}</h2>
          <p className="text-sm text-[var(--color-ink-muted)]">
            {t("calendarColorDesc")}
          </p>
        </header>
        <CalendarColorPicker />
      </section>
      <section className="space-y-4">
        <header className="space-y-1">
          <h2 className="text-lg font-bold tracking-tight">{t("language")}</h2>
          <p className="text-sm text-[var(--color-ink-muted)]">
            {t("languageDesc")}
          </p>
        </header>
        <div className="flex gap-3 flex-wrap">
          <button
            type="button"
            onClick={() => setLocale("ja")}
            className="px-4 py-2 rough text-sm font-medium hover:bg-[var(--color-surface-alt)] transition-colors"
          >
            {t("japanese")}
          </button>
          <button
            type="button"
            onClick={() => setLocale("en")}
            className="px-4 py-2 rough text-sm font-medium hover:bg-[var(--color-surface-alt)] transition-colors"
          >
            {t("english")}
          </button>
        </div>
      </section>
    </div>
  );
}
