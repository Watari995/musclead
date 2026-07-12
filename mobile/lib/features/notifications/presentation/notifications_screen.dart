import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:intl/intl.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/tab_page.dart';
import '../data/notification_dtos.dart';
import '../data/notification_repository.dart';
import '../../../l10n/app_localizations.dart';

class NotificationsScreen extends ConsumerWidget {
  const NotificationsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final notifications = ref.watch(notificationsProvider);
    return TabPage(
      title: l.notificationTitle,
      onRefresh: () => ref.refresh(notificationsProvider.future),
      children: [
        AsyncValueView<GetNotificationsResponse>(
          value: notifications,
          onRetry: () => ref.invalidate(notificationsProvider),
          data: (data) {
            if (data.notifications.isEmpty) {
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 48),
                child: Center(
                  child: Text(
                    l.notificationEmpty,
                    style: TextStyle(color: context.tokens.muted),
                  ),
                ),
              );
            }
            return Column(
              children: data.notifications
                  .map((n) => _NotificationTile(notification: n))
                  .toList(),
            );
          },
        ),
      ],
    );
  }
}

class _NotificationTile extends ConsumerWidget {
  const _NotificationTile({required this.notification});

  final NotificationDto notification;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final achieved = notification.metadata['achieved'] as bool? ?? false;
    final label = notification.notificationType == 'weekly_goal'
        ? (achieved
              ? l.notificationWeeklyGoalAchieved
              : l.notificationWeeklyGoalCheck)
        : l.notificationTitle;
    final dateStr = DateFormat(
      'M/d HH:mm',
    ).format(notification.createdAt.toLocal());

    return GestureDetector(
      onTap: () => context.push('/notifications/${notification.id}'),
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        child: RoughBox(
          radius: BorderRadius.circular(12),
          fill: notification.isRead ? t.paper : t.paperAlt,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  notification.notificationType == 'weekly_goal' ? '🏆' : '🔔',
                  style: const TextStyle(fontSize: 20),
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        label,
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: notification.isRead
                              ? FontWeight.normal
                              : FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        dateStr,
                        style: TextStyle(fontSize: 12, color: t.muted),
                      ),
                    ],
                  ),
                ),
                if (!notification.isRead)
                  Container(
                    margin: const EdgeInsets.only(top: 4),
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: Colors.red,
                      shape: BoxShape.circle,
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
