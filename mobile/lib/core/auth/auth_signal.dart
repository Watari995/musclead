import 'package:flutter/foundation.dart';

/// dio の認証層（低レベル）と認証状態コントローラ（高レベル）を疎結合にするシグナル。
///
/// access token の refresh も失敗してセッションが切れたとき [expire] が呼ばれ、
/// 認証コントローラがこれを購読して未認証状態へ遷移する。dio → controller の
/// 循環依存を避けるための仕掛け。
class AuthSignal extends ChangeNotifier {
  void expire() => notifyListeners();
}
