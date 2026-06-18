"use client";

import { useState } from "react";
import {
  useMealTemplatesQuery,
  useCreateMealTemplateMutation,
  useDeleteMealTemplateMutation,
} from "@/features/meal/api/meal_templates";
import { mealTypeLabel } from "@/features/meal/model/meal";
import type { MealTemplate } from "@/features/meal/model/meal_template";
import { Card, SectionTitle, ErrorText, Button, Label, TextInput, NumberField } from "@/shared/ui";

type Props = {
  onSelect: (t: MealTemplate) => void;
};

export function MealTemplateSection({ onSelect }: Props) {
  const query = useMealTemplatesQuery();
  const [showForm, setShowForm] = useState(false);

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <SectionTitle>食事テンプレート</SectionTitle>
        <button
          type="button"
          onClick={() => setShowForm((v) => !v)}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] underline"
        >
          {showForm ? "キャンセル" : "+ 新規"}
        </button>
      </div>

      {showForm && (
        <CreateTemplateForm onCreated={() => setShowForm(false)} />
      )}

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {query.data && query.data.length === 0 && !showForm && (
        <Card className="p-4 text-center text-sm text-[var(--color-ink-muted)]">
          テンプレートはまだありません
        </Card>
      )}
      {query.data && query.data.length > 0 && (
        <div className="grid gap-2">
          {query.data.map((t) => (
            <TemplateCard key={t.id} template={t} onSelect={onSelect} />
          ))}
        </div>
      )}
    </div>
  );
}

function TemplateCard({
  template,
  onSelect,
}: {
  template: MealTemplate;
  onSelect: (t: MealTemplate) => void;
}) {
  const del = useDeleteMealTemplateMutation();

  return (
    <Card className="p-3 flex items-center gap-3 hover:bg-[var(--color-surface-alt)] transition-colors">
      <button
        type="button"
        onClick={() => onSelect(template)}
        className="flex-1 text-left min-w-0"
      >
        <div className="font-medium text-sm truncate">{template.name}</div>
        <div className="flex gap-3 mt-0.5 text-xs text-[var(--color-ink-muted)]">
          <span>{mealTypeLabel(template.mealType)}</span>
          <span className="font-medium text-[var(--color-ink)]">{template.calories} kcal</span>
          <span>P {template.proteinG}g</span>
          <span>F {template.fatG}g</span>
          <span>C {template.carbohydrateG}g</span>
        </div>
      </button>
      <button
        type="button"
        disabled={del.isPending}
        onClick={() => {
          if (confirm(`「${template.name}」を削除しますか?`)) {
            del.mutate(template.id);
          }
        }}
        className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] shrink-0 disabled:opacity-50"
      >
        削除
      </button>
    </Card>
  );
}

type FormState = {
  name: string;
  meal_type: string;
  calories: number | undefined;
  protein_g: number | undefined;
  fat_g: number | undefined;
  carbohydrate_g: number | undefined;
};

function CreateTemplateForm({ onCreated }: { onCreated: () => void }) {
  const [form, setForm] = useState<FormState>({
    name: "",
    meal_type: "breakfast",
    calories: undefined,
    protein_g: undefined,
    fat_g: undefined,
    carbohydrate_g: undefined,
  });
  const mutation = useCreateMealTemplateMutation();

  const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    mutation.mutate(
      {
        name: form.name,
        meal_type: form.meal_type,
        calories: form.calories ?? 0,
        protein_g: form.protein_g,
        fat_g: form.fat_g,
        carbohydrate_g: form.carbohydrate_g,
      },
      { onSuccess: onCreated },
    );
  };

  return (
    <Card className="p-4">
      <form className="space-y-3" onSubmit={handleSubmit}>
        <Label label="テンプレート名">
          <TextInput
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="例: プロテインシェイク"
            required
          />
        </Label>
        <Label label="種類">
          <select
            value={form.meal_type}
            onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
            className="block w-full h-11 px-3 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] focus:outline-none focus:border-[var(--color-ink)]"
          >
            <option value="breakfast">朝食</option>
            <option value="lunch">昼食</option>
            <option value="dinner">夕食</option>
            <option value="snack">間食</option>
          </select>
        </Label>
        <div className="grid grid-cols-2 gap-3">
          <Label label="カロリー (kcal)">
            <NumberField
              min={0}
              placeholder="0"
              value={form.calories}
              onChange={(v) => setForm({ ...form, calories: v })}
            />
          </Label>
          <Label label="タンパク質 (g)">
            <NumberField
              step="0.1"
              min={0}
              placeholder="0"
              value={form.protein_g}
              onChange={(v) => setForm({ ...form, protein_g: v })}
            />
          </Label>
          <Label label="脂質 (g)">
            <NumberField
              step="0.1"
              min={0}
              placeholder="0"
              value={form.fat_g}
              onChange={(v) => setForm({ ...form, fat_g: v })}
            />
          </Label>
          <Label label="炭水化物 (g)">
            <NumberField
              step="0.1"
              min={0}
              placeholder="0"
              value={form.carbohydrate_g}
              onChange={(v) => setForm({ ...form, carbohydrate_g: v })}
            />
          </Label>
        </div>
        {mutation.isError && (
          <ErrorText>{(mutation.error as Error).message}</ErrorText>
        )}
        <Button type="submit" fullWidth disabled={mutation.isPending}>
          {mutation.isPending ? "保存中…" : "保存する"}
        </Button>
      </form>
    </Card>
  );
}
