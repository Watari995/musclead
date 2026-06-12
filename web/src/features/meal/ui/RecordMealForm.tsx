"use client";

import { useRef, useState } from "react";
import type { RecordMealRequest } from "@/shared/api/client";
import {
  useRecordMealMutation,
  useUploadMealPhotoMutation,
} from "@/features/meal/api/meals";
import { toLocalInput } from "@/features/meal/model/meal";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

const MAX_PHOTOS = 5;
const ACCEPT_TYPES = ["image/jpeg", "image/png", "image/webp"];

type LocalPhoto = {
  file: File;
  previewURL: string;
};

export function RecordMealForm() {
  const [form, setForm] = useState<RecordMealRequest>(initialForm);
  const [photos, setPhotos] = useState<LocalPhoto[]>([]);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const recordMutation = useRecordMealMutation();
  const uploadMutation = useUploadMealPhotoMutation();

  const isPending = recordMutation.isPending || uploadMutation.isPending;

  const handlePickPhotos = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUploadError(null);
    const files = Array.from(e.target.files ?? []);
    e.target.value = "";
    if (files.length === 0) return;

    const next: LocalPhoto[] = [];
    for (const file of files) {
      if (!ACCEPT_TYPES.includes(file.type)) {
        setUploadError("JPEG / PNG / WebP のみアップロード可能です");
        return;
      }
      if (file.size > 10 * 1024 * 1024) {
        setUploadError("ファイルサイズは 10 MB 以下にしてください");
        return;
      }
      if (photos.length + next.length >= MAX_PHOTOS) {
        setUploadError(`写真は最大 ${MAX_PHOTOS} 枚まで追加できます`);
        break;
      }
      next.push({ file, previewURL: URL.createObjectURL(file) });
    }
    setPhotos([...photos, ...next]);
  };

  const handleRemovePhoto = (idx: number) => {
    const target = photos[idx];
    if (target) URL.revokeObjectURL(target.previewURL);
    setPhotos(photos.filter((_, i) => i !== idx));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setUploadError(null);

    // photo を全部並列アップロード → path を集める
    let imagePaths: string[] = [];
    if (photos.length > 0) {
      try {
        const results = await Promise.all(
          photos.map((p) => uploadMutation.mutateAsync({ file: p.file })),
        );
        imagePaths = results.map((r) => r.path);
      } catch (err) {
        const msg = err instanceof Error ? err.message : "写真のアップロードに失敗しました";
        setUploadError(msg);
        return;
      }
    }

    // 空文字のメモは BE で String1000 validation NG になるので、 undefined にして送らない
    const trimmedMemo = (form.memo ?? "").trim();
    recordMutation.mutate(
      {
        ...form,
        eaten_at: new Date(form.eaten_at!).toISOString(),
        calories: form.calories ?? 0,
        protein_g: form.protein_g ?? 0,
        fat_g: form.fat_g ?? 0,
        carbohydrate_g: form.carbohydrate_g ?? 0,
        memo: trimmedMemo === "" ? undefined : trimmedMemo,
        photos: imagePaths.map((path, i) => ({
          image_path: path,
          display_order: i,
        })),
      },
      {
        onSuccess: () => {
          photos.forEach((p) => URL.revokeObjectURL(p.previewURL));
          setPhotos([]);
          setForm(initialForm());
        },
      },
    );
  };

  return (
    <Card className="p-4 sm:p-5">
      <form className="space-y-4" onSubmit={handleSubmit}>
        <Label label="種類">
          <select
            value={form.meal_type}
            onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
            className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] focus:outline-none focus:border-[var(--color-ink)]"
          >
            <option value="breakfast">朝食</option>
            <option value="lunch">昼食</option>
            <option value="dinner">夕食</option>
            <option value="snack">間食</option>
          </select>
        </Label>
        <Label label="日時">
          <TextInput
            type="datetime-local"
            value={form.eaten_at}
            onChange={(e) => setForm({ ...form, eaten_at: e.target.value })}
            required
          />
        </Label>
        <div className="grid grid-cols-2 gap-3">
          <NumField
            label="カロリー (kcal)"
            value={form.calories}
            onChange={(v) => setForm({ ...form, calories: v })}
          />
          <NumField
            label="タンパク質 (g)"
            step="0.1"
            value={form.protein_g}
            onChange={(v) => setForm({ ...form, protein_g: v })}
          />
          <NumField
            label="脂質 (g)"
            step="0.1"
            value={form.fat_g}
            onChange={(v) => setForm({ ...form, fat_g: v })}
          />
          <NumField
            label="炭水化物 (g)"
            step="0.1"
            value={form.carbohydrate_g}
            onChange={(v) => setForm({ ...form, carbohydrate_g: v })}
          />
        </div>
        <Label label="メモ">
          <textarea
            value={form.memo ?? ""}
            onChange={(e) => setForm({ ...form, memo: e.target.value })}
            rows={2}
            className="block w-full px-3 py-2 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] focus:outline-none focus:border-[var(--color-ink)]"
          />
        </Label>

        {/* 写真 */}
        <div className="space-y-2">
          <span className="block text-xs text-[var(--color-ink-muted)]">
            写真({photos.length}/{MAX_PHOTOS})
          </span>
          {photos.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {photos.map((p, i) => (
                <div key={p.previewURL} className="relative">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={p.previewURL}
                    alt=""
                    className="w-20 h-20 rounded-md object-cover border border-[var(--color-line)]"
                  />
                  <button
                    type="button"
                    onClick={() => handleRemovePhoto(i)}
                    aria-label="削除"
                    disabled={isPending}
                    className="absolute -top-2 -right-2 w-5 h-5 rounded-full bg-black/70 text-white text-xs flex items-center justify-center hover:bg-black disabled:opacity-50"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={isPending || photos.length >= MAX_PHOTOS}
            className="text-xs text-[var(--color-ink)] underline disabled:opacity-50 disabled:no-underline"
          >
            写真を追加
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept={ACCEPT_TYPES.join(",")}
            multiple
            className="hidden"
            onChange={handlePickPhotos}
          />
        </div>

        {uploadError && <ErrorText>{uploadError}</ErrorText>}
        {recordMutation.isError && (
          <ErrorText>{(recordMutation.error as Error).message}</ErrorText>
        )}
        <Button type="submit" fullWidth disabled={isPending}>
          {isPending ? "記録中…" : "記録する"}
        </Button>
      </form>
    </Card>
  );
}

function NumField({
  label,
  value,
  onChange,
  step = "1",
}: {
  label: string;
  value: number | undefined;
  onChange: (v: number | undefined) => void;
  step?: string;
}) {
  return (
    <Label label={label}>
      <TextInput
        type="number"
        step={step}
        min={0}
        placeholder="0"
        value={value ?? ""}
        onChange={(e) => {
          const raw = e.target.value;
          onChange(raw === "" ? undefined : Number(raw));
        }}
      />
    </Label>
  );
}

function initialForm(): RecordMealRequest {
  return {
    meal_type: "breakfast",
    eaten_at: toLocalInput(new Date()),
    memo: "",
    photos: [],
  };
}
