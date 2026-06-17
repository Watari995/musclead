import 'package:freezed_annotation/freezed_annotation.dart';

part 'pagination.freezed.dart';
part 'pagination.g.dart';

@freezed
abstract class PaginationDto with _$PaginationDto {
  const factory PaginationDto({
    @Default(1) int currentPage,
    @Default(0) int itemsPerPage,
    @Default(0) int totalItems,
    @Default(0) int totalPages,
  }) = _PaginationDto;

  factory PaginationDto.fromJson(Map<String, dynamic> json) =>
      _$PaginationDtoFromJson(json);
}
