import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'meal_dtos.dart';

class MealRepository {
  MealRepository(this._dio);

  final Dio _dio;

  Future<List<MealDto>> list({int limit = 50, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/meals',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListMealsResponse.fromJson(res.data!).meals;
      });

  Future<void> record(RecordMealRequest request) =>
      guardApi(() => _dio.post<void>('/meals', data: request.toJson()));

  Future<void> delete(String id) =>
      guardApi(() => _dio.delete<void>('/meals/$id'));
}

final mealRepositoryProvider = Provider<MealRepository>(
  (ref) => MealRepository(ref.watch(dioProvider)),
);

final mealsProvider = FutureProvider<List<MealDto>>(
  (ref) => ref.watch(mealRepositoryProvider).list(),
);
