"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useDeleteWeightMutation } from "@/features/weight/api/weights";
import {
  formatWeightDateTime,
  type Weight,
} from "@/features/weight/model/weight";
import { EditWeightModal } from "./EditWeightModal";

export function WeightRow({ weight }: { weight: Weight }) {
  const t = useTranslations("weights");
  const tCommon = useTranslations("common");
  const [editing, setEditing] = useState(false);
  const del = useDeleteWeightMutation();

  return (
    <>
      <li className="p-4 flex items-start gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <span className="text-base font-bold tracking-tight">
              {weight.weightKg} kg
            </span>
            <span className="text-xs text-[var(--color-ink-muted)]">
              {formatWeightDateTime(weight.measuredAt)}
            </span>
          </div>
          {(weight.bodyFatPercentage || weight.skeletalMuscleKg) && (
            <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1 text-xs text-[var(--color-ink-muted)]">
              {weight.bodyFatPercentage && (
                <span>{t("bodyFatLabel", { value: weight.bodyFatPercentage })}</span>
              )}
              {weight.skeletalMuscleKg && (
                <span>{t("muscleMassLabel", { value: weight.skeletalMuscleKg })}</span>
              )}
            </div>
          )}
        </div>
        <div className="flex flex-col items-end gap-1 shrink-0">
          <button
            type="button"
            onClick={() => setEditing(true)}
            className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-ink)]"
          >
            {tCommon("edit")}
          </button>
          <button
            type="button"
            onClick={() => {
              if (confirm(t("deleteConfirm"))) {
                del.mutate(weight.id);
              }
            }}
            disabled={del.isPending}
            className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)]"
          >
            {tCommon("delete")}
          </button>
        </div>
      </li>
      {editing && (
        <EditWeightModal weight={weight} onClose={() => setEditing(false)} />
      )}
    </>
  );
}
