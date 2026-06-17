import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../theme/app_tokens.dart';
import '../theme/glass.dart';

/// 下部にフローティングのガラス製タブバーを持つホームシェル。
/// `extendBody: true` でコンテンツがバーの背後に回り込み、ガラス越しに透ける。
class HomeShell extends StatelessWidget {
  const HomeShell({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  static const _tabs = [
    (icon: Icons.restaurant_outlined, active: Icons.restaurant, label: '食事'),
    (icon: Icons.fitness_center_outlined, active: Icons.fitness_center, label: 'トレーニング'),
    (icon: Icons.monitor_weight_outlined, active: Icons.monitor_weight, label: '体重'),
    (icon: Icons.person_outline, active: Icons.person, label: 'マイページ'),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      extendBody: true,
      body: navigationShell,
      bottomNavigationBar: _GlassTabBar(
        currentIndex: navigationShell.currentIndex,
        tabs: _tabs,
        onTap: (i) => navigationShell.goBranch(
          i,
          initialLocation: i == navigationShell.currentIndex,
        ),
      ),
    );
  }
}

typedef _TabSpec = ({IconData icon, IconData active, String label});

class _GlassTabBar extends StatelessWidget {
  const _GlassTabBar({
    required this.currentIndex,
    required this.tabs,
    required this.onTap,
  });

  final int currentIndex;
  final List<_TabSpec> tabs;
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
                for (var i = 0; i < tabs.length; i++)
                  Expanded(
                    child: _TabItem(
                      spec: tabs[i],
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
  const _TabItem({required this.spec, required this.selected, required this.onTap});

  final _TabSpec spec;
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
            child: Icon(selected ? spec.active : spec.icon, size: 23, color: color),
          ),
          const SizedBox(height: 2),
          Text(
            spec.label,
            style: TextStyle(fontSize: 10, fontWeight: FontWeight.w600, color: color),
          ),
        ],
      ),
    );
  }
}
