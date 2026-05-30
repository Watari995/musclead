"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { apiClient, type MealDTO, type RecordMealRequest } from "@/api/client";
import { useUserId } from "@/lib/auth";

const MEALS_QUERY_KEY = ["meals"] as const;

export default function MealsPage() {
  const router = useRouter();
  const { userId, ready } = useUserId();

  useEffect(() => {
    if (ready && !userId) router.replace("/login");
  }, [ready, userId, router]);

  const query = useQuery({
    queryKey: MEALS_QUERY_KEY,
    enabled: Boolean(userId),
    queryFn: async () => {
      const { data, error, response } = await apiClient.GET("/meals", {
        params: { query: { limit: 50, offset: 0 } },
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return data?.meals ?? [];
    },
  });

  if (!ready || !userId) return null;

  return (
    <div className="space-y-6">
      <RecordMealForm />
      <section>
        <h2 className="text-lg font-bold mb-3">食事一覧</h2>
        {query.isLoading && <p className="text-sm text-slate-500">読み込み中…</p>}
        {query.isError && (
          <p className="text-sm text-red-600">{(query.error as Error).message}</p>
        )}
        {query.data && query.data.length === 0 && (
          <p className="text-sm text-slate-500">まだ食事が記録されていません。</p>
        )}
        {query.data && query.data.length > 0 && (
          <ul className="space-y-2">
            {query.data.map((m) => (
              <MealCard key={m.id} meal={m} />
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}

function MealCard({ meal }: { meal: MealDTO }) {
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
    <li className="bg-white border border-slate-200 rounded p-4 flex items-start justify-between gap-4">
      <div className="space-y-1">
        <div className="flex items-center gap-2">
          <span className="text-sm font-bold">{labelOfType(meal.meal_type)}</span>
          <span className="text-xs text-slate-500">{formatDateTime(meal.eaten_at)}</span>
        </div>
        {meal.memo && <p className="text-sm text-slate-700">{meal.memo}</p>}
        <div className="text-xs text-slate-500 flex flex-wrap gap-x-3">
          <span>カロリー: {meal.calories ?? 0} kcal</span>
          <span>P: {meal.protein_g ?? "0"}g</span>
          <span>F: {meal.fat_g ?? "0"}g</span>
          <span>C: {meal.carbohydrate_g ?? "0"}g</span>
        </div>
      </div>
      <button
        type="button"
        onClick={() => del.mutate()}
        disabled={del.isPending}
        className="text-sm text-red-600 hover:text-red-800 disabled:opacity-50"
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
    <section className="bg-white border border-slate-200 rounded p-4">
      <h2 className="text-lg font-bold mb-3">食事を記録</h2>
      <form
        className="grid grid-cols-2 gap-3"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate({
            ...form,
            eaten_at: new Date(form.eaten_at!).toISOString(),
          });
        }}
      >
        <label className="col-span-2 block">
          <span className="text-sm font-medium text-slate-700">種類</span>
          <select
            value={form.meal_type}
            onChange={(e) => setForm({ ...form, meal_type: e.target.value })}
            className="mt-1 block w-full rounded border border-slate-300 px-3 py-2"
          >
            <option value="breakfast">朝食</option>
            <option value="lunch">昼食</option>
            <option value="dinner">夕食</option>
            <option value="snack">間食</option>
          </select>
        </label>
        <label className="col-span-2 block">
          <span className="text-sm font-medium text-slate-700">日時</span>
          <input
            type="datetime-local"
            value={form.eaten_at}
            onChange={(e) => setForm({ ...form, eaten_at: e.target.value })}
            className="mt-1 block w-full rounded border border-slate-300 px-3 py-2"
            required
          />
        </label>
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
        <label className="col-span-2 block">
          <span className="text-sm font-medium text-slate-700">メモ</span>
          <textarea
            value={form.memo ?? ""}
            onChange={(e) => setForm({ ...form, memo: e.target.value })}
            rows={2}
            className="mt-1 block w-full rounded border border-slate-300 px-3 py-2"
          />
        </label>
        {mutation.isError && (
          <p className="col-span-2 text-sm text-red-600">
            {(mutation.error as Error).message}
          </p>
        )}
        <button
          type="submit"
          disabled={mutation.isPending}
          className="col-span-2 rounded bg-slate-900 text-white py-2 hover:bg-slate-700 disabled:opacity-50"
        >
          {mutation.isPending ? "記録中…" : "記録"}
        </button>
      </form>
    </section>
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
    <label className="block">
      <span className="text-sm font-medium text-slate-700">{label}</span>
      <input
        type="number"
        step={step}
        min={0}
        value={value ?? 0}
        onChange={(e) => onChange(Number(e.target.value))}
        className="mt-1 block w-full rounded border border-slate-300 px-3 py-2"
      />
    </label>
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
    case "breakfast":
      return "朝食";
    case "lunch":
      return "昼食";
    case "dinner":
      return "夕食";
    case "snack":
      return "間食";
    default:
      return t ?? "";
  }
}
