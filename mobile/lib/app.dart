import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/theme/theme_controller.dart';
import 'features/auth/application/auth_controller.dart';
import 'features/user/data/user_repository.dart';

class MuscleadApp extends ConsumerWidget {
  const MuscleadApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final accent = ref.watch(accentProvider);
    final mode = ref.watch(themeModeProvider);
    final router = ref.watch(goRouterProvider);

    // 認証後にサーバーの外観設定（preferences.theme）を適用する。
    // 未認証時は meProvider を起動しない（401 を避ける）ため authenticated で限定。
    // ref.listen は meProvider が変化した時のみ発火する（AsyncLoading→AsyncData など）。
    // rebuild ごとに発火しないため、ユーザーがローカルで変更した設定を上書きしない。
    if (ref.watch(authControllerProvider) == AuthStatus.authenticated) {
      ref.listen(meProvider, (_, next) {
        final theme = next.asData?.value.preferences?.theme;
        if (theme != null) {
          ref.read(themeModeProvider.notifier).hydrate(theme);
        }
      });
    }

    return MaterialApp.router(
      title: 'musclead',
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(Brightness.light, accent),
      darkTheme: buildAppTheme(Brightness.dark, accent),
      themeMode: mode,
      routerConfig: router,
    );
  }
}
