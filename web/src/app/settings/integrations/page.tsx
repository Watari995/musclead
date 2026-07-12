"use client";

import { useState } from "react";
import { useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
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
  const t = useTranslations("integrations");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleConnect = async () => {
    setLoading(true);
    setError(null);
    try {
      const url = await fetchAuthURL();
      window.location.href = url;
    } catch {
      setError(t("connectFailed"));
      setLoading(false);
    }
  };

  return (
    <section className="space-y-4">
      <header className="space-y-1">
        <h2 className="font-hand text-2xl">{t("title")}</h2>
        <p className="text-sm text-[var(--color-ink-muted)]">
          {t("desc")}
        </p>
      </header>

      {connected && (
        <p className="text-sm text-green-600">
          {t("healthPlanetConnected")}
        </p>
      )}

      <Card className="p-4">
        <div className="flex items-center justify-between gap-4">
          <div className="space-y-1">
            <p className="text-sm font-semibold">Tanita HealthPlanet</p>
            <p className="text-xs text-[var(--color-ink-muted)]">
              {t("healthPlanetDesc")}
            </p>
          </div>
          <Button
            type="button"
            variant="ghost"
            onClick={handleConnect}
            disabled={loading}
          >
            {loading ? t("moving") : connected ? t("reconnect") : t("connect")}
          </Button>
        </div>
        {error && (
          <p className="mt-2 text-xs text-red-500">{error}</p>
        )}
      </Card>
    </section>
  );
}
