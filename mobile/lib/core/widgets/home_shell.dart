import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../theme/app_tokens.dart';
import '../theme/glass.dart';
import '../../l10n/app_localizations.dart';

/// 下部にフローティングのガラス製タブバーを持つホームシェル。
/// `extendBody: true` でコンテンツがバーの背後に回り込み、ガラス越しに透ける。
class HomeShell extends StatelessWidget {
  const HomeShell({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  static const _tabIcons = [
    (icon: Icons.calendar_month_outlined, active: Icons.calendar_month),
    (icon: Icons.restaurant_outlined, active: Icons.restaurant),
    (icon: Icons.fitness_center_outlined, active: Icons.fitness_center),
    (icon: Icons.monitor_weight_outlined, active: Icons.monitor_weight),
    (icon: Icons.person_outline, active: Icons.person),
  ];

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final labels = [
      l.navHome,
      l.navMeals,
      l.navTraining,
      l.navWeight,
      l.navProfile,
    ];

    return Scaffold(
      extendBody: true,
      body: navigationShell,
      bottomNavigationBar: _GlassTabBar(
        currentIndex: navigationShell.currentIndex,
        icons: _tabIcons,
        labels: labels,
        onTap: (i) => navigationShell.goBranch(
          i,
          initialLocation: i == navigationShell.currentIndex,
        ),
      ),
    );
  }
}

typedef _IconSpec = ({IconData icon, IconData active});

class _GlassTabBar extends StatelessWidget {
  const _GlassTabBar({
    required this.currentIndex,
    required this.icons,
    required this.labels,
    required this.onTap,
  });

  final int currentIndex;
  final List<_IconSpec> icons;
  final List<String> labels;
  final ValueChanged<int> onTap;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(14, 0, 14, 12),
      child: SafeArea(
        top: false,
        child: GlassSurface(
          borderRadius: BorderRadius.circular(30),
          child: SizedBox(
            height: 62,
            child: Row(
              children: [
                for (var i = 0; i < icons.length; i++)
                  Expanded(
                    child: _TabItem(
                      spec: icons[i],
                      label: labels[i],
                      selected: i == currentIndex,
                      onTap: () => onTap(i),
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class _TabItem extends StatelessWidget {
  const _TabItem({
    required this.spec,
    required this.label,
    required this.selected,
    required this.onTap,
  });

  final _IconSpec spec;
  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final color = selected ? t.accent : t.subtle;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(16),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 46,
            height: 30,
            decoration: BoxDecoration(
              color: selected ? t.accentWeak : Colors.transparent,
              borderRadius: BorderRadius.circular(15),
            ),
            child: Icon(
              selected ? spec.active : spec.icon,
              size: 23,
              color: color,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            label,
            style: TextStyle(
              fontSize: 10,
              fontWeight: FontWeight.w600,
              color: color,
            ),
          ),
        ],
      ),
    );
  }
}
