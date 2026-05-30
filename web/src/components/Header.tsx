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
    <header className="border-b border-slate-200 bg-white">
      <div className="w-full max-w-3xl mx-auto px-4 py-3 flex items-center justify-between">
        <Link href="/" className="font-bold text-lg">
          💪 musclead
        </Link>
        {ready && (
          <nav className="flex items-center gap-4 text-sm">
            {userId ? (
              <>
                <Link href="/meals" className="hover:underline">
                  食事
                </Link>
                <span className="text-slate-500 font-mono text-xs">
                  {userId.slice(0, 8)}…
                </span>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="text-slate-600 hover:text-slate-900"
                >
                  ログアウト
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="hover:underline">
                  ログイン
                </Link>
                <Link
                  href="/register"
                  className="rounded bg-slate-900 text-white px-3 py-1 hover:bg-slate-700"
                >
                  新規登録
                </Link>
              </>
            )}
          </nav>
        )}
      </div>
    </header>
  );
}
