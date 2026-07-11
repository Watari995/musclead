"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { SectionTitle } from "@/shared/ui";
import { useNotificationsQuery } from "@/features/notification/api/notifications";
import type { NotificationDTO } from "@/features/notification/api/notifications";

function formatDate(iso: string) {
  return new Date(iso).toLocaleString("ja-JP", {
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function notificationLabel(n: NotificationDTO, t: (key: string) => string): string {
  if (n.notification_type === "weekly_goal") {
    const achieved = n.metadata["achieved"] as boolean | undefined;
    return achieved ? t("weeklyGoalAchieved") : t("checkWeeklyGoal");
  }
  return t("title");
}

export default function NotificationsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const t = useTranslations("notifications");
  const tc = useTranslations("common");
  const { data, isLoading, isError } = useNotificationsQuery();

  if (!ready || !token) return null;

  return (
    <div className="space-y-6">
      <SectionTitle>{t("title")}</SectionTitle>

      {isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      )}
      {isError && (
        <p className="text-sm text-red-500">{t("loadFailed")}</p>
      )}

      {data && data.notifications.length === 0 && (
        <p className="text-sm text-[var(--color-ink-muted)]">{t("noNotifications")}</p>
      )}

      {data && data.notifications.length > 0 && (
        <ul className="divide-y divide-[var(--color-line)] border border-[var(--color-line)] rounded-xl overflow-hidden">
          {data.notifications.map((n) => (
            <li key={n.id}>
              <Link
                href={`/notifications/${n.id}`}
                className={`flex items-start gap-3 px-4 py-3 hover:bg-[var(--color-surface-alt)] transition-colors ${
                  !n.is_read ? "bg-[var(--color-surface-alt)]" : ""
                }`}
              >
                <span className="mt-1 shrink-0">
                  {n.notification_type === "weekly_goal" ? "🏆" : "🔔"}
                </span>
                <div className="flex-1 min-w-0">
                  <p className={`text-sm ${!n.is_read ? "font-bold" : ""} text-[var(--color-ink)]`}>
                    {notificationLabel(n, t)}
                  </p>
                  <p className="text-xs text-[var(--color-ink-muted)] mt-0.5">
                    {formatDate(n.created_at)}
                  </p>
                </div>
                {!n.is_read && (
                  <span className="mt-1 w-2 h-2 rounded-full bg-red-500 shrink-0" />
                )}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
