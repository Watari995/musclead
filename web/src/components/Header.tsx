"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { clearStoredUserId, useUserId } from "@/lib/auth";

export function Header() {
  const router = useRouter();
  const { userId, ready } = useUserId();

  const handleLogout = () => {
    clearStoredUserId();
    router.push("/login");
  };

  return (
    <header className="border-b border-[var(--color-line)] bg-white sticky top-0 z-10">
      <div className="w-full max-w-5xl mx-auto px-5 sm:px-6 h-14 flex items-center justify-between">
        <Link
          href="/"
          className="font-bold text-lg tracking-tight text-[var(--color-ink)]"
        >
          musclead
        </Link>
        {ready && userId && (
          <nav className="hidden sm:flex items-center gap-6 text-sm text-[var(--color-ink)]">
            <Link href="/meals" className="hover:opacity-60 transition-opacity">
              食事
            </Link>
            <Link
              href="/meals"
              className="hover:opacity-60 transition-opacity text-[var(--color-ink-muted)]"
            >
              トレーニング
            </Link>
            <Link
              href="/meals"
              className="hover:opacity-60 transition-opacity text-[var(--color-ink-muted)]"
            >
              体重
            </Link>
          </nav>
        )}
        {ready && (
          <div className="flex items-center gap-3 text-sm">
            {userId ? (
              <button
                type="button"
                onClick={handleLogout}
                className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
              >
                ログアウト
              </button>
            ) : (
              <>
                <Link
                  href="/login"
                  className="text-[var(--color-ink)] hover:opacity-60"
                >
                  ログイン
                </Link>
                <Link
                  href="/register"
                  className="bg-[var(--color-ink)] text-white px-4 h-9 inline-flex items-center rounded-md text-sm font-medium hover:opacity-90"
                >
                  新規登録
                </Link>
              </>
            )}
          </div>
        )}
      </div>
    </header>
  );
}
