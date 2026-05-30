"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { setStoredUserId } from "@/lib/auth";

export default function LoginPage() {
  const router = useRouter();
  const [userId, setUserId] = useState("");
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = userId.trim();
    if (!isUUID(trimmed)) {
      setError("UserID は UUID 形式で入力してください");
      return;
    }
    setStoredUserId(trimmed);
    router.replace("/meals");
  };

  return (
    <div className="max-w-md mx-auto bg-white rounded-lg shadow-sm border border-slate-200 p-6">
      <h1 className="text-xl font-bold mb-4">ログイン</h1>
      <p className="text-sm text-slate-600 mb-4">
        現状は UserID(UUID)を入力してログインします。 認証は今後実装予定。
      </p>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="block">
          <span className="text-sm font-medium text-slate-700">UserID</span>
          <input
            type="text"
            value={userId}
            onChange={(e) => {
              setUserId(e.target.value);
              setError(null);
            }}
            placeholder="00000000-0000-0000-0000-000000000000"
            className="mt-1 block w-full rounded border border-slate-300 px-3 py-2 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-slate-400"
            required
          />
        </label>
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          className="w-full rounded bg-slate-900 text-white py-2 hover:bg-slate-700"
        >
          ログイン
        </button>
      </form>
      <p className="mt-4 text-sm text-slate-600">
        アカウントをお持ちでないですか?{" "}
        <Link href="/register" className="text-blue-600 hover:underline">
          新規登録
        </Link>
      </p>
    </div>
  );
}

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
function isUUID(v: string) {
  return UUID_RE.test(v);
}
