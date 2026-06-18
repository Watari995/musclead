import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/auth/token_store.dart';
import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import '../../user/data/user_dtos.dart';
import 'auth_dtos.dart';

class AuthRepository {
  AuthRepository(this._dio, this._tokenStore);

  final Dio _dio;
  final TokenStore _tokenStore;

  Future<void> login(String email, String password) => guardApi(() async {
    final res = await _dio.post<Map<String, dynamic>>(
      '/auth/login',
      data: LoginRequest(email: email, password: password).toJson(),
    );
    final token = AccessTokenResponse.fromJson(res.data!);
    await _tokenStore.writeAccessToken(token.accessToken);
  });

  Future<void> register({
    required String name,
    required String email,
    required String password,
    String? birthday,
  }) => guardApi(() async {
    await _dio.post<Map<String, dynamic>>(
      '/users',
      data: RegisterRequest(
        name: name,
        email: email,
        password: password,
        birthday: birthday,
      ).toJson(),
    );
  });

  Future<void> logout() => guardApi(() async {
    try {
      await _dio.post<void>('/auth/logout');
    } finally {
      await _tokenStore.clear();
    }
  });
}

final authRepositoryProvider = Provider<AuthRepository>(
  (ref) =>
      AuthRepository(ref.watch(dioProvider), ref.watch(tokenStoreProvider)),
);
