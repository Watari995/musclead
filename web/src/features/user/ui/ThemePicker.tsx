"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import {
  useUpdatePreferencesMutation,
  type Theme,
} from "@/features/user/api/user";
import { ErrorText } from "@/shared/ui";

const THEME_VALUES: Theme[] = ["light", "dark", "system"];

// 各 theme の固定パレット(プレビュー描画時に inline で適用、 現在の theme に
// 左右されないようにするため)
type Palette = {
  surface: string;
  surfaceAlt: string;
  line: string;
  ink: string;
  inkMuted: string;
};

const LIGHT_PALETTE: Palette = {
  surface: "#ffffff",
  surfaceAlt: "#f6f6f6",
  line: "#e5e7eb",
  ink: "#111111",
  inkMuted: "#6b7280",
};

const DARK_PALETTE: Palette = {
  surface: "#0d1117",
  surfaceAlt: "#161b22",
  line: "#30363d",
  ink: "#c9d1d9",
  inkMuted: "#8b949e",
};

export function ThemePicker() {
  const t = useTranslations("appearance");
  const { theme, setTheme } = useTheme();
  const mutation = useUpdatePreferencesMutation();

  // SSR では theme が undefined なので、 mount 後に hydrate 完了状態を出す。
  // next-themes 公式の推奨パターン。
  const [mounted, setMounted] = useState(false);
  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => setMounted(true), []);

  const current: Theme | undefined = mounted
    ? ((theme as Theme | undefined) ?? "system")
    : undefined;

  const change = (next: Theme) => {
    if (next === current) return;
    setTheme(next); // 楽観的 UI: next-themes 即時反映
    mutation.mutate({ theme: next });
    // 失敗時はサーバ値が refetch で同期されるので、 そこで自動 rollback
  };

  const options = THEME_VALUES.map((value) => ({
    value,
    label: t(value),
    description: t(`${value}Desc` as "lightDesc" | "darkDesc" | "systemDesc"),
  }));

  return (
    <div className="space-y-3">
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        {options.map((opt) => (
          <ThemeCard
            key={opt.value}
            value={opt.value}
            label={opt.label}
            description={opt.description}
            selected={current === opt.value}
            disabled={mutation.isPending}
            onClick={() => change(opt.value)}
          />
        ))}
      </div>
      {mutation.isError && (
        <ErrorText>{(mutation.error as Error).message}</ErrorText>
      )}
    </div>
  );
}

function ThemeCard({
  value,
  label,
  description,
  selected,
  disabled,
  onClick,
}: {
  value: Theme;
  label: string;
  description: string;
  selected: boolean;
  disabled: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      aria-pressed={selected}
      className={`text-left rounded-lg border-2 transition-colors overflow-hidden ${
        selected
          ? "border-[var(--color-ink)]"
          : "border-[var(--color-line)] hover:border-[var(--color-ink-muted)]"
      } disabled:opacity-60 disabled:cursor-not-allowed`}
    >
      <ThemePreview value={value} />
      <div className="px-3 py-2.5 bg-[var(--color-surface-alt)] border-t border-[var(--color-line)]">
        <div className="flex items-center gap-2">
          <RadioIndicator selected={selected} />
          <span className="text-sm font-medium">{label}</span>
        </div>
        <p className="text-xs text-[var(--color-ink-muted)] mt-1 pl-6">
          {description}
        </p>
      </div>
    </button>
  );
}

// 各 theme のミニ画面プレビュー。 system は light/dark の左右分割で表現。
function ThemePreview({ value }: { value: Theme }) {
  if (value === "system") {
    return (
      <div className="grid grid-cols-2 h-28">
        <PreviewPane palette={LIGHT_PALETTE} />
        <PreviewPane palette={DARK_PALETTE} />
      </div>
    );
  }
  const palette = value === "dark" ? DARK_PALETTE : LIGHT_PALETTE;
  return (
    <div className="h-28">
      <PreviewPane palette={palette} />
    </div>
  );
}

function PreviewPane({ palette }: { palette: Palette }) {
  return (
    <div
      className="h-full px-3 pt-3 overflow-hidden"
      style={{ background: palette.surface }}
    >
      {/* ヘッダ風 */}
      <div
        className="h-3 w-12 rounded-sm mb-2"
        style={{ background: palette.ink }}
      />
      {/* テキスト行 */}
      <div
        className="h-2 w-full rounded-sm mb-1.5"
        style={{ background: palette.inkMuted, opacity: 0.4 }}
      />
      <div
        className="h-2 w-3/4 rounded-sm mb-3"
        style={{ background: palette.inkMuted, opacity: 0.4 }}
      />
      {/* カード風 */}
      <div
        className="rounded border h-9"
        style={{
          background: palette.surfaceAlt,
          borderColor: palette.line,
        }}
      />
    </div>
  );
}

function RadioIndicator({ selected }: { selected: boolean }) {
  return (
    <span
      aria-hidden
      className={`inline-block w-4 h-4 rounded-full border-2 transition-colors ${
        selected
          ? "border-[var(--color-ink)] bg-[var(--color-ink)]"
          : "border-[var(--color-line)] bg-transparent"
      }`}
    >
      {selected && (
        <span
          className="block w-1.5 h-1.5 rounded-full mx-auto mt-[3px]"
          style={{ background: "var(--color-surface)" }}
        />
      )}
    </span>
  );
}
