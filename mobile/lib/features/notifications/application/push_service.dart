import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// FCM プッシュ通知サービスの雛形 — **Phase 3**。
///
/// 前提（あなたのコンソール作業）:
/// - Firebase プロジェクト + APNs 認証鍵の設定
/// - バックエンドに device token 登録エンドポイントと配信処理を追加（現状未実装）
///
/// 通知基盤が整うまでアプリ起動時には呼ばない。Phase 3 で有効化する足場。
class PushService {
  PushService(this._messaging);

  final FirebaseMessaging _messaging;

  /// 権限要求 → FCM トークン取得 → 前面メッセージ購読。トークンはサーバ登録する。
  Future<String?> initAndGetToken() async {
    final settings = await _messaging.requestPermission();
    if (settings.authorizationStatus == AuthorizationStatus.denied) {
      return null;
    }
    FirebaseMessaging.onMessage.listen((message) {
      debugPrint('FCM foreground: ${message.notification?.title}');
    });
    final token = await _messaging.getToken();
    // TODO(phase3): token をサーバの device-token 登録 API に送る
    return token;
  }
}

final pushServiceProvider = Provider<PushService>(
  (ref) => PushService(FirebaseMessaging.instance),
);
