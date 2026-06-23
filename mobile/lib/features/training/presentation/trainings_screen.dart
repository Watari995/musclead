import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';
import 'training_detail_screen.dart';

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
        Row(
          children: [
            Expanded(
              child: AppButton(
                label: '種目',
                icon: Icons.fitness_center,
                variant: AppButtonVariant.glass,
                onPressed: () => context.push('/exercises'),
              ),
            ),
            const SizedBox(width: 10),
            Expanded(
              child: AppButton(
                label: '記録',
                icon: Icons.show_chart,
                variant: AppButtonVariant.glass,
                onPressed: () => context.push('/records'),
              ),
            ),
            const SizedBox(width: 10),
            Expanded(
              child: AppButton(
                label: 'ルーティン',
                icon: Icons.list_alt,
                variant: AppButtonVariant.glass,
                onPressed: () => context.push('/routines'),
              ),
            ),
          ],
        ),
        const SizedBox(height: 10),
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

  String? _duration() {
    final end = training.endedAt;
    if (end == null) return null;
    final mins = end.difference(training.startedAt).inMinutes;
    return '$mins分';
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final count = training.exercises.length;
    return AppListRow(
      onTap: () => Navigator.of(context, rootNavigator: true).push(
        MaterialPageRoute<void>(
          builder: (_) => TrainingDetailScreen(training: training),
        ),
      ),
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
          if (_duration() != null)
            Text(_duration()!, style: TextStyle(fontSize: 12, color: t.muted)),
          const SizedBox(width: 6),
          Icon(Icons.chevron_right, color: t.subtle),
        ],
      ),
    );
  }
}
