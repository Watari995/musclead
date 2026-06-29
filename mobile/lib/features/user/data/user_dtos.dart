import 'package:freezed_annotation/freezed_annotation.dart';

part 'user_dtos.freezed.dart';
part 'user_dtos.g.dart';

@freezed
abstract class UserDto with _$UserDto {
  const factory UserDto({
    required String id,
    required String name,
    required String email,
    String? birthday,
    String? profileImageUrl,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _UserDto;

  factory UserDto.fromJson(Map<String, dynamic> json) =>
      _$UserDtoFromJson(json);
}

@freezed
abstract class PreferencesDto with _$PreferencesDto {
  const factory PreferencesDto({
    @Default('system') String theme,
    @JsonKey(name: 'training_color') @Default('#4A90E2') String trainingColor,
    @JsonKey(name: 'meal_color') @Default('#7ED321') String mealColor,
    @JsonKey(name: 'weight_color') @Default('#FF6B6B') String weightColor,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _PreferencesDto;

  factory PreferencesDto.fromJson(Map<String, dynamic> json) =>
      _$PreferencesDtoFromJson(json);
}

@freezed
abstract class MeResponse with _$MeResponse {
  const factory MeResponse({
    required UserDto user,
    PreferencesDto? preferences,
  }) = _MeResponse;

  factory MeResponse.fromJson(Map<String, dynamic> json) =>
      _$MeResponseFromJson(json);
}

@freezed
abstract class RegisterRequest with _$RegisterRequest {
  const factory RegisterRequest({
    required String name,
    required String email,
    required String password,
    String? birthday,
  }) = _RegisterRequest;

  factory RegisterRequest.fromJson(Map<String, dynamic> json) =>
      _$RegisterRequestFromJson(json);
}

@freezed
abstract class RegisterResponse with _$RegisterResponse {
  const factory RegisterResponse({String? userId}) = _RegisterResponse;

  factory RegisterResponse.fromJson(Map<String, dynamic> json) =>
      _$RegisterResponseFromJson(json);
}
