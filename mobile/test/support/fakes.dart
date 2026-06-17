import 'package:musclead/core/auth/token_store.dart';

/// テスト用のインメモリ TokenStore（Keychain を使わない）。
class FakeTokenStore implements TokenStore {
  String? _token;

  @override
  Future<String?> readAccessToken() async => _token;

  @override
  Future<void> writeAccessToken(String token) async => _token = token;

  @override
  Future<void> clear() async => _token = null;
}
