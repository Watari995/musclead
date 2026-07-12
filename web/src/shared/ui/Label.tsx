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
      <span className="font-hand block text-base text-[var(--color-ink)] mb-1.5">
        {label}
      </span>
      {children}
    </label>
  );
}
