import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../data/exercise_dtos.dart';
import '../data/training_repository.dart';

/// 種目の管理（一覧 / 追加 / 削除）。
class ExercisesScreen extends ConsumerWidget {
  const ExercisesScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final exercises = ref.watch(exercisesProvider);
    return Scaffold(
      appBar: AppBar(
        title: const Text('種目'),
        backgroundColor: Colors.transparent,
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _createDialog(context, ref),
          ),
        ],
      ),
      body: SafeArea(
        child: AsyncValueView<List<ExerciseDto>>(
          value: exercises,
          onRetry: () => ref.invalidate(exercisesProvider),
          data: (list) {
            if (list.isEmpty) {
              return const Center(child: Text('種目がありません。右上の + で追加'));
            }
            return ListView(
              padding: const EdgeInsets.all(16),
              children: [
                AppListBox(
                  children: [
                    for (final ex in list)
                      AppListRow(
                        child: Row(
                          children: [
                            Expanded(
                              child: Text(
                                ex.name,
                                style: const TextStyle(
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ),
                            IconButton(
                              icon: Icon(
                                Icons.delete_outline,
                                size: 20,
                                color: context.tokens.subtle,
                              ),
                              onPressed: () => _delete(context, ref, ex),
                            ),
                          ],
                        ),
                      ),
                  ],
                ),
              ],
            );
          },
        ),
      ),
    );
  }

  Future<void> _createDialog(BuildContext context, WidgetRef ref) async {
    final controller = TextEditingController();
    final name = await showDialog<String>(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('種目を追加'),
        content: TextField(
          controller: controller,
          autofocus: true,
          decoration: const InputDecoration(hintText: 'ベンチプレス'),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(),
            child: const Text('キャンセル'),
          ),
          FilledButton(
            onPressed: () =>
                Navigator.of(dialogContext).pop(controller.text.trim()),
            child: const Text('追加'),
          ),
        ],
      ),
    );
    if (name == null || name.isEmpty) return;
    try {
      await ref.read(trainingRepositoryProvider).createExercise(name);
      ref.invalidate(exercisesProvider);
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
        ).showSnackBar(const SnackBar(content: Text('追加に失敗しました')));
      }
    }
  }

  Future<void> _delete(
    BuildContext context,
    WidgetRef ref,
    ExerciseDto ex,
  ) async {
    try {
      await ref.read(trainingRepositoryProvider).deleteExercise(ex.id);
      ref.invalidate(exercisesProvider);
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
