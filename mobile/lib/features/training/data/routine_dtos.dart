import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';

part 'routine_dtos.freezed.dart';
part 'routine_dtos.g.dart';

@freezed
abstract class RoutineExerciseDto with _$RoutineExerciseDto {
  const factory RoutineExerciseDto({
    required String id,
    required String exerciseId,
    String? exerciseName,
    @Default(0) int displayOrder,
  }) = _RoutineExerciseDto;

  factory RoutineExerciseDto.fromJson(Map<String, dynamic> json) =>
      _$RoutineExerciseDtoFromJson(json);
}

@freezed
abstract class RoutineDto with _$RoutineDto {
  const factory RoutineDto({
    required String id,
    required String name,
    String? userId,
    @Default(<RoutineExerciseDto>[]) List<RoutineExerciseDto> routineExercises,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _RoutineDto;

  factory RoutineDto.fromJson(Map<String, dynamic> json) =>
      _$RoutineDtoFromJson(json);
}

@freezed
abstract class ListRoutinesResponse with _$ListRoutinesResponse {
  const factory ListRoutinesResponse({
    @Default(<RoutineDto>[]) List<RoutineDto> routines,
    PaginationDto? pagination,
  }) = _ListRoutinesResponse;

  factory ListRoutinesResponse.fromJson(Map<String, dynamic> json) =>
      _$ListRoutinesResponseFromJson(json);
}

@freezed
abstract class UpsertRoutineExerciseRequest with _$UpsertRoutineExerciseRequest {
  const factory UpsertRoutineExerciseRequest({
    required String exerciseId,
    required int displayOrder,
  }) = _UpsertRoutineExerciseRequest;

  factory UpsertRoutineExerciseRequest.fromJson(Map<String, dynamic> json) =>
      _$UpsertRoutineExerciseRequestFromJson(json);
}

@freezed
abstract class UpsertRoutineRequest with _$UpsertRoutineRequest {
  const factory UpsertRoutineRequest({
    required String name,
    @Default(<UpsertRoutineExerciseRequest>[])
    List<UpsertRoutineExerciseRequest> exercises,
  }) = _UpsertRoutineRequest;

  factory UpsertRoutineRequest.fromJson(Map<String, dynamic> json) =>
      _$UpsertRoutineRequestFromJson(json);
}
