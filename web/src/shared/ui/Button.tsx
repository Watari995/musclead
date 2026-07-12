import type { ComponentPropsWithoutRef } from "react";

type ButtonProps = ComponentPropsWithoutRef<"button"> & {
  variant?: "primary" | "ghost" | "danger";
  fullWidth?: boolean;
};

export function Button({
  children,
  variant = "primary",
  fullWidth = false,
  className = "",
  ...props
}: ButtonProps) {
  const base =
    "rough h-11 px-5 text-sm font-medium inline-flex items-center justify-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed";
  const variants = {
    primary:
      "bg-[var(--color-ink)] text-[var(--color-surface)] hover:opacity-90 active:opacity-80",
    ghost:
      "bg-[var(--color-surface)] text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)]",
    danger:
      "bg-[var(--color-surface)] text-[var(--color-accent)] hover:bg-[var(--color-surface-alt)]",
  };
  return (
    <button
      {...props}
      className={`${base} ${variants[variant]} ${fullWidth ? "w-full" : ""} ${className}`}
    >
      {children}
    </button>
  );
}
