"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button, ErrorText, Label, NumberField, TextInput } from "@/shared/ui";
import { useCreateFoodProductMutation } from "../api/food_products";
import type { FoodProduct } from "../model/food_product";

type Props = {
  initialBarcode?: string;
  onSuccess: (food: FoodProduct) => void;
  onCancel: () => void;
};

export function FoodRegisterModal({ initialBarcode, onSuccess, onCancel }: Props) {
  const t = useTranslations("food");
  const tCommon = useTranslations("common");
  const [barcode, setBarcode] = useState(initialBarcode ?? "");
  const [name, setName] = useState("");
  const [calories, setCalories] = useState<number | undefined>();
  const [proteinG, setProteinG] = useState<number | undefined>();
  const [fatG, setFatG] = useState<number | undefined>();
  const [carbohydrateG, setCarbohydrateG] = useState<number | undefined>();

  const mutation = useCreateFoodProductMutation();

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!name.trim() || calories == null) return;
    try {
      const id = await mutation.mutateAsync({
        barcode: barcode.trim() || undefined,
        name: name.trim(),
        calories,
        protein_g: proteinG != null ? proteinG.toString() : undefined,
        fat_g: fatG != null ? fatG.toString() : undefined,
        carbohydrate_g: carbohydrateG != null ? carbohydrateG.toString() : undefined,
      });
      onSuccess({
        id,
        barcode: barcode.trim() || null,
        name: name.trim(),
        calories,
        proteinG: proteinG != null ? proteinG.toString() : null,
        fatG: fatG != null ? fatG.toString() : null,
        carbohydrateG: carbohydrateG != null ? carbohydrateG.toString() : null,
        registerSource: "user",
      });
    } catch {
      // エラーは mutation.isError で表示
    }
  };

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-end sm:items-center justify-center z-50 p-4"
      onClick={(e) => e.target === e.currentTarget && onCancel()}
    >
      <div className="bg-[var(--color-surface)] rounded-xl p-6 w-full max-w-sm space-y-4">
        <h2 className="text-base font-bold">{t("registerTitle")}</h2>
        <form className="space-y-3" onSubmit={handleSubmit}>
          {initialBarcode && (
            <Label label={t("barcodeLabel")}>
              <TextInput
                value={barcode}
                onChange={(e) => setBarcode(e.target.value)}
                placeholder="4901085615881"
              />
            </Label>
          )}
          <Label label={t("foodName")}>
            <TextInput
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t("foodNamePlaceholder")}
            />
          </Label>
          <Label label={t("calories")}>
            <NumberField
              min={0}
              placeholder="0"
              value={calories}
              onChange={setCalories}
            />
          </Label>
          <div className="grid grid-cols-3 gap-2">
            <Label label="P (g)">
              <NumberField step="0.1" min={0} placeholder="0" value={proteinG} onChange={setProteinG} />
            </Label>
            <Label label="F (g)">
              <NumberField step="0.1" min={0} placeholder="0" value={fatG} onChange={setFatG} />
            </Label>
            <Label label="C (g)">
              <NumberField step="0.1" min={0} placeholder="0" value={carbohydrateG} onChange={setCarbohydrateG} />
            </Label>
          </div>
          {mutation.isError && (
            <ErrorText>{(mutation.error as Error).message}</ErrorText>
          )}
          <div className="flex gap-2 pt-1">
            <Button type="button" variant="ghost" fullWidth onClick={onCancel}>
              {tCommon("cancel")}
            </Button>
            <Button type="submit" fullWidth disabled={mutation.isPending}>
              {mutation.isPending ? t("registering") : t("register")}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
