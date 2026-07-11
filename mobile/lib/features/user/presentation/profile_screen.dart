import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../../core/providers/core_providers.dart';
import '../../../core/providers/locale_provider.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/theme_controller.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../../auth/application/auth_controller.dart';
import '../data/user_dtos.dart';
import '../data/user_repository.dart';
import '../../../l10n/app_localizations.dart';

class ProfileScreen extends ConsumerWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    ref.listen<Uri?>(oauthCallbackProvider, (_, uri) {
      if (uri == null || !context.mounted) return;
      final connected = uri.queryParameters['connected'] == 'true';
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(connected ? l.profileHpConnected : l.profileHpFailed),
        ),
      );
      if (connected) ref.invalidate(meProvider);
      ref.read(oauthCallbackProvider.notifier).set(null);
    });

    final me = ref.watch(meProvider);
    return TabPage(
      title: l.profileTitle,
      onRefresh: () async => ref.invalidate(meProvider),
      children: [
        AsyncValueView<MeResponse>(
          value: me,
          onRetry: () => ref.invalidate(meProvider),
          data: (data) =>
              _ProfileBody(user: data.user, preferences: data.preferences),
        ),
      ],
    );
  }
}

class _ProfileBody extends ConsumerWidget {
  const _ProfileBody({required this.user, this.preferences});

