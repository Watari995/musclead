import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../data/notification_repository.dart';

/// FCM プッシュ通知サービス。
///
/// `app.dart` でログイン済み判定後に `initAndGetToken()` を呼び、
/// device token をサーバに登録する。配信(outbox relay)はバックエンド側で処理される。
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
