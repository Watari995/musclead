"use client";

import { useEffect, useState } from "react";
import Cropper, { type Area } from "react-easy-crop";
import { Button } from "@/shared/ui";
import { cropImageToBlob } from "./cropImage";

/**
 * 画像クロップモーダル(GitHub 風)。
 * - 1:1 アスペクトの正方形 crop(表示は丸マスク)
 * - ズーム/ドラッグ操作
 * - 「適用」 で 512x512 JPEG の Blob を返す
 */
export function CropModal({
  imageSrc,
  open,
  onApply,
  onCancel,
}: {
  imageSrc: string | null;
  open: boolean;
  onApply: (blob: Blob) => void;
  onCancel: () => void;
}) {
  const [crop, setCrop] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null);
  const [working, setWorking] = useState(false);

  // モーダルが開く / 画像が変わるたびに crop 状態をリセット
  // (React 19 公式パターン: prop 変化で派生状態をリセットする)
  const [lastImageSrc, setLastImageSrc] = useState(imageSrc);
  if (open && lastImageSrc !== imageSrc) {
    setLastImageSrc(imageSrc);
    setCrop({ x: 0, y: 0 });
    setZoom(1);
    setCroppedAreaPixels(null);
  }

  // ESC キーでキャンセル
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onCancel();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onCancel]);

  // 背景スクロールロック
  useEffect(() => {
    if (!open) return;
    const original = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = original;
    };
  }, [open]);

  if (!open || !imageSrc) return null;

  const handleApply = async () => {
    if (!croppedAreaPixels) return;
    try {
      setWorking(true);
      const blob = await cropImageToBlob(imageSrc, croppedAreaPixels);
      onApply(blob);
    } catch (e) {
      console.error("crop failed", e);
      setWorking(false);
    }
  };

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="プロフィール画像をクロップ"
      className="fixed inset-0 z-50 flex items-center justify-center"
    >
      <button
        type="button"
        aria-label="キャンセル"
        onClick={onCancel}
        className="absolute inset-0 bg-black/50"
      />
      <div className="relative bg-[var(--color-surface)] rounded-lg shadow-xl w-full max-w-md mx-4 overflow-hidden">
        <div className="px-5 py-4 border-b border-[var(--color-line)]">
          <h2 className="text-base font-bold tracking-tight">
            プロフィール画像を編集
          </h2>
        </div>

        {/* Crop エリア */}
        <div className="relative h-80 bg-[var(--color-surface-alt)]">
          <Cropper
            image={imageSrc}
            crop={crop}
            zoom={zoom}
            aspect={1}
            cropShape="round"
            showGrid={false}
            onCropChange={setCrop}
            onZoomChange={setZoom}
            onCropComplete={(_, pixels) => setCroppedAreaPixels(pixels)}
          />
        </div>

        {/* ズームスライダー */}
        <div className="px-5 py-4 space-y-2 border-t border-[var(--color-line)]">
          <label className="block text-xs text-[var(--color-ink-muted)]">
            ズーム
          </label>
          <input
            type="range"
            min={1}
            max={3}
            step={0.01}
            value={zoom}
            onChange={(e) => setZoom(Number(e.target.value))}
            className="w-full"
            disabled={working}
          />
        </div>

        {/* アクション */}
        <div className="flex gap-3 px-5 py-4 border-t border-[var(--color-line)]">
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={working}
          >
            キャンセル
          </Button>
          <Button
            type="button"
            onClick={handleApply}
            disabled={working || !croppedAreaPixels}
            className="flex-1"
          >
            {working ? "処理中…" : "適用"}
          </Button>
        </div>
      </div>
    </div>
  );
}
