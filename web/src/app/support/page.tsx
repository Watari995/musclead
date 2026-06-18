import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "サポート | musclead",
  description: "musclead のサポートページ",
};

export default function SupportPage() {
  return (
    <div className="max-w-2xl mx-auto py-8 space-y-8 text-[var(--color-ink)]">
      <h1 className="text-2xl font-bold">サポート</h1>

      <section className="space-y-3">
        <p className="text-[var(--color-ink-muted)]">
          musclead をご利用いただきありがとうございます。
          ご質問・ご要望・不具合報告は以下よりお問い合わせください。
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">お問い合わせ</h2>
        <div className="rounded-xl border border-[var(--color-line)] bg-[var(--color-surface-alt)] p-6 space-y-3">
          <p className="text-sm text-[var(--color-ink-muted)]">メールでのお問い合わせ</p>
          <a
            href="mailto:support@musclead.app"
            className="text-[var(--color-accent)] font-medium text-lg break-all"
          >
            support@musclead.app
          </a>
          <p className="text-xs text-[var(--color-ink-muted)]">
            通常2〜3営業日以内に返信いたします。
          </p>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold">よくある質問</h2>
        <div className="space-y-4">
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">アカウントを削除したい</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              上記メールアドレスに「アカウント削除希望」とご連絡ください。
              7日以内にすべての個人情報を削除いたします。
            </p>
          </div>
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">パスワードを忘れた</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              現在パスワードリセット機能を実装中です。お手数ですがサポートメールにご連絡ください。
            </p>
          </div>
          <div className="border-b border-[var(--color-line)] pb-4">
            <p className="font-medium mb-1">アプリがクラッシュする・正しく動作しない</p>
            <p className="text-sm text-[var(--color-ink-muted)]">
              発生状況（操作手順・端末・OSバージョン）をメールでお知らせください。
            </p>
          </div>
        </div>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">関連リンク</h2>
        <ul className="space-y-2">
          <li>
            <a
              href="/privacy"
              className="text-[var(--color-accent)] underline underline-offset-2 text-sm"
            >
              プライバシーポリシー
            </a>
          </li>
        </ul>
      </section>
    </div>
  );
}
