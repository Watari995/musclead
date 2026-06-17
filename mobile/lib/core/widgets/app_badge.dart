import 'package:flutter/material.dart';

import '../theme/app_tokens.dart';

enum BadgeTone { neutral, accent, gold }

class AppBadge extends StatelessWidget {
  const AppBadge(this.label, {super.key, this.tone = BadgeTone.neutral});

  final String label;
  final BadgeTone tone;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final (bg, fg) = switch (tone) {
      BadgeTone.neutral => (t.accentWeak.withValues(alpha: 0), t.muted),
      BadgeTone.accent => (t.accentWeak, t.accent),
      BadgeTone.gold => (t.goldWeak, t.gold),
    };
    final background = tone == BadgeTone.neutral
        ? context.colors.onSurface.withValues(alpha: 0.06)
        : bg;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 3),
      decoration: BoxDecoration(
        color: background,
        borderRadius: BorderRadius.circular(999),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: fg),
      ),
    );
  }
}
