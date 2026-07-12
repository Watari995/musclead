import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

export const metadata: Metadata = {
  title: "サポート | musclead",
  description: "musclead のサポートページ",
};

export default async function SupportPage() {
  const t = await getTranslations("support");

  return (
    <div className="max-w-2xl mx-auto py-8 space-y-8 text-[var(--color-ink)]">
      <h1 className="text-2xl font-bold">{t("title")}</h1>

      <section className="space-y-3">
        <p className="text-[var(--color-ink-muted)]">
          {t("intro")}
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">{t("contactTitle")}</h2>
        <div className="rounded-xl border border-[var(--color-line)] bg-[var(--color-surface-alt)] p-6 space-y-3">
          <p className="text-sm text-[var(--color-ink-muted)]">{t("contactDesc")}</p>
          <a
            href="mailto:support@musclead.app"
            className="text-[var(--color-accent)] font-medium text-lg break-all"
          >
            support@musclead.app
          </a>
          <p className="text-xs text-[var(--color-ink-muted)]">
            {t("contactReply")}
          </p>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">{t("faqTitle")}</h2>
        <div className="space-y-4">
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">{t("faq1Q")}</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              {t("faq1A")}
            </p>
          </div>
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">{t("faq2Q")}</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              {t("faq2A")}
            </p>
          </div>
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">{t("faq3Q")}</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              {t("faq3A")}
            </p>
          </div>
        </div>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("relatedTitle")}</h2>
        <ul className="space-y-2">
          <li>
            <a
              href="/privacy"
              className="text-[var(--color-accent)] underline underline-offset-2 text-sm"
            >
              {t("privacyPolicy")}
            </a>
          </li>
        </ul>
      </section>
    </div>
  );
}
