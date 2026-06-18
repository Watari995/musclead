import 'package:cookie_jar/cookie_jar.dart';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '../api/api_client.dart';
import '../api/auth_interceptor.dart';
import '../auth/auth_signal.dart';
import '../auth/token_store.dart';

final secureStorageProvider = Provider<FlutterSecureStorage>(
  (ref) => const FlutterSecureStorage(),
);

final tokenStoreProvider = Provider<TokenStore>(
  (ref) => TokenStore(ref.watch(secureStorageProvider)),
);

final authSignalProvider = Provider<AuthSignal>((ref) {
  final signal = AuthSignal();
  ref.onDispose(signal.dispose);
  return signal;
});

/// 実体は main() で PersistCookieJar（永続）に override する。
final cookieJarProvider = Provider<CookieJar>(
  (ref) => throw UnimplementedError(
    'cookieJarProvider must be overridden in main()',
  ),
);

/// 認証付き Dio。refresh と失敗時処理は controller を参照せず自己完結させ循環依存を避ける。
final dioProvider = Provider<Dio>((ref) {
  final dio = buildDio(cookieJar: ref.watch(cookieJarProvider));
  final tokenStore = ref.watch(tokenStoreProvider);
  final signal = ref.watch(authSignalProvider);

  dio.interceptors.add(
    AuthInterceptor(
      tokenStore: tokenStore,
      dio: dio,
      onRefresh: () async {
        try {
          final res = await dio.post<Map<String, dynamic>>('/auth/refresh');
          final token = res.data?['access_token'] as String?;
          if (token == null) return false;
          await tokenStore.writeAccessToken(token);
          return true;
        } catch (_) {
          return false;
        }
      },
      onAuthFailure: () async {
        await tokenStore.clear();
        signal.expire();
      },
    ),
  );
  return dio;
});
