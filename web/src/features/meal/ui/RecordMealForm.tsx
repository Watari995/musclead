"use client";

import { useEffect, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import type { RecordMealRequest } from "@/shared/api/client";
import {
  useRecordMealMutation,
  useUploadMealPhotoMutation,
} from "@/features/meal/api/meals";
import { toLocalInput } from "@/features/meal/model/meal";
import type { MealTemplate } from "@/features/meal/model/meal_template";
import { FoodSearchSection } from "@/features/food/ui/FoodSearchSection";
import { FoodRegisterModal } from "@/features/food/ui/FoodRegisterModal";
import type { FoodProduct } from "@/features/food/model/food_product";
import { Button, Card, ErrorText, Label, NumberField, TextInput } from "@/shared/ui";

const MAX_PHOTOS = 5;
const ACCEPT_TYPES = ["image/jpeg", "image/png", "image/webp"];

type LocalPhoto = {
  file: File;
  previewURL: string;
};

type BaseNutrients = {
  calories: number;
  proteinG: number | undefined;
  fatG: number | undefined;
  carbohydrateG: number | undefined;
};

type Props = {
  prefill?: MealTemplate | null;
  onPrefillConsumed?: () => void;
  onSuccess?: () => void;
};

export function RecordMealForm({ prefill, onPrefillConsumed, onSuccess }: Props) {
  const t = useTranslations("meals");
  const tCommon = useTranslations("common");
  const [form, setForm] = useState<RecordMealRequest>(initialForm);
  const [showRegisterModal, setShowRegisterModal] = useState(false);
  const [registerBarcode, setRegisterBarcode] = useState<string | undefined>();
  const [selectedFoodId, setSelectedFoodId] = useState<string | undefined>();
  const [baseNutrients, setBaseNutrients] = useState<BaseNutrients | undefined>();
  const [servingCount, setServingCount] = useState<number>(1);

  useEffect(() => {
    if (!prefill) return;
    setForm((prev) => ({
      ...initialForm(),
      meal_type: prefill.mealType,
      eaten_at: prev.eaten_at,
      calories: prefill.calories,
      protein_g: parseFloat(prefill.proteinG) || undefined,
      fat_g: parseFloat(prefill.fatG) || undefined,
      carbohydrate_g: parseFloat(prefill.carbohydrateG) || undefined,
    }));
    onPrefillConsumed?.();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [prefill]);

  const [photos, setPhotos] = useState<LocalPhoto[]>([]);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const recordMutation = useRecordMealMutation();
  const uploadMutation = useUploadMealPhotoMutation();

  const isPending = recordMutation.isPending || uploadMutation.isPending;

  const handleFoodSelect = (food: FoodProduct) => {
    const base: BaseNutrients = {
      calories: food.calories,
      proteinG: food.proteinG ? parseFloat(food.proteinG) : undefined,
      fatG: food.fatG ? parseFloat(food.fatG) : undefined,
      carbohydrateG: food.carbohydrateG ? parseFloat(food.carbohydrateG) : undefined,
    };
    setSelectedFoodId(food.id);
    setBaseNutrients(base);
    setServingCount(1);
    setForm((prev) => ({
      ...prev,
      calories: base.calories,
      protein_g: base.proteinG,
      fat_g: base.fatG,
      carbohydrate_g: base.carbohydrateG,
      memo: (prev.memo ?? "").trim() === "" ? food.name : prev.memo,
    }));
  };

  const handleServingCountChange = (v: number | undefined) => {
    const s = v ?? 1;
    setServingCount(s);
    if (baseNutrients) {
      setForm((prev) => ({
        ...prev,
        calories: Math.round(baseNutrients.calories * s),
        protein_g: baseNutrients.proteinG !== undefined ? baseNutrients.proteinG * s : undefined,
        fat_g: baseNutrients.fatG !== undefined ? baseNutrients.fatG * s : undefined,
        carbohydrate_g: baseNutrients.carbohydrateG !== undefined ? baseNutrients.carbohydrateG * s : undefined,
      }));
    }
  };

  const handlePickPhotos = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUploadError(null);
    const files = Array.from(e.target.files ?? []);
    e.target.value = "";
    if (files.length === 0) return;

    const next: LocalPhoto[] = [];
    for (const file of files) {
      if (!ACCEPT_TYPES.includes(file.type)) {
        setUploadError(tCommon("uploadError.type"));
        return;
      }
      if (file.size > 10 * 1024 * 1024) {
        setUploadError(tCommon("uploadError.size"));
        return;
      }
      if (photos.length + next.length >= MAX_PHOTOS) {
        setUploadError(tCommon("uploadError.max", { max: MAX_PHOTOS }));
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

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setUploadError(null);

    let imagePaths: string[] = [];
    if (photos.length > 0) {
      try {
        const results = await Promise.all(
          photos.map((p) => uploadMutation.mutateAsync({ file: p.file })),
        );
        imagePaths = results.map((r) => r.path);
      } catch (err) {
        const msg = err instanceof Error ? err.message : tCommon("uploadError.fail");
        setUploadError(msg);
        return;
      }
    }

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
        food_product_id: selectedFoodId,
        serving_count: servingCount,
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
          setSelectedFoodId(undefined);
          setBaseNutrients(undefined);
          setServingCount(1);
          onSuccess?.();
        },
      },
    );
  };

  return (
    <Card className="p-4 sm:p-5">
      <div className="space-y-4">
        <FoodSearchSection
          onSelect={handleFoodSelect}
          onNotFound={(barcode) => {
            setRegisterBarcode(barcode);
            setShowRegisterModal(true);
          }}
        />
        {baseNutrients && (
          <Label label={t("servingCount")}>
            <NumberField
              step="0.5"
              min={0.5}
              placeholder="1"
              value={servingCount}
              onChange={handleServingCountChange}
            />
          </Label>
        )}

        <form className="space-y-4" onSubmit={handleSubmit}>
          <Label label={t("mealType")}>
            <select
              value={form.meal_type}
              onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
              className="block w-full h-11 px-3 rough bg-[var(--color-surface)] focus:outline-none focus:[--rough-color:var(--color-accent)]"
            >
              <option value="breakfast">{t("breakfast")}</option>
              <option value="lunch">{t("lunch")}</option>
              <option value="dinner">{t("dinner")}</option>
              <option value="snack">{t("snack")}</option>
            </select>
          </Label>
          <Label label={tCommon("dateTime")}>
            <TextInput
              type="datetime-local"
              value={form.eaten_at}
              onChange={(e) => setForm({ ...form, eaten_at: e.target.value })}
              required
            />
          </Label>
          <div className="grid grid-cols-2 gap-3">
            <NumField
              label={t("calories")}
              value={form.calories}
              onChange={(v) => setForm({ ...form, calories: v })}
            />
            <NumField
              label={t("protein")}
              step="0.1"
              value={form.protein_g}
              onChange={(v) => setForm({ ...form, protein_g: v })}
            />
            <NumField
              label={t("fat")}
              step="0.1"
              value={form.fat_g}
              onChange={(v) => setForm({ ...form, fat_g: v })}
            />
            <NumField
              label={t("carbs")}
              step="0.1"
              value={form.carbohydrate_g}
              onChange={(v) => setForm({ ...form, carbohydrate_g: v })}
            />
          </div>
          <Label label={tCommon("memo")}>
            <textarea
              value={form.memo ?? ""}
              onChange={(e) => setForm({ ...form, memo: e.target.value })}
              rows={2}
              className="block w-full px-3 py-2 rough bg-[var(--color-surface)] focus:outline-none focus:[--rough-color:var(--color-accent)]"
            />
          </Label>

          <div className="space-y-2">
            <span className="block text-xs text-[var(--color-ink-muted)]">
              {tCommon("photo")}({photos.length}/{MAX_PHOTOS})
            </span>
            {photos.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {photos.map((p, i) => (
                  <div key={p.previewURL} className="relative">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={p.previewURL}
                      alt=""
                      className="rough w-20 h-20 object-cover"
                    />
                    <button
                      type="button"
                      onClick={() => handleRemovePhoto(i)}
                      aria-label={tCommon("delete")}
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
              {tCommon("addPhoto")}
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
            {isPending ? tCommon("recording") : tCommon("record")}
          </Button>
        </form>
      </div>

      {showRegisterModal && (
        <FoodRegisterModal
          initialBarcode={registerBarcode}
          onSuccess={(food) => {
            setShowRegisterModal(false);
            handleFoodSelect(food);
          }}
          onCancel={() => setShowRegisterModal(false)}
        />
      )}
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
      <NumberField
        step={step}
        min={0}
        placeholder="0"
        value={value}
        onChange={onChange}
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
