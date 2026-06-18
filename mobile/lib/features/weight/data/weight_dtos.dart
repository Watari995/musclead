import 'package:decimal/decimal.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../core/data/pagination.dart';
import '../../../core/util/datetime_converter.dart';
import '../../../core/util/decimal_converter.dart';

part 'weight_dtos.freezed.dart';
part 'weight_dtos.g.dart';

@freezed
abstract class WeightDto with _$WeightDto {
  const factory WeightDto({
    required String id,
    @DecimalStringConverter() required Decimal weightKg,
    required DateTime measuredAt,
    @NullableDecimalStringConverter() Decimal? bodyFatPercentage,
    @NullableDecimalStringConverter() Decimal? skeletalMuscleKg,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _WeightDto;

  factory WeightDto.fromJson(Map<String, dynamic> json) =>
      _$WeightDtoFromJson(json);
}

@freezed
abstract class ListWeightsResponse with _$ListWeightsResponse {
  const factory ListWeightsResponse({
    @Default(<WeightDto>[]) List<WeightDto> weights,
    PaginationDto? pagination,
  }) = _ListWeightsResponse;

  factory ListWeightsResponse.fromJson(Map<String, dynamic> json) =>
      _$ListWeightsResponseFromJson(json);
}

@freezed
abstract class TimeseriesWeightsResponse with _$TimeseriesWeightsResponse {
  const factory TimeseriesWeightsResponse({
    String? period,
    @Default(<WeightDto>[]) List<WeightDto> weights,
  }) = _TimeseriesWeightsResponse;

  factory TimeseriesWeightsResponse.fromJson(Map<String, dynamic> json) =>
      _$TimeseriesWeightsResponseFromJson(json);
}

@freezed
abstract class UpsertWeightRequest with _$UpsertWeightRequest {
  const factory UpsertWeightRequest({
    @DecimalStringConverter() required Decimal weightKg,
    @UtcDateTimeConverter() required DateTime measuredAt,
    @NullableDecimalStringConverter() Decimal? bodyFatPercentage,
    @NullableDecimalStringConverter() Decimal? skeletalMuscleKg,
  }) = _UpsertWeightRequest;

  factory UpsertWeightRequest.fromJson(Map<String, dynamic> json) =>
      _$UpsertWeightRequestFromJson(json);
}
