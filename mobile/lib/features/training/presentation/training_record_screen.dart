import 'package:decimal/decimal.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_card.dart';
import '../data/exercise_dtos.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';

class _SetDraft {
  final weight = TextEditingController();
  final reps = TextEditingController();
  void dispose() {
    weight.dispose();
    reps.dispose();
  }
}

class _ExerciseDraft {
  _ExerciseDraft(this.exerciseId, this.name);
  final String exerciseId;
  final String name;
  final List<_SetDraft> sets = [_SetDraft()];
  void dispose() {
    for (final s in sets) {
      s.dispose();
    }
  }
}

/// トレーニング記録画面（Training > 種目 > セット の入れ子を入力して保存）。
class TrainingRecordScreen extends ConsumerStatefulWidget {
  const TrainingRecordScreen({super.key, this.initialExercises = const []});

  /// ルーティンから開始する場合の初期種目（exerciseId + 表示名）。
  final List<({String exerciseId, String name})> initialExercises;

  @override
  ConsumerState<TrainingRecordScreen> createState() =>
      _TrainingRecordScreenState();
}

class _TrainingRecordScreenState extends ConsumerState<TrainingRecordScreen> {
  final DateTime _startedAt = DateTime.now();
  final List<_ExerciseDraft> _exercises = [];
  bool _saving = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    for (final e in widget.initialExercises) {
      _exercises.add(_ExerciseDraft(e.exerciseId, e.name));
    }
  }

  @override
  void dispose() {
    for (final e in _exercises) {
      e.dispose();
    }
    super.dispose();
  }

  void _addExercise(ExerciseDto ex) {
    if (_exercises.any((e) => e.exerciseId == ex.id)) return;
    setState(() => _exercises.add(_ExerciseDraft(ex.id, ex.name)));
  }

  void _removeExercise(int i) => setState(() {
    _exercises[i].dispose();
    _exercises.removeAt(i);
  });

  void _addSet(int i) => setState(() => _exercises[i].sets.add(_SetDraft()));

  void _removeSet(int ei, int si) => setState(() {
    _exercises[ei].sets[si].dispose();
    _exercises[ei].sets.removeAt(si);
  });

  Future<void> _save() async {
    final reqExercises = <RecordTrainingExerciseRequest>[];
    for (var i = 0; i < _exercises.length; i++) {
      final e = _exercises[i];
      final sets = <RecordTrainingSetRequest>[];
      var n = 1;
      for (final s in e.sets) {
        final w = Decimal.tryParse(s.weight.text.trim());
        final r = int.tryParse(s.reps.text.trim());
        if (w == null || r == null) continue;
        sets.add(
          RecordTrainingSetRequest(setNumber: n++, weightKg: w, reps: r),
        );
      }
      if (sets.isEmpty) continue;
      reqExercises.add(
        RecordTrainingExerciseRequest(
          exerciseId: e.exerciseId,
          displayOrder: i,
          sets: sets,
        ),
      );
    }
    if (reqExercises.isEmpty) {
      setState(() => _error = '有効なセット（重量と回数）を1つ以上入力してください');
      return;
    }
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      await ref
          .read(trainingRepositoryProvider)
          .recordTraining(
            RecordTrainingRequest(
              startedAt: _startedAt,
              endedAt: DateTime.now(),
              exercises: reqExercises,
            ),
          );
      ref.invalidate(trainingsProvider);
      if (mounted) context.go('/trainings');
    } on Failure catch (f) {
      setState(() => _error = f.message);
    } catch (_) {
      setState(() => _error = '保存に失敗しました');
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  Future<void> _pickExercise() async {
    final selected = await showModalBottomSheet<ExerciseDto>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (_) => const _ExercisePicker(),
    );
    if (selected != null) _addExercise(selected);
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Scaffold(
      appBar: AppBar(
        title: const Text('トレーニング記録'),
        backgroundColor: Colors.transparent,
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
          children: [
            for (var i = 0; i < _exercises.length; i++) _exerciseCard(i),
            const SizedBox(height: 4),
            AppButton(
              label: '種目を追加',
              icon: Icons.add,
              variant: AppButtonVariant.glass,
              onPressed: _pickExercise,
            ),
            if (_error != null) ...[
              const SizedBox(height: 12),
              Text(_error!, style: TextStyle(color: t.accent, fontSize: 13)),
            ],
          ],
        ),
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 20),
        child: SafeArea(
          top: false,
          child: AppButton(label: '保存', loading: _saving, onPressed: _save),
        ),
      ),
    );
  }

  Widget _exerciseCard(int i) {
    final e = _exercises[i];
    final t = context.tokens;
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: AppCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    e.name,
                    style: const TextStyle(
                      fontWeight: FontWeight.w700,
                      fontSize: 15,
                    ),
                  ),
                ),
                IconButton(
                  icon: Icon(Icons.close, size: 18, color: t.subtle),
                  onPressed: () => _removeExercise(i),
                ),
              ],
            ),
            for (var si = 0; si < e.sets.length; si++) _setRow(i, si),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                onPressed: () => _addSet(i),
                icon: const Icon(Icons.add, size: 16),
                label: const Text('セット追加'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _setRow(int ei, int si) {
    final sets = _exercises[ei].sets;
    final s = sets[si];
    final t = context.tokens;
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          SizedBox(
            width: 26,
            child: Text('${si + 1}', style: TextStyle(color: t.muted)),
          ),
          Expanded(child: _numField(s.weight, 'kg')),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 8),
            child: Text('×'),
          ),
          Expanded(child: _numField(s.reps, '回')),
          IconButton(
            icon: Icon(Icons.remove_circle_outline, size: 18, color: t.subtle),
            onPressed: sets.length > 1 ? () => _removeSet(ei, si) : null,
          ),
        ],
      ),
    );
  }

  Widget _numField(TextEditingController c, String suffix) => TextField(
    controller: c,
    keyboardType: const TextInputType.numberWithOptions(decimal: true),
    decoration: InputDecoration(
      isDense: true,
      suffixText: suffix,
      border: const OutlineInputBorder(),
      contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
    ),
  );
}

