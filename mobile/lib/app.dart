import 'dart:async';

import 'package:app_links/app_links.dart';
import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import 'core/providers/core_providers.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/theme/theme_controller.dart';
import 'features/auth/application/auth_controller.dart';
import 'features/notifications/application/push_service.dart';
import 'features/user/data/user_repository.dart';

class MuscleadApp extends ConsumerStatefulWidget {
  const MuscleadApp({super.key});

  @override
  ConsumerState<MuscleadApp> createState() => _MuscleadAppState();
}

class _MuscleadAppState extends ConsumerState<MuscleadApp> {
  StreamSubscription<Uri>? _linkSub;
  bool _pushTokenRegistered = false;

  @override
  void initState() {
    super.initState();
    _linkSub = AppLinks().uriLinkStream.listen((uri) {
      if (uri.scheme == 'musclead' && uri.host == 'oauth') {
        ref.read(oauthCallbackProvider.notifier).set(uri);
      }
    });
  }

  @override
  void dispose() {
    _linkSub?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final accent = ref.watch(accentProvider);
    final mode = ref.watch(themeModeProvider);
    final router = ref.watch(goRouterProvider);

    if (ref.watch(authControllerProvider) == AuthStatus.authenticated) {
      ref.listen(meProvider, (_, next) {
        final theme = next.asData?.value.preferences?.theme;
        if (theme != null) {
          ref.read(themeModeProvider.notifier).hydrate(theme);
        }
      });
      if (!_pushTokenRegistered) {
        _pushTokenRegistered = true;
        ref.read(pushServiceProvider).initAndGetToken();
      }
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
