import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_badge.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../data/subscription_dtos.dart';
import '../data/subscription_repository.dart';

/// プラン（読み取り専用）。購入導線は持たない（App Store 3.1.1）。
class PlanScreen extends ConsumerWidget {
  const PlanScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final sub = ref.watch(subscriptionProvider);
    return Scaffold(
      appBar: AppBar(
        title: const Text('プラン'),
        backgroundColor: Colors.transparent,
      ),
      body: SafeArea(
        child: AsyncValueView<GetSubscriptionResponse>(
          value: sub,
          onRetry: () => ref.invalidate(subscriptionProvider),
          data: (s) => ListView(
            padding: const EdgeInsets.all(16),
            children: [
              AppCard(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(
                          s.isPro ? 'Pro' : 'Free',
                          style: const TextStyle(
                            fontSize: 24,
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                        const SizedBox(width: 10),
                        AppBadge(
                          s.isPro ? '有効' : '無料プラン',
                          tone: s.isPro ? BadgeTone.accent : BadgeTone.neutral,
                        ),
                      ],
                    ),
                    if (s.expiresAt != null)
                      Padding(
                        padding: const EdgeInsets.only(top: 8),
                        child: Text(
                          '有効期限: ${dateJpLong(s.expiresAt!)}',
                          style: TextStyle(color: context.tokens.muted),
                        ),
                      ),
                  ],
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'プランの購入・変更はウェブからお手続きください。',
                style: TextStyle(fontSize: 13, color: context.tokens.muted),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