/// 種目選択シート（既存から選ぶ / 新規作成）。
class _ExercisePicker extends ConsumerStatefulWidget {
  const _ExercisePicker();

  @override
  ConsumerState<_ExercisePicker> createState() => _ExercisePickerState();
}

class _ExercisePickerState extends ConsumerState<_ExercisePicker> {
  final _name = TextEditingController();
  bool _creating = false;

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> _create() async {
    final name = _name.text.trim();
    if (name.isEmpty) return;
    setState(() => _creating = true);
    try {
      await ref.read(trainingRepositoryProvider).createExercise(name);
      ref.invalidate(exercisesProvider);
      _name.clear();
    } catch (_) {
      // 一覧は AsyncValue 側でエラー表示
    } finally {
      if (mounted) setState(() => _creating = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final exercises = ref.watch(exercisesProvider);
    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(16, 4, 16, 16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                '種目を選択',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _name,
                      decoration: const InputDecoration(
                        hintText: '新規種目名',
                        isDense: true,
                        border: OutlineInputBorder(),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  AppButton(
                    label: '作成',
                    expand: false,
                    loading: _creating,
                    onPressed: _create,
                  ),
                ],
              ),
              const SizedBox(height: 12),
              ConstrainedBox(
                constraints: BoxConstraints(
                  maxHeight: MediaQuery.of(context).size.height * 0.45,
                ),
                child: exercises.when(
                  data: (list) => list.isEmpty
                      ? const Padding(
                          padding: EdgeInsets.all(20),
                          child: Text('種目がありません。上で作成してください'),
                        )
                      : ListView(
                          shrinkWrap: true,
                          children: [
                            for (final ex in list)
                              ListTile(
                                title: Text(ex.name),
                                onTap: () => Navigator.of(context).pop(ex),
                              ),
                          ],
                        ),
                  loading: () => const Padding(
                    padding: EdgeInsets.all(20),
                    child: Center(child: CircularProgressIndicator()),
                  ),
                  error: (e, _) => const Padding(
                    padding: EdgeInsets.all(20),
                    child: Text('読み込みに失敗しました'),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
