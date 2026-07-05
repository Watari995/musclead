import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../data/notification_repository.dart';

/// FCM プッシュ通知サービスの雛形 — **Phase 3**。
///
/// device token のサーバ登録までは実装済み。配信処理(outbox relay)は
/// まだバックエンド未実装のため、通知基盤が整うまでアプリ起動時には呼ばない。
class PushService {
  PushService(this._messaging, this._notificationRepository);

  final FirebaseMessaging _messaging;
  final NotificationRepository _notificationRepository;

  /// 権限要求 → FCM トークン取得 → サーバ登録 → 前面メッセージ購読。
  Future<String?> initAndGetToken() async {
    final settings = await _messaging.requestPermission();
    if (settings.authorizationStatus == AuthorizationStatus.denied) {
      return null;
    }
    FirebaseMessaging.onMessage.listen((message) {
      debugPrint('FCM foreground: ${message.notification?.title}');
    });
    final token = await _messaging.getToken();
    if (token != null) {
      await _notificationRepository.registerDeviceToken(
        token: token,
        platform: defaultTargetPlatform == TargetPlatform.iOS
            ? 'ios'
            : 'android',
      );
    }
    return token;
  }
}

final pushServiceProvider = Provider<PushService>(
  (ref) => PushService(
    FirebaseMessaging.instance,
    ref.watch(notificationRepositoryProvider),
  ),
);
