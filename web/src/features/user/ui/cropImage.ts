/**
 * react-easy-crop の croppedAreaPixels を使って、 元画像から正方形の crop image を生成する。
 * 出力サイズは 512x512、 形式は JPEG(quality 0.9)。
 *
 * 戻り値: Blob + width/height(デバッグ用)
 */
export type PixelCrop = {
  x: number;
  y: number;
  width: number;
  height: number;
};

const OUTPUT_SIZE = 512;
const JPEG_QUALITY = 0.9;

export async function cropImageToBlob(
  imageSrc: string,
  pixelCrop: PixelCrop,
): Promise<Blob> {
  const image = await loadImage(imageSrc);

  const canvas = document.createElement("canvas");
  canvas.width = OUTPUT_SIZE;
  canvas.height = OUTPUT_SIZE;
  const ctx = canvas.getContext("2d");
  if (!ctx) throw new Error("failed to get 2d context");

  ctx.drawImage(
    image,
    pixelCrop.x,
    pixelCrop.y,
    pixelCrop.width,
    pixelCrop.height,
    0,
    0,
    OUTPUT_SIZE,
    OUTPUT_SIZE,
  );

  return new Promise<Blob>((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (!blob) return reject(new Error("toBlob returned null"));
        resolve(blob);
      },
      "image/jpeg",
      JPEG_QUALITY,
    );
  });
}

function loadImage(src: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve(img);
    img.onerror = () => reject(new Error("failed to load image"));
    img.src = src;
  });
}
