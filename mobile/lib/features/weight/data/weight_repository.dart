import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'weight_dtos.dart';

class WeightRepository {
  WeightRepository(this._dio);

  final Dio _dio;

  Future<List<WeightDto>> list({int limit = 60, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/weights',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListWeightsResponse.fromJson(res.data!).weights;
      });

  Future<void> upsert(UpsertWeightRequest request) =>
      guardApi(() => _dio.post<void>('/weights', data: request.toJson()));

  Future<void> delete(String id) =>
      guardApi(() => _dio.delete<void>('/weights/$id'));
}

final weightRepositoryProvider = Provider<WeightRepository>(
  (ref) => WeightRepository(ref.watch(dioProvider)),
);

/// 直近の体重記録（新しい順で返るため、グラフ用に昇順へ並べ替えて使う）。
final weightsProvider = FutureProvider<List<WeightDto>>(
  (ref) => ref.watch(weightRepositoryProvider).list(),
);
