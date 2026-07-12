"use client";

import Link from "next/link";
import { useTranslations } from "next-intl";
import { useDeleteMealMutation } from "@/features/meal/api/meals";
import {
  formatMealDateTime,
  mealTypeEmoji,
  mealTypeLabelKey,
  type Meal,
} from "@/features/meal/model/meal";

export function MealRow({ meal }: { meal: Meal }) {
  const t = useTranslations("meals");
  const tCommon = useTranslations("common");
  const del = useDeleteMealMutation();

  const firstPhoto = meal.photos[0];

  return (
    <li className="p-4 flex items-start gap-4">
      {firstPhoto ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={firstPhoto.imageURL}
          alt=""
          className="rough w-14 h-14 shrink-0 object-cover"
        />
      ) : (
        <div className="rough w-14 h-14 shrink-0 bg-[var(--color-surface-alt)] flex items-center justify-center text-xl">
          {mealTypeEmoji(meal.type)}
        </div>
      )}
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {t(mealTypeLabelKey(meal.type))}
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
          <span className="text-[var(--color-macro-protein)]">
            P {meal.proteinG}g
          </span>
          <span className="text-[var(--color-macro-fat)]">
            F {meal.fatG}g
          </span>
          <span className="text-[var(--color-macro-carb)]">
            C {meal.carbohydrateG}g
          </span>
        </div>
        {meal.photos.length > 1 && (
          <div className="mt-2 flex gap-1">
            {meal.photos.slice(1, 4).map((p) => (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                key={p.imageURL}
                src={p.imageURL}
                alt=""
                className="rough w-10 h-10 object-cover"
              />
            ))}
            {meal.photos.length > 4 && (
              <div className="rough w-10 h-10 bg-[var(--color-surface-alt)] flex items-center justify-center text-xs text-[var(--color-ink-muted)]">
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
          {tCommon("edit")}
        </Link>
        <button
          type="button"
          onClick={() => {
            if (confirm(t("deleteConfirm", { mealType: t(mealTypeLabelKey(meal.type)) }))) {
              del.mutate(meal.id);
            }
          }}
          disabled={del.isPending}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] transition-colors"
        >
          {tCommon("delete")}
        </button>
      </div>
    </li>
  );
}
