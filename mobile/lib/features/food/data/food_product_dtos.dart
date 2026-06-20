import 'package:freezed_annotation/freezed_annotation.dart';

part 'food_product_dtos.freezed.dart';
part 'food_product_dtos.g.dart';

@freezed
abstract class FoodProductDto with _$FoodProductDto {
  const factory FoodProductDto({
    required String id,
    String? barcode,
    required String name,
    required int calories,
    String? proteinG,
    String? fatG,
    String? carbohydrateG,
    required String registerSource,
  }) = _FoodProductDto;

  factory FoodProductDto.fromJson(Map<String, dynamic> json) =>
      _$FoodProductDtoFromJson(json);
}

@freezed
abstract class SearchFoodProductsResponse with _$SearchFoodProductsResponse {
  const factory SearchFoodProductsResponse({
    @Default(<FoodProductDto>[]) List<FoodProductDto> foodProducts,
  }) = _SearchFoodProductsResponse;

  factory SearchFoodProductsResponse.fromJson(Map<String, dynamic> json) =>
      _$SearchFoodProductsResponseFromJson(json);
}

@freezed
abstract class CreateFoodProductRequest with _$CreateFoodProductRequest {
  const factory CreateFoodProductRequest({
    String? barcode,
    required String name,
    required int calories,
    String? proteinG,
    String? fatG,
    String? carbohydrateG,
  }) = _CreateFoodProductRequest;

  factory CreateFoodProductRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateFoodProductRequestFromJson(json);
}

@freezed
abstract class CreateFoodProductResponse with _$CreateFoodProductResponse {
  const factory CreateFoodProductResponse({String? foodProductId}) =
      _CreateFoodProductResponse;

  factory CreateFoodProductResponse.fromJson(Map<String, dynamic> json) =>
      _$CreateFoodProductResponseFromJson(json);
}
