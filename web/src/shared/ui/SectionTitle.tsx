import type { ReactNode } from "react";

export function SectionTitle({ children }: { children: ReactNode }) {
  return (
    <h2 className="font-hand text-2xl text-[var(--color-ink)] mb-3">
      {children}
    </h2>
  );
}
