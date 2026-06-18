import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/async_value_view.dart';
import '../data/exercise_dtos.dart';
import '../data/training_repository.dart';

/// 種目の管理（一覧 / 追加 / 削除 / 並び替え）。
class ExercisesScreen extends ConsumerStatefulWidget {
  const ExercisesScreen({super.key});

  @override
  ConsumerState<ExercisesScreen> createState() => _ExercisesScreenState();
}

class _ExercisesScreenState extends ConsumerState<ExercisesScreen> {
  List<ExerciseDto>? _items;

  bool _sameIds(List<ExerciseDto> a, List<ExerciseDto> b) {
    if (a.length != b.length) return false;
    final ids = b.map((e) => e.id).toSet();
    return a.every((e) => ids.contains(e.id));
  }

  @override
  Widget build(BuildContext context) {
    final async = ref.watch(exercisesProvider);
    return Scaffold(
      appBar: AppBar(
        title: const Text('種目'),
        backgroundColor: Colors.transparent,
        actions: [
          IconButton(icon: const Icon(Icons.add), onPressed: _createDialog),
        ],
      ),
      body: SafeArea(
        child: AsyncValueView<List<ExerciseDto>>(
          value: async,
          onRetry: () => ref.invalidate(exercisesProvider),
          data: (list) {
            // 種目の集合が変わったら（追加/削除）ローカルを再同期。
            if (_items == null || !_sameIds(_items!, list)) {
              _items = List.of(list);
            }
            final items = _items!;
            if (items.isEmpty) {
              return const Center(child: Text('種目がありません。右上の + で追加'));
            }
            return Column(
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(18, 12, 18, 2),
                  child: Row(
                    children: [
                      Icon(
                        Icons.swap_vert,
                        size: 16,
                        color: context.tokens.subtle,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '右の ☰ をドラッグで並び替え',
                        style: TextStyle(
                          fontSize: 12,
                          color: context.tokens.muted,
                        ),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: ReorderableListView.builder(
                    padding: const EdgeInsets.fromLTRB(16, 4, 16, 16),
                    itemCount: items.length,
                    onReorder: _onReorder,
                    buildDefaultDragHandles: false,
                    proxyDecorator: (child, index, animation) =>
                        Material(color: Colors.transparent, child: child),
                    itemBuilder: (context, i) => _tile(items[i], i),
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }

  void _onReorder(int oldIndex, int newIndex) {
    setState(() {
      final items = _items!;
      var target = newIndex;
      if (target > oldIndex) target -= 1;
      final moved = items.removeAt(oldIndex);
      items.insert(target, moved);
    });
    ref
        .read(trainingRepositoryProvider)
        .reorderExercises(_items!.map((e) => e.id).toList())
        .catchError((Object _) {
          if (mounted) {
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(const SnackBar(content: Text('並び替えの保存に失敗しました')));
          }
        });
  }

  Widget _tile(ExerciseDto ex, int index) {
    final t = context.tokens;
    return Padding(
      key: ValueKey(ex.id),
      padding: const EdgeInsets.only(bottom: 8),
      child: Container(
        decoration: BoxDecoration(
          color: context.colors.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: t.border),
        ),
        child: ListTile(
          // 波紋(ink)を角丸にクリップし、長押し時に四角い角が出るのを防ぐ。
          shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(Radius.circular(12)),
          ),
          title: Text(
            ex.name,
            style: const TextStyle(fontWeight: FontWeight.w500),
          ),
          trailing: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              IconButton(
                icon: Icon(Icons.delete_outline, size: 20, color: t.subtle),
                onPressed: () => _delete(ex),
              ),
              ReorderableDragStartListener(
                index: index,
                child: Padding(
                  padding: const EdgeInsets.all(8),
                  child: Icon(Icons.drag_handle, color: t.subtle),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _createDialog() async {
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
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(f.message)));
      }
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('追加に失敗しました')));
      }
    }
  }

  Future<void> _delete(ExerciseDto ex) async {
    try {
      await ref.read(trainingRepositoryProvider).deleteExercise(ex.id);
      ref.invalidate(exercisesProvider);
    } on Failure catch (f) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(f.message)));
      }
    } catch (_) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text('削除に失敗しました')));
      }
    }
  }
}
