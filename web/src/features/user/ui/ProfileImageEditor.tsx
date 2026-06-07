"use client";

import { useRef, useState } from "react";
import {
  useUpdateUserMutation,
  useUploadProfileImageMutation,
} from "@/features/user/api/user";
import { Avatar } from "./Avatar";
import { CropModal } from "./CropModal";

/**
 * GitHub 風プロフィール画像エディタ:
 * - 大きい丸アバター表示
 * - アバタークリックでファイル選択
 * - 選択 → CropModal → 適用で S3 にアップロード → PATCH で確定
 * - 「削除」 で default 復帰
 */
export function ProfileImageEditor({
  imageURL,
  displayName,
}: {
  imageURL: string;
  displayName: string;
}) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [pickedImageSrc, setPickedImageSrc] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const uploadMutation = useUploadProfileImageMutation();
  const updateMutation = useUpdateUserMutation();

  const isPending = uploadMutation.isPending || updateMutation.isPending;
  // default 画像の時は削除ボタンを出さない(削除して default に戻す対象が無い)
  const isDefault = imageURL.endsWith("/profiles/default.png");

  const handlePickClick = () => {
    setError(null);
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    // 同じファイルを再選択しても onChange 発火させるためにリセット
    e.target.value = "";
    if (!file) return;

    if (!["image/jpeg", "image/png", "image/webp"].includes(file.type)) {
      setError("JPEG / PNG / WebP のみアップロード可能です");
      return;
    }
    if (file.size > 10 * 1024 * 1024) {
      setError("ファイルサイズは 10 MB 以下にしてください");
      return;
    }

    const url = URL.createObjectURL(file);
    setPickedImageSrc(url);
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    if (pickedImageSrc) {
      URL.revokeObjectURL(pickedImageSrc);
      setPickedImageSrc(null);
    }
  };

  const handleApply = async (blob: Blob) => {
    setError(null);
    try {
      const { path } = await uploadMutation.mutateAsync({
        blob,
        contentType: "image/jpeg",
      });
      await updateMutation.mutateAsync({ profile_image_path: path });
      closeModal();
    } catch (e) {
      const msg = e instanceof Error ? e.message : "アップロードに失敗しました";
      setError(msg);
      // モーダルは開いたままにしてもう一度試せるようにする
    }
  };

  const handleRemove = async () => {
    if (!confirm("プロフィール画像を削除しますか?")) return;
    setError(null);
    try {
      await updateMutation.mutateAsync({ profile_image_path: null });
    } catch (e) {
      const msg = e instanceof Error ? e.message : "削除に失敗しました";
      setError(msg);
    }
  };

  return (
    <div className="flex flex-col items-center gap-3">
      {/* アバター + 編集 overlay */}
      <button
        type="button"
        onClick={handlePickClick}
        disabled={isPending}
        aria-label="プロフィール画像を編集"
        className="relative group rounded-full focus:outline-none focus:ring-2 focus:ring-[var(--color-ink)] focus:ring-offset-2 disabled:opacity-60"
      >
        <Avatar
          src={imageURL}
          alt={displayName}
          size="w-28 h-28"
          className="border-2"
        />
        <span className="absolute inset-0 rounded-full bg-black/0 group-hover:bg-black/40 transition-colors flex items-center justify-center">
          <span className="text-white text-xs font-medium opacity-0 group-hover:opacity-100 transition-opacity">
            編集
          </span>
        </span>
      </button>

      {/* 隠しファイル入力 */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/webp"
        className="hidden"
        onChange={handleFileChange}
      />

      {/* 削除ボタン(default 画像の時は隠す)*/}
      {!isDefault && (
        <button
          type="button"
          onClick={handleRemove}
          disabled={isPending}
          className="text-xs text-[var(--color-ink-muted)] hover:text-red-600 disabled:opacity-50"
        >
          プロフィール画像を削除
        </button>
      )}

      {/* エラー表示 */}
      {error && (
        <p className="text-xs text-red-600 text-center max-w-xs">{error}</p>
      )}

      {/* クロップモーダル */}
      <CropModal
        imageSrc={pickedImageSrc}
        open={modalOpen}
        onApply={handleApply}
        onCancel={closeModal}
      />
    </div>
  );
}