  final UserDto user;
  final PreferencesDto? preferences;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final mode = ref.watch(themeModeProvider);
    final accent = ref.watch(accentProvider);
    final locale = ref.watch(localeProvider);
    final versionText = ref
        .watch(packageInfoProvider)
        .when(
          data: (info) => '${info.version} (${info.buildNumber})',
          loading: () => '...',
          error: (_, _) => '-',
        );

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
        SectionTitle(l.profileAccount),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => _showThemeSheet(context, ref, mode),
              child: _row(
                context,
                l.profileAppearance,
                value: _modeLabel(l, mode),
              ),
            ),
            AppListRow(
              onTap: () => _showAccentSheet(context, ref),
              child: Row(
                children: [
                  Text(
                    l.profileThemeColor,
                    style: const TextStyle(fontWeight: FontWeight.w500),
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
            AppListRow(
              onTap: () => _showLanguageSheet(context, ref, locale),
              child: _row(
                context,
                l.profileLanguage,
                value: locale.languageCode == 'ja'
                    ? l.profileLanguageJa
                    : l.profileLanguageEn,
              ),
            ),
          ],
        ),
        if (preferences != null) ...[
          SectionTitle(l.profileCalendar),
          _CalendarColorSection(prefs: preferences!),
        ],
        SectionTitle(l.profileGoal),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => context.push('/weekly-goal'),
              child: _row(context, l.profileWeeklyGoal),
            ),
          ],
        ),
        SectionTitle(l.profileIntegrations),
        AppListBox(
          children: [
            AppListRow(
              onTap: () => _connectHealthPlanet(context, ref),
              child: _row(context, l.profileTanitaHp),
            ),
          ],
        ),
        SectionTitle(l.profileOther),
        AppListBox(
          children: [
            AppListRow(
              child: _row(
                context,
                l.profileVersion,
                value: versionText,
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
                  l.profileLogout,
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
                  l.profileDeleteAccount,
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
    final l = AppLocalizations.of(context)!;
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
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(l.profileHpStartFailed)));
      }
    }
  }

  Future<void> _pickAndUpload(BuildContext context, WidgetRef ref) async {
    final l = AppLocalizations.of(context)!;
    final file = await ImagePicker().pickImage(
      source: ImageSource.gallery,
      maxWidth: 512,
      maxHeight: 512,
      imageQuality: 80,
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
        ).showSnackBar(SnackBar(content: Text(l.profileAvatarSuccess)));
      }
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(l.profileAvatarFailed)));
      }
    }
  }

  void _showAccentSheet(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
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
              Text(
                l.profileThemeColor,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w800,
                ),
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

  void _showLanguageSheet(BuildContext context, WidgetRef ref, Locale current) {
    final l = AppLocalizations.of(context)!;
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            for (final loc in const [Locale('ja'), Locale('en')])
              ListTile(
                title: Text(
                  loc.languageCode == 'ja'
                      ? l.profileLanguageJa
                      : l.profileLanguageEn,
                ),
                trailing: loc.languageCode == current.languageCode
                    ? Icon(Icons.check, color: sheetContext.tokens.accent)
                    : null,
                onTap: () {
                  ref.read(localeProvider.notifier).set(loc);
                  Navigator.of(sheetContext).pop();
                },
              ),
          ],
        ),
      ),
    );
  }

  Future<void> _confirmDelete(
    BuildContext context,
    WidgetRef ref,
    String userId,
  ) async {
    final l = AppLocalizations.of(context)!;
    final ok = await showDialog<bool>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: Text(l.profileDeleteAccount),
        content: Text(l.profileDeleteAccountMsg),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(false),
            child: Text(l.commonCancel),
          ),
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(true),
            child: Text(
              l.commonDeleteOk,
              style: TextStyle(color: context.tokens.accent),
            ),
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
        ).showSnackBar(SnackBar(content: Text(l.profileDeleteFailed)));
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

  static String _modeLabel(AppLocalizations l, ThemeMode m) => switch (m) {
    ThemeMode.system => l.profileThemeSystem,
    ThemeMode.light => l.profileThemeLight,
    ThemeMode.dark => l.profileThemeDark,
  };

  void _showThemeSheet(BuildContext context, WidgetRef ref, ThemeMode current) {
    final l = AppLocalizations.of(context)!;
    showModalBottomSheet<void>(
      context: context,
      builder: (sheetContext) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            for (final m in ThemeMode.values)
              ListTile(
                title: Text(_modeLabel(l, m)),
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

class _CalendarColorSection extends ConsumerStatefulWidget {
  const _CalendarColorSection({required this.prefs});
  final PreferencesDto prefs;

  @override
  ConsumerState<_CalendarColorSection> createState() =>
      _CalendarColorSectionState();
}

class _CalendarColorSectionState extends ConsumerState<_CalendarColorSection> {
  bool _saving = false;

  static const _presets = [
    Color(0xFF4A90E2),
    Color(0xFF7ED321),
    Color(0xFFFF6B6B),
    Color(0xFFF5A623),
    Color(0xFFBD10E0),
    Color(0xFF50E3C2),
    Color(0xFFB8E986),
    Color(0xFF9013FE),
  ];

  String _colorToHex(Color c) {
    final r = (c.r * 255.0).round().clamp(0, 255);
    final g = (c.g * 255.0).round().clamp(0, 255);
    final b = (c.b * 255.0).round().clamp(0, 255);
    return '#${r.toRadixString(16).padLeft(2, '0')}${g.toRadixString(16).padLeft(2, '0')}${b.toRadixString(16).padLeft(2, '0')}';
  }

  Color _parseColor(String hex) {
    try {
      return Color(int.parse('FF${hex.replaceAll('#', '')}', radix: 16));
    } catch (_) {
      return const Color(0xFF4A90E2);
    }
  }

  Future<void> _updateColor(String field, Color color) async {
    setState(() => _saving = true);
    try {
      final repo = ref.read(userRepositoryProvider);
      await repo.updateCalendarColors(
        trainingColor: field == 'training' ? _colorToHex(color) : null,
        mealColor: field == 'meal' ? _colorToHex(color) : null,
        weightColor: field == 'weight' ? _colorToHex(color) : null,
      );
      ref.invalidate(meProvider);
    } catch (_) {
      // ignore - UI stays optimistic
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final trainingColor = _parseColor(widget.prefs.trainingColor);
    final mealColor = _parseColor(widget.prefs.mealColor);
    final weightColor = _parseColor(widget.prefs.weightColor);

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l.profileCalendarColors,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: t.muted,
            ),
          ),
          const SizedBox(height: 12),
          _ColorRow(
            label: l.calendarTraining,
            currentColor: trainingColor,
            presets: _presets,
            onSelect: (c) => _updateColor('training', c),
            disabled: _saving,
          ),
          const SizedBox(height: 10),
          _ColorRow(
            label: l.calendarMeal,
            currentColor: mealColor,
            presets: _presets,
            onSelect: (c) => _updateColor('meal', c),
            disabled: _saving,
          ),
          const SizedBox(height: 10),
          _ColorRow(
            label: l.calendarWeight,
            currentColor: weightColor,
            presets: _presets,
            onSelect: (c) => _updateColor('weight', c),
            disabled: _saving,
          ),
        ],
      ),
    );
  }
}

class _ColorRow extends StatelessWidget {
  const _ColorRow({
    required this.label,
    required this.currentColor,
    required this.presets,
    required this.onSelect,
    required this.disabled,
  });

  final String label;
  final Color currentColor;
  final List<Color> presets;
  final ValueChanged<Color> onSelect;
  final bool disabled;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Row(
      children: [
        Container(
          width: 18,
          height: 18,
          decoration: BoxDecoration(
            color: currentColor,
            shape: BoxShape.circle,
            border: Border.all(color: t.border),
          ),
        ),
        const SizedBox(width: 8),
        SizedBox(
          width: 72,
          child: Text(label, style: const TextStyle(fontSize: 13)),
        ),
        Expanded(
          child: Wrap(
            spacing: 6,
            children: presets.map((c) {
              final selected = c.toARGB32() == currentColor.toARGB32();
              return GestureDetector(
                onTap: disabled ? null : () => onSelect(c),
                child: Container(
                  width: 24,
                  height: 24,
                  decoration: BoxDecoration(
                    color: c,
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: selected
                          ? Theme.of(context).colorScheme.onSurface
                          : Colors.transparent,
                      width: selected ? 2.5 : 0,
                    ),
                  ),
                ),
              );
            }).toList(),
          ),
        ),
      ],
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
    Widget placeholder() => Container(width: 52, height: 52, color: t.border);
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
