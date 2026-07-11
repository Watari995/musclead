import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_card.dart';
import '../../../l10n/app_localizations.dart';
import '../data/exercise_dtos.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';
import 'training_record_screen.dart';

/// トレーニング詳細（種目 > セットの読み取り表示）。
class TrainingDetailScreen extends ConsumerWidget {
  const TrainingDetailScreen({super.key, required this.training});

  final TrainingDto training;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final exList =
        ref.watch(exercisesProvider).asData?.value ?? const <ExerciseDto>[];
    final names = {for (final e in exList) e.id: e.name};
    final minutes = training.endedAt?.difference(training.startedAt).inMinutes;

    return Scaffold(
      appBar: AppBar(
        title: Text(l.trainingDetailTitle(mdWeekday(training.startedAt))),
        backgroundColor: Colors.transparent,
        actions: [
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () => Navigator.of(context).push(
              MaterialPageRoute<void>(
                builder: (_) => TrainingRecordScreen(editingTraining: training),
              ),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () => _confirmDelete(context, ref),
          ),
        ],
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
          children: [
            AppCard(
              child: Row(
                children: [
                  Expanded(
                    child: _stat(
                      context,
                      minutes == null ? '-' : '$minutes',
                      l.trainingMinutes,
                      l.trainingHoursLabel,
                    ),
                  ),
                  Container(width: 1, height: 36, color: t.border),
                  Expanded(
                    child: _stat(
                      context,
                      '${training.exercises.length}',
                      '',
                      l.trainingExercisesLabel,
                    ),
                  ),
                ],
              ),
            ),
            if (training.memo?.isNotEmpty == true)
              Padding(
                padding: const EdgeInsets.fromLTRB(4, 12, 4, 0),
                child: Text(
                  l.trainingMemoPrefix(training.memo!),
                  style: TextStyle(fontSize: 13, color: t.muted),
                ),
              ),
            const SizedBox(height: 12),
            for (final ex in training.exercises)
              Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: _ExerciseCard(
                  name: names[ex.exerciseId] ?? l.trainingExerciseDefault,
                  exercise: ex,
                ),
              ),
          ],
        ),
      ),
    );
  }

  Future<void> _confirmDelete(BuildContext context, WidgetRef ref) async {
    final l = AppLocalizations.of(context)!;
    final ok = await showDialog<bool>(
      context: context,
      builder: (dialogContext) {
        final dl = AppLocalizations.of(dialogContext)!;
        return AlertDialog(
          title: Text(dl.trainingDeleteConfirmTitle),
          content: Text(dl.trainingDeleteConfirmMsg),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(false),
              child: Text(dl.commonCancel),
            ),
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(true),
              child: Text(
                dl.commonDeleteOk,
                style: TextStyle(color: context.tokens.accent),
              ),
            ),
          ],
        );
      },
    );
    if (ok != true) return;
    try {
      await ref.read(trainingRepositoryProvider).deleteTraining(training.id);
      ref.invalidate(trainingsProvider);
      if (context.mounted) Navigator.of(context).pop();
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(l.commonDeleteFailed)));
      }
    }
  }

  Widget _stat(BuildContext context, String value, String unit, String label) {
    final t = context.tokens;
    return Column(
      children: [
        Text.rich(
          TextSpan(
            text: value,
            style: const TextStyle(fontSize: 19, fontWeight: FontWeight.w800),
            children: [
              TextSpan(text: unit, style: const TextStyle(fontSize: 12)),
            ],
          ),
        ),
        Text(label, style: TextStyle(fontSize: 11, color: t.muted)),
      ],
    );
  }
}

class _ExerciseCard extends StatelessWidget {
  const _ExerciseCard({required this.name, required this.exercise});

  final String name;
  final TrainingExerciseDto exercise;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            name,
            style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15),
          ),
          const SizedBox(height: 6),
          for (final s in exercise.sets)
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 3),
              child: Row(
                children: [
                  SizedBox(
                    width: 24,
                    child: Text(
                      '${s.setNumber}',
                      style: TextStyle(color: t.muted),
                    ),
                  ),
                  Text(
                    '${s.weightKg} kg',
                    style: const TextStyle(fontWeight: FontWeight.w700),
                  ),
                  const Spacer(),
                  Text(
                    l.trainingRepsUnit(s.reps),
                    style: TextStyle(color: t.muted),
                  ),
                ],
              ),
            ),
        ],
      ),
    );
  }
}
