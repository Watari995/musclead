import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/error/failure.dart';
import '../../../core/providers/core_providers.dart';
import 'food_product_dtos.dart';

class FoodProductRepository {
  FoodProductRepository(this._dio);
  final Dio _dio;

  Future<List<FoodProductDto>> searchByName(String name) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/food_products',
          queryParameters: {'q': name},
        );
        return SearchFoodProductsResponse.fromJson(res.data!).foodProducts;
      });

  Future<List<FoodProductDto>> searchByBarcode(String barcode) async {
    try {
      return await guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/food_products/barcode/$barcode',
        );
        return SearchFoodProductsResponse.fromJson(res.data!).foodProducts;
      });
    } on ApiFailure catch (f) {
      if (f.statusCode == 404) return [];
      rethrow;
    }
  }

  Future<String> create(CreateFoodProductRequest request) =>
      guardApi(() async {
        final res = await _dio.post<Map<String, dynamic>>(
          '/food_products',
          data: request.toJson(),
        );
        return CreateFoodProductResponse.fromJson(res.data!).foodProductId ?? '';
      });
}

final foodProductRepositoryProvider = Provider<FoodProductRepository>(
  (ref) => FoodProductRepository(ref.watch(dioProvider)),
);
