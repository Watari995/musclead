import { forwardRef, type ComponentPropsWithoutRef, type ReactNode } from "react";

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

export function SectionTitle({ children }: { children: ReactNode }) {
  return (
    <h2 className="text-base font-bold tracking-tight text-[var(--color-ink)] mb-3">
      {children}
    </h2>
  );
}

export const TextInput = forwardRef<
  HTMLInputElement,
  ComponentPropsWithoutRef<"input">
>(function TextInput({ className = "", ...props }, ref) {
  return (
    <input
      ref={ref}
      {...props}
      className={`block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-white text-[var(--color-ink)] placeholder:text-[var(--color-ink-muted)] focus:outline-none focus:border-[var(--color-ink)] transition-colors ${className}`}
    />
  );
});

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
    "h-11 px-5 rounded-md text-sm font-medium inline-flex items-center justify-center transition-colors disabled:opacity-50 disabled:cursor-not-allowed";
  const variants = {
    primary:
      "bg-[var(--color-ink)] text-white hover:opacity-90 active:opacity-80",
    ghost:
      "bg-white text-[var(--color-ink)] border border-[var(--color-line)] hover:bg-[var(--color-surface-alt)]",
    danger:
      "bg-white text-[var(--color-accent)] border border-[var(--color-line)] hover:bg-[var(--color-surface-alt)]",
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

export function ErrorText({ children }: { children: ReactNode }) {
  return <p className="text-sm text-[var(--color-accent)]">{children}</p>;
}
