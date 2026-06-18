import 'package:dio/dio.dart';

import '../auth/token_store.dart';

/// 各リクエストに Bearer access token を付与し、401 時に refresh → 元リクエスト再試行を行う。
///
/// - `/auth/*` は対象外（refresh 自体や login に bearer を付けない）。
/// - 並行して複数の 401 が出ても refresh は 1 回に集約（スタンピード防止）。
/// - refresh も失敗したら [onAuthFailure] を呼ぶ。
class AuthInterceptor extends Interceptor {
  AuthInterceptor({
    required this.tokenStore,
    required this.dio,
    required this.onRefresh,
    required this.onAuthFailure,
  });

  final TokenStore tokenStore;
  final Dio dio;
  final Future<bool> Function() onRefresh;
  final void Function() onAuthFailure;

  Future<bool>? _refreshing;
  static const _retriedKey = 'retried';

  bool _isAuthPath(String path) => path.contains('/auth/');

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (!_isAuthPath(options.path)) {
      final token = await tokenStore.readAccessToken();
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }
    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final status = err.response?.statusCode;
    final opts = err.requestOptions;
    final alreadyRetried = opts.extra[_retriedKey] == true;

    if (status != 401 || _isAuthPath(opts.path) || alreadyRetried) {
      return handler.next(err);
    }

    final refreshed = await (_refreshing ??= _runRefresh());
    _refreshing = null;

    if (!refreshed) {
      onAuthFailure();
      return handler.next(err);
    }

    try {
      final token = await tokenStore.readAccessToken();
      final retryOpts = opts..extra[_retriedKey] = true;
      if (token != null) {
        retryOpts.headers['Authorization'] = 'Bearer $token';
      }
      final response = await dio.fetch<dynamic>(retryOpts);
      return handler.resolve(response);
    } on DioException catch (e) {
      return handler.next(e);
    }
  }

  Future<bool> _runRefresh() async {
    try {
      return await onRefresh();
    } catch (_) {
      return false;
    }
  }
}
