import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'notification_dtos.dart';

class NotificationRepository {
  NotificationRepository(this._dio);

  final Dio _dio;

  Future<GetNotificationsResponse> listNotifications() => guardApi(() async {
    final res = await _dio.get<Map<String, dynamic>>('/notifications');
    return GetNotificationsResponse.fromJson(res.data!);
  });

  Future<void> markAsRead(String id) =>
      guardApi(() => _dio.put<void>('/notifications/$id/read'));

  Future<void> registerDeviceToken({
    required String token,
    required String platform,
  }) => guardApi(
    () => _dio.post<void>(
      '/device-tokens',
      data: {'token': token, 'platform': platform},
    ),
  );

  Future<WeeklyGoalDto?> getWeeklyGoal() => guardApi(() async {
    final res = await _dio.get<Map<String, dynamic>?>('/users/me/weekly-goal');
    if (res.data == null) return null;
    return WeeklyGoalDto.fromJson(res.data!);
  });

  Future<WeeklyGoalDto> upsertWeeklyGoal({
    int? trainingCount,
    int? calorieAverage,
    double? weightChangeKg,
  }) => guardApi(() async {
    final res = await _dio.put<Map<String, dynamic>>(
      '/users/me/weekly-goal',
      data: {
        'training_count': trainingCount,
        'calorie_average': calorieAverage,
        'weight_change_kg': weightChangeKg,
      },
    );
    return WeeklyGoalDto.fromJson(res.data!);
  });
}

final notificationRepositoryProvider = Provider<NotificationRepository>(
  (ref) => NotificationRepository(ref.watch(dioProvider)),
);

final notificationsProvider = FutureProvider<GetNotificationsResponse>((
  ref,
) async {
  return ref.watch(notificationRepositoryProvider).listNotifications();
});

final weeklyGoalProvider = FutureProvider<WeeklyGoalDto?>((ref) async {
  return ref.watch(notificationRepositoryProvider).getWeeklyGoal();
});
