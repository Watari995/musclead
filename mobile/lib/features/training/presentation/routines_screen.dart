import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../l10n/app_localizations.dart';
import '../data/routine_dtos.dart';
import '../data/training_repository.dart';
import 'training_record_screen.dart';

/// ルーティンの一覧 / 削除 / 並び替え / タップで記録開始。
class RoutinesScreen extends ConsumerStatefulWidget {
  const RoutinesScreen({super.key});

  @override
  ConsumerState<RoutinesScreen> createState() => _RoutinesScreenState();
}

class _RoutinesScreenState extends ConsumerState<RoutinesScreen> {
  List<RoutineDto>? _items;

  bool _sameIds(List<RoutineDto> a, List<RoutineDto> b) {
    if (a.length != b.length) return false;
    final ids = b.map((e) => e.id).toSet();
    return a.every((e) => ids.contains(e.id));
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final async = ref.watch(routinesProvider);
    return Scaffold(
      appBar: AppBar(
        title: Text(l.trainingRoutinesLabel),
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
          value: async,
          onRetry: () => ref.invalidate(routinesProvider),
          data: (list) {
            if (_items == null || !_sameIds(_items!, list)) {
              _items = List.of(list);
            }
            final items = _items!;
            if (items.isEmpty) {
              return Center(child: Text(l.trainingRoutinesEmpty));
            }
            return ReorderableListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: items.length,
              onReorder: _onReorder,
              buildDefaultDragHandles: false,
              proxyDecorator: (child, index, animation) =>
                  Material(color: Colors.transparent, child: child),
              itemBuilder: (context, i) => _tile(items[i], i),
            );
          },
        ),
      ),
    );
  }

  void _onReorder(int oldIndex, int newIndex) {
    final l = AppLocalizations.of(context)!;
    setState(() {
      final items = _items!;
      var target = newIndex;
      if (target > oldIndex) target -= 1;
      final moved = items.removeAt(oldIndex);
      items.insert(target, moved);
    });
    ref
        .read(trainingRepositoryProvider)
        .reorderRoutines(_items!.map((e) => e.id).toList())
        .catchError((Object _) {
          if (mounted) {
            ScaffoldMessenger.of(
              context,
            ).showSnackBar(SnackBar(content: Text(l.commonReorderFailed)));
          }
        });
  }

  Widget _tile(RoutineDto r, int index) {
    final t = context.tokens;
    return Padding(
      key: ValueKey(r.id),
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
          onTap: () {
            final l = AppLocalizations.of(context)!;
            Navigator.of(context, rootNavigator: true).push(
              MaterialPageRoute<void>(
                builder: (_) => TrainingRecordScreen(
                  initialExercises: [
                    for (final e in r.routineExercises)
                      (exerciseId: e.exerciseId, name: e.exerciseName ?? l.trainingExerciseDefault),
                  ],
                ),
              ),
            );
          },
          title: Text(
            r.name,
            style: const TextStyle(fontWeight: FontWeight.w600),
          ),
          subtitle: Builder(
            builder: (context) {
              final l = AppLocalizations.of(context)!;
              return Text(
                l.trainingRoutineTapHint(r.routineExercises.length),
                style: TextStyle(fontSize: 12, color: t.muted),
              );
            },
          ),
          trailing: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              IconButton(
                icon: Icon(Icons.delete_outline, size: 20, color: t.subtle),
                onPressed: () => _delete(r),
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

  Future<void> _delete(RoutineDto r) async {
    try {
      await ref.read(trainingRepositoryProvider).deleteRoutine(r.id);
      ref.invalidate(routinesProvider);
    } on Failure catch (f) {
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(f.message)));
      }
    } catch (_) {
      if (mounted) {
        final l = AppLocalizations.of(context)!;
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(l.commonDeleteFailed)));
      }
    }
  }
}
