import 'package:decimal/decimal.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';
import '../../../core/util/datetime_converter.dart';
import '../../../core/util/decimal_converter.dart';

part 'training_dtos.freezed.dart';
part 'training_dtos.g.dart';

@freezed
abstract class TrainingSetDto with _$TrainingSetDto {
  const factory TrainingSetDto({
    required String id,
    required int setNumber,
    @DecimalStringConverter() required Decimal weightKg,
    required int reps,
    int? restSeconds,
    String? memo,
  }) = _TrainingSetDto;

  factory TrainingSetDto.fromJson(Map<String, dynamic> json) =>
      _$TrainingSetDtoFromJson(json);
}

@freezed
abstract class TrainingExerciseDto with _$TrainingExerciseDto {
  const factory TrainingExerciseDto({
    required String id,
    required String exerciseId,
    @Default(0) int displayOrder,
    int? restSeconds,
    String? memo,
    @Default(<TrainingSetDto>[]) List<TrainingSetDto> sets,
  }) = _TrainingExerciseDto;

  factory TrainingExerciseDto.fromJson(Map<String, dynamic> json) =>
      _$TrainingExerciseDtoFromJson(json);
}

@freezed
abstract class TrainingDto with _$TrainingDto {
  const factory TrainingDto({
    required String id,
    required DateTime startedAt,
    DateTime? endedAt,
    String? memo,
    String? userId,
    @Default(<TrainingExerciseDto>[]) List<TrainingExerciseDto> exercises,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _TrainingDto;

  factory TrainingDto.fromJson(Map<String, dynamic> json) =>
      _$TrainingDtoFromJson(json);
}

@freezed
abstract class ListTrainingsResponse with _$ListTrainingsResponse {
  const factory ListTrainingsResponse({
    @Default(<TrainingDto>[]) List<TrainingDto> trainings,
    PaginationDto? pagination,
  }) = _ListTrainingsResponse;

  factory ListTrainingsResponse.fromJson(Map<String, dynamic> json) =>
      _$ListTrainingsResponseFromJson(json);
}

@freezed
abstract class RecordTrainingSetRequest with _$RecordTrainingSetRequest {
  const factory RecordTrainingSetRequest({
    required int setNumber,
    @DecimalStringConverter() required Decimal weightKg,
    required int reps,
    int? restSeconds,
    String? memo,
  }) = _RecordTrainingSetRequest;

  factory RecordTrainingSetRequest.fromJson(Map<String, dynamic> json) =>
      _$RecordTrainingSetRequestFromJson(json);
}

@freezed
abstract class RecordTrainingExerciseRequest with _$RecordTrainingExerciseRequest {
  const factory RecordTrainingExerciseRequest({
    required String exerciseId,
    required int displayOrder,
    int? restSeconds,
    String? memo,
    @Default(<RecordTrainingSetRequest>[]) List<RecordTrainingSetRequest> sets,
  }) = _RecordTrainingExerciseRequest;

  factory RecordTrainingExerciseRequest.fromJson(Map<String, dynamic> json) =>
      _$RecordTrainingExerciseRequestFromJson(json);
}

@freezed
abstract class RecordTrainingRequest with _$RecordTrainingRequest {
  const factory RecordTrainingRequest({
    @UtcDateTimeConverter() required DateTime startedAt,
    @NullableUtcDateTimeConverter() DateTime? endedAt,
    String? memo,
    @Default(<RecordTrainingExerciseRequest>[])
    List<RecordTrainingExerciseRequest> exercises,
  }) = _RecordTrainingRequest;

  factory RecordTrainingRequest.fromJson(Map<String, dynamic> json) =>
      _$RecordTrainingRequestFromJson(json);
}
