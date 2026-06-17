import '../api/dio_error_mapper.dart';
import '../error/failure.dart';

/// リポジトリの API 呼び出しを包み、例外を必ず [Failure] に正規化する。
Future<T> guardApi<T>(Future<T> Function() request) async {
  try {
    return await request();
  } on Failure {
    rethrow;
  } catch (e) {
    throw mapDioException(e);
  }
}
