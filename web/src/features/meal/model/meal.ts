import type { MealDTO } from "@/shared/api/client";

export type MealType = "breakfast" | "lunch" | "dinner" | "snack";

export type Meal = {
  id: string;
  type: MealType;
  eatenAt: string; // ISO
  calories: number;
  proteinG: string;
  fatG: string;
  carbohydrateG: string;
  memo: string;
};

export function toMeal(dto: MealDTO): Meal {
  return {
    id: dto.id ?? "",
    type: (dto.meal_type as MealType) ?? "breakfast",
    eatenAt: dto.eaten_at ?? "",
    calories: dto.calories ?? 0,
    proteinG: dto.protein_g ?? "0",
    fatG: dto.fat_g ?? "0",
    carbohydrateG: dto.carbohydrate_g ?? "0",
    memo: dto.memo ?? "",
  };
}

export function mealTypeLabel(t: MealType): string {
  switch (t) {
    case "breakfast":
      return "朝食";
    case "lunch":
      return "昼食";
    case "dinner":
      return "夕食";
    case "snack":
      return "間食";
  }
}

export function mealTypeEmoji(t: MealType): string {
  switch (t) {
    case "breakfast":
      return "🍳";
    case "lunch":
      return "🍱";
    case "dinner":
      return "🍽️";
    case "snack":
      return "🍎";
  }
}

export function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function formatMealDateTime(iso: string): string {
  if (!iso) return "";
  return new Date(iso).toLocaleString("ja-JP", {
    dateStyle: "short",
    timeStyle: "short",
  });
}
