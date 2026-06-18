import 'package:flutter/material.dart';

import 'app_tokens.dart';

/// brightness と accent から [ThemeData] を構築する。
///
/// 配色はニュートラル基調 + 単一 accent。Material コンポーネントテーマは
/// 最小限にとどめ、見た目は core/widgets のカスタム widget で統一する
/// （Flutter 側 component-theme API の変更に強くするため）。
ThemeData buildAppTheme(Brightness brightness, Color accent) {
  final dark = brightness == Brightness.dark;
  final tokens = AppTokens.build(brightness, accent);

  final scheme = ColorScheme.fromSeed(seedColor: accent, brightness: brightness)
      .copyWith(
        primary: accent,
        onPrimary: Colors.white,
        surface: dark ? const Color(0xFF1C1C1E) : Colors.white,
        onSurface: dark ? const Color(0xFFF5F5F7) : const Color(0xFF1C1C1E),
        outline: tokens.border,
      );

  return ThemeData(
    useMaterial3: true,
    brightness: brightness,
    colorScheme: scheme,
    scaffoldBackgroundColor: dark ? Colors.black : const Color(0xFFF2F2F2),
    extensions: <ThemeExtension<dynamic>>[tokens],
  );
}
