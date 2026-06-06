"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import {
  useMeQuery,
  useUpdateUserMutation,
  type UpdateUserBody,
} from "@/features/user/api/user";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  Button,
  Card,
  ErrorText,
  Label,
  SectionTitle,
  TextInput,
} from "@/shared/ui";

export default function ProfilePage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const meQuery = useMeQuery(Boolean(token));

  if (!ready || !token) return null;
  if (meQuery.isLoading) {
    return <p className="text-sm text-[var(--color-ink-muted)]">読み込み中…</p>;
  }
  if (meQuery.isError) {
    return <ErrorText>{(meQuery.error as Error).message}</ErrorText>;
  }
  if (!meQuery.data) return null;

  return (
    <ProfileForm
      initialName={meQuery.data.name ?? ""}
      initialBirthday={meQuery.data.birthday ?? ""}
      email={meQuery.data.email ?? ""}
      onCancel={() => router.back()}
    />
  );
}

function ProfileForm({
  initialName,
  initialBirthday,
  email,
  onCancel,
}: {
  initialName: string;
  initialBirthday: string;
  email: string;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initialName);
  const [birthday, setBirthday] = useState(initialBirthday);
  const [done, setDone] = useState(false);
  const mutation = useUpdateUserMutation();

  // 差分を計算: 変わってないキーは body に含めない(= サーバーは未送信扱いで更新しない)
  const body = buildBody(initialName, initialBirthday, name, birthday);
  const hasChange = Object.keys(body).length > 0;
  const nameInvalid = name.trim().length === 0;

  return (
    <div className="space-y-6">
      <SectionTitle>プロフィール</SectionTitle>

      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          if (!hasChange || nameInvalid) return;
          mutation.mutate(body, {
            onSuccess: () => {
              setDone(true);
            },
          });
        }}
      >
        <Card className="p-5 space-y-4">
          <Label label="メールアドレス">
            <TextInput value={email} disabled readOnly />
          </Label>
          <p className="text-xs text-[var(--color-ink-muted)] -mt-2">
            メールアドレスは現状変更できません
          </p>

          <Label label="名前">
            <TextInput
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              maxLength={50}
              disabled={mutation.isPending}
            />
          </Label>

          <Label label="誕生日(空欄で削除)">
            <TextInput
              type="date"
              value={birthday}
              onChange={(e) => setBirthday(e.target.value)}
              disabled={mutation.isPending}
            />
          </Label>
        </Card>

        {mutation.isError && (
          <ErrorText>{(mutation.error as Error).message}</ErrorText>
        )}
        {done && !mutation.isPending && !mutation.isError && (
          <p className="text-sm text-[var(--color-ink-muted)]">保存しました</p>
        )}

        <div className="flex gap-3">
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={mutation.isPending}
          >
            キャンセル
          </Button>
          <Button
            type="submit"
            disabled={mutation.isPending || !hasChange || nameInvalid}
            className="flex-1"
          >
            {mutation.isPending ? "保存中…" : "保存"}
          </Button>
        </div>
      </form>
    </div>
  );
}

// 変更されたキーだけを含む PATCH body を生成。
// birthday は空欄 → null(クリア)、 値あり → 文字列。
// 元と同じならキー自体を含めない(サーバーは「未送信」 と判定し更新しない)。
function buildBody(
  initialName: string,
  initialBirthday: string,
  name: string,
  birthday: string,
): UpdateUserBody {
  const body: UpdateUserBody = {};
  if (name !== initialName) {
    body.name = name;
  }
  if (birthday !== initialBirthday) {
    body.birthday = birthday === "" ? null : birthday;
  }
  return body;
}
