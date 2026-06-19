import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/foundation.dart';

/// Firebase（Analytics）を初期化する。エラー監視は Sentry が担うため
/// Crashlytics のエラーフックは設定しない。
///
/// GoogleService-Info.plist / firebase_options が未設定の環境では例外を
/// 握りつぶして起動を継続する（dev での利便性 + CI/test での安全性）。
Future<void> initFirebase() async {
  try {
    await Firebase.initializeApp();
  } catch (e) {
    debugPrint('Firebase 初期化をスキップ（未設定）: $e');
  }
}
