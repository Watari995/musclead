"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAccessToken } from "@/shared/auth/access-token";
import { SectionTitle } from "@/shared/ui";
import {
  useNotificationsQuery,
  useReadNotificationMutation,
} from "@/features/notification/api/notifications";

function WeeklyGoalDetail({ metadata }: { metadata: Record<string, unknown> }) {
  const achieved = metadata["achieved"] as boolean | undefined;
  const trainingGoal = metadata["training_goal"] as number | null | undefined;
  const trainingActual = metadata["training_actual"] as number | undefined;
  const calorieGoal = metadata["calorie_goal"] as number | null | undefined;
  const calorieActual = metadata["calorie_actual"] as number | null | undefined;
  const weightGoal = metadata["weight_goal"] as number | null | undefined;
  const weightActual = metadata["weight_actual"] as number | null | undefined;

  return (
    <div className="space-y-4">
      <p className="text-lg font-bold">
        {achieved ? "今週の目標を達成しました！ 🎉" : "今週の目標結果"}
      </p>
      <ul className="space-y-2 text-sm">
        {trainingGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">トレーニング</span>
            <span>
              {trainingActual ?? 0} / {trainingGoal} 回
              {(trainingActual ?? 0) >= trainingGoal ? " ✅" : " ❌"}
            </span>
          </li>
        )}
        {calorieGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">平均カロリー</span>
            <span>
              {calorieActual != null ? Math.round(calorieActual) : "—"} / {calorieGoal} kcal
              {calorieActual != null && calorieActual <= calorieGoal ? " ✅" : " ❌"}
            </span>
          </li>
        )}
        {weightGoal != null && (
          <li className="flex justify-between border-b border-[var(--color-line)] pb-2">
            <span className="text-[var(--color-ink-muted)]">体重変化</span>
            <span>
              {weightActual != null ? `${weightActual > 0 ? "+" : ""}${weightActual}` : "—"} / 目標 {weightGoal > 0 ? "+" : ""}{weightGoal} kg
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
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }

  if (!notification) {
    return <p className="text-sm text-red-500">通知が見つかりません</p>;
  }

  return (
    <div className="space-y-6">
      <SectionTitle>通知詳細</SectionTitle>
      <div className="border border-[var(--color-line)] rounded-xl p-4 space-y-3">
        {notification.notification_type === "weekly_goal" && (
          <WeeklyGoalDetail metadata={notification.metadata} />
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
        ← 通知一覧に戻る
      </button>
    </div>
  );
}
