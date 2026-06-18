import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:musclead/core/api/dio_error_mapper.dart';
import 'package:musclead/core/error/failure.dart';

void main() {
  RequestOptions opts() => RequestOptions(path: '/x');

  DioException dio(DioExceptionType type, {Response<dynamic>? response}) =>
      DioException(requestOptions: opts(), type: type, response: response);

  test('接続エラー → NetworkFailure', () {
    expect(
      dio(DioExceptionType.connectionError).let(mapDioException),
      isA<NetworkFailure>(),
    );
  });

  test('タイムアウト → TimeoutFailure', () {
    expect(
      mapDioException(dio(DioExceptionType.receiveTimeout)),
      isA<TimeoutFailure>(),
    );
  });

  test('401 → UnauthorizedFailure', () {
    final f = mapDioException(
      dio(
        DioExceptionType.badResponse,
        response: Response<dynamic>(
          requestOptions: opts(),
          statusCode: 401,
          data: {
            'error': {'code': 'general.unauthorized_error', 'message': '要ログイン'},
          },
        ),
      ),
    );
    expect(f, isA<UnauthorizedFailure>());
  });

  test('エラー封筒 → ApiFailure(code/status/data を保持)', () {
    final f = mapDioException(
      dio(
        DioExceptionType.badResponse,
        response: Response<dynamic>(
          requestOptions: opts(),
          statusCode: 403,
          data: {
            'error': {
              'code': 'training.routine_limit_reached_error',
              'message': '上限です',
              'data': {'max': 10},
            },
          },
        ),
      ),
    );
    expect(f, isA<ApiFailure>());
    final api = f as ApiFailure;
    expect(api.code, 'training.routine_limit_reached_error');
    expect(api.statusCode, 403);
    expect(api.message, '上限です');
    expect(api.data?['max'], 10);
  });

  test('Dio 以外の例外 → UnknownFailure', () {
    expect(mapDioException(Exception('boom')), isA<UnknownFailure>());
  });
}

extension<T> on T {
  R let<R>(R Function(T) f) => f(this);
}
