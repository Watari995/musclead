import 'package:flutter/material.dart';

import '../theme/app_tokens.dart';
import '../theme/sketchy.dart';

/// 手描き風の二重ストローク輪郭を持つカード。
class AppCard extends StatelessWidget {
  const AppCard({super.key, required this.child, this.padding});

  final Widget child;
  final EdgeInsetsGeometry? padding;

  @override
  Widget build(BuildContext context) {
    return RoughBox(padding: padding ?? const EdgeInsets.all(15), child: child);
  }
}

/// 罫線で区切られた行リスト(iOS グループ表示)。
class AppListBox extends StatelessWidget {
  const AppListBox({super.key, required this.children});

  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final rows = <Widget>[];
    for (var i = 0; i < children.length; i++) {
      if (i > 0) rows.add(Divider(height: 1, thickness: 1.4, color: t.hairline));
      rows.add(children[i]);
    }
    return RoughBox(clipBehavior: Clip.antiAlias, child: Column(children: rows));
  }
}

/// AppListBox 内のタップ可能な行。
class AppListRow extends StatelessWidget {
  const AppListRow({
    super.key,
    required this.child,
    this.onTap,
    this.padding = const EdgeInsets.symmetric(horizontal: 15, vertical: 13),
  });

  final Widget child;
  final VoidCallback? onTap;
  final EdgeInsetsGeometry padding;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        child: Padding(padding: padding, child: child),
      ),
    );
  }
}
