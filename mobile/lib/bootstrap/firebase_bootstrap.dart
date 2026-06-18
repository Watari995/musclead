import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_crashlytics/firebase_crashlytics.dart';
import 'package:flutter/foundation.dart';

/// Firebase（Crashlytics / Analytics）を初期化する。
///
/// GoogleService-Info.plist / firebase_options が未設定の環境では例外を
/// 握りつぶして起動を継続する（dev での利便性 + CI/test での安全性）。
/// 本番では `flutterfire configure` で設定を生成すること。
Future<void> initFirebase() async {
  try {
    await Firebase.initializeApp();
    final crashlytics = FirebaseCrashlytics.instance;
    await crashlytics.setCrashlyticsCollectionEnabled(!kDebugMode);
    FlutterError.onError = crashlytics.recordFlutterFatalError;
    PlatformDispatcher.instance.onError = (error, stack) {
      crashlytics.recordError(error, stack, fatal: true);
      return true;
    };
  } catch (e) {
    debugPrint('Firebase 初期化をスキップ（未設定）: $e');
  }
}
