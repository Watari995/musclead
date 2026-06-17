import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/providers/core_providers.dart';
import '../data/auth_repository.dart';

enum AuthStatus { unknown, authenticated, unauthenticated }

final authControllerProvider =
    NotifierProvider<AuthController, AuthStatus>(AuthController.new);

/// 認証状態を保持。dio 層からの [AuthSignal]（refresh も失敗＝セッション切れ）を購読し、
/// 未認証へ遷移する。go_router の redirect がこの状態を見て遷移先を決める。
class AuthController extends Notifier<AuthStatus> {
  @override
  AuthStatus build() {
    final signal = ref.watch(authSignalProvider);
    signal.addListener(_onExpired);
    ref.onDispose(() => signal.removeListener(_onExpired));
    Future.microtask(_bootstrap);
    return AuthStatus.unknown;
  }

  void _onExpired() => state = AuthStatus.unauthenticated;

  Future<void> _bootstrap() async {
    final token = await ref.read(tokenStoreProvider).readAccessToken();
    state = token == null ? AuthStatus.unauthenticated : AuthStatus.authenticated;
  }

  Future<void> login(String email, String password) async {
    await ref.read(authRepositoryProvider).login(email, password);
    state = AuthStatus.authenticated;
  }

  Future<void> register({
    required String name,
    required String email,
    required String password,
    String? birthday,
  }) async {
    await ref.read(authRepositoryProvider).register(
          name: name,
          email: email,
          password: password,
          birthday: birthday,
        );
    await login(email, password);
  }

  Future<void> logout() async {
    await ref.read(authRepositoryProvider).logout();
    state = AuthStatus.unauthenticated;
  }
}
