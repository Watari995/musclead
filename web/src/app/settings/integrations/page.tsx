"use client";

import { useState } from "react";
import { useSearchParams } from "next/navigation";
import { Button, Card } from "@/shared/ui";
import { getAccessToken } from "@/shared/auth/access-token";

const BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

async function fetchAuthURL(): Promise<string> {
  const token = getAccessToken();
  const res = await fetch(`${BASE_URL}/integrations/healthplanet/auth`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error("認証URLの取得に失敗しました");
  const data = await res.json();
  return data.url as string;
}

export default function IntegrationsSettingsPage() {
  const searchParams = useSearchParams();
  const connected = searchParams.get("connected") === "true";

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleConnect = async () => {
    setLoading(true);
    setError(null);
    try {
      const url = await fetchAuthURL();
      window.location.href = url;
    } catch {
      setError("連携の開始に失敗しました。もう一度お試しください。");
      setLoading(false);
    }
  };

  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="text-lg font-bold tracking-tight">連携サービス</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          外部サービスと連携して、データを自動で同期します。
        </p>
      </header>

      {connected && (
        <p className="text-sm text-green-600">
          HealthPlanet との連携が完了しました。
        </p>
      )}

      <Card className="p-4">
        <div className="flex items-center justify-between gap-4">
          <div className="space-y-1">
            <p className="text-sm font-semibold">Tanita HealthPlanet</p>
            <p className="text-xs text-[var(--color-ink-muted)]">
              体重・体脂肪率・骨格筋量を自動で同期します。
            </p>
          </div>
          <Button
            type="button"
            variant="ghost"
            onClick={handleConnect}
            disabled={loading}
          >
            {loading ? "移動中…" : connected ? "再連携する" : "連携する"}
          </Button>
        </div>
        {error && (
          <p className="mt-2 text-xs text-red-500">{error}</p>
        )}
      </Card>
    </section>
  );
}
