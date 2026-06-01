"use client";

import { useMutation } from "@tanstack/react-query";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { loginRequest } from "@/api/auth";
import { setAccessToken } from "@/lib/access-token";
import { Button, Card, ErrorText, Label, TextInput } from "@/components/ui";

export default function LoginPage() {
  const router = useRouter();
  const [form, setForm] = useState({ email: "", password: "" });

  const mutation = useMutation({
    mutationFn: async (input: { email: string; password: string }) =>
      loginRequest(input.email, input.password),
    onSuccess: (tokens) => {
      setAccessToken(tokens.access_token);
      router.replace("/meals");
    },
  });

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold tracking-tight mb-6">ログイン</h1>
      <Card className="p-6">
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            mutation.mutate(form);
          }}
        >
          <Label label="メールアドレス">
            <TextInput
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              autoComplete="email"
            />
          </Label>
          <Label label="パスワード">
            <TextInput
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              required
              autoComplete="current-password"
            />
          </Label>
          {mutation.isError && <ErrorText>{mutation.error.message}</ErrorText>}
          <Button type="submit" fullWidth disabled={mutation.isPending}>
            {mutation.isPending ? "ログイン中…" : "ログイン"}
          </Button>
        </form>
      </Card>
      <p className="mt-6 text-sm text-[var(--color-ink-muted)] text-center">
        アカウントをお持ちでないですか?{" "}
        <Link
          href="/register"
          className="text-[var(--color-ink)] font-medium hover:opacity-60"
        >
          新規登録
        </Link>
      </p>
    </div>
  );
}
