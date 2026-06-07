"use client";

import { useTheme } from "next-themes";
import { useEffect } from "react";
import { usePreferencesQuery } from "@/features/user/api/user";
import { useAccessToken } from "@/shared/auth/access-token";

// サーバが保持している theme を取得して next-themes に反映する。
// ログインしている間はサーバ値を真実とし、 別デバイスで変更された場合も
// 次回 me 取得時に追従する。 ローカル切替は ThemePicker 側で楽観的に
// next-themes を更新しつつ PATCH を投げる。
export function ThemePreferenceSync() {
  const { token } = useAccessToken();
  const prefsQuery = usePreferencesQuery(Boolean(token));
  const { setTheme, theme } = useTheme();
  const serverTheme = prefsQuery.data?.theme;

  useEffect(() => {
    if (!serverTheme) return;
    if (serverTheme !== theme) {
      setTheme(serverTheme);
    }
    // サーバ値 (= 真実) が変わった時のみ next-themes に流す。 theme 自体は
    // 依存に入れない(無限ループ防止)。
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [serverTheme]);

  return null;
}
