import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../data/routine_dtos.dart';
import '../data/training_repository.dart';
import 'training_record_screen.dart';

/// ルーティンの一覧 / 削除 / タップで記録開始。
class RoutinesScreen extends ConsumerWidget {
  const RoutinesScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final routines = ref.watch(routinesProvider);
    return Scaffold(
      appBar: AppBar(
        title: const Text('ルーティン'),
        backgroundColor: Colors.transparent,
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/routines/new'),
          ),
        ],
      ),
      body: SafeArea(
        child: AsyncValueView<List<RoutineDto>>(
          value: routines,
          onRetry: () => ref.invalidate(routinesProvider),
          data: (list) {
            if (list.isEmpty) {
              return const Center(child: Text('ルーティンがありません。右上の + で作成'));
            }
            return ListView(
              padding: const EdgeInsets.all(16),
              children: [
                AppListBox(
                  children: [for (final r in list) _row(context, ref, r)],
                ),
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _row(BuildContext context, WidgetRef ref, RoutineDto r) {
    return AppListRow(
      onTap: () => Navigator.of(context, rootNavigator: true).push(
        MaterialPageRoute<void>(
          builder: (_) => TrainingRecordScreen(
            initialExercises: [
              for (final e in r.routineExercises)
                (exerciseId: e.exerciseId, name: e.exerciseName ?? '種目'),
            ],
          ),
        ),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  r.name,
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                Text(
                  '${r.routineExercises.length} 種目 ・ タップで記録開始',
                  style: TextStyle(fontSize: 12, color: context.tokens.muted),
                ),
              ],
            ),
          ),
          IconButton(
            icon: Icon(
              Icons.delete_outline,
              size: 20,
              color: context.tokens.subtle,
            ),
            onPressed: () => _delete(context, ref, r),
          ),
        ],
      ),
    );
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    RoutineDto r,
  ) async {
    try {
      await ref.read(trainingRepositoryProvider).deleteRoutine(r.id);
      ref.invalidate(routinesProvider);
    } on Failure catch (f) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(f.message)));
      }
    } catch (_) {
      if (context.mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('削除に失敗しました')));
      }
    }
  }
}
