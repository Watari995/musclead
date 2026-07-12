import type { ReactNode } from "react";

export function Card({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div className={`rough bg-[var(--color-surface)] ${className}`}>
      {children}
    </div>
  );
}
