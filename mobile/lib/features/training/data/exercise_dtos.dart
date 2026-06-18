import 'package:decimal/decimal.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';
import '../../../core/util/decimal_converter.dart';

part 'exercise_dtos.freezed.dart';
part 'exercise_dtos.g.dart';

@freezed
abstract class ExerciseDto with _$ExerciseDto {
  const factory ExerciseDto({
    required String id,
    required String name,
    @Default(0) int displayOrder,
    String? userId,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _ExerciseDto;

  factory ExerciseDto.fromJson(Map<String, dynamic> json) =>
      _$ExerciseDtoFromJson(json);
}

@freezed
abstract class BestSetDto with _$BestSetDto {
  const factory BestSetDto({
    required String exerciseId,
    @DecimalStringConverter() required Decimal weightKg,
    required int reps,
    String? trainingId,
    DateTime? performedAt,
  }) = _BestSetDto;

  factory BestSetDto.fromJson(Map<String, dynamic> json) =>
      _$BestSetDtoFromJson(json);
}

@freezed
abstract class ListExercisesResponse with _$ListExercisesResponse {
  const factory ListExercisesResponse({
    @Default(<ExerciseDto>[]) List<ExerciseDto> exercises,
    PaginationDto? pagination,
  }) = _ListExercisesResponse;

  factory ListExercisesResponse.fromJson(Map<String, dynamic> json) =>
      _$ListExercisesResponseFromJson(json);
}

@freezed
abstract class ListBestSetsResponse with _$ListBestSetsResponse {
  const factory ListBestSetsResponse({
    @Default(<BestSetDto>[]) List<BestSetDto> bestSets,
  }) = _ListBestSetsResponse;

  factory ListBestSetsResponse.fromJson(Map<String, dynamic> json) =>
      _$ListBestSetsResponseFromJson(json);
}

@freezed
abstract class UpsertExerciseRequest with _$UpsertExerciseRequest {
  const factory UpsertExerciseRequest({required String name}) =
      _UpsertExerciseRequest;

  factory UpsertExerciseRequest.fromJson(Map<String, dynamic> json) =>
      _$UpsertExerciseRequestFromJson(json);
}
