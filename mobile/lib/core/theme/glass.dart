import 'dart:ui';

import 'package:flutter/material.dart';

import 'app_tokens.dart';

/// Liquid Glass の操作レイヤー（タブバー / ナビバー / FAB）に使う半透明サーフェス。
///
/// 影は clip の外側、ブラー + 半透明 tint は内側に置く。
class GlassSurface extends StatelessWidget {
  const GlassSurface({
    super.key,
    required this.child,
    this.borderRadius,
    this.blur = 24,
    this.padding,
    this.showBorder = true,
    this.shadow = true,
  });

  final Widget child;
  final BorderRadius? borderRadius;
  final double blur;
  final EdgeInsetsGeometry? padding;
  final bool showBorder;
  final bool shadow;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final radius = borderRadius ?? BorderRadius.circular(20);

    return DecoratedBox(
      decoration: BoxDecoration(
        borderRadius: radius,
        boxShadow: shadow
            ? [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.16),
                  blurRadius: 24,
                  offset: const Offset(0, 8),
                ),
              ]
            : null,
      ),
      child: ClipRRect(
        borderRadius: radius,
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: t.glassTint,
              borderRadius: radius,
              border: showBorder ? Border.all(color: t.glassBorder) : null,
            ),
            child: Padding(padding: padding ?? EdgeInsets.zero, child: child),
          ),
        ),
      ),
    );
  }
}
