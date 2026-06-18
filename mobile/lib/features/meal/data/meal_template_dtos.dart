import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';

part 'meal_template_dtos.freezed.dart';
part 'meal_template_dtos.g.dart';

@freezed
abstract class MealTemplateDto with _$MealTemplateDto {
  const factory MealTemplateDto({
    required String id,
    required String name,
    @Default(0) int displayOrder,
    required String mealType,
    @Default(0) int calories,
    String? proteinG,
    String? fatG,
    String? carbohydrateG,
    String? userId,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _MealTemplateDto;

  factory MealTemplateDto.fromJson(Map<String, dynamic> json) =>
      _$MealTemplateDtoFromJson(json);
}

@freezed
abstract class ListMealTemplatesResponse with _$ListMealTemplatesResponse {
  const factory ListMealTemplatesResponse({
    @Default(<MealTemplateDto>[]) List<MealTemplateDto> mealTemplates,
    PaginationDto? pagination,
  }) = _ListMealTemplatesResponse;

  factory ListMealTemplatesResponse.fromJson(Map<String, dynamic> json) =>
      _$ListMealTemplatesResponseFromJson(json);
}

@freezed
abstract class UpsertMealTemplateRequest with _$UpsertMealTemplateRequest {
  const factory UpsertMealTemplateRequest({
    required String name,
    required String mealType,
    required int calories,
    double? proteinG,
    double? fatG,
    double? carbohydrateG,
  }) = _UpsertMealTemplateRequest;

  factory UpsertMealTemplateRequest.fromJson(Map<String, dynamic> json) =>
      _$UpsertMealTemplateRequestFromJson(json);
}

@freezed
abstract class ReorderMealTemplatesRequest with _$ReorderMealTemplatesRequest {
  const factory ReorderMealTemplatesRequest({
    @Default(<String>[]) List<String> mealTemplateIds,
  }) = _ReorderMealTemplatesRequest;

  factory ReorderMealTemplatesRequest.fromJson(Map<String, dynamic> json) =>
      _$ReorderMealTemplatesRequestFromJson(json);
}
