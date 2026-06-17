import 'package:dio/dio.dart';

import '../error/failure.dart';

/// DioException（やその他例外）を UI 向けの [Failure] に変換する。
/// バックエンドのエラー封筒 `{error:{code,message,data}}` を解釈する。
Failure mapDioException(Object error) {
  if (error is! DioException) return const UnknownFailure();
  switch (error.type) {
    case DioExceptionType.connectionTimeout:
    case DioExceptionType.sendTimeout:
    case DioExceptionType.receiveTimeout:
      return const TimeoutFailure();
    case DioExceptionType.connectionError:
      return const NetworkFailure();
    case DioExceptionType.badResponse:
      return _fromResponse(error.response);
    case DioExceptionType.cancel:
    case DioExceptionType.badCertificate:
    case DioExceptionType.unknown:
      return const UnknownFailure();
  }
}

Failure _fromResponse(Response<dynamic>? res) {
  final status = res?.statusCode;
  final data = res?.data;
  String? code;
  String? message;
  Map<String, dynamic>? extra;

  if (data is Map && data['error'] is Map) {
    final err = (data['error'] as Map).cast<String, dynamic>();
    code = err['code'] as String?;
    message = err['message'] as String?;
    extra = err['data'] is Map
        ? (err['data'] as Map).cast<String, dynamic>()
        : null;
  }

  if (status == 401) return UnauthorizedFailure(message ?? 'ログインが必要です');
  return ApiFailure(
    message: message ?? 'エラーが発生しました (HTTP $status)',
    code: code ?? 'unknown',
    statusCode: status,
    data: extra,
  );
}
