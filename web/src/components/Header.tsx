"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { logoutRequest } from "@/api/auth";
import { apiClient, type UserDTO } from "@/api/client";
import { clearAccessToken, useAccessToken } from "@/lib/access-token";

const ME_QUERY_KEY = ["me"] as const;

export function Header() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { token, ready } = useAccessToken();

  const meQuery = useQuery({
    queryKey: ME_QUERY_KEY,
    enabled: Boolean(token),
    queryFn: async () => {
      const { data, error, response } = await apiClient.GET("/users/me");
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data as UserDTO;
    },
  });

  const handleLogout = async () => {
    await logoutRequest();
    clearAccessToken();
    queryClient.removeQueries({ queryKey: ME_QUERY_KEY });
    router.replace("/login");
  };

  const loggedIn = Boolean(token);

  return (
    <header className="border-b border-[var(--color-line)] bg-white sticky top-0 z-10">
      <div className="w-full max-w-5xl mx-auto px-5 sm:px-6 h-14 flex items-center justify-between">
        <Link
          href="/"
          className="font-bold text-lg tracking-tight text-[var(--color-ink)]"
        >
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
                  className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
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
