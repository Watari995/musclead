"use client";

import Link from "next/link";
import { useDeleteMealMutation } from "@/features/meal/api/meals";
import {
  formatMealDateTime,
  mealTypeEmoji,
  mealTypeLabel,
  type Meal,
} from "@/features/meal/model/meal";

export function MealRow({ meal }: { meal: Meal }) {
  const del = useDeleteMealMutation();

  const firstPhoto = meal.photos[0];

  return (
    <li className="p-4 flex items-start gap-4">
      {firstPhoto ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={firstPhoto.imageURL}
          alt=""
          className="w-14 h-14 shrink-0 rounded-md object-cover border border-[var(--color-line)]"
        />
      ) : (
        <div className="w-14 h-14 shrink-0 rounded-md bg-[var(--color-surface-alt)] flex items-center justify-center text-xl">
          {mealTypeEmoji(meal.type)}
        </div>
      )}
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
        {meal.photos.length > 1 && (
          <div className="mt-2 flex gap-1">
            {meal.photos.slice(1, 4).map((p) => (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                key={p.imageURL}
                src={p.imageURL}
                alt=""
                className="w-10 h-10 rounded object-cover border border-[var(--color-line)]"
              />
            ))}
            {meal.photos.length > 4 && (
              <div className="w-10 h-10 rounded bg-[var(--color-surface-alt)] flex items-center justify-center text-xs text-[var(--color-ink-muted)]">
                +{meal.photos.length - 4}
              </div>
            )}
          </div>
        )}
      </div>
      <div className="flex flex-col items-end gap-2 shrink-0">
        <Link
          href={`/meals/${meal.id}/edit`}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] transition-colors"
        >
          編集
        </Link>
        <button
          type="button"
          onClick={() => {
            if (confirm(`${mealTypeLabel(meal.type)} の記録を削除しますか?`)) {
              del.mutate(meal.id);
            }
          }}
          disabled={del.isPending}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] transition-colors"
        >
          削除
        </button>
      </div>
    </li>
  );
}
