"use client";

import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { useLogoutMutation } from "@/features/user/api/user";
import { Button } from "@/shared/ui";

// 認証ガード / サイドバーは settings/layout.tsx が持つ。 ここは内容のみ。
export default function AccountSettingsPage() {
  const router = useRouter();
  const logout = useLogoutMutation();
  const t = useTranslations("account");

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace("/login"),
    });
  };

  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="font-hand text-2xl">{t("title")}</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          {t("signOutDesc")}
        </p>
      </header>
      <Button
        type="button"
        variant="ghost"
        onClick={handleLogout}
        disabled={logout.isPending}
      >
        {logout.isPending ? t("signingOut") : t("signOut")}
      </Button>
    </section>
  );
}
