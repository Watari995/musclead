import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/theme/theme_controller.dart';

class MuscleadApp extends ConsumerWidget {
  const MuscleadApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final accent = ref.watch(accentProvider);
    final mode = ref.watch(themeModeProvider);
    final router = ref.watch(goRouterProvider);

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
