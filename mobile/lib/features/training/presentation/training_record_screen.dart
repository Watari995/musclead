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
  _SetDraft({String weight = '', String reps = ''})
    : weight = TextEditingController(text: weight),
      reps = TextEditingController(text: reps);

  final TextEditingController weight;
  final TextEditingController reps;

  void dispose() {
    weight.dispose();
    reps.dispose();
  }
}

class _ExerciseDraft {
  _ExerciseDraft(this.exerciseId, this.name);

  String exerciseId;
  String name;
  final List<_SetDraft> sets = [_SetDraft()];
  final TextEditingController memo = TextEditingController();

  void dispose() {
    for (final s in sets) {
      s.dispose();
    }
    memo.dispose();
  }
}

/// トレーニング記録画面（Training > 種目 > セット の入れ子を入力して保存）。
/// [editingTraining] を渡すと既存記録の編集モードになる。
class TrainingRecordScreen extends ConsumerStatefulWidget {
  const TrainingRecordScreen({
    super.key,
    this.initialExercises = const [],
    this.editingTraining,
  });

  /// ルーティンから開始する場合の初期種目（exerciseId + 表示名）。
  final List<({String exerciseId, String name})> initialExercises;

  /// 編集モード時の既存トレーニングデータ。
  final TrainingDto? editingTraining;

  @override
  ConsumerState<TrainingRecordScreen> createState() =>
      _TrainingRecordScreenState();
}

class _TrainingRecordScreenState extends ConsumerState<TrainingRecordScreen> {
  late DateTime _startedAt;
  final List<_ExerciseDraft> _exercises = [];
  final TextEditingController _memo = TextEditingController();
  Map<String, BestSetDto> _bestSets = {};
  bool _saving = false;
  String? _error;

  bool get _isEditing => widget.editingTraining != null;

  @override
  void initState() {
    super.initState();
    _startedAt = widget.editingTraining?.startedAt ?? DateTime.now();
    if (_isEditing) {
      _initFromTraining(widget.editingTraining!);
    } else {
      for (final e in widget.initialExercises) {
        _exercises.add(_ExerciseDraft(e.exerciseId, e.name));
      }
    }
    if (_exercises.isNotEmpty) _loadBestSets();
  }

  void _initFromTraining(TrainingDto training) {
    _memo.text = training.memo ?? '';
    final exList = ref.read(exercisesProvider).asData?.value ?? [];
    final names = {for (final e in exList) e.id: e.name};
    for (final ex in training.exercises) {
      final draft = _ExerciseDraft(
        ex.exerciseId,
        names[ex.exerciseId] ?? '種目',
      );
      draft.sets.first.dispose();
      draft.sets.clear();
      for (final s in ex.sets) {
        draft.sets.add(
          _SetDraft(
            weight: s.weightKg.toString(),
            reps: s.reps.toString(),
          ),
        );
      }
      if (draft.sets.isEmpty) draft.sets.add(_SetDraft());
      draft.memo.text = ex.memo ?? '';
      _exercises.add(draft);
    }
  }

  @override
  void dispose() {
    for (final e in _exercises) {
      e.dispose();
    }
    _memo.dispose();
    super.dispose();
  }

  Future<void> _loadBestSets() async {
    final ids = _exercises.map((e) => e.exerciseId).toList();
    if (ids.isEmpty) return;
    try {
      final list = await ref.read(trainingRepositoryProvider).bestSets(ids);
      if (!mounted) return;
      setState(() => _bestSets = {for (final b in list) b.exerciseId: b});
    } catch (_) {
      // 自己ベストは取得できなくても致命的でないため無視
    }
  }

  void _addExercise(ExerciseDto ex) {
    if (_exercises.any((e) => e.exerciseId == ex.id)) return;
    setState(() => _exercises.add(_ExerciseDraft(ex.id, ex.name)));
    _loadBestSets();
  }

  void _removeExercise(int i) => setState(() {
    _exercises[i].dispose();
    _exercises.removeAt(i);
  });

  /// セット追加時は直前セットの kg / 回数 を引き継ぐ。
  void _addSet(int i) => setState(() {
    final sets = _exercises[i].sets;
    final last = sets.isNotEmpty ? sets.last : null;
    // kg は引き継ぐが、回数は引き継がない。
    sets.add(_SetDraft(weight: last?.weight.text ?? ''));
  });

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
          memo: e.memo.text.trim().isEmpty ? null : e.memo.text.trim(),
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
      final req = RecordTrainingRequest(
        startedAt: _startedAt,
        endedAt: _isEditing ? widget.editingTraining!.endedAt : DateTime.now(),
        memo: _memo.text.trim().isEmpty ? null : _memo.text.trim(),
        exercises: reqExercises,
      );
      if (_isEditing) {
        await ref
            .read(trainingRepositoryProvider)
            .updateTraining(widget.editingTraining!.id, req);
      } else {
        await ref.read(trainingRepositoryProvider).recordTraining(req);
      }
      ref.invalidate(trainingsProvider);
      if (mounted) {
        if (_isEditing) {
          Navigator.of(context, rootNavigator: true)
              .popUntil((route) => route.isFirst);
        } else {
          context.go('/trainings');
        }
      }
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

