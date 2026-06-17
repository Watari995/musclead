import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'user_dtos.dart';

class UserRepository {
  UserRepository(this._dio);

  final Dio _dio;

  Future<MeResponse> me() => guardApi(() async {
    final res = await _dio.get<Map<String, dynamic>>('/users/me');
    return MeResponse.fromJson(res.data!);
  });

  /// テーマ設定の更新（Patch-string: set + value）。
  Future<void> updateTheme(String theme) => guardApi(() async {
    await _dio.patch<Map<String, dynamic>>(
      '/users/me/preferences',
      data: {
        'theme': {'set': true, 'value': theme},
      },
    );
  });

  Future<({String uploadUrl, String path})> profileImagePresignedUrl(
    String contentType,
  ) => guardApi(() async {
    final res = await _dio.post<Map<String, dynamic>>(
      '/users/me/profile-image/presigned-url',
      data: {'content_type': contentType},
    );
    return (
      uploadUrl: res.data!['url'] as String,
      path: res.data!['path'] as String,
    );
  });
}

final userRepositoryProvider = Provider<UserRepository>(
  (ref) => UserRepository(ref.watch(dioProvider)),
);

/// 現在ログイン中ユーザー（プロフィール + 設定）。
final meProvider = FutureProvider<MeResponse>(
  (ref) => ref.watch(userRepositoryProvider).me(),
);
