import { forwardRef, type ComponentPropsWithoutRef } from "react";

export const TextInput = forwardRef<
  HTMLInputElement,
  ComponentPropsWithoutRef<"input">
>(function TextInput({ className = "", ...props }, ref) {
  return (
    <input
      ref={ref}
      {...props}
      className={`block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] text-[var(--color-ink)] placeholder:text-[var(--color-ink-muted)] focus:outline-none focus:border-[var(--color-ink)] transition-colors ${className}`}
    />
  );
});
