import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';

class TrainingsScreen extends ConsumerWidget {
  const TrainingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final trainings = ref.watch(trainingsProvider);
    return TabPage(
      title: 'トレーニング',
      trailing: IconButton(
        icon: Icon(Icons.add, color: context.tokens.accent),
        onPressed: () => context.push('/trainings/new'),
      ),
      onRefresh: () => ref.refresh(trainingsProvider.future),
      children: [
        AsyncValueView<List<TrainingDto>>(
          value: trainings,
          onRetry: () => ref.invalidate(trainingsProvider),
          data: (list) {
            if (list.isEmpty) {
              return const Padding(
                padding: EdgeInsets.symmetric(vertical: 40),
                child: Center(child: Text('まだ記録がありません')),
              );
            }
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                const SectionTitle('履歴'),
                AppListBox(
                  children: [for (final tr in list) _TrainingRow(training: tr)],
                ),
              ],
            );
          },
        ),
      ],
    );
  }
}

class _TrainingRow extends StatelessWidget {
  const _TrainingRow({required this.training});

  final TrainingDto training;

  String _duration() {
    final end = training.endedAt;
    if (end == null) return '進行中';
    final mins = end.difference(training.startedAt).inMinutes;
    return '$mins分';
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final count = training.exercises.length;
    return AppListRow(
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  mdWeekday(training.startedAt),
                  style: const TextStyle(
                    fontWeight: FontWeight.w700,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '$count 種目',
                  style: TextStyle(fontSize: 12, color: t.muted),
                ),
              ],
            ),
          ),
          Text(_duration(), style: TextStyle(fontSize: 12, color: t.muted)),
          const SizedBox(width: 6),
          Icon(Icons.chevron_right, color: t.subtle),
        ],
      ),
    );
  }
}
