"use client";

import { useState } from "react";
import type { RecordMealRequest } from "@/shared/api/client";
import { useRecordMealMutation } from "@/features/meal/api/meals";
import { toLocalInput } from "@/features/meal/model/meal";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

export function RecordMealForm() {
  const [form, setForm] = useState<RecordMealRequest>(initialForm);
  const mutation = useRecordMealMutation();

  return (
    <Card className="p-5">
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate(
            {
              ...form,
              eaten_at: new Date(form.eaten_at!).toISOString(),
            },
            { onSuccess: () => setForm(initialForm()) },
          );
        }}
      >
        <Label label="種類">
          <select
            value={form.meal_type}
            onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
            className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-white focus:outline-none focus:border-[var(--color-ink)]"
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
            className="block w-full px-3 py-2 rounded-md border border-[var(--color-line)] bg-white focus:outline-none focus:border-[var(--color-ink)]"
          />
        </Label>
        {mutation.isError && (
          <ErrorText>{(mutation.error as Error).message}</ErrorText>
        )}
        <Button type="submit" fullWidth disabled={mutation.isPending}>
          {mutation.isPending ? "記録中…" : "記録する"}
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
  onChange: (v: number) => void;
  step?: string;
}) {
  return (
    <Label label={label}>
      <TextInput
        type="number"
        step={step}
        min={0}
        value={value ?? 0}
        onChange={(e) => onChange(Number(e.target.value))}
      />
    </Label>
  );
}

function initialForm(): RecordMealRequest {
  return {
    meal_type: "breakfast",
    eaten_at: toLocalInput(new Date()),
    calories: 0,
    protein_g: 0,
    fat_g: 0,
    carbohydrate_g: 0,
    memo: "",
  };
}
