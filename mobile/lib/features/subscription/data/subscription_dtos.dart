import 'package:freezed_annotation/freezed_annotation.dart';

part 'subscription_dtos.freezed.dart';
part 'subscription_dtos.g.dart';

@freezed
abstract class GetSubscriptionResponse with _$GetSubscriptionResponse {
  const factory GetSubscriptionResponse({
    @Default(false) bool isPro,
    @Default('free') String plan,
    DateTime? expiresAt,
  }) = _GetSubscriptionResponse;

  factory GetSubscriptionResponse.fromJson(Map<String, dynamic> json) =>
      _$GetSubscriptionResponseFromJson(json);
}

@freezed
abstract class SubscribeResponse with _$SubscribeResponse {
  const factory SubscribeResponse({String? checkoutUrl}) = _SubscribeResponse;

  factory SubscribeResponse.fromJson(Map<String, dynamic> json) =>
      _$SubscribeResponseFromJson(json);
}

@freezed
abstract class CreatePortalSessionResponse with _$CreatePortalSessionResponse {
  const factory CreatePortalSessionResponse({String? portalUrl}) =
      _CreatePortalSessionResponse;

  factory CreatePortalSessionResponse.fromJson(Map<String, dynamic> json) =>
      _$CreatePortalSessionResponseFromJson(json);
}
