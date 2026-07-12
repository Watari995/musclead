import 'package:flutter/material.dart';

import '../theme/app_tokens.dart';
import '../theme/sketchy.dart';

enum AppButtonVariant { primary, glass, text }

/// プレビューの主CTA(accent 塗り) / セカンダリ(手描き輪郭) / テキストボタンに対応。
class AppButton extends StatelessWidget {
  const AppButton({
    super.key,
    required this.label,
    this.onPressed,
    this.variant = AppButtonVariant.primary,
    this.icon,
    this.loading = false,
    this.expand = true,
  });

  final String label;
  final VoidCallback? onPressed;
  final AppButtonVariant variant;
  final IconData? icon;
  final bool loading;
  final bool expand;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final enabled = onPressed != null && !loading;

    final isDark = Theme.of(context).brightness == Brightness.dark;
    final fg = switch (variant) {
      AppButtonVariant.primary => isDark ? t.paper : Colors.white,
      AppButtonVariant.glass => t.ink,
      AppButtonVariant.text => t.accent,
    };

    final content = Row(
      mainAxisSize: expand ? MainAxisSize.max : MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        if (loading)
          SizedBox(
            width: 18,
            height: 18,
            child: CircularProgressIndicator(strokeWidth: 2, color: fg),
          )
        else ...[
          if (icon != null) ...[
            Icon(icon, size: 19, color: fg),
            const SizedBox(width: 8),
          ],
          Text(
            label,
            style: TextStyle(
              color: fg,
              fontSize: 15,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ],
    );

    final inner = Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: enabled ? onPressed : null,
        borderRadius: BorderRadius.circular(14),
        child: SizedBox(
          height: 48,
          child: Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: content,
            ),
          ),
        ),
      ),
    );

    final opacity = enabled ? 1.0 : 0.5;
    final Widget button = switch (variant) {
      AppButtonVariant.primary => Opacity(
        opacity: opacity,
        child: RoughBox(
          radius: BorderRadius.circular(14),
          fill: t.accent,
          color: t.accent,
          clipBehavior: Clip.antiAlias,
          child: inner,
        ),
      ),
      AppButtonVariant.glass => Opacity(
        opacity: opacity,
        child: RoughBox(
          radius: BorderRadius.circular(14),
          clipBehavior: Clip.antiAlias,
          child: inner,
        ),
      ),
      AppButtonVariant.text => Opacity(opacity: opacity, child: inner),
    };

    return expand ? SizedBox(width: double.infinity, child: button) : button;
  }
}
