import type { Metadata } from "next";
import { getTranslations } from "next-intl/server";

export const metadata: Metadata = {
  title: "プライバシーポリシー | musclead",
  description: "musclead のプライバシーポリシー",
};

export default async function PrivacyPage() {
  const t = await getTranslations("privacy");

  return (
    <div className="max-w-2xl mx-auto py-8 space-y-8 text-[var(--color-ink)]">
      <h1 className="text-2xl font-bold">{t("title")}</h1>
      <p className="text-sm text-[var(--color-ink-muted)]">{t("lastUpdated")}</p>

      <section className="space-y-3">
        <p>{t("intro")}</p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section1Title")}</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          {(t.raw("section1Items") as string[]).map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section2Title")}</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          {(t.raw("section2Items") as string[]).map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section3Title")}</h2>
        <p className="text-[var(--color-ink-muted)]">{t("section3Text")}</p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section4Title")}</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          {(t.raw("section4Items") as string[]).map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section5Title")}</h2>
        <p className="text-[var(--color-ink-muted)]">{t("section5Text")}</p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">{t("section6Title")}</h2>
        <p className="text-[var(--color-ink-muted)]">
          {t("section6Text")}{" "}
          <a
            href="/support"
            className="text-[var(--color-accent)] underline underline-offset-2"
          >
            {t("section6Link")}
          </a>{" "}
          {t("section6TextAfter")}
        </p>
      </section>
    </div>
  );
}
