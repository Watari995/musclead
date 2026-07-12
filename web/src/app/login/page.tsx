"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { useLoginMutation } from "@/features/user/api/user";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

export default function LoginPage() {
  const router = useRouter();
  const [form, setForm] = useState({ email: "", password: "" });
  const mutation = useLoginMutation();
  const t = useTranslations("login");

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold tracking-tight mb-6">{t("title")}</h1>
      <Card className="p-6">
        <form
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            mutation.mutate(form, {
              onSuccess: () => router.replace("/meals"),
            });
          }}
        >
          <Label label={t("email")}>
            <TextInput
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
              autoComplete="email"
            />
          </Label>
          <Label label={t("password")}>
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
            {mutation.isPending ? t("submitting") : t("submit")}
          </Button>
        </form>
      </Card>
      <p className="mt-6 text-sm text-[var(--color-ink-muted)] text-center">
        {t("noAccount")}{" "}
        <Link
          href="/register"
          className="text-[var(--color-ink)] font-medium hover:opacity-60"
        >
          {t("register")}
        </Link>
      </p>
    </div>
  );
}
