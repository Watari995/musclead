import 'package:flutter/material.dart';

/// タブ配下の画面共通スキャフォールド。iOS ラージタイトル + 引っ張って更新 +
/// フローティングタブバー分の下部余白。
class TabPage extends StatelessWidget {
  const TabPage({
    super.key,
    required this.title,
    required this.children,
    this.subtitle,
    this.trailing,
    this.onRefresh,
  });

  final String title;
  final String? subtitle;
  final Widget? trailing;
  final Future<void> Function()? onRefresh;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final list = ListView(
      padding: const EdgeInsets.fromLTRB(18, 8, 18, 124),
      children: [
        Padding(
          padding: EdgeInsets.only(top: 8, bottom: subtitle == null ? 8 : 2),
          child: Row(
            children: [
              Expanded(
                child: Text(
                  title,
                  style: const TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.w800,
                    letterSpacing: -0.4,
                  ),
                ),
              ),
              ?trailing,
            ],
          ),
        ),
        if (subtitle != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Text(
              subtitle!,
              style: TextStyle(
                fontSize: 13,
                fontWeight: FontWeight.w600,
                color: Theme.of(
                  context,
                ).colorScheme.onSurface.withValues(alpha: 0.55),
              ),
            ),
          ),
        const SizedBox(height: 6),
        ...children,
      ],
    );

    return SafeArea(
      bottom: false,
      child: onRefresh == null
          ? list
          : RefreshIndicator(onRefresh: onRefresh!, child: list),
    );
  }
}
