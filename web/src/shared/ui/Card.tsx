import type { ReactNode } from "react";

export function Card({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`bg-white border border-[var(--color-line)] rounded-lg ${className}`}
    >
      {children}
    </div>
  );
}
