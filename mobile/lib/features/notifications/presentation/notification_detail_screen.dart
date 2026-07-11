import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../l10n/app_localizations.dart';
import '../data/notification_dtos.dart';
import '../data/notification_repository.dart';

class NotificationDetailScreen extends ConsumerStatefulWidget {
  const NotificationDetailScreen({super.key, required this.id});

  final String id;

  @override
  ConsumerState<NotificationDetailScreen> createState() =>
      _NotificationDetailScreenState();
}

class _NotificationDetailScreenState
    extends ConsumerState<NotificationDetailScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref
          .read(notificationRepositoryProvider)
          .markAsRead(widget.id)
          .then((_) => ref.invalidate(notificationsProvider));
    });
  }

  @override
  Widget build(BuildContext context) {
    final notifications = ref.watch(notificationsProvider);

    final notification = notifications.asData?.value.notifications
        .where((n) => n.id == widget.id)
        .firstOrNull;

    return Scaffold(
      appBar: AppBar(
        title: Text(AppLocalizations.of(context)!.notificationDetailTitle),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: notification == null
          ? const Center(child: CircularProgressIndicator(strokeWidth: 2.4))
          : _NotificationBody(notification: notification),
    );
  }
}

class _NotificationBody extends StatelessWidget {
  const _NotificationBody({required this.notification});

  final NotificationDto notification;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final m = notification.metadata;
    final l = AppLocalizations.of(context)!;

    return ListView(
      padding: const EdgeInsets.all(18),
      children: [
        if (notification.notificationType == 'weekly_goal') ...[
          _WeeklyGoalDetail(metadata: m),
        ] else ...[
          Text(l.notificationTitle, style: TextStyle(fontSize: 16, color: t.subtle)),
        ],
        const SizedBox(height: 16),
        Text(
          notification.createdAt.toLocal().toString().substring(0, 16),
          style: TextStyle(fontSize: 12, color: t.muted),
        ),
      ],
    );
  }
}

class _WeeklyGoalDetail extends StatelessWidget {
  const _WeeklyGoalDetail({required this.metadata});

  final Map<String, dynamic> metadata;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final achieved = metadata['achieved'] as bool? ?? false;
    final trainingGoal = metadata['training_goal'] as num?;
    final trainingActual = metadata['training_actual'] as num?;
    final calorieGoal = metadata['calorie_goal'] as num?;
    final calorieActual = metadata['calorie_actual'] as num?;
    final weightGoal = metadata['weight_goal'] as num?;
    final weightActual = metadata['weight_actual'] as num?;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Builder(builder: (context) {
          final l = AppLocalizations.of(context)!;
          return Text(
            achieved ? l.notificationGoalAchievedDetail : l.notificationGoalResult,
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
              color: t.muted,
            ),
          );
        }),
        const SizedBox(height: 16),
        if (trainingGoal != null)
          Builder(builder: (context) {
            final l = AppLocalizations.of(context)!;
            return _GoalRow(
              label: l.notificationTraining,
              value: AppLocalizations.of(context)!.notificationTrainingActualGoal(
                (trainingActual ?? 0).toInt(), trainingGoal.toInt()),
              achieved: (trainingActual?.toInt() ?? 0) >= trainingGoal.toInt(),
            );
          }),
        if (calorieGoal != null)
          Builder(builder: (context) {
            final l = AppLocalizations.of(context)!;
            return _GoalRow(
              label: l.notificationCalorie,
              value:
                  '${calorieActual != null ? calorieActual.round() : '—'} / $calorieGoal kcal',
              achieved: calorieActual != null && calorieActual <= calorieGoal,
            );
          }),
        if (weightGoal != null)
          Builder(builder: (context) {
            final l = AppLocalizations.of(context)!;
            return _GoalRow(
              label: l.notificationWeightChange,
              value:
                  '${weightActual != null ? (weightActual > 0 ? '+' : '') + weightActual.toStringAsFixed(1) : '—'} kg',
              achieved:
                weightActual != null &&
                weightGoal.isNegative == weightActual.isNegative &&
                weightActual.abs() >= weightGoal.abs(),
            );
          }),
      ],
    );
  }
}

class _GoalRow extends StatelessWidget {
  const _GoalRow({
    required this.label,
    required this.value,
    required this.achieved,
  });

  final String label;
  final String value;
  final bool achieved;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: TextStyle(fontSize: 12, color: t.muted)),
                Text(value, style: const TextStyle(fontSize: 14)),
              ],
            ),
          ),
          Text(achieved ? '✅' : '❌', style: const TextStyle(fontSize: 18)),
        ],
      ),
    );
  }
}
