"use client";

import {
  useEffect,
  useId,
  useRef,
  useState,
  type ReactNode,
} from "react";

type Props = {
  trigger: (props: {
    open: boolean;
    onClick: () => void;
    "aria-expanded": boolean;
    "aria-controls": string;
  }) => ReactNode;
  children: ReactNode;
  align?: "start" | "end";
};

// 軽量 Popover: Radix を入れずに済む程度のシンプル要件 (single anchor,
// click-outside + Escape で閉じる, 縦方向は常に下開き) に絞った実装。
// ビューポート衝突や focus trap などはサポートしない。
export function Popover({ trigger, children, align = "start" }: Props) {
  const [open, setOpen] = useState(false);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const contentId = useId();

  useEffect(() => {
    if (!open) return;
    const onPointer = (e: MouseEvent | TouchEvent) => {
      if (!wrapperRef.current?.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onPointer);
    document.addEventListener("touchstart", onPointer);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("touchstart", onPointer);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <div ref={wrapperRef} className="relative">
      {trigger({
        open,
        onClick: () => setOpen((o) => !o),
        "aria-expanded": open,
        "aria-controls": contentId,
      })}
      {open && (
        <div
          id={contentId}
          role="dialog"
          className={`absolute top-full mt-2 z-20 bg-white rounded-lg shadow-lg border border-[var(--color-line)] ${
            align === "end" ? "right-0" : "left-0"
          }`}
        >
          {children}
        </div>
      )}
    </div>
  );
}
