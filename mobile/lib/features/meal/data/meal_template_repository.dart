import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'meal_template_dtos.dart';

class MealTemplateRepository {
  MealTemplateRepository(this._dio);

  final Dio _dio;

  Future<List<MealTemplateDto>> list({int limit = 100, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/meal_templates',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListMealTemplatesResponse.fromJson(res.data!).mealTemplates;
      });

  Future<void> create(UpsertMealTemplateRequest request) => guardApi(
    () => _dio.post<void>('/meal_templates', data: request.toJson()),
  );

  Future<void> update(String id, UpsertMealTemplateRequest request) =>
      guardApi(
        () => _dio.put<void>('/meal_templates/$id', data: request.toJson()),
      );

  Future<void> delete(String id) =>
      guardApi(() => _dio.delete<void>('/meal_templates/$id'));

  Future<void> reorder(ReorderMealTemplatesRequest request) => guardApi(
    () => _dio.post<void>('/meal_templates/reorder', data: request.toJson()),
  );
}

final mealTemplateRepositoryProvider = Provider<MealTemplateRepository>(
  (ref) => MealTemplateRepository(ref.watch(dioProvider)),
);

final mealTemplatesProvider = FutureProvider<List<MealTemplateDto>>(
  (ref) => ref.watch(mealTemplateRepositoryProvider).list(),
);
