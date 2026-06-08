"use client";

import { useState } from "react";
import {
  useUpdateWeightMutation,
  type UpsertWeightRequest,
} from "@/features/weight/api/weights";
import { toLocalInput, type Weight } from "@/features/weight/model/weight";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

type FormState = {
  weight_kg: string;
  body_fat_percentage: string;
  skeletal_muscle_kg: string;
  measured_at: string;
};

export function EditWeightModal({
  weight,
  onClose,
}: {
  weight: Weight;
  onClose: () => void;
}) {
  const [form, setForm] = useState<FormState>(() => ({
    weight_kg: weight.weightKg,
    body_fat_percentage: weight.bodyFatPercentage ?? "",
    skeletal_muscle_kg: weight.skeletalMuscleKg ?? "",
    measured_at: toLocalInput(new Date(weight.measuredAt)),
  }));
  const updateMutation = useUpdateWeightMutation();

  return (
    <div
      role="dialog"
      aria-modal="true"
      className="fixed inset-0 z-40 flex items-center justify-center p-4"
    >
      <button
        type="button"
        aria-label="閉じる"
        onClick={onClose}
        className="absolute inset-0 bg-black/40"
      />
      <Card className="relative w-full max-w-md p-5 z-50">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-bold tracking-tight">体重を編集</h2>
          <button
            type="button"
            onClick={onClose}
            aria-label="閉じる"
            className="text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] text-lg leading-none"
          >
            ×
          </button>
        </div>
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            const body: UpsertWeightRequest = {
              weight_kg: form.weight_kg,
              measured_at: new Date(form.measured_at).toISOString(),
            };
            if (form.body_fat_percentage.trim() !== "") {
              body.body_fat_percentage = form.body_fat_percentage;
            }
            if (form.skeletal_muscle_kg.trim() !== "") {
              body.skeletal_muscle_kg = form.skeletal_muscle_kg;
            }
            updateMutation.mutate(
              { id: weight.id, body },
              {
                onSuccess: () => onClose(),
              },
            );
          }}
        >
          <Label label="体重 (kg)">
            <TextInput
              type="number"
              step="0.01"
              min={0}
              value={form.weight_kg}
              onChange={(e) =>
                setForm({ ...form, weight_kg: e.target.value })
              }
              required
            />
          </Label>
          <Label label="日時">
            <TextInput
              type="datetime-local"
              value={form.measured_at}
              onChange={(e) =>
                setForm({ ...form, measured_at: e.target.value })
              }
              required
            />
          </Label>
          <Label label="体脂肪率 (%) ※任意">
            <TextInput
              type="number"
              step="0.01"
              min={0}
              max={100}
              value={form.body_fat_percentage}
              onChange={(e) =>
                setForm({ ...form, body_fat_percentage: e.target.value })
              }
            />
          </Label>
          <Label label="骨格筋量 (kg) ※任意">
            <TextInput
              type="number"
              step="0.01"
              min={0}
              value={form.skeletal_muscle_kg}
              onChange={(e) =>
                setForm({ ...form, skeletal_muscle_kg: e.target.value })
              }
            />
          </Label>
          {updateMutation.isError && (
            <ErrorText>{(updateMutation.error as Error).message}</ErrorText>
          )}
          <div className="flex gap-2">
            <button
              type="button"
              onClick={onClose}
              disabled={updateMutation.isPending}
              className="flex-1 h-10 rounded-md border border-[var(--color-line)] text-sm text-[var(--color-ink)] hover:bg-[var(--color-surface-alt)] disabled:opacity-50"
            >
              キャンセル
            </button>
            <Button
              type="submit"
              fullWidth
              disabled={updateMutation.isPending}
            >
              {updateMutation.isPending ? "更新中…" : "更新する"}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
