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
import '../../../l10n/app_localizations.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';
import 'training_detail_screen.dart';

class TrainingsScreen extends ConsumerWidget {
  const TrainingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final trainings = ref.watch(trainingsProvider);
    return TabPage(
      title: l.trainingTitle,
      trailing: IconButton(
        icon: Icon(Icons.add, color: context.tokens.accent),
        onPressed: () => context.push('/trainings/new'),
      ),
      onRefresh: () => ref.refresh(trainingsProvider.future),
      children: [
        Wrap(
          spacing: 10,
          runSpacing: 10,
          children: [
            AppButton(
              label: l.trainingExercisesLabel,
              icon: Icons.fitness_center,
              variant: AppButtonVariant.glass,
              onPressed: () => context.push('/exercises'),
            ),
            AppButton(
              label: l.trainingRecordsLabel,
              icon: Icons.show_chart,
              variant: AppButtonVariant.glass,
              onPressed: () => context.push('/records'),
            ),
            AppButton(
              label: l.trainingRoutinesLabel,
              icon: Icons.list_alt,
              variant: AppButtonVariant.glass,
              onPressed: () => context.push('/routines'),
            ),
          ],
        ),
        const SizedBox(height: 10),
        AsyncValueView<List<TrainingDto>>(
          value: trainings,
          onRetry: () => ref.invalidate(trainingsProvider),
          data: (list) {
            if (list.isEmpty) {
              return Padding(
                padding: const EdgeInsets.symmetric(vertical: 40),
                child: Center(child: Text(l.trainingNoRecords)),
              );
            }
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                SectionTitle(l.trainingHistory),
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

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final count = training.exercises.length;
    final end = training.endedAt;
    final duration = end != null
        ? l.trainingDuration(end.difference(training.startedAt).inMinutes)
        : null;
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
                  l.trainingExerciseCount(count),
                  style: TextStyle(fontSize: 12, color: t.muted),
                ),
              ],
            ),
          ),
          if (duration != null)
            Text(duration, style: TextStyle(fontSize: 12, color: t.muted)),
          const SizedBox(width: 6),
          Icon(Icons.chevron_right, color: t.subtle),
        ],
      ),
    );
  }
}
