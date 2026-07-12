import 'package:flutter/material.dart';

import 'app_tokens.dart';

/// brightness と accent から [ThemeData] を構築する。
///
/// 配色は手描き（Excalidraw 風）の紙・インク基調 + 単一 accent。Material コンポーネント
/// テーマは最小限にとどめ、見た目は core/widgets のカスタム widget で統一する
/// （Flutter 側 component-theme API の変更に強くするため）。
ThemeData buildAppTheme(Brightness brightness, Color accent) {
  final dark = brightness == Brightness.dark;
  final tokens = AppTokens.build(brightness, accent);

  final scheme = ColorScheme.fromSeed(seedColor: accent, brightness: brightness)
      .copyWith(
        primary: accent,
        // モノクロ既定では accent が ink/paper と同色になりうるため、
        // 常に白固定ではなく paper/white を brightness で切り替えてコントラストを保つ。
        onPrimary: dark ? tokens.paper : Colors.white,
        surface: tokens.paper,
        onSurface: tokens.ink,
        surfaceContainerHighest: tokens.paperAlt,
        outline: tokens.border,
      );

  final base = ThemeData(
    useMaterial3: true,
    brightness: brightness,
    colorScheme: scheme,
  );

  // 本文/UI テキストは Architects Daughter を既定にし、和文は Yomogi にフォールバック
  // させる。見出し等で手描き感を強めたい箇所は個別に `fontFamily: 'Caveat'` を指定する。
  return base.copyWith(
    scaffoldBackgroundColor: tokens.paper,
    textTheme: base.textTheme.apply(
      fontFamily: 'Architects Daughter',
      fontFamilyFallback: const ['Yomogi'],
    ),
    primaryTextTheme: base.primaryTextTheme.apply(
      fontFamily: 'Architects Daughter',
      fontFamilyFallback: const ['Yomogi'],
    ),
    extensions: <ThemeExtension<dynamic>>[tokens],
  );
}
