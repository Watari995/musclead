"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { setStoredUserId } from "@/lib/auth";
import { Button, Card, ErrorText, Label, TextInput } from "@/components/ui";

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
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold tracking-tight mb-2">ログイン</h1>
      <p className="text-sm text-[var(--color-ink-muted)] mb-6">
        現状は UserID(UUID)を入力してログインします。 認証は今後実装予定。
      </p>
      <Card className="p-6">
        <form className="space-y-4" onSubmit={handleSubmit}>
          <Label label="UserID">
            <TextInput
              type="text"
              value={userId}
              onChange={(e) => {
                setUserId(e.target.value);
                setError(null);
              }}
              placeholder="00000000-0000-0000-0000-000000000000"
              className="font-mono text-sm"
              required
            />
          </Label>
          {error && <ErrorText>{error}</ErrorText>}
          <Button type="submit" fullWidth>
            ログイン
          </Button>
        </form>
      </Card>
      <p className="mt-6 text-sm text-[var(--color-ink-muted)] text-center">
        アカウントをお持ちでないですか?{" "}
        <Link
          href="/register"
          className="text-[var(--color-ink)] font-medium hover:opacity-60"
        >
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
