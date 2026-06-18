import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "プライバシーポリシー | musclead",
  description: "musclead のプライバシーポリシー",
};

export default function PrivacyPage() {
  return (
    <div className="max-w-2xl mx-auto py-8 space-y-8 text-[var(--color-ink)]">
      <h1 className="text-2xl font-bold">プライバシーポリシー</h1>
      <p className="text-sm text-[var(--color-ink-muted)]">最終更新日: 2026年6月18日</p>

      <section className="space-y-3">
        <p>
          musclead（以下「本サービス」）は、ユーザーの個人情報保護を重要視し、
          適切に取り扱います。本ポリシーでは、収集する情報とその利用目的を説明します。
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">1. 収集する情報</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          <li>メールアドレスおよび氏名（アカウント登録時）</li>
          <li>体重・体脂肪率・筋肉量などの身体測定データ</li>
          <li>トレーニング記録（種目・セット数・重量・回数）</li>
          <li>食事記録（カロリー・栄養素・食事メモ）</li>
          <li>プロフィール画像（任意）</li>
          <li>アプリの利用状況（クラッシュレポート・使用統計）</li>
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">2. 情報の利用目的</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          <li>本サービスの提供・改善・サポート</li>
          <li>トレーニング・食事・体重の記録・分析機能の提供</li>
          <li>アカウント管理および認証</li>
          <li>アプリの安定性向上（Crashlytics によるクラッシュ分析）</li>
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">3. 第三者への提供</h2>
        <p className="text-[var(--color-ink-muted)]">
          収集した個人情報は、法令に基づく場合を除き、第三者へ提供しません。
          本サービスは広告目的での情報共有を行いません。
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">4. 利用する外部サービス</h2>
        <ul className="list-disc list-inside space-y-2 text-[var(--color-ink-muted)]">
          <li>Firebase（クラッシュレポート・利用統計）</li>
        </ul>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">5. データの保管・削除</h2>
        <p className="text-[var(--color-ink-muted)]">
          ユーザーデータはサービス提供に必要な期間のみ保管します。
          アカウント削除をご希望の場合は、サポートページからお問い合わせください。
          削除依頼から7日以内に個人情報を削除します。
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold">6. お問い合わせ</h2>
        <p className="text-[var(--color-ink-muted)]">
          プライバシーに関するご質問は{" "}
          <a
            href="/support"
            className="text-[var(--color-accent)] underline underline-offset-2"
          >
            サポートページ
          </a>{" "}
          からお問い合わせください。
        </p>
      </section>
    </div>
  );
}
