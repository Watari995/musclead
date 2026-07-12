import 'package:flutter/material.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';

class SplashScreen extends StatelessWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RoughBox(
              fill: t.accent,
              color: t.accent,
              radius: BorderRadius.circular(18),
              child: SizedBox(
                width: 64,
                height: 64,
                child: Icon(
                  Icons.fitness_center,
                  color: context.colors.onPrimary,
                  size: 34,
                ),
              ),
            ),
            const SizedBox(height: 20),
            const SizedBox(
              width: 22,
              height: 22,
              child: CircularProgressIndicator(strokeWidth: 2.2),
            ),
          ],
        ),
      ),
    );
  }
}
