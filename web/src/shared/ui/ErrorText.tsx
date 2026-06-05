import type { ReactNode } from "react";

export function ErrorText({ children }: { children: ReactNode }) {
  return <p className="text-sm text-[var(--color-accent)]">{children}</p>;
}
