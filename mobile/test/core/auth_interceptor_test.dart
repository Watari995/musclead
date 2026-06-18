import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:musclead/core/api/auth_interceptor.dart';

import '../support/fakes.dart';

/// access token の有効性をモデル化したスタブ。
/// `Bearer new` のみ 200、それ以外（old/無し）は 401 を返す。
/// `/auth/refresh` は新トークンを返し回数を数える。
class _StubAdapter implements HttpClientAdapter {
  int refreshCount = 0;

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) async {
    if (options.path.contains('/auth/refresh')) {
      refreshCount++;
      return _json('{"access_token":"new"}', 200);
    }
    final auth = options.headers['Authorization'];
    if (auth == 'Bearer new') return _json('{"ok":true}', 200);
    return _json(
      '{"error":{"code":"general.unauthorized_error","message":"x"}}',
      401,
    );
  }

  @override
  void close({bool force = false}) {}
}

ResponseBody _json(String body, int code) => ResponseBody.fromString(
  body,
  code,
  headers: {
    Headers.contentTypeHeader: ['application/json'],
  },
);

Dio _buildDio(FakeTokenStore store, {void Function()? onAuthFailure}) {
  final dio = Dio(BaseOptions(baseUrl: 'https://api.test'));
  dio.interceptors.add(
    AuthInterceptor(
      tokenStore: store,
      dio: dio,
      onRefresh: () async {
        final r = await dio.post<Map<String, dynamic>>('/auth/refresh');
        await store.writeAccessToken(r.data!['access_token'] as String);
        return true;
      },
      onAuthFailure: onAuthFailure ?? () {},
    ),
  );
  dio.httpClientAdapter = _StubAdapter();
  return dio;
}

void main() {
  test('401 → refresh → 新トークンで再試行して成功', () async {
    final store = FakeTokenStore();
    await store.writeAccessToken('old');
    final dio = _buildDio(store);

    final res = await dio.get<Map<String, dynamic>>('/users/me');

    expect(res.statusCode, 200);
    expect(res.data!['ok'], true);
    expect(await store.readAccessToken(), 'new');
    expect((dio.httpClientAdapter as _StubAdapter).refreshCount, 1);
  });

  test('並行 401 は refresh を 1 回に集約（スタンピード防止）', () async {
    final store = FakeTokenStore();
    await store.writeAccessToken('old');
    final dio = _buildDio(store);

    final results = await Future.wait([
      dio.get<Map<String, dynamic>>('/users/me'),
      dio.get<Map<String, dynamic>>('/users/me'),
      dio.get<Map<String, dynamic>>('/users/me'),
    ]);

    for (final r in results) {
      expect(r.statusCode, 200);
    }
    expect((dio.httpClientAdapter as _StubAdapter).refreshCount, 1);
  });
}
