import 'package:flutter/material.dart';

/// 既定のアクセント（モノクロ = ink と同色）。実行時に [accentProvider] で差し替え可能。
const Color kBrandAccent = Color(0xFF1E1E1E);

/// デザインプレビュー（mobile/preview）と 1:1 対応するデザイントークン。
///
/// 手描き（Excalidraw 風）の紙（paper）・インク（ink）基調 + 単一 [accent]。
/// アクセントは 1 つの [accent] から派生（[accentWeak]）。
/// PFC・自己ベスト（[gold]）は意味を持つデータ配色として accent から独立。
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
    required this.paper,
    required this.paperAlt,
    required this.ink,
    required this.border,
    required this.hairline,
    required this.muted,
    required this.subtle,
  });

  final Color accent;
  final Color accentWeak;
  final Color gold;
  final Color goldWeak;
  final Color macroP;
  final Color macroF;
  final Color macroC;

  /// 紙の背景色（canvas/card 共通）。カードは [paper] 塗り + [border] の手描き輪郭で区別する。
  final Color paper;

  /// [paper] よりわずかに濃い代替面（アバター・写真プレースホルダ等）。
  final Color paperAlt;

  /// 本文・アイコンのインク色。
  final Color ink;

  /// 手描き輪郭（RoughBox）のストローク色。
  final Color border;
  final Color hairline;
  final Color muted;
  final Color subtle;

  /// brightness と accent から実トークンを組み立てる。
  ///
  /// [kBrandAccent](モノクロ既定)は brightness に応じて ink 自体に解決する。
  /// 固定 hex のままだとダークモードで「ink と同じ暗色の accent」が
  /// 暗い紙の上に沈んで見えなくなるため。
  factory AppTokens.build(Brightness brightness, Color accent) {
    final dark = brightness == Brightness.dark;
    final ink = dark ? const Color(0xFFE8E6DE) : const Color(0xFF1E1E1E);
    final resolvedAccent = accent.toARGB32() == kBrandAccent.toARGB32()
        ? ink
        : accent;
    return AppTokens(
      accent: resolvedAccent,
      accentWeak: resolvedAccent.withValues(alpha: dark ? 0.22 : 0.14),
      gold: dark ? const Color(0xFFE3B341) : const Color(0xFFC8861A),
      goldWeak: (dark ? const Color(0xFFE3B341) : const Color(0xFFC8861A))
          .withValues(alpha: 0.16),
      macroP: dark ? const Color(0xFF5AA7F6) : const Color(0xFF1971C2),
      macroF: dark ? const Color(0xFFF3A94B) : const Color(0xFFE8830A),
      macroC: dark ? const Color(0xFF5ED676) : const Color(0xFF2F9E44),
      paper: dark ? const Color(0xFF17181C) : const Color(0xFFFAFAF8),
      paperAlt: dark ? const Color(0xFF212227) : const Color(0xFFEEF0EA),
      ink: ink,
      border: ink,
      hairline: dark ? const Color(0xFF37383D) : const Color(0xFFC7C7BD),
      muted: dark ? const Color(0xFF8D8E88) : const Color(0xFF68685F),
      subtle: dark ? const Color(0xFF6B6C66) : const Color(0xFFA4A49C),
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
    Color? paper,
    Color? paperAlt,
    Color? ink,
    Color? border,
    Color? hairline,
    Color? muted,
    Color? subtle,
  }) {
    return AppTokens(
      accent: accent ?? this.accent,
      accentWeak: accentWeak ?? this.accentWeak,
      gold: gold ?? this.gold,
      goldWeak: goldWeak ?? this.goldWeak,
      macroP: macroP ?? this.macroP,
      macroF: macroF ?? this.macroF,
      macroC: macroC ?? this.macroC,
      paper: paper ?? this.paper,
      paperAlt: paperAlt ?? this.paperAlt,
      ink: ink ?? this.ink,
      border: border ?? this.border,
      hairline: hairline ?? this.hairline,
      muted: muted ?? this.muted,
      subtle: subtle ?? this.subtle,
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
      paper: Color.lerp(paper, other.paper, t)!,
      paperAlt: Color.lerp(paperAlt, other.paperAlt, t)!,
      ink: Color.lerp(ink, other.ink, t)!,
      border: Color.lerp(border, other.border, t)!,
      hairline: Color.lerp(hairline, other.hairline, t)!,
      muted: Color.lerp(muted, other.muted, t)!,
      subtle: Color.lerp(subtle, other.subtle, t)!,
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
