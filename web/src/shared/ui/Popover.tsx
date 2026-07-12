"use client";

import {
  useEffect,
  useId,
  useLayoutEffect,
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

// viewport の左右端からこの距離は必ず空ける(モバイルで端ギリギリにしない)
const VIEWPORT_MARGIN = 8;
// trigger と content の隙間 (mt-2 相当)
const TRIGGER_GAP = 8;

// 軽量 Popover: シンプル要件 (single anchor, click-outside / Escape で閉じる,
// 縦方向は常に下開き) を満たしつつ、 viewport の左右端からはみ出さないように
// 開いた直後に位置を補正する。 補正後の位置はビューポート基準 (fixed) なので、
// scroll / resize では追従させずに閉じる。
export function Popover({ trigger, children, align = "start" }: Props) {
  const [open, setOpen] = useState(false);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const contentId = useId();

  useEffect(() => {
    if (!open) return;
    const onPointer = (e: MouseEvent | TouchEvent) => {
      const target = e.target as Node;
      if (
        !wrapperRef.current?.contains(target) &&
        !contentRef.current?.contains(target)
      ) {
        setOpen(false);
      }
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    const close = () => setOpen(false);
    document.addEventListener("mousedown", onPointer);
    document.addEventListener("touchstart", onPointer);
    document.addEventListener("keydown", onKey);
    window.addEventListener("scroll", close, true);
    window.addEventListener("resize", close);
    return () => {
      document.removeEventListener("mousedown", onPointer);
      document.removeEventListener("touchstart", onPointer);
      document.removeEventListener("keydown", onKey);
      window.removeEventListener("scroll", close, true);
      window.removeEventListener("resize", close);
    };
  }, [open]);

  // 初回 paint で flicker しないよう、 content は invisible で mount し、
  // 位置を確定した直後 (同 commit phase 内) に表示する
  useLayoutEffect(() => {
    if (!open) return;
    const trigger = wrapperRef.current?.getBoundingClientRect();
    const node = contentRef.current;
    if (!trigger || !node) return;
    const content = node.getBoundingClientRect();
    const vw = document.documentElement.clientWidth;

    // 1) align に従って希望位置を決める
    let left =
      align === "end" ? trigger.right - content.width : trigger.left;
    // 2) 右端衝突 → 左に詰める
    if (left + content.width > vw - VIEWPORT_MARGIN) {
      left = vw - VIEWPORT_MARGIN - content.width;
    }
    // 3) 左端衝突 → 右に詰める (2 の結果より優先)
    if (left < VIEWPORT_MARGIN) left = VIEWPORT_MARGIN;

    node.style.left = `${left}px`;
    node.style.top = `${trigger.bottom + TRIGGER_GAP}px`;
    node.style.maxWidth = `calc(100vw - ${VIEWPORT_MARGIN * 2}px)`;
    node.style.visibility = "visible";
  }, [open, align]);

  return (
    <div ref={wrapperRef} className="relative inline-block">
      {trigger({
        open,
        onClick: () => setOpen((o) => !o),
        "aria-expanded": open,
        "aria-controls": contentId,
      })}
      {open && (
        <div
          ref={contentRef}
          id={contentId}
          role="dialog"
          className="rough fixed invisible z-20 bg-[var(--color-surface)]"
        >
          {children}
        </div>
      )}
    </div>
  );
}
