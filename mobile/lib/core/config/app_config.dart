import 'package:flutter/foundation.dart';

/// アプリの実行環境設定（flavor）。
///
/// FLAVOR 未指定なら release ビルド=prod / debug=dev（Xcode Archive はそのまま本番）。
/// `--dart-define=FLAVOR=prod|dev` で明示切替、`API_BASE_URL=...` で base URL 上書き。
enum Flavor { dev, prod }

class AppConfig {
  const AppConfig({required this.flavor, required this.apiBaseUrl});

  final Flavor flavor;
  final String apiBaseUrl;

  bool get isProd => flavor == Flavor.prod;

  static const String _flavorEnv = String.fromEnvironment('FLAVOR');
  static const String _baseUrlEnv = String.fromEnvironment('API_BASE_URL');

  static final AppConfig current = _resolve();

  static AppConfig _resolve() {
    final isProd = _flavorEnv.isEmpty ? kReleaseMode : _flavorEnv == 'prod';
    final defaultBase = isProd
        ? 'https://api.musclead.com'
        : 'http://localhost:8080';
    return AppConfig(
      flavor: isProd ? Flavor.prod : Flavor.dev,
      apiBaseUrl: _baseUrlEnv.isNotEmpty ? _baseUrlEnv : defaultBase,
    );
  }
}
