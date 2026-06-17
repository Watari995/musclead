/// UI に提示するためのドメインエラー。API 例外はリポジトリで [Failure] に変換する。
sealed class Failure {
  const Failure(this.message);
  final String message;
}

class NetworkFailure extends Failure {
  const NetworkFailure([super.message = 'ネットワーク接続を確認してください']);
}

class TimeoutFailure extends Failure {
  const TimeoutFailure([super.message = '通信がタイムアウトしました']);
}

class UnauthorizedFailure extends Failure {
  const UnauthorizedFailure([super.message = 'ログインが必要です']);
}

/// バックエンドのエラー封筒 `{error:{code,message,data}}` に対応。
class ApiFailure extends Failure {
  const ApiFailure({
    required String message,
    required this.code,
    this.statusCode,
    this.data,
  }) : super(message);

  final String code;
  final int? statusCode;
  final Map<String, dynamic>? data;
}

class UnknownFailure extends Failure {
  const UnknownFailure([super.message = '予期しないエラーが発生しました']);
}
