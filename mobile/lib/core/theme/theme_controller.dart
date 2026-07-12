import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/core_providers.dart';
import 'app_tokens.dart';

/// 設定画面で選べるアクセント候補（プレビューのスウォッチ準拠）。
/// Excalidraw のストロークパレットに合わせた 5 色。既定はモノクロ（[kBrandAccent] = ink）。
const List<Color> kAccentPresets = [
  kBrandAccent, // モノクロ(既定)
  Color(0xFFE03131), // レッド
  Color(0xFF1971C2), // ブルー
  Color(0xFF2F9E44), // グリーン
  Color(0xFFF08C00), // オレンジ
];

/// アクセントカラー。1 トークンで全テーマを駆動し、secure storage に永続化する。
final accentProvider = NotifierProvider<AccentController, Color>(
  AccentController.new,
);

class AccentController extends Notifier<Color> {
  static const _key = 'accent_color';

  @override
  Color build() {
    _load();
    return kBrandAccent;
  }

  Future<void> _load() async {
    final v = await ref.read(secureStorageProvider).read(key: _key);
    final parsed = v == null ? null : int.tryParse(v);
    if (parsed != null) state = Color(parsed);
  }

  Future<void> set(Color color) async {
    state = color;
    await ref
        .read(secureStorageProvider)
        .write(key: _key, value: color.toARGB32().toString());
  }
}

/// ライト/ダーク/システム。ユーザー設定（preferences.theme）と同期させる。
final themeModeProvider = NotifierProvider<ThemeModeController, ThemeMode>(
  ThemeModeController.new,
);

class ThemeModeController extends Notifier<ThemeMode> {
  @override
  ThemeMode build() => ThemeMode.system;

  void set(ThemeMode mode) => state = mode;

  /// サーバーの preferences.theme（'system' / 'light' / 'dark'）から復元する。
  /// 外観はローカル保存しないため、起動時に必ずサーバー値へ同期する。
  void hydrate(String theme) {
    state = switch (theme) {
      'light' => ThemeMode.light,
      'dark' => ThemeMode.dark,
      _ => ThemeMode.system,
    };
  }
}
