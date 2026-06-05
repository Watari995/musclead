"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { Button } from "./Button";

export type ConfirmOptions = {
  title: string;
  description?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  destructive?: boolean;
};

type Confirm = (options: ConfirmOptions) => Promise<boolean>;

const ConfirmContext = createContext<Confirm | null>(null);

export function useConfirm(): Confirm {
  const fn = useContext(ConfirmContext);
  if (!fn) {
    throw new Error("useConfirm must be used within <ConfirmProvider>");
  }
  return fn;
}

type Pending = {
  options: ConfirmOptions;
  resolve: (result: boolean) => void;
};

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [pending, setPending] = useState<Pending | null>(null);

  const confirm = useCallback<Confirm>((options) => {
    return new Promise<boolean>((resolve) => {
      setPending({ options, resolve });
    });
  }, []);

  const close = useCallback(
    (result: boolean) => {
      setPending((prev) => {
        prev?.resolve(result);
        return null;
      });
    },
    [],
  );

  return (
    <ConfirmContext.Provider value={confirm}>
      {children}
      {pending && (
        <ConfirmDialog
          {...pending.options}
          onConfirm={() => close(true)}
          onCancel={() => close(false)}
        />
      )}
    </ConfirmContext.Provider>
  );
}

type DialogProps = ConfirmOptions & {
  onConfirm: () => void;
  onCancel: () => void;
};

function ConfirmDialog({
  title,
  description,
  confirmLabel = "OK",
  cancelLabel = "キャンセル",
  destructive,
  onConfirm,
  onCancel,
}: DialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    dialogRef.current?.showModal();
  }, []);

  // Esc / dialog 内部 close() を全部 cancel に集約。
  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;
    const onClose = () => onCancel();
    dialog.addEventListener("close", onClose);
    return () => dialog.removeEventListener("close", onClose);
  }, [onCancel]);

  // 背景クリック (= ::backdrop) で閉じる。 dialog 要素本体上に
  // クリックが落ちた = backdrop クリック。
  const onBackdropClick = (e: React.MouseEvent<HTMLDialogElement>) => {
    if (e.target === dialogRef.current) onCancel();
  };

  return (
    <dialog
      ref={dialogRef}
      onClick={onBackdropClick}
      // Tailwind v4 の preflight が <dialog> の margin auto を潰すので、
      // m-auto + fixed inset-0 でブラウザデフォルトの中央配置を復元する。
      className="fixed inset-0 m-auto rounded-lg border border-[var(--color-line)] bg-white p-0 w-[min(420px,calc(100vw-2rem))] backdrop:bg-black/40"
    >
      <div className="p-6 space-y-4">
        <h2 className="text-base font-bold tracking-tight text-[var(--color-ink)]">
          {title}
        </h2>
        {description && (
          <p className="text-sm text-[var(--color-ink)] whitespace-pre-line">
            {description}
          </p>
        )}
        <div className="flex gap-3 justify-end pt-2">
          <Button type="button" variant="ghost" onClick={onCancel}>
            {cancelLabel}
          </Button>
          <Button
            type="button"
            variant={destructive ? "danger" : "primary"}
            onClick={onConfirm}
          >
            {confirmLabel}
          </Button>
        </div>
      </div>
    </dialog>
  );
}
