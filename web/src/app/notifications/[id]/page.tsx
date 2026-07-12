"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { SectionTitle } from "@/shared/ui";
import {
  useNotificationsQuery,
  useReadNotificationMutation,
} from "@/features/notification/api/notifications";

function WeeklyGoalDetail({ metadata, t }: { metadata: Record<string, unknown>; t: (key: string) => string }) {
  const achieved = metadata["achieved"] as boolean | undefined;
  const trainingGoal = metadata["training_goal"] as number | null | undefined;
  const trainingActual = metadata["training_actual"] as number | undefined;
  const calorieGoal = metadata["calorie_goal"] as number | null | undefined;
  const calorieActual = metadata["calorie_actual"] as number | null | undefined;
  const weightGoal = metadata["weight_goal"] as number | null | undefined;
  const weightActual = metadata["weight_actual"] as number | null | undefined;

  return (
    <div className="space-y-4">
      <p className="font-hand text-2xl">
        {achieved ? t("weeklyGoalAchievedBanner") : t("weeklyGoalResults")}
      </p>
      <ul className="space-y-2 text-sm">
        {trainingGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">{t("training")}</span>
            <span>
              {trainingActual ?? 0} / {trainingGoal} {t("timesUnit")}
              {(trainingActual ?? 0) >= trainingGoal ? " ✅" : " ❌"}
            </span>
          </li>
        )}
        {calorieGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">{t("avgCalories")}</span>
            <span>
              {calorieActual != null ? Math.round(calorieActual) : "—"} / {calorieGoal} {t("kcalUnit")}
              {calorieActual != null && calorieActual <= calorieGoal ? " ✅" : " ❌"}
            </span>
          </li>
        )}
        {weightGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">{t("weightChange")}</span>
            <span>
              {weightActual != null ? `${weightActual > 0 ? "+" : ""}${weightActual}` : "—"} / {t("goal")} {weightGoal > 0 ? "+" : ""}{weightGoal} {t("kgUnit")}
            </span>
          </li>
        )}
      </ul>
    </div>
  );
}

export default function NotificationDetailPage() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const { token, ready } = useAccessToken();
  const t = useTranslations("notifications");
  const tc = useTranslations("common");
  const { data } = useNotificationsQuery();
  const readMutation = useReadNotificationMutation();

  const notification = data?.notifications.find((n) => n.id === params.id);

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  useEffect(() => {
    if (notification && !notification.is_read) {
      readMutation.mutate(notification.id);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [notification?.id, notification?.is_read]);

  if (!ready || !token) return null;

  if (!data) {
    return <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>;
  }

  if (!notification) {
    return <p className="text-sm text-red-500">{t("notFound")}</p>;
  }

  return (
    <div className="space-y-6">
      <SectionTitle>{t("notificationDetail")}</SectionTitle>
      <div className="rough p-4 space-y-3">
        {notification.notification_type === "weekly_goal" && (
          <WeeklyGoalDetail metadata={notification.metadata} t={t} />
        )}
        <p className="text-xs text-[var(--color-ink-muted)]">
          {new Date(notification.created_at).toLocaleString("ja-JP")}
        </p>
      </div>
      <button
        type="button"
        onClick={() => router.back()}
        className="text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
      >
        {t("backToList")}
      </button>
    </div>
  );
}
