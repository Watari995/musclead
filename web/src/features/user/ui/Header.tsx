"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useAccessToken } from "@/shared/auth/access-token";
import { useMeQuery } from "@/features/user/api/user";
import { Avatar } from "@/features/user/ui/Avatar";
import { NotificationBell } from "@/features/notification/ui/NotificationBell";

const NAV_ITEMS = [
  { href: "/meals", key: "meals" },
  { href: "/trainings", key: "trainings" },
  { href: "/exercises", key: "exercises" },
  { href: "/routines", key: "routines" },
  { href: "/weights", key: "weights" },
  { href: "/settings", key: "settings" },
] as const;

export function Header() {
  const t = useTranslations("nav");
  const pathname = usePathname();
  const { token, ready } = useAccessToken();
  const loggedIn = Boolean(token);

  const meQuery = useMeQuery(loggedIn);
  const [menuOpen, setMenuOpen] = useState(false);

  // ルート遷移で自動的にメニューを閉じる (React 19: prop 変化で派生状態をリセットする公式パターン)
  const [lastPathname, setLastPathname] = useState(pathname);
  if (lastPathname !== pathname) {
    setLastPathname(pathname);
    if (menuOpen) setMenuOpen(false);
  }

  // メニュー展開中は背景スクロールをロック
  useEffect(() => {
    if (!menuOpen) return;
    const original = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = original;
    };
  }, [menuOpen]);

  return (
    <header className="border-b border-[var(--color-line)] bg-[var(--color-surface)] sticky top-0 z-30">
      <div className="w-full max-w-5xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between gap-3">
        <Link
          href="/"
          className="flex items-center gap-2 font-bold text-lg tracking-tight text-[var(--color-ink)] hover:opacity-80 transition-opacity min-w-0"
        >
          <Image
            src="/icon.png"
            alt=""
            width={28}
            height={28}
            className="rounded-full shrink-0"
            priority
          />
          <span className="truncate">musclead</span>
        </Link>

        {ready && loggedIn && (
          <nav className="hidden sm:flex items-center gap-6 text-sm text-[var(--color-ink)]">
            {NAV_ITEMS.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="hover:opacity-60 transition-opacity"
              >
                {t(item.key)}
              </Link>
            ))}
          </nav>
        )}

        {ready && (
          <div className="flex items-center gap-3 text-sm">
            {loggedIn ? (
              <>
                <NotificationBell />
                {meQuery.data?.name && (
                  <Link
                    href="/profile"
                    className="hidden sm:flex items-center gap-2 text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] max-w-[14rem]"
                  >
                    {meQuery.data.profile_image_url && (
                      <Avatar
                        src={meQuery.data.profile_image_url}
                        alt={meQuery.data.name}
                        size="w-7 h-7"
                      />
                    )}
                    <span className="truncate">{meQuery.data.name}</span>
                  </Link>
                )}
                {/* mobile: ハンバーガーの左にアバターのみ表示(name はドロワー内) */}
                {meQuery.data?.profile_image_url && (
                  <Link
                    href="/profile"
                    aria-label={t("profile")}
                    className="sm:hidden inline-flex items-center"
                  >
                    <Avatar
                      src={meQuery.data.profile_image_url}
                      alt={meQuery.data.name ?? ""}
                      size="w-8 h-8"
                    />
                  </Link>
                )}
                <button
                  type="button"
                  onClick={() => setMenuOpen(true)}
                  aria-label={t("openMenu")}
                  aria-expanded={menuOpen}
                  aria-controls="mobile-nav"
                  className="sm:hidden inline-flex items-center justify-center w-10 h-10 -mr-2 rounded-md text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)]"
                >
                  <HamburgerIcon />
                </button>
              </>
            ) : (
              <>
                <Link
                  href="/login"
                  className="text-[var(--color-ink)] hover:opacity-60"
                >
                  {t("login")}
                </Link>
                <Link
                  href="/register"
                  className="bg-[var(--color-ink)] text-[var(--color-surface)] px-4 h-9 inline-flex items-center rounded-md text-sm font-medium hover:opacity-90 whitespace-nowrap"
                >
                  {t("register")}
                </Link>
              </>
            )}
          </div>
        )}
      </div>

      {ready && loggedIn && (
        <MobileMenu
          open={menuOpen}
          userName={meQuery.data?.name ?? ""}
          pathname={pathname}
          onClose={() => setMenuOpen(false)}
        />
      )}
    </header>
  );
}

function MobileMenu({
  open,
  userName,
  pathname,
  onClose,
}: {
  open: boolean;
  userName: string;
  pathname: string;
  onClose: () => void;
}) {
  const t = useTranslations("nav");

  return (
    <div
      id="mobile-nav"
      role="dialog"
      aria-modal="true"
      aria-hidden={!open}
      className={`sm:hidden fixed inset-0 z-40 ${open ? "" : "pointer-events-none"}`}
    >
      <button
        type="button"
        aria-label={t("closeMenu")}
        onClick={onClose}
        className={`absolute inset-0 bg-black/40 transition-opacity ${
          open ? "opacity-100" : "opacity-0"
        }`}
      />
      <div
        className={`absolute top-0 right-0 h-full w-72 max-w-[85vw] bg-[var(--color-surface)] shadow-xl flex flex-col transition-transform duration-200 ${
          open ? "translate-x-0" : "translate-x-full"
        }`}
      >
        <div className="h-14 px-5 flex items-center justify-between border-b border-[var(--color-line)]">
          {userName ? (
            <Link
              href="/profile"
              onClick={onClose}
              className="text-sm font-bold tracking-tight truncate hover:opacity-60"
            >
              {userName}
            </Link>
          ) : (
            <span className="text-sm font-bold tracking-tight truncate">
              {t("menu")}
            </span>
          )}
          <button
            type="button"
            onClick={onClose}
            aria-label={t("closeMenu")}
            className="inline-flex items-center justify-center w-10 h-10 -mr-2 rounded-md text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)]"
          >
            <CloseIcon />
          </button>
        </div>
        <nav className="flex-1 overflow-y-auto py-2">
          {NAV_ITEMS.map((item) => {
            const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={onClose}
                className={`block px-5 py-3 text-base border-l-2 ${
                  active
                    ? "border-[var(--color-ink)] font-bold bg-[var(--color-surface-alt)]"
                    : "border-transparent text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)]"
                }`}
              >
                {t(item.key)}
              </Link>
            );
          })}
        </nav>
      </div>
    </div>
  );
}

function HamburgerIcon() {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      aria-hidden="true"
    >
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  );
}

function CloseIcon() {
  return (
    <svg
      width="22"
      height="22"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      aria-hidden="true"
    >
      <line x1="6" y1="6" x2="18" y2="18" />
      <line x1="18" y1="6" x2="6" y2="18" />
    </svg>
  );
}
