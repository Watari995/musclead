import type { MealTemplateDTO } from "@/shared/api/client";
import type { MealType } from "./meal";

export type MealTemplate = {
  id: string;
  name: string;
  displayOrder: number;
  mealType: MealType;
  calories: number;
  proteinG: string;
  fatG: string;
  carbohydrateG: string;
};

export function toMealTemplate(dto: MealTemplateDTO): MealTemplate {
  return {
    id: dto.id ?? "",
    name: dto.name ?? "",
    displayOrder: dto.display_order ?? 0,
    mealType: (dto.meal_type as MealType) ?? "breakfast",
    calories: dto.calories ?? 0,
    proteinG: dto.protein_g ?? "0",
    fatG: dto.fat_g ?? "0",
    carbohydrateG: dto.carbohydrate_g ?? "0",
  };
}
