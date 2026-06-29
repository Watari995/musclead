import 'package:freezed_annotation/freezed_annotation.dart';

part 'calendar_dtos.freezed.dart';
part 'calendar_dtos.g.dart';

@freezed
abstract class MonthlySummaryDayDto with _$MonthlySummaryDayDto {
  const factory MonthlySummaryDayDto({
    required String date,
    @JsonKey(name: 'has_training') @Default(false) bool hasTraining,
    @JsonKey(name: 'has_meal') @Default(false) bool hasMeal,
    @JsonKey(name: 'has_weight') @Default(false) bool hasWeight,
  }) = _MonthlySummaryDayDto;

  factory MonthlySummaryDayDto.fromJson(Map<String, dynamic> json) =>
      _$MonthlySummaryDayDtoFromJson(json);
}

@freezed
abstract class GetMonthlySummaryResponse with _$GetMonthlySummaryResponse {
  const factory GetMonthlySummaryResponse({
    @Default(<MonthlySummaryDayDto>[]) List<MonthlySummaryDayDto> days,
  }) = _GetMonthlySummaryResponse;

  factory GetMonthlySummaryResponse.fromJson(Map<String, dynamic> json) =>
      _$GetMonthlySummaryResponseFromJson(json);
}

@freezed
abstract class CalendarTrainingSummaryDto with _$CalendarTrainingSummaryDto {
  const factory CalendarTrainingSummaryDto({
    @JsonKey(name: 'training_id') required String trainingId,
    @JsonKey(name: 'started_at') required DateTime startedAt,
    @JsonKey(name: 'ended_at') DateTime? endedAt,
    @JsonKey(name: 'exercise_count') @Default(0) int exerciseCount,
    @JsonKey(name: 'set_count') @Default(0) int setCount,
  }) = _CalendarTrainingSummaryDto;

  factory CalendarTrainingSummaryDto.fromJson(Map<String, dynamic> json) =>
      _$CalendarTrainingSummaryDtoFromJson(json);
}

@freezed
abstract class CalendarMealSummaryDto with _$CalendarMealSummaryDto {
  const factory CalendarMealSummaryDto({
    @JsonKey(name: 'meal_id') required String mealId,
    @JsonKey(name: 'meal_type') required String mealType,
    @JsonKey(name: 'eaten_at') required String eatenAt,
    @Default(0) int calories,
    @JsonKey(name: 'protein_g') String? proteinG,
    @JsonKey(name: 'fat_g') String? fatG,
    @JsonKey(name: 'carbohydrate_g') String? carbohydrateG,
  }) = _CalendarMealSummaryDto;

  factory CalendarMealSummaryDto.fromJson(Map<String, dynamic> json) =>
      _$CalendarMealSummaryDtoFromJson(json);
}

@freezed
abstract class CalendarWeightSummaryDto with _$CalendarWeightSummaryDto {
  const factory CalendarWeightSummaryDto({
    @JsonKey(name: 'weight_id') required String weightId,
    @JsonKey(name: 'weight_kg') required String weightKg,
    @JsonKey(name: 'body_fat_percentage') String? bodyFatPercentage,
    @JsonKey(name: 'skeletal_muscle_kg') String? skeletalMuscleKg,
    @JsonKey(name: 'measured_at') required String measuredAt,
  }) = _CalendarWeightSummaryDto;

  factory CalendarWeightSummaryDto.fromJson(Map<String, dynamic> json) =>
      _$CalendarWeightSummaryDtoFromJson(json);
}

@freezed
abstract class GetDailySummaryResponse with _$GetDailySummaryResponse {
  const factory GetDailySummaryResponse({
    @Default(<CalendarTrainingSummaryDto>[])
    List<CalendarTrainingSummaryDto> trainings,
    @Default(<CalendarMealSummaryDto>[]) List<CalendarMealSummaryDto> meals,
    @Default(<CalendarWeightSummaryDto>[])
    List<CalendarWeightSummaryDto> weights,
  }) = _GetDailySummaryResponse;

  factory GetDailySummaryResponse.fromJson(Map<String, dynamic> json) =>
      _$GetDailySummaryResponseFromJson(json);
}
