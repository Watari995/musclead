import 'package:freezed_annotation/freezed_annotation.dart';

part 'notification_dtos.freezed.dart';
part 'notification_dtos.g.dart';

@freezed
abstract class NotificationDto with _$NotificationDto {
  const factory NotificationDto({
    required String id,
    @JsonKey(name: 'notification_type') required String notificationType,
    required Map<String, dynamic> metadata,
    @JsonKey(name: 'is_read') required bool isRead,
    @JsonKey(name: 'read_at') DateTime? readAt,
    @JsonKey(name: 'created_at') required DateTime createdAt,
  }) = _NotificationDto;

  factory NotificationDto.fromJson(Map<String, dynamic> json) =>
      _$NotificationDtoFromJson(json);
}

@freezed
abstract class GetNotificationsResponse with _$GetNotificationsResponse {
  const factory GetNotificationsResponse({
    required List<NotificationDto> notifications,
    @JsonKey(name: 'unread_count') required int unreadCount,
  }) = _GetNotificationsResponse;

  factory GetNotificationsResponse.fromJson(Map<String, dynamic> json) =>
      _$GetNotificationsResponseFromJson(json);
}

@freezed
abstract class WeeklyGoalDto with _$WeeklyGoalDto {
  const factory WeeklyGoalDto({
    @JsonKey(name: 'training_count') int? trainingCount,
    @JsonKey(name: 'calorie_average') int? calorieAverage,
    @JsonKey(name: 'weight_change_kg') double? weightChangeKg,
    @JsonKey(name: 'created_at') DateTime? createdAt,
    @JsonKey(name: 'updated_at') DateTime? updatedAt,
  }) = _WeeklyGoalDto;

  factory WeeklyGoalDto.fromJson(Map<String, dynamic> json) =>
      _$WeeklyGoalDtoFromJson(json);
}
