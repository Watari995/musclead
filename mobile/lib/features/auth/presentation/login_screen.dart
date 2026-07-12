import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../application/auth_controller.dart';
import '../../../l10n/app_localizations.dart';

class LoginScreen extends HookConsumerWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final email = useTextEditingController();
    final password = useTextEditingController();
    final loading = useState(false);
    final error = useState<String?>(null);
    final t = context.tokens;

    Future<void> submit() async {
      if (loading.value) return;
      loading.value = true;
      error.value = null;
      try {
        await ref
            .read(authControllerProvider.notifier)
            .login(email.text.trim(), password.text);
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = l.loginFailed;
      } finally {
        if (context.mounted) loading.value = false;
      }
    }

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 40),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Center(
                  child: RoughBox(
                    fill: t.accent,
                    color: t.accent,
                    radius: BorderRadius.circular(18),
                    child: SizedBox(
                      width: 62,
                      height: 62,
                      child: Icon(
                        Icons.fitness_center,
                        color: context.colors.onPrimary,
                        size: 34,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 18),
                const Center(
                  child: Text(
                    'musclead',
                    style: TextStyle(
                      fontFamily: 'Architects Daughter',
                      fontSize: 26,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                const SizedBox(height: 6),
                Center(
                  child: Text(
                    l.loginTagline,
                    style: TextStyle(fontSize: 13, color: t.muted),
                  ),
                ),
                const SizedBox(height: 34),
                AppTextField(
                  label: l.loginEmailLabel,
                  controller: email,
                  hint: 'you@example.com',
                  keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  autofillHints: const [AutofillHints.email],
                ),
                const SizedBox(height: 14),
                AppTextField(
                  label: l.loginPasswordLabel,
                  controller: password,
                  hint: '••••••••',
                  obscureText: true,
                  textInputAction: TextInputAction.done,
                  autofillHints: const [AutofillHints.password],
                  onSubmitted: (_) => submit(),
                ),
                if (error.value != null) ...[
                  const SizedBox(height: 12),
                  Text(
                    error.value!,
                    style: TextStyle(color: t.accent, fontSize: 13),
                  ),
                ],
                const SizedBox(height: 22),
                AppButton(
                  label: l.loginTitle,
                  loading: loading.value,
                  onPressed: submit,
                ),
                const SizedBox(height: 14),
                AppButton(
                  label: l.loginCreateAccountBtn,
                  variant: AppButtonVariant.text,
                  onPressed: () => context.go('/register'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
