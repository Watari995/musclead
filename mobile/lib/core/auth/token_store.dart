import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// access token を Keychain（iOS）に保存する。
/// refresh token は HttpOnly Cookie として cookie_jar 側で永続化する。
class TokenStore {
  TokenStore(this._storage);

  final FlutterSecureStorage _storage;
  static const _accessKey = 'access_token';

  Future<String?> readAccessToken() => _storage.read(key: _accessKey);

  Future<void> writeAccessToken(String token) =>
      _storage.write(key: _accessKey, value: token);

  Future<void> clear() => _storage.delete(key: _accessKey);
}
