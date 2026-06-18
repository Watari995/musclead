import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../application/auth_controller.dart';

class LoginScreen extends HookConsumerWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
        // 認証状態の変化で go_router が自動的に /meals へ遷移する
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = 'ログインに失敗しました';
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
                  child: Container(
                    width: 62,
                    height: 62,
                    decoration: BoxDecoration(
                      color: t.accent,
                      borderRadius: BorderRadius.circular(18),
                    ),
                    child: const Icon(
                      Icons.fitness_center,
                      color: Colors.white,
                      size: 34,
                    ),
                  ),
                ),
                const SizedBox(height: 18),
                const Center(
                  child: Text(
                    'musclead',
                    style: TextStyle(fontSize: 25, fontWeight: FontWeight.w800),
                  ),
                ),
                const SizedBox(height: 6),
                Center(
                  child: Text(
                    '筋トレ・食事・体重を一元管理',
                    style: TextStyle(fontSize: 13, color: t.muted),
                  ),
                ),
                const SizedBox(height: 34),
                AppTextField(
                  label: 'メールアドレス',
                  controller: email,
                  hint: 'you@example.com',
                  keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  autofillHints: const [AutofillHints.email],
                ),
                const SizedBox(height: 14),
                AppTextField(
                  label: 'パスワード',
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
                  label: 'ログイン',
                  loading: loading.value,
                  onPressed: submit,
                ),
                const SizedBox(height: 14),
                AppButton(
                  label: 'アカウントを作成',
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
