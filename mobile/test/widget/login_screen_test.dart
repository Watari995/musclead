import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:musclead/core/theme/app_theme.dart';
import 'package:musclead/core/theme/app_tokens.dart';
import 'package:musclead/features/auth/presentation/login_screen.dart';
import 'package:musclead/l10n/app_localizations.dart';

void main() {
  testWidgets('LoginScreen がロゴ・入力欄・ボタンを表示する', (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        child: MaterialApp(
          theme: buildAppTheme(Brightness.light, kBrandAccent),
          locale: const Locale('ja'),
          localizationsDelegates: const [
            AppLocalizations.delegate,
            GlobalMaterialLocalizations.delegate,
            GlobalWidgetsLocalizations.delegate,
          ],
          supportedLocales: AppLocalizations.supportedLocales,
          home: const LoginScreen(),
        ),
      ),
    );
    await tester.pumpAndSettle();

    expect(find.text('musclead'), findsOneWidget);
    expect(find.byType(TextField), findsNWidgets(2));
    expect(find.text('ログイン'), findsWidgets);
    expect(find.text('アカウントを作成'), findsOneWidget);
  });
}
