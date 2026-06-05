"use client";

import { useDeleteMealMutation } from "@/features/meal/api/meals";
import {
  formatMealDateTime,
  mealTypeEmoji,
  mealTypeLabel,
  type Meal,
} from "@/features/meal/model/meal";

export function MealRow({ meal }: { meal: Meal }) {
  const del = useDeleteMealMutation();

  return (
    <li className="p-4 flex items-start gap-4">
      <div className="w-14 h-14 shrink-0 rounded-md bg-[var(--color-surface-alt)] flex items-center justify-center text-xl">
        {mealTypeEmoji(meal.type)}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {mealTypeLabel(meal.type)}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {formatMealDateTime(meal.eatenAt)}
          </span>
        </div>
        {meal.memo && (
          <p className="mt-1 text-sm text-[var(--color-ink)] line-clamp-2">
            {meal.memo}
          </p>
        )}
        <div className="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-[var(--color-ink-muted)]">
          <span className="font-medium text-[var(--color-ink)]">
            {meal.calories} kcal
          </span>
          <span>P {meal.proteinG}g</span>
          <span>F {meal.fatG}g</span>
          <span>C {meal.carbohydrateG}g</span>
        </div>
      </div>
      <button
        type="button"
        onClick={() => del.mutate(meal.id)}
        disabled={del.isPending}
        className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] shrink-0"
      >
        削除
      </button>
    </li>
  );
}
