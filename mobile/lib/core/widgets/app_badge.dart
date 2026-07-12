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
    final (bg, border, fg) = switch (tone) {
      BadgeTone.neutral => (Colors.transparent, t.hairline, t.muted),
      BadgeTone.accent => (t.accentWeak, Colors.transparent, t.accent),
      BadgeTone.gold => (t.goldWeak, Colors.transparent, t.gold),
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 9, vertical: 3),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(999),
        border: Border.all(color: border, width: 1.2),
      ),
      child: Text(
        label,
        style: TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: fg),
      ),
    );
  }
}
