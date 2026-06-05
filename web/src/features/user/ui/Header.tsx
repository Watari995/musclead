"use client";

import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAccessToken } from "@/shared/auth/access-token";
import { useLogoutMutation, useMeQuery } from "@/features/user/api/user";

export function Header() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const loggedIn = Boolean(token);

  const meQuery = useMeQuery(loggedIn);
  const logout = useLogoutMutation();

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace("/login"),
    });
  };

  return (
    <header className="border-b border-[var(--color-line)] bg-white sticky top-0 z-10">
      <div className="w-full max-w-5xl mx-auto px-5 sm:px-6 h-14 flex items-center justify-between">
        <Link
          href="/"
          className="flex items-center gap-2 font-bold text-lg tracking-tight text-[var(--color-ink)] hover:opacity-80 transition-opacity"
        >
          <Image
            src="/icon.png"
            alt=""
            width={28}
            height={28}
            className="rounded-full"
            priority
          />
          musclead
        </Link>
        {ready && loggedIn && (
          <nav className="hidden sm:flex items-center gap-6 text-sm text-[var(--color-ink)]">
            <Link href="/meals" className="hover:opacity-60 transition-opacity">
              食事
            </Link>
            <Link
              href="/trainings"
              className="hover:opacity-60 transition-opacity"
            >
              トレーニング
            </Link>
            <Link
              href="/exercises"
              className="hover:opacity-60 transition-opacity"
            >
              種目
            </Link>
            <Link
              href="/routines"
              className="hover:opacity-60 transition-opacity"
            >
              ルーティン
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
            {loggedIn ? (
              <>
                {meQuery.data?.name && (
                  <span className="hidden sm:inline text-[var(--color-ink-muted)]">
                    {meQuery.data.name}
                  </span>
                )}
                <button
                  type="button"
                  onClick={handleLogout}
                  disabled={logout.isPending}
                  className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] disabled:opacity-50"
                >
                  ログアウト
                </button>
              </>
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
