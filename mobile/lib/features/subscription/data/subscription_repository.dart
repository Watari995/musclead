import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'subscription_dtos.dart';

class SubscriptionRepository {
  SubscriptionRepository(this._dio);

  final Dio _dio;

  /// 現在のサブスク状態（is_pro / plan / expires_at）を取得。
  /// iOS では購入導線は持たず、この状態のみを参照する（App Store 3.1.1 対応）。
  Future<GetSubscriptionResponse> get() => guardApi(() async {
    final res = await _dio.get<Map<String, dynamic>>('/purchase/subscription');
    return GetSubscriptionResponse.fromJson(res.data!);
  });
}

final subscriptionRepositoryProvider = Provider<SubscriptionRepository>(
  (ref) => SubscriptionRepository(ref.watch(dioProvider)),
);

final subscriptionProvider = FutureProvider<GetSubscriptionResponse>(
  (ref) => ref.watch(subscriptionRepositoryProvider).get(),
);
