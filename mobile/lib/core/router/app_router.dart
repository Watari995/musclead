import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/application/auth_controller.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/auth/presentation/register_screen.dart';
import '../../features/auth/presentation/splash_screen.dart';
import '../../features/meal/presentation/meals_screen.dart';
import '../../features/training/presentation/trainings_screen.dart';
import '../../features/user/presentation/profile_screen.dart';
import '../../features/weight/presentation/weights_screen.dart';
import '../widgets/home_shell.dart';

final goRouterProvider = Provider<GoRouter>((ref) {
  // 認証状態の変化を GoRouter の refreshListenable に橋渡しする。
  final refresh = ValueNotifier<AuthStatus>(AuthStatus.unknown);
  ref.listen<AuthStatus>(
    authControllerProvider,
    (_, next) => refresh.value = next,
    fireImmediately: true,
  );
  ref.onDispose(refresh.dispose);

  return GoRouter(
    initialLocation: '/splash',
    refreshListenable: refresh,
    redirect: (context, state) {
      final status = refresh.value;
      final loc = state.matchedLocation;
      final atAuth = loc == '/login' || loc == '/register';
      final atSplash = loc == '/splash';

      if (status == AuthStatus.unknown) return atSplash ? null : '/splash';
      if (status == AuthStatus.unauthenticated) return atAuth ? null : '/login';
      // authenticated
      if (atAuth || atSplash) return '/meals';
      return null;
    },
    routes: [
      GoRoute(path: '/splash', builder: (_, _) => const SplashScreen()),
      GoRoute(path: '/login', builder: (_, _) => const LoginScreen()),
      GoRoute(path: '/register', builder: (_, _) => const RegisterScreen()),
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            HomeShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [GoRoute(path: '/meals', builder: (_, _) => const MealsScreen())],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(path: '/trainings', builder: (_, _) => const TrainingsScreen())
            ],
          ),
          StatefulShellBranch(
            routes: [GoRoute(path: '/weights', builder: (_, _) => const WeightsScreen())],
          ),
          StatefulShellBranch(
            routes: [GoRoute(path: '/profile', builder: (_, _) => const ProfileScreen())],
          ),
        ],
      ),
    ],
  );
});
