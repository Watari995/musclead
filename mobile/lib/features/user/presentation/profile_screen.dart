import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/theme_controller.dart';
import '../../../core/widgets/app_badge.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../../auth/application/auth_controller.dart';
import '../../subscription/data/subscription_repository.dart';
import '../data/user_dtos.dart';
import '../data/user_repository.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final me = ref.watch(meProvider);
    return TabPage(
      title: 'マイページ',
      onRefresh: () async {
        ref.invalidate(meProvider);
        ref.invalidate(subscriptionProvider);
      },
      children: [
        AsyncValueView<MeResponse>(
          value: me,
          onRetry: () => ref.invalidate(meProvider),
          data: (data) => _ProfileBody(user: data.user),
        ),
      ],
    );
  }
}

class _ProfileBody extends ConsumerWidget {
  const _ProfileBody({required this.user});

  final UserDto user;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = context.tokens;
    final sub = ref.watch(subscriptionProvider);
    final isPro = sub.asData?.value.isPro ?? false;
    final mode = ref.watch(themeModeProvider);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        AppCard(
          child: Row(
            children: [
              CircleAvatar(
                radius: 23,
                backgroundColor: t.accent,
                child: Text(
                  user.name.isNotEmpty ? user.name.characters.first : '?',
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w700,
                    fontSize: 16,
                  ),
                ),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      user.name,
                      style: const TextStyle(
                        fontWeight: FontWeight.w700,
                        fontSize: 16,
                      ),
                    ),
                    Text(
                      user.email,
                      style: TextStyle(fontSize: 12, color: t.muted),
                    ),
                  ],
                ),
              ),
              AppBadge(
                isPro ? 'Pro' : 'Free',
                tone: isPro ? BadgeTone.accent : BadgeTone.neutral,
              ),
            ],
          ),
        ),
        const SectionTitle('アカウント'),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => _showThemeSheet(context, ref, mode),
              child: _row(context, '外観', value: _modeLabel(mode)),
            ),
            AppListRow(
              onTap: () => context.push('/plan'),
              child: _row(context, 'プラン', value: isPro ? 'Pro' : 'Free'),
            ),
          ],
        ),
        const SectionTitle('その他'),
        AppListBox(
          children: [
            AppListRow(
              child: _row(context, 'バージョン', value: '1.0.0 (1)', chevron: false),
            ),
          ],
        ),
        const SizedBox(height: 18),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => ref.read(authControllerProvider.notifier).logout(),
              child: Center(
                child: Text(
                  'ログアウト',
                  style: TextStyle(
                    color: context.colors.onSurface,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ],
        ),
        const SizedBox(height: 10),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => _confirmDelete(context, ref, user.id),
              child: Center(
                child: Text(
                  'アカウントを削除',
                  style: TextStyle(
                    color: context.tokens.accent,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Future<void> _confirmDelete(
    BuildContext context,
    WidgetRef ref,
    String userId,
  ) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('アカウントを削除'),
        content: const Text('アカウントとすべてのデータが削除されます。この操作は取り消せません。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: const Text('キャンセル'),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: Text('削除する', style: TextStyle(color: context.tokens.accent)),
          ),
        ],
      ),
    );
    if (ok != true) return;
    try {
      await ref.read(userRepositoryProvider).deleteAccount(userId);
      await ref.read(authControllerProvider.notifier).logout();
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('削除に失敗しました')));
      }
    }
  }

  Widget _row(
    BuildContext context,
    String label, {
    String? value,
    bool chevron = true,
  }) {
    final t = context.tokens;
    return Row(
      children: [
        Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
        const Spacer(),
        if (value != null)
          Text(value, style: TextStyle(fontSize: 13, color: t.muted)),
        if (chevron) ...[
          const SizedBox(width: 6),
          Icon(Icons.chevron_right, color: t.subtle, size: 18),
        ],
      ],
    );
  }

  static String _modeLabel(ThemeMode m) => switch (m) {
    ThemeMode.system => 'システム',
    ThemeMode.light => 'ライト',
    ThemeMode.dark => 'ダーク',
  };

  void _showThemeSheet(BuildContext context, WidgetRef ref, ThemeMode current) {
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            for (final m in ThemeMode.values)
              ListTile(
                title: Text(_modeLabel(m)),
                trailing: m == current
                    ? Icon(Icons.check, color: sheetContext.tokens.accent)
                    : null,
                onTap: () {
                  ref.read(themeModeProvider.notifier).set(m);
                  final apiTheme = switch (m) {
                    ThemeMode.system => 'system',
                    ThemeMode.light => 'light',
                    ThemeMode.dark => 'dark',
                  };
                  // ベストエフォートでサーバ設定も更新
                  ref
                      .read(userRepositoryProvider)
                      .updateTheme(apiTheme)
                      .ignore();
                  Navigator.of(sheetContext).pop();
                },
              ),
          ],
        ),
      ),
    );
  }
}
