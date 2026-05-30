"use client";

import { useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { apiClient, type RegisterRequest } from "@/api/client";
import { setStoredUserId } from "@/lib/auth";

export default function RegisterPage() {
  const router = useRouter();
  const [form, setForm] = useState<RegisterRequest>({
    name: "",
    email: "",
    password: "",
    birthday: "",
  });

  const mutation = useMutation({
    mutationFn: async (body: RegisterRequest) => {
      const { data, error, response } = await apiClient.POST("/users", {
        body,
      });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      if (!data?.user_id) throw new Error("user_id がレスポンスに含まれません");
      return data.user_id;
    },
    onSuccess: (userId) => {
      setStoredUserId(userId);
      router.replace("/meals");
    },
  });

  return (
    <div className="max-w-md mx-auto bg-white rounded-lg shadow-sm border border-slate-200 p-6">
      <h1 className="text-xl font-bold mb-4">新規登録</h1>
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          mutation.mutate(form);
        }}
      >
        <Field
          label="名前"
          value={form.name ?? ""}
          onChange={(v) => setForm({ ...form, name: v })}
          required
        />
        <Field
          label="メールアドレス"
          type="email"
          value={form.email ?? ""}
          onChange={(v) => setForm({ ...form, email: v })}
          required
        />
        <Field
          label="パスワード"
          type="password"
          value={form.password ?? ""}
          onChange={(v) => setForm({ ...form, password: v })}
          required
        />
        <Field
          label="誕生日 (YYYY-MM-DD)"
          type="date"
          value={form.birthday ?? ""}
          onChange={(v) => setForm({ ...form, birthday: v })}
          required
        />
        {mutation.isError && (
          <p className="text-sm text-red-600">{mutation.error.message}</p>
        )}
        <button
          type="submit"
          disabled={mutation.isPending}
          className="w-full rounded bg-slate-900 text-white py-2 hover:bg-slate-700 disabled:opacity-50"
        >
          {mutation.isPending ? "登録中…" : "登録"}
        </button>
      </form>
      <p className="mt-4 text-sm text-slate-600">
        既にアカウントをお持ちですか?{" "}
        <Link href="/login" className="text-blue-600 hover:underline">
          ログイン
        </Link>
      </p>
    </div>
  );
}

function Field({
  label,
  value,
  onChange,
  type = "text",
  required,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  type?: string;
  required?: boolean;
}) {
  return (
    <label className="block">
      <span className="text-sm font-medium text-slate-700">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        className="mt-1 block w-full rounded border border-slate-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-slate-400"
      />
    </label>
  );
}
