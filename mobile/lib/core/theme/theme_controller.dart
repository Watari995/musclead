import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../providers/core_providers.dart';
import 'app_tokens.dart';

/// 設定画面で選べるアクセント候補（プレビューのスウォッチ準拠）。
const List<Color> kAccentPresets = [
  kBrandAccent, // ブランド赤
  Color(0xFF0A84FF), // ブルー
  Color(0xFF00BCD4), // ターコイズ
  Color(0xFF30C759), // グリーン
  Color(0xFFFF9F0A), // オレンジ
  Color(0xFFBF5AF2), // パープル
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
}
