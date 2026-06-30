"use client";

import Link from "next/link";
import { useNotificationsQuery } from "../api/notifications";

export function NotificationBell() {
  const { data } = useNotificationsQuery();
  const unread = data?.unread_count ?? 0;

  return (
    <Link
      href="/notifications"
      aria-label={`通知${unread > 0 ? `（未読${unread}件）` : ""}`}
      className="relative inline-flex items-center justify-center w-9 h-9 rounded-md text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)] transition-colors"
    >
      <BellIcon />
      {unread > 0 && (
        <span className="absolute top-1 right-1 min-w-[16px] h-4 px-1 rounded-full bg-red-500 text-white text-[10px] font-bold leading-4 text-center">
          {unread > 99 ? "99+" : unread}
        </span>
      )}
    </Link>
  );
}

function BellIcon() {
  return (
    <svg
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
      <path d="M13.73 21a2 2 0 0 1-3.46 0" />
    </svg>
  );
}
