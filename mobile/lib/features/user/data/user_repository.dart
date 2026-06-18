import 'dart:typed_data';

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

  /// アカウント削除（App Store のアカウント削除要件に対応）。
  Future<void> deleteAccount(String userId) =>
      guardApi(() => _dio.delete<void>('/users/$userId'));

  /// テーマ設定の更新。サーバの Patch[T] は plain 値を受ける（{theme:"dark"}）。
  Future<void> updateTheme(String theme) => guardApi(
    () => _dio.patch<void>('/users/me/preferences', data: {'theme': theme}),
  );

  /// プロフィール画像アップロード:
  /// 1) presigned URL 取得 → 2) 署名URLへ直接 PUT → 3) PATCH /users/me で path 紐付け。
  Future<void> uploadProfileImage(Uint8List bytes, String contentType) =>
      guardApi(() async {
        final presigned = await _dio.post<Map<String, dynamic>>(
          '/users/me/profile-image/presigned-url',
          data: {'content_type': contentType},
        );
        final url = presigned.data!['url'] as String;
        final path = presigned.data!['path'] as String;
        // 署名済み URL へ直接 PUT（認証 interceptor を通さない素の Dio）。
        await Dio().put<void>(
          url,
          data: Stream<List<int>>.fromIterable([bytes]),
          options: Options(
            headers: {
              Headers.contentTypeHeader: contentType,
              Headers.contentLengthHeader: bytes.length,
            },
          ),
        );
        await _dio.patch<void>('/users/me', data: {'profile_image_path': path});
      });
}

final userRepositoryProvider = Provider<UserRepository>(
  (ref) => UserRepository(ref.watch(dioProvider)),
);

/// 現在ログイン中ユーザー（プロフィール + 設定）。
final meProvider = FutureProvider<MeResponse>(
  (ref) => ref.watch(userRepositoryProvider).me(),
);
