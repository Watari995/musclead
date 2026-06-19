import 'package:cookie_jar/cookie_jar.dart';
import 'package:flutter/widgets.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sentry_flutter/sentry_flutter.dart';

import 'app.dart';
import 'bootstrap/firebase_bootstrap.dart';
import 'core/providers/core_providers.dart';

const _sentryDsn = String.fromEnvironment('SENTRY_DSN');

Future<void> main() async {
  await SentryFlutter.init(
    (options) {
      options.dsn = _sentryDsn;
      options.tracesSampleRate = 1.0;
      options.environment = const String.fromEnvironment(
        'FLAVOR',
        defaultValue: 'dev',
      );
    },
    appRunner: () async {
      WidgetsFlutterBinding.ensureInitialized();

      final supportDir = await getApplicationSupportDirectory();
      final cookieJar = PersistCookieJar(
        storage: FileStorage('${supportDir.path}/.cookies'),
      );

      await initFirebase();

      runApp(
        ProviderScope(
          overrides: [cookieJarProvider.overrideWithValue(cookieJar)],
          child: const MuscleadApp(),
        ),
      );
    },
  );
}
