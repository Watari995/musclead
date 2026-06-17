import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app_tokens.dart';

/// アクセントカラー。1 つのソースから全テーマを駆動（後から差し替え可能）。
/// 既定はブランド赤。`ref.read(accentProvider.notifier).set(color)` で全画面に即反映。
final accentProvider = NotifierProvider<AccentController, Color>(
  AccentController.new,
);

class AccentController extends Notifier<Color> {
  @override
  Color build() => kBrandAccent;

  void set(Color color) => state = color;
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
