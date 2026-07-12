import { forwardRef, type ComponentPropsWithoutRef } from "react";

export const TextInput = forwardRef<
  HTMLInputElement,
  ComponentPropsWithoutRef<"input">
>(function TextInput({ className = "", ...props }, ref) {
  return (
    <input
      ref={ref}
      {...props}
      className={`rough block w-full h-11 px-3 bg-[var(--color-surface)] text-[var(--color-ink)] placeholder:text-[var(--color-ink-muted)] focus:outline-none focus:[--rough-color:var(--color-accent)] transition-colors ${className}`}
    />
  );
});
