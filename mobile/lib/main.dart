import 'package:cookie_jar/cookie_jar.dart';
import 'package:flutter/widgets.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:path_provider/path_provider.dart';

import 'app.dart';
import 'bootstrap/firebase_bootstrap.dart';
import 'core/providers/core_providers.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // refresh token（HttpOnly Cookie）をアプリサンドボックス内に永続化する。
  final supportDir = await getApplicationSupportDirectory();
  final cookieJar = PersistCookieJar(
    storage: FileStorage('${supportDir.path}/.cookies'),
  );

  // Firebase（Crashlytics/Analytics）。未設定環境でも起動を止めない。
  await initFirebase();

  runApp(
    ProviderScope(
      overrides: [cookieJarProvider.overrideWithValue(cookieJar)],
      child: const MuscleadApp(),
    ),
  );
}
