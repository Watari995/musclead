import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../theme/app_tokens.dart';
import 'app_button.dart';

/// AsyncValue を loading / error / data に振り分けて描画する共通 view。
class AsyncValueView<T> extends StatelessWidget {
  const AsyncValueView({
    super.key,
    required this.value,
    required this.data,
    this.onRetry,
  });

  final AsyncValue<T> value;
  final Widget Function(T data) data;
  final VoidCallback? onRetry;

  @override
  Widget build(BuildContext context) {
    return value.when(
      skipLoadingOnReload: true,
      skipLoadingOnRefresh: true,
      data: data,
      loading: () => const Center(
        child: Padding(
          padding: EdgeInsets.all(40),
          child: CircularProgressIndicator(strokeWidth: 2.4),
        ),
      ),
      error: (e, _) => Center(
        child: Padding(
          padding: const EdgeInsets.all(28),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.error_outline, color: context.tokens.muted, size: 34),
              const SizedBox(height: 10),
              Text(
                _message(e),
                textAlign: TextAlign.center,
                style: TextStyle(color: context.tokens.muted),
              ),
              if (onRetry != null) ...[
                const SizedBox(height: 16),
                AppButton(
                  label: '再試行',
                  variant: AppButtonVariant.glass,
                  expand: false,
                  onPressed: onRetry,
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  String _message(Object e) {
    final s = e.toString();
    return s.length > 120 ? '読み込みに失敗しました' : s;
  }
}
