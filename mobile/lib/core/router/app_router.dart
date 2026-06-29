import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/application/auth_controller.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/calendar/presentation/calendar_screen.dart';
import '../../features/auth/presentation/register_screen.dart';
import '../../features/auth/presentation/splash_screen.dart';
import '../../features/meal/data/meal_dtos.dart';
import '../../features/meal/data/meal_template_dtos.dart';
import '../../features/meal/presentation/meal_record_screen.dart';
import '../../features/meal/presentation/meals_screen.dart';
import '../../features/training/presentation/exercise_records_screen.dart';
import '../../features/training/presentation/exercises_screen.dart';
import '../../features/training/presentation/routine_create_screen.dart';
import '../../features/training/presentation/routines_screen.dart';
import '../../features/training/presentation/training_record_screen.dart';
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
    onException: (context, state, router) {
      // musclead://oauth/callback 等のカスタム URL スキームは go_router では
      // ルート解決できないが、app_links が処理するため無視して問題ない。
      if (state.uri.scheme == 'musclead') return;
      router.go('/calendar');
    },
    redirect: (context, state) {
      final status = refresh.value;
      final loc = state.matchedLocation;
      final atAuth = loc == '/login' || loc == '/register';
      final atSplash = loc == '/splash';

      if (status == AuthStatus.unknown) return atSplash ? null : '/splash';
      if (status == AuthStatus.unauthenticated) return atAuth ? null : '/login';
      // authenticated
      if (atAuth || atSplash) return '/calendar';
      return null;
    },
    routes: [
      GoRoute(path: '/splash', builder: (_, _) => const SplashScreen()),
      GoRoute(path: '/login', builder: (_, _) => const LoginScreen()),
      GoRoute(path: '/register', builder: (_, _) => const RegisterScreen()),
      // 食事記録・編集はタブシェルの上にフルスクリーンで push する
      GoRoute(
        path: '/meals/new',
        builder: (_, state) {
          final extra = state.extra;
          return MealRecordScreen(
            fromTemplate: extra is MealTemplateDto ? extra : null,
          );
        },
      ),
      GoRoute(
        path: '/meals/:id/edit',
        builder: (_, state) => MealRecordScreen(
          existing: state.extra is MealDto ? state.extra as MealDto : null,
        ),
      ),
      // トレーニング記録はタブシェルの上にフルスクリーンで push する
      GoRoute(
        path: '/trainings/new',
        builder: (_, _) => const TrainingRecordScreen(),
      ),
      GoRoute(path: '/exercises', builder: (_, _) => const ExercisesScreen()),
      GoRoute(
        path: '/records',
        builder: (_, _) => const ExerciseRecordsScreen(),
      ),
      GoRoute(path: '/routines', builder: (_, _) => const RoutinesScreen()),
      GoRoute(
        path: '/routines/new',
        builder: (_, _) => const RoutineCreateScreen(),
      ),
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) =>
            HomeShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/calendar',
                builder: (_, _) => const CalendarScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(path: '/meals', builder: (_, _) => const MealsScreen()),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/trainings',
                builder: (_, _) => const TrainingsScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/weights',
                builder: (_, _) => const WeightsScreen(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: '/profile',
                builder: (_, _) => const ProfileScreen(),
              ),
            ],
          ),
        ],
      ),
    ],
  );
});
