import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../data/routine_dtos.dart';
import '../data/training_repository.dart';
import '../../../l10n/app_localizations.dart';

/// ルーティン作成（名前 + 種目を順番に選択）。
class RoutineCreateScreen extends ConsumerStatefulWidget {
  const RoutineCreateScreen({super.key});

  @override
  ConsumerState<RoutineCreateScreen> createState() =>
      _RoutineCreateScreenState();
}

class _RoutineCreateScreenState extends ConsumerState<RoutineCreateScreen> {
  final _name = TextEditingController();
  final List<String> _selected = []; // 選択順 = displayOrder
  bool _saving = false;
  String? _error;

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  void _toggle(String id) => setState(() {
    if (_selected.contains(id)) {
      _selected.remove(id);
    } else {
      _selected.add(id);
    }
  });

  Future<void> _save() async {
    final l = AppLocalizations.of(context)!;
    final name = _name.text.trim();
    if (name.isEmpty) {
      setState(() => _error = l.commonNameRequired);
      return;
    }
    if (_selected.isEmpty) {
      setState(() => _error = l.trainingExerciseRequired);
      return;
    }
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      await ref
          .read(trainingRepositoryProvider)
          .createRoutine(
            UpsertRoutineRequest(
              name: name,
              exercises: [
                for (var i = 0; i < _selected.length; i++)
                  UpsertRoutineExerciseRequest(
                    exerciseId: _selected[i],
                    displayOrder: i,
                  ),
              ],
            ),
          );
      ref.invalidate(routinesProvider);
      if (mounted) context.pop();
    } on Failure catch (f) {
      setState(() => _error = f.message);
    } catch (_) {
      if (mounted) {
        final l2 = AppLocalizations.of(context)!;
        setState(() => _error = l2.commonSaveFailed);
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final exercises = ref.watch(exercisesProvider);
    return Scaffold(
      appBar: AppBar(
        title: Text(l.trainingNewRoutine),
        backgroundColor: Colors.transparent,
      ),
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
          children: [
            AppTextField(
              label: l.trainingRoutineName,
              controller: _name,
              hint: l.trainingRoutineHint,
            ),
            const SizedBox(height: 16),
            Text(
              l.trainingSelectExercisesHint,
              style: TextStyle(
                fontWeight: FontWeight.w700,
                color: t.muted,
                fontSize: 12.5,
              ),
            ),
            const SizedBox(height: 4),
            exercises.when(
              data: (list) => list.isEmpty
                  ? Padding(
                      padding: const EdgeInsets.all(16),
                      child: Text(l.trainingNoExercisesYet),
                    )
                  : Column(
                      children: [
                        for (final ex in list)
                          CheckboxListTile(
                            value: _selected.contains(ex.id),
                            title: Text(ex.name),
                            dense: true,
                            controlAffinity: ListTileControlAffinity.leading,
                            onChanged: (_) => _toggle(ex.id),
                          ),
                      ],
                    ),
              loading: () => const Padding(
                padding: EdgeInsets.all(16),
                child: Center(child: CircularProgressIndicator()),
              ),
              error: (e, _) => Padding(
                padding: const EdgeInsets.all(16),
                child: Text(l.trainingExerciseLoadFailed),
              ),
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
          child: AppButton(
            label: l.commonCreate,
            loading: _saving,
            onPressed: _save,
          ),
        ),
      ),
    );
  }
}
