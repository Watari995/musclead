import type { ReactNode } from "react";

export function Label({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <label className="block">
      <span className="block text-xs font-medium text-[var(--color-ink-muted)] mb-1.5">
        {label}
      </span>
      {children}
    </label>
  );
}
