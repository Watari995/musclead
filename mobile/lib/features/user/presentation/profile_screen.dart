import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import 'package:package_info_plus/package_info_plus.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../../core/providers/core_providers.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/theme_controller.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../../auth/application/auth_controller.dart';
import '../data/user_dtos.dart';
import '../data/user_repository.dart';

final _appVersionProvider = FutureProvider<String>((ref) async {
  final info = await PackageInfo.fromPlatform();
  return '${info.version} (${info.buildNumber})';
});

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.listen<Uri?>(oauthCallbackProvider, (_, uri) {
      if (uri == null || !context.mounted) return;
      final connected = uri.queryParameters['connected'] == 'true';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(
            connected ? 'HealthPlanet を連携しました' : 'HealthPlanet 連携に失敗しました',
          ),
        ),
      );
      if (connected) ref.invalidate(meProvider);
      ref.read(oauthCallbackProvider.notifier).set(null);
    });

    final me = ref.watch(meProvider);
    return TabPage(
      title: 'マイページ',
      onRefresh: () async => ref.invalidate(meProvider),
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
    final mode = ref.watch(themeModeProvider);
    final accent = ref.watch(accentProvider);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        AppCard(
          child: Row(
            children: [
              _Avatar(user: user, onTap: () => _pickAndUpload(context, ref)),
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
              onTap: () => _showAccentSheet(context, ref),
              child: Row(
                children: [
                  const Text(
                    'テーマカラー',
                    style: TextStyle(fontWeight: FontWeight.w500),
                  ),
                  const Spacer(),
                  Container(
                    width: 18,
                    height: 18,
                    decoration: BoxDecoration(
                      color: accent,
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 6),
                  Icon(Icons.chevron_right, color: t.subtle, size: 18),
                ],
              ),
            ),
          ],
        ),
        const SectionTitle('連携サービス'),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => _connectHealthPlanet(context, ref),
              child: _row(context, 'Tanita HealthPlanet'),
            ),
          ],
        ),
        const SectionTitle('その他'),
        AppListBox(
          children: [
            AppListRow(
              child: _row(
                context,
                'バージョン',
                value: ref
                    .watch(_appVersionProvider)
                    .when(
                      data: (v) => v,
                      loading: () => '...',
                      error: (e, _) => '-',
                    ),
                chevron: false,
              ),
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

  Future<void> _connectHealthPlanet(BuildContext context, WidgetRef ref) async {
    try {
      final dio = ref.read(dioProvider);
      final res = await dio.get<Map<String, dynamic>>(
        '/integrations/healthplanet/auth',
        queryParameters: {'redirect_url': 'musclead://oauth/callback'},
      );
      final url = res.data?['url'] as String?;
      if (url == null) throw Exception('url not found');
      await launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication);
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('HealthPlanet 連携の開始に失敗しました')),
        );
      }
    }
  }

  Future<void> _pickAndUpload(BuildContext context, WidgetRef ref) async {
    final file = await ImagePicker().pickImage(
      source: ImageSource.gallery,
      maxWidth: 512,
      maxHeight: 512,
      imageQuality: 80, // iOS では JPEG に再エンコードされる
    );
    if (file == null) return;
    final bytes = await file.readAsBytes();
    try {
      await ref
          .read(userRepositoryProvider)
          .uploadProfileImage(bytes, 'image/jpeg');
      ref.invalidate(meProvider);
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('プロフィール画像を更新しました')));
      }
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('画像の更新に失敗しました')));
      }
    }
  }

  void _showAccentSheet(BuildContext context, WidgetRef ref) {
    final current = ref.read(accentProvider);
    showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (sheetContext) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(20, 4, 20, 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'テーマカラー',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 18),
              Wrap(
                spacing: 18,
                runSpacing: 18,
                children: [
                  for (final c in kAccentPresets)
                    GestureDetector(
                      onTap: () {
                        ref.read(accentProvider.notifier).set(c);
                        Navigator.of(sheetContext).pop();
                      },
                      child: Container(
                        width: 48,
                        height: 48,
                        decoration: BoxDecoration(
                          color: c,
                          shape: BoxShape.circle,
                          border: Border.all(
                            color: c.toARGB32() == current.toARGB32()
                                ? sheetContext.colors.onSurface
                                : Colors.transparent,
                            width: 3,
                          ),
                        ),
                        child: c.toARGB32() == current.toARGB32()
                            ? const Icon(Icons.check, color: Colors.white)
                            : null,
                      ),
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
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

class _Avatar extends StatelessWidget {
  const _Avatar({required this.user, required this.onTap});

  final UserDto user;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final url = user.profileImageUrl;
    final initial = user.name.isNotEmpty ? user.name.characters.first : '?';
    // 読み込み中の中立プレースホルダ（色付きイニシャルを出さないことでチラつきを防ぐ）。
    Widget placeholder() => Container(width: 52, height: 52, color: t.border);
    // 通常はサーバーが必ずデフォルト画像 URL を返すため、これは読込失敗時のみ表示。
    Widget fallback() => Container(
      width: 52,
      height: 52,
      color: t.accent,
      alignment: Alignment.center,
      child: Text(
        initial,
        style: const TextStyle(
          color: Colors.white,
          fontWeight: FontWeight.w700,
          fontSize: 18,
        ),
      ),
    );

    return GestureDetector(
      onTap: onTap,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          ClipOval(
            child: (url != null && url.isNotEmpty)
                ? CachedNetworkImage(
                    imageUrl: url,
                    width: 52,
                    height: 52,
                    fit: BoxFit.cover,
                    // フェードを無効化し、キャッシュ済みなら即時表示（チラつき防止）。
                    fadeInDuration: Duration.zero,
                    fadeOutDuration: Duration.zero,
                    placeholder: (_, _) => placeholder(),
                    errorWidget: (_, _, _) => fallback(),
                  )
                : fallback(),
          ),
          Positioned(
            right: -2,
            bottom: -2,
            child: Container(
              padding: const EdgeInsets.all(4),
              decoration: BoxDecoration(
                color: t.accent,
                shape: BoxShape.circle,
                border: Border.all(color: context.colors.surface, width: 2),
              ),
              child: const Icon(
                Icons.camera_alt,
                size: 11,
                color: Colors.white,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
