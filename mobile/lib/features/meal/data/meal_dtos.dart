import 'package:decimal/decimal.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';
import '../../../core/util/datetime_converter.dart';
import '../../../core/util/decimal_converter.dart';

part 'meal_dtos.freezed.dart';
part 'meal_dtos.g.dart';

@freezed
abstract class PhotoDto with _$PhotoDto {
  const factory PhotoDto({
    @Default(0) int displayOrder,
    String? imagePath,
    String? imageUrl,
  }) = _PhotoDto;

  factory PhotoDto.fromJson(Map<String, dynamic> json) =>
      _$PhotoDtoFromJson(json);
}

@freezed
abstract class MealDto with _$MealDto {
  const factory MealDto({
    required String id,
    required DateTime eatenAt,
    required String mealType,
    @Default(0) int calories,
    @NullableDecimalStringConverter() Decimal? proteinG,
    @NullableDecimalStringConverter() Decimal? fatG,
    @NullableDecimalStringConverter() Decimal? carbohydrateG,
    String? memo,
    @Default(<PhotoDto>[]) List<PhotoDto> photos,
    String? userId,
    DateTime? createdAt,
    DateTime? updatedAt,
    String? foodProductId,
    String? servingCount,
  }) = _MealDto;

  factory MealDto.fromJson(Map<String, dynamic> json) =>
      _$MealDtoFromJson(json);
}

@freezed
abstract class ListMealsResponse with _$ListMealsResponse {
  const factory ListMealsResponse({
    @Default(<MealDto>[]) List<MealDto> meals,
    PaginationDto? pagination,
  }) = _ListMealsResponse;

  factory ListMealsResponse.fromJson(Map<String, dynamic> json) =>
      _$ListMealsResponseFromJson(json);
}

@freezed
abstract class MealPhotoInput with _$MealPhotoInput {
  const factory MealPhotoInput({
    required int displayOrder,
    required String imagePath,
  }) = _MealPhotoInput;

  factory MealPhotoInput.fromJson(Map<String, dynamic> json) =>
      _$MealPhotoInputFromJson(json);
}

@freezed
abstract class RecordMealRequest with _$RecordMealRequest {
  const factory RecordMealRequest({
    @UtcDateTimeConverter() required DateTime eatenAt,
    required String mealType,
    required int calories,
    double? proteinG,
    double? fatG,
    double? carbohydrateG,
    String? memo,
    @Default(<MealPhotoInput>[]) List<MealPhotoInput> photos,
    String? foodProductId,
    double? servingCount,
  }) = _RecordMealRequest;

  factory RecordMealRequest.fromJson(Map<String, dynamic> json) =>
      _$RecordMealRequestFromJson(json);
}