  /// 既存カードの種目を別の種目に切り替える（入力済みのセットは保持）。
  Future<void> _switchExercise(int i) async {
    final selected = await showModalBottomSheet<ExerciseDto>(
      context: context,
      isScrollControlled: true,
      showDragHandle: true,
      builder: (_) => const _ExercisePicker(),
    );
    if (selected == null) return;
    setState(() {
      _exercises[i].exerciseId = selected.id;
      _exercises[i].name = selected.name;
    });
    _loadBestSets();
  }

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? '記録を編集' : 'トレーニング記録'),
        backgroundColor: Colors.transparent,
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        behavior: HitTestBehavior.translucent,
        child: SafeArea(
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
              const SizedBox(height: 12),
              _overallMemoCard(),
              if (_error != null) ...[
                const SizedBox(height: 12),
                Text(_error!, style: TextStyle(color: t.accent, fontSize: 13)),
              ],
            ],
          ),
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

  Widget _overallMemoCard() {
    final t = context.tokens;
    return AppCard(
      padding: EdgeInsets.zero,
      child: Theme(
        data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
        child: ExpansionTile(
          tilePadding: const EdgeInsets.symmetric(horizontal: 15),
          childrenPadding: const EdgeInsets.fromLTRB(15, 0, 15, 14),
          title: Text(
            'メモ（全体）',
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: t.muted,
              fontSize: 14,
            ),
          ),
          children: [
            TextField(
              controller: _memo,
              maxLines: 3,
              decoration: const InputDecoration(
                hintText: '今日のコンディション など',
                border: OutlineInputBorder(),
                isDense: true,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _exerciseCard(int i) {
    final e = _exercises[i];
    final t = context.tokens;
    final best = _bestSets[e.exerciseId];
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: AppCard(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: InkWell(
                    onTap: () => _switchExercise(i),
                    borderRadius: BorderRadius.circular(8),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: Row(
                        children: [
                          Flexible(
                            child: Text(
                              e.name,
                              style: const TextStyle(
                                fontWeight: FontWeight.w700,
                                fontSize: 15,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                          const SizedBox(width: 4),
                          Icon(Icons.unfold_more, size: 16, color: t.subtle),
                        ],
                      ),
                    ),
                  ),
                ),
                IconButton(
                  icon: Icon(Icons.close, size: 18, color: t.subtle),
                  onPressed: () => _removeExercise(i),
                ),
              ],
            ),
            if (best != null) _bestSetLine(best),
            const SizedBox(height: 4),
            for (var si = 0; si < e.sets.length; si++) _setRow(i, si),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                onPressed: () => _addSet(i),
                icon: const Icon(Icons.add, size: 16),
                label: const Text('セット追加'),
              ),
            ),
            _exerciseMemo(e),
          ],
        ),
      ),
    );
  }

  Widget _bestSetLine(BestSetDto best) {
    final t = context.tokens;
    final d = best.performedAt?.toLocal();
    final date = d == null ? '' : ' (${d.year}/${d.month}/${d.day})';
    return Padding(
      padding: const EdgeInsets.only(top: 2, bottom: 4),
      child: Text.rich(
        TextSpan(
          children: [
            TextSpan(
              text: '★ ',
              style: TextStyle(color: t.gold),
            ),
            TextSpan(
              text: '最高記録 ',
              style: TextStyle(color: t.muted, fontWeight: FontWeight.w600),
            ),
            TextSpan(
              text: '${best.weightKg}kg × ${best.reps}回$date',
              style: TextStyle(color: t.muted),
            ),
          ],
        ),
        style: const TextStyle(fontSize: 12),
      ),
    );
  }

  Widget _exerciseMemo(_ExerciseDraft e) {
    final t = context.tokens;
    return Theme(
      data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
      child: ExpansionTile(
        tilePadding: EdgeInsets.zero,
        childrenPadding: const EdgeInsets.only(bottom: 8),
        title: Text('メモ', style: TextStyle(fontSize: 13, color: t.muted)),
        children: [
          TextField(
            controller: e.memo,
            maxLines: 2,
            decoration: const InputDecoration(
              hintText: 'フォーム・感覚 など',
              border: OutlineInputBorder(),
              isDense: true,
            ),
          ),
        ],
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
