/// アプリの実行環境設定（flavor）。
///
/// 切替は `--dart-define=FLAVOR=prod`（既定 dev）。
/// base URL は `--dart-define=API_BASE_URL=...` で上書き可能。
enum Flavor { dev, prod }

class AppConfig {
  const AppConfig({required this.flavor, required this.apiBaseUrl});

  final Flavor flavor;
  final String apiBaseUrl;

  bool get isProd => flavor == Flavor.prod;

  static const String _flavorEnv = String.fromEnvironment('FLAVOR', defaultValue: 'dev');
  static const String _baseUrlEnv = String.fromEnvironment('API_BASE_URL');

  static final AppConfig current = _resolve();

  static AppConfig _resolve() {
    final isProd = _flavorEnv == 'prod';
    final defaultBase = isProd ? 'https://api.musclead.com' : 'http://localhost:8080';
    return AppConfig(
      flavor: isProd ? Flavor.prod : Flavor.dev,
      apiBaseUrl: _baseUrlEnv.isNotEmpty ? _baseUrlEnv : defaultBase,
    );
  }
}
