"use client";

import { useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { loginRequest } from "@/api/auth";
import { apiClient, type RegisterRequest } from "@/api/client";
import { setAccessToken } from "@/lib/access-token";
import { Button, Card, ErrorText, Label, TextInput } from "@/components/ui";

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
      const { error, response } = await apiClient.POST("/users", { body });
      if (error) throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      return loginRequest(body.email ?? "", body.password ?? "");
    },
    onSuccess: (tokens) => {
      setAccessToken(tokens.access_token);
      router.replace("/meals");
    },
  });

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold tracking-tight mb-6">新規登録</h1>
      <Card className="p-6">
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            mutation.mutate(form);
          }}
        >
          <Label label="名前">
            <TextInput
              value={form.name ?? ""}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
              autoComplete="name"
            />
          </Label>
          <Label label="メールアドレス">
            <TextInput
              type="email"
              value={form.email ?? ""}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              autoComplete="email"
            />
          </Label>
          <Label label="パスワード">
            <TextInput
              type="password"
              value={form.password ?? ""}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              required
              autoComplete="new-password"
            />
          </Label>
          <Label label="誕生日">
            <TextInput
              type="date"
              value={form.birthday ?? ""}
              onChange={(e) => setForm({ ...form, birthday: e.target.value })}
              required
            />
          </Label>
          {mutation.isError && <ErrorText>{mutation.error.message}</ErrorText>}
          <Button type="submit" fullWidth disabled={mutation.isPending}>
            {mutation.isPending ? "登録中…" : "登録する"}
          </Button>
        </form>
      </Card>
      <p className="mt-6 text-sm text-[var(--color-ink-muted)] text-center">
        既にアカウントをお持ちですか?{" "}
        <Link
          href="/login"
          className="text-[var(--color-ink)] font-medium hover:opacity-60"
        >
          ログイン
        </Link>
      </p>
    </div>
  );
}
