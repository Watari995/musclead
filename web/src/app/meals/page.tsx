"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { apiClient, type MealDTO, type RecordMealRequest } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/shared/ui";

const MEALS_QUERY_KEY = ["meals"] as const;

export default function MealsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useQuery({
    queryKey: MEALS_QUERY_KEY,
    enabled: Boolean(token),
    queryFn: async () => {
      const { data, error, response } = await apiClient.GET("/meals", {
        params: { query: { limit: 50, offset: 0 } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data?.meals ?? [];
    },
  });

  if (!ready || !token) return null;

  return (
    <div className="grid gap-8 lg:grid-cols-[1fr_360px]">
      <section>
        <SectionTitle>食事一覧</SectionTitle>
        {query.isLoading && (
          <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>
        )}
        {query.isError && <ErrorText>{(query.error as Error).message}</ErrorText>}
        {query.data && query.data.length === 0 && (
          <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
            まだ食事が記録されていません。
          </Card>
        )}
        {query.data && query.data.length > 0 && (
          <ul className="divide-y divide-[var(--color-line)] border border-[var(--color-line)] rounded-lg overflow-hidden bg-white">
            {query.data.map((m) => (
              <MealRow key={m.id} meal={m} />
            ))}
          </ul>
        )}
      </section>
      <aside>
        <SectionTitle>食事を記録</SectionTitle>
        <RecordMealForm />
      </aside>
    </div>
  );
}

function MealRow({ meal }: { meal: MealDTO }) {
  const queryClient = useQueryClient();
  const del = useMutation({
    mutationFn: async () => {
      const { error, response } = await apiClient.DELETE("/meals/{id}", {
        params: { path: { id: meal.id! } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: MEALS_QUERY_KEY }),
  });

  return (
    <li className="p-4 flex items-start gap-4">
      <div className="w-14 h-14 shrink-0 rounded-md bg-[var(--color-surface-alt)] flex items-center justify-center text-xl">
        {emojiOfType(meal.meal_type)}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {labelOfType(meal.meal_type)}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {formatDateTime(meal.eaten_at)}
          </span>
        </div>
        {meal.memo && (
          <p className="mt-1 text-sm text-[var(--color-ink)] line-clamp-2">
            {meal.memo}
          </p>
        )}
        <div className="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-[var(--color-ink-muted)]">
          <span className="font-medium text-[var(--color-ink)]">
            {meal.calories ?? 0} kcal
          </span>
          <span>P {meal.protein_g ?? "0"}g</span>
          <span>F {meal.fat_g ?? "0"}g</span>
          <span>C {meal.carbohydrate_g ?? "0"}g</span>
        </div>
      </div>
      <button
        type="button"
        onClick={() => del.mutate()}
        disabled={del.isPending}
        className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] shrink-0"
      >
        削除
      </button>
    </li>
  );
}

function RecordMealForm() {
  const queryClient = useQueryClient();
  const [form, setForm] = useState<RecordMealRequest>(initialForm());

  const mutation = useMutation({
    mutationFn: async (body: RecordMealRequest) => {
      const { data, error, response } = await apiClient.POST("/meals", { body });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data;
    },
    onSuccess: () => {
      setForm(initialForm());
      queryClient.invalidateQueries({ queryKey: MEALS_QUERY_KEY });
    },
  });

  return (
    <Card className="p-5">
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate({
            ...form,
            eaten_at: new Date(form.eaten_at!).toISOString(),
          });
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

function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function formatDateTime(iso?: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  return d.toLocaleString("ja-JP", { dateStyle: "short", timeStyle: "short" });
}

function labelOfType(t?: string): string {
  switch (t) {
    case "breakfast": return "朝食";
    case "lunch": return "昼食";
    case "dinner": return "夕食";
    case "snack": return "間食";
    default: return t ?? "";
  }
}

function emojiOfType(t?: string): string {
  switch (t) {
    case "breakfast": return "🍳";
    case "lunch": return "🍱";
    case "dinner": return "🍽️";
    case "snack": return "🍎";
    default: return "🍴";
  }
}
