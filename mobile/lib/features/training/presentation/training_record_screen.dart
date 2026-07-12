import 'package:decimal/decimal.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_card.dart';
import '../data/exercise_dtos.dart';
import '../data/training_dtos.dart';
import '../data/training_repository.dart';
import '../../../l10n/app_localizations.dart';

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
    this.focusExerciseId,
  });

  final List<({String exerciseId, String name})> initialExercises;
  final TrainingDto? editingTraining;

  /// 指定すると、この種目のカードまで初期表示時にスクロールする。
  final String? focusExerciseId;

  @override
  ConsumerState<TrainingRecordScreen> createState() =>
      _TrainingRecordScreenState();
}

class _TrainingRecordScreenState extends ConsumerState<TrainingRecordScreen> {
  late DateTime _startedAt;
  DateTime? _endedAt;
  final List<_ExerciseDraft> _exercises = [];
  final Map<int, GlobalKey> _exerciseKeys = {};
  final TextEditingController _memo = TextEditingController();
  Map<String, BestSetDto> _bestSets = {};
  Map<String, LastSessionSetsByExerciseDto> _lastSessionSets = {};
  bool _saving = false;
  String? _error;

  bool get _isEditing => widget.editingTraining != null;

  @override
  void initState() {
    super.initState();
    _startedAt = widget.editingTraining?.startedAt ?? DateTime.now();
    _endedAt = widget.editingTraining?.endedAt;
    if (_isEditing) {
      _initFromTraining(widget.editingTraining!);
    } else {
      for (final e in widget.initialExercises) {
        _exercises.add(_ExerciseDraft(e.exerciseId, e.name));
      }
    }
    if (_exercises.isNotEmpty) {
      _loadBestSets();
      _loadLastSessionSets();
    }
    if (widget.focusExerciseId != null) {
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => _scrollToFocusedExercise(),
      );
    }
  }

  void _scrollToFocusedExercise() {
    final index = _exercises.indexWhere(
      (e) => e.exerciseId == widget.focusExerciseId,
    );
    final targetContext = _exerciseKeys[index]?.currentContext;
    if (targetContext == null) return;
    Scrollable.ensureVisible(
      targetContext,
      duration: const Duration(milliseconds: 300),
      alignment: 0.1,
    );
  }

  void _initFromTraining(TrainingDto training) {
    _memo.text = training.memo ?? '';
    final exList = ref.read(exercisesProvider).asData?.value ?? [];
    final names = {for (final e in exList) e.id: e.name};
    for (final ex in training.exercises) {
      final draft = _ExerciseDraft(
        ex.exerciseId,
        names[ex.exerciseId] ?? ex.exerciseId,
      );
      draft.sets.first.dispose();
      draft.sets.clear();
      for (final s in ex.sets) {
        draft.sets.add(
          _SetDraft(weight: s.weightKg.toString(), reps: s.reps.toString()),
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
    } catch (_) {}
  }

  Future<void> _loadLastSessionSets() async {
    final ids = _exercises.map((e) => e.exerciseId).toList();
    if (ids.isEmpty) return;
    try {
      final map = await ref
          .read(trainingRepositoryProvider)
          .lastSessionSets(ids);
      if (!mounted) return;
      setState(() => _lastSessionSets = map);
    } catch (_) {}
  }

  void _addExercise(ExerciseDto ex) {
    if (_exercises.any((e) => e.exerciseId == ex.id)) return;
    setState(() => _exercises.add(_ExerciseDraft(ex.id, ex.name)));
    _loadBestSets();
    _loadLastSessionSets();
  }

  void _removeExercise(int i) => setState(() {
    _exercises[i].dispose();
    _exercises.removeAt(i);
  });

  void _addSet(int i) => setState(() {
    final sets = _exercises[i].sets;
    final last = sets.isNotEmpty ? sets.last : null;
    sets.add(_SetDraft(weight: last?.weight.text ?? ''));
  });

  void _removeSet(int ei, int si) => setState(() {
    _exercises[ei].sets[si].dispose();
    _exercises[ei].sets.removeAt(si);
  });

  Future<void> _save() async {
    final l = AppLocalizations.of(context)!;
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
      setState(() => _error = l.trainingExerciseRequired2);
      return;
    }
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      final req = RecordTrainingRequest(
        startedAt: _startedAt,
        endedAt: _endedAt,
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
          Navigator.of(
            context,
            rootNavigator: true,
          ).popUntil((route) => route.isFirst);
        } else {
          context.go('/trainings');
        }
      }
    } on Failure catch (f) {
      setState(() => _error = f.message);
    } catch (_) {
      if (mounted) {
        setState(() => _error = AppLocalizations.of(context)!.commonSaveFailed);
      }
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
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? l.trainingEditTitle : l.trainingRecordTitle),
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
                label: l.trainingExercisesLabel,
                icon: Icons.add,
                variant: AppButtonVariant.glass,
                onPressed: _pickExercise,
              ),
              const SizedBox(height: 12),
              _overallMemoCard(l),
              const SizedBox(height: 12),
              _endedAtCard(l),
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
          child: AppButton(
            label: l.trainingSaveBtn,
            loading: _saving,
            onPressed: _save,
          ),
        ),
      ),
    );
  }

  Widget _overallMemoCard(AppLocalizations l) {
    final t = context.tokens;
    return AppCard(
      padding: EdgeInsets.zero,
      child: Theme(
        data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
        child: ExpansionTile(
          tilePadding: const EdgeInsets.symmetric(horizontal: 15),
          childrenPadding: const EdgeInsets.fromLTRB(15, 0, 15, 14),
          title: Text(
            l.trainingMemoAll,
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: t.muted,
              fontSize: 14,
            ),
          ),
          children: [
            RoughTextField(
              controller: _memo,
              hint: l.trainingMemoAllHint,
              maxLines: 3,
            ),
          ],
        ),
      ),
    );
  }

  Widget _endedAtCard(AppLocalizations l) {
    final t = context.tokens;
    final label = _endedAt == null
        ? l.trainingEndTimeEmpty
        : l.trainingEndTimeSet(
            '${_endedAt!.toLocal().hour.toString().padLeft(2, '0')}:${_endedAt!.toLocal().minute.toString().padLeft(2, '0')}',
          );
    return AppCard(
      child: InkWell(
        onTap: () async {
          final now = DateTime.now();
          final initial = _endedAt ?? now;
          final picked = await showTimePicker(
            context: context,
            initialTime: TimeOfDay(hour: initial.hour, minute: initial.minute),
          );
          if (picked == null) return;
          final base = _endedAt ?? now;
          setState(() {
            _endedAt = DateTime(
              base.year,
              base.month,
              base.day,
              picked.hour,
              picked.minute,
            );
          });
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 14),
          child: Row(
            children: [
              Icon(Icons.timer_off_outlined, size: 18, color: t.muted),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  label,
                  style: TextStyle(
                    fontWeight: FontWeight.w600,
                    color: _endedAt == null ? t.muted : null,
                    fontSize: 14,
                  ),
                ),
              ),
              if (_endedAt != null)
                GestureDetector(
                  onTap: () => setState(() => _endedAt = null),
                  child: Icon(Icons.close, size: 16, color: t.subtle),
                ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _exerciseCard(int i) {
    final l = AppLocalizations.of(context)!;
    final e = _exercises[i];
    final t = context.tokens;
    final best = _bestSets[e.exerciseId];
    final lastSession = _lastSessionSets[e.exerciseId];
    return Padding(
      key: _exerciseKeys.putIfAbsent(i, () => GlobalKey()),
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
                                fontFamily: 'Caveat',
                                fontSize: 19,
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
            if (best != null) _bestSetLine(best, l),
            if (lastSession != null && lastSession.sets.isNotEmpty)
              _lastSessionLine(lastSession, l),
            const SizedBox(height: 4),
            for (var si = 0; si < e.sets.length; si++) _setRow(i, si, l),
            Align(
              alignment: Alignment.centerLeft,
              child: TextButton.icon(
                onPressed: () => _addSet(i),
                icon: const Icon(Icons.add, size: 16),
                label: Text(l.trainingAddSet),
              ),
            ),
            _exerciseMemo(e, l),
          ],
        ),
      ),
    );
  }

  Widget _lastSessionLine(
    LastSessionSetsByExerciseDto lastSession,
    AppLocalizations l,
  ) {
    final t = context.tokens;
    final d = lastSession.performedAt.toLocal();
    final date = '${d.year}/${d.month}/${d.day}';
    final setsText = lastSession.sets
        .map((s) => '${s.setNumber}. ${s.weightKg}kg×${s.reps}')
        .join('  ');
    return Padding(
      padding: const EdgeInsets.only(top: 2, bottom: 2),
      child: Text.rich(
        TextSpan(
          children: [
            TextSpan(
              text: '📅 ${l.trainingPrevSession}',
              style: TextStyle(color: t.muted, fontWeight: FontWeight.w600),
            ),
            TextSpan(
              text: '($date)  ',
              style: TextStyle(color: t.muted),
            ),
            TextSpan(
              text: setsText,
              style: TextStyle(color: t.muted),
            ),
          ],
        ),
        style: const TextStyle(fontSize: 12),
      ),
    );
  }

  Widget _bestSetLine(BestSetDto best, AppLocalizations l) {
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
              text: l.trainingBestRecord,
              style: TextStyle(color: t.muted, fontWeight: FontWeight.w600),
            ),
            TextSpan(
              text: '${best.weightKg}kg × ${best.reps}${l.trainingReps}$date',
              style: TextStyle(color: t.muted),
            ),
          ],
        ),
        style: const TextStyle(fontSize: 12),
      ),
    );
  }

  Widget _exerciseMemo(_ExerciseDraft e, AppLocalizations l) {
    final t = context.tokens;
    return Theme(
      data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
      child: ExpansionTile(
        tilePadding: EdgeInsets.zero,
        childrenPadding: const EdgeInsets.only(bottom: 8),
        title: Text(
          l.trainingExerciseMemo,
          style: TextStyle(fontSize: 13, color: t.muted),
        ),
        children: [
          RoughTextField(
            controller: e.memo,
            hint: l.trainingExerciseMemoHint,
            maxLines: 2,
          ),
        ],
      ),
    );
  }

  Widget _setRow(int ei, int si, AppLocalizations l) {
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
          Expanded(child: _numField(s.reps, l.trainingReps)),
          IconButton(
            icon: Icon(Icons.remove_circle_outline, size: 18, color: t.subtle),
            onPressed: sets.length > 1 ? () => _removeSet(ei, si) : null,
          ),
        ],
      ),
    );
  }

  Widget _numField(TextEditingController c, String suffix) => RoughTextField(
    controller: c,
    keyboardType: const TextInputType.numberWithOptions(decimal: true),
    suffixText: suffix,
    textAlign: TextAlign.center,
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
    final l = AppLocalizations.of(context)!;
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
              Text(
                l.trainingSelectExercise,
                style: const TextStyle(fontFamily: 'Caveat', fontSize: 24),
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: RoughTextField(
                      controller: _name,
                      hint: l.trainingNewExerciseName,
                    ),
                  ),
                  const SizedBox(width: 8),
                  AppButton(
                    label: l.commonCreate,
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
                      ? Padding(
                          padding: const EdgeInsets.all(20),
                          child: Text(l.trainingNoExercisesCreate),
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
                  error: (e, _) => Padding(
                    padding: const EdgeInsets.all(20),
                    child: Text(l.commonLoadFailed),
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
