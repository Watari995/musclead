"use client";

/**
 * 丸いプロフィールアバター。 画像 URL を受け取って表示する。
 * サイズは Tailwind class で指定(例: "w-10 h-10")。
 * 読み込み失敗時は default が表示される(サーバー側で必ず URL が返るため、 実質常に画像)。
 */
export function Avatar({
  src,
  alt,
  size = "w-10 h-10",
  className = "",
}: {
  src: string;
  alt: string;
  size?: string;
  className?: string;
}) {
  return (
    // eslint-disable-next-line @next/next/no-img-element -- S3 直配信のため Image コンポーネント非対応
    <img
      src={src}
      alt={alt}
      className={`${size} rough rough-pill object-cover bg-[var(--color-surface-alt)] ${className}`}
    />
  );
}
