import type { ReactNode } from "react";

export function SectionTitle({ children }: { children: ReactNode }) {
  return (
    <h2 className="text-base font-bold tracking-tight text-[var(--color-ink)] mb-3">
      {children}
    </h2>
  );
}
