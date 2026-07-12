import 'dart:math';

import 'package:flutter/material.dart';

import 'app_tokens.dart';

/// 手描き(Excalidraw 風)の二重ストローク角丸矩形を描く [CustomPainter]。
///
/// 同じ角丸矩形の輪郭を、わずかにランダムな位置ゆらぎを与えて 2 回描画することで
/// ラフスケッチのような線の揺れを表現する。web モックアップの
/// `.rough` / `.rough::before` (二重の角丸ボーダー + ランダム回転) と同じ発想を
/// Canvas 描画に置き換えたもの。
class SketchyBorderPainter extends CustomPainter {
  const SketchyBorderPainter({
    required this.color,
    required this.radius,
    this.strokeWidth = 1.75,
    this.seed = 7,
  });

  final Color color;
  final BorderRadius radius;
  final double strokeWidth;
  final int seed;

  @override
  void paint(Canvas canvas, Size size) {
    final rnd = Random(seed);
    // 内側のくっきりした 1 本目。
    _paintStroke(canvas, size, rnd, opacityFactor: 1, jitter: 0.6, inset: 0);
    // 外側にわずかにずれた 2 本目(薄め)。手描きの重ね書き感を出す。
    _paintStroke(
      canvas,
      size,
      rnd,
      opacityFactor: 0.32,
      jitter: 1.5,
      inset: -2.2,
    );
  }

  void _paintStroke(
    Canvas canvas,
    Size size,
    Random rnd, {
    required double opacityFactor,
    required double jitter,
    required double inset,
  }) {
    final rect = Rect.fromLTWH(
      inset,
      inset,
      size.width - inset * 2,
      size.height - inset * 2,
    );
    if (rect.width <= 0 || rect.height <= 0) return;
    final rrect = radius.toRRect(rect);
    final path = _jitteredPath(rrect, rnd, jitter);
    final paint = Paint()
      ..color = color.withValues(alpha: color.a * opacityFactor)
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round
      ..strokeJoin = StrokeJoin.round;
    canvas.drawPath(path, paint);
  }

  /// 角丸矩形の輪郭を一定間隔でサンプリングし、各点を少しだけランダムにずらして
  /// 滑らかな曲線で結び直す。これが手描き線のゆらぎになる。
  Path _jitteredPath(RRect rrect, Random rnd, double jitter) {
    const steps = 48;
    final metrics = (Path()..addRRect(rrect)).computeMetrics().first;
    final points = <Offset>[];
    for (var i = 0; i < steps; i++) {
      final t = metrics.length * i / steps;
      final tangent = metrics.getTangentForOffset(t);
      if (tangent == null) continue;
      final dx = (rnd.nextDouble() - 0.5) * jitter * 2;
      final dy = (rnd.nextDouble() - 0.5) * jitter * 2;
      points.add(tangent.position + Offset(dx, dy));
    }
    if (points.isEmpty) return Path()..addRRect(rrect);

    final path = Path()..moveTo(points.first.dx, points.first.dy);
    for (var i = 1; i < points.length; i++) {
      final prev = points[i - 1];
      final curr = points[i];
      final mid = Offset((prev.dx + curr.dx) / 2, (prev.dy + curr.dy) / 2);
      path.quadraticBezierTo(prev.dx, prev.dy, mid.dx, mid.dy);
    }
    path.close();
    return path;
  }

  @override
  bool shouldRepaint(covariant SketchyBorderPainter oldDelegate) {
    return color != oldDelegate.color ||
        radius != oldDelegate.radius ||
        strokeWidth != oldDelegate.strokeWidth ||
        seed != oldDelegate.seed;
  }
}

/// 手描き風の輪郭 + 塗りで子要素を囲むラッパー。
///
/// `AppCard` / `AppListBox` / `AppButton` / `AppTextField` など、枠線を持つ
/// 共通ウィジェットはすべてこれを土台にする。既定値は [BuildContext.tokens] から
/// 取るため、通常は `color` / `fill` を指定しなくてよい。
class RoughBox extends StatelessWidget {
  const RoughBox({
    super.key,
    required this.child,
    this.color,
    this.fill,
    this.radius = const BorderRadius.all(Radius.circular(14)),
    this.strokeWidth = 1.75,
    this.padding,
    this.clipBehavior = Clip.none,
    this.seed = 7,
  });

  /// 円/ピル型(タブのアクティブ・スクリブルなど)に使う既定の丸め。
  static const BorderRadius pill = BorderRadius.all(Radius.circular(999));

  final Widget child;
  final Color? color;
  final Color? fill;
  final BorderRadius radius;
  final double strokeWidth;
  final EdgeInsetsGeometry? padding;
  final Clip clipBehavior;
  final int seed;

  @override
  Widget build(BuildContext context) {
    final tokens = context.tokens;
    final strokeColor = color ?? tokens.border;
    final fillColor = fill ?? tokens.paper;

    return CustomPaint(
      foregroundPainter: SketchyBorderPainter(
        color: strokeColor,
        radius: radius,
        strokeWidth: strokeWidth,
        seed: seed,
      ),
      child: Container(
        padding: padding,
        clipBehavior: clipBehavior,
        decoration: BoxDecoration(color: fillColor, borderRadius: radius),
        child: child,
      ),
    );
  }
}
