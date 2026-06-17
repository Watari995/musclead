import 'package:flutter/material.dart';

/// 既定のアクセント（ブランド赤）。実行時に [accentProvider] で差し替え可能。
const Color kBrandAccent = Color(0xFFEF3A3A);

/// デザインプレビュー（mobile/preview）と 1:1 対応するデザイントークン。
///
/// アクセントは 1 つの [accent] から派生（[accentWeak]）。
/// PFC は意味を持つデータ配色、自己ベストは [gold]。
/// ガラス系トークンは操作レイヤー（タブ/ナビ/FAB）専用。
@immutable
class AppTokens extends ThemeExtension<AppTokens> {
  const AppTokens({
    required this.accent,
    required this.accentWeak,
    required this.gold,
    required this.goldWeak,
    required this.macroP,
    required this.macroF,
    required this.macroC,
    required this.border,
    required this.hairline,
    required this.muted,
    required this.subtle,
    required this.glassTint,
    required this.glassBorder,
    required this.glassHighlight,
  });

  final Color accent;
  final Color accentWeak;
  final Color gold;
  final Color goldWeak;
  final Color macroP;
  final Color macroF;
  final Color macroC;
  final Color border;
  final Color hairline;
  final Color muted;
  final Color subtle;
  final Color glassTint;
  final Color glassBorder;
  final Color glassHighlight;

  /// brightness と accent から実トークンを組み立てる。
  factory AppTokens.build(Brightness brightness, Color accent) {
    final dark = brightness == Brightness.dark;
    return AppTokens(
      accent: accent,
      accentWeak: accent.withValues(alpha: dark ? 0.22 : 0.14),
      gold: dark ? const Color(0xFFE3B341) : const Color(0xFFC8861A),
      goldWeak: (dark ? const Color(0xFFE3B341) : const Color(0xFFC8861A))
          .withValues(alpha: 0.16),
      macroP: dark ? const Color(0xFF6EA0FF) : const Color(0xFF4A7FE0),
      macroF: dark ? const Color(0xFFE6B552) : const Color(0xFFD99528),
      macroC: dark ? const Color(0xFF5FCE82) : const Color(0xFF3FA364),
      border: dark
          ? Colors.white.withValues(alpha: 0.14)
          : Colors.black.withValues(alpha: 0.10),
      hairline: dark
          ? Colors.white.withValues(alpha: 0.10)
          : Colors.black.withValues(alpha: 0.08),
      muted: dark ? const Color(0xFF98989D) : const Color(0xFF6B6B70),
      subtle: dark ? const Color(0xFF5B5B60) : const Color(0xFFA0A0A5),
      glassTint: dark
          ? const Color(0xFF2C2C2E).withValues(alpha: 0.5)
          : Colors.white.withValues(alpha: 0.55),
      glassBorder: dark
          ? Colors.white.withValues(alpha: 0.16)
          : Colors.white.withValues(alpha: 0.70),
      glassHighlight: dark
          ? Colors.white.withValues(alpha: 0.18)
          : Colors.white.withValues(alpha: 0.85),
    );
  }

  @override
  AppTokens copyWith({
    Color? accent,
    Color? accentWeak,
    Color? gold,
    Color? goldWeak,
    Color? macroP,
    Color? macroF,
    Color? macroC,
    Color? border,
    Color? hairline,
    Color? muted,
    Color? subtle,
    Color? glassTint,
    Color? glassBorder,
    Color? glassHighlight,
  }) {
    return AppTokens(
      accent: accent ?? this.accent,
      accentWeak: accentWeak ?? this.accentWeak,
      gold: gold ?? this.gold,
      goldWeak: goldWeak ?? this.goldWeak,
      macroP: macroP ?? this.macroP,
      macroF: macroF ?? this.macroF,
      macroC: macroC ?? this.macroC,
      border: border ?? this.border,
      hairline: hairline ?? this.hairline,
      muted: muted ?? this.muted,
      subtle: subtle ?? this.subtle,
      glassTint: glassTint ?? this.glassTint,
      glassBorder: glassBorder ?? this.glassBorder,
      glassHighlight: glassHighlight ?? this.glassHighlight,
    );
  }

  @override
  AppTokens lerp(ThemeExtension<AppTokens>? other, double t) {
    if (other is! AppTokens) return this;
    return AppTokens(
      accent: Color.lerp(accent, other.accent, t)!,
      accentWeak: Color.lerp(accentWeak, other.accentWeak, t)!,
      gold: Color.lerp(gold, other.gold, t)!,
      goldWeak: Color.lerp(goldWeak, other.goldWeak, t)!,
      macroP: Color.lerp(macroP, other.macroP, t)!,
      macroF: Color.lerp(macroF, other.macroF, t)!,
      macroC: Color.lerp(macroC, other.macroC, t)!,
      border: Color.lerp(border, other.border, t)!,
      hairline: Color.lerp(hairline, other.hairline, t)!,
      muted: Color.lerp(muted, other.muted, t)!,
      subtle: Color.lerp(subtle, other.subtle, t)!,
      glassTint: Color.lerp(glassTint, other.glassTint, t)!,
      glassBorder: Color.lerp(glassBorder, other.glassBorder, t)!,
      glassHighlight: Color.lerp(glassHighlight, other.glassHighlight, t)!,
    );
  }
}

/// `Theme.of(context).tokens` で参照できるショートカット。
extension AppTokensX on ThemeData {
  AppTokens get tokens => extension<AppTokens>()!;
}

extension AppTokensContextX on BuildContext {
  AppTokens get tokens => Theme.of(this).extension<AppTokens>()!;
  ColorScheme get colors => Theme.of(this).colorScheme;
}
