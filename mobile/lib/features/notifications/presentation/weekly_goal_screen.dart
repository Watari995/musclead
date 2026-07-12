import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/widgets/app_button.dart';
import '../../../l10n/app_localizations.dart';
import '../data/notification_repository.dart';

class WeeklyGoalScreen extends HookConsumerWidget {
  const WeeklyGoalScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = context.tokens;
    final goalAsync = ref.watch(weeklyGoalProvider);

    final trainingCtrl = useTextEditingController();
    final calorieCtrl = useTextEditingController();
    final weightCtrl = useTextEditingController();
    final initialized = useState(false);
    final saving = useState(false);
    final saved = useState(false);

    goalAsync.whenData((goal) {
      if (!initialized.value) {
        trainingCtrl.text = goal?.trainingCount?.toString() ?? '';
        calorieCtrl.text = goal?.calorieAverage?.toString() ?? '';
        weightCtrl.text = goal?.weightChangeKg?.toString() ?? '';
        initialized.value = true;
      }
    });

    Future<void> save() async {
      saving.value = true;
      try {
        await ref
            .read(notificationRepositoryProvider)
            .upsertWeeklyGoal(
              trainingCount: trainingCtrl.text.isEmpty
                  ? null
                  : int.tryParse(trainingCtrl.text),
              calorieAverage: calorieCtrl.text.isEmpty
                  ? null
                  : int.tryParse(calorieCtrl.text),
              weightChangeKg: weightCtrl.text.isEmpty
                  ? null
                  : double.tryParse(weightCtrl.text),
            );
        ref.invalidate(weeklyGoalProvider);
        saved.value = true;
        await Future.delayed(const Duration(seconds: 2));
        saved.value = false;
      } finally {
        saving.value = false;
      }
    }

    final l = AppLocalizations.of(context)!;
    return Scaffold(
      appBar: AppBar(
        title: Text(l.weeklyGoalTitle),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: ListView(
        padding: const EdgeInsets.all(18),
        children: [
          Text(
            l.weeklyGoalDescription,
            style: TextStyle(fontSize: 13, color: t.muted),
          ),
          const SizedBox(height: 24),
          _GoalField(
            label: l.weeklyGoalTrainingCount,
            unit: l.weeklyGoalTrainingUnit,
            controller: trainingCtrl,
            keyboardType: TextInputType.number,
          ),
          const SizedBox(height: 16),
          _GoalField(
            label: l.weeklyGoalCalorieGoal,
            unit: l.weeklyGoalCalorieUnit,
            controller: calorieCtrl,
            keyboardType: TextInputType.number,
          ),
          const SizedBox(height: 16),
          _GoalField(
            label: l.weeklyGoalWeightChange,
            unit: l.weeklyGoalWeightUnit,
            controller: weightCtrl,
            keyboardType: const TextInputType.numberWithOptions(
              signed: true,
              decimal: true,
            ),
          ),
          const SizedBox(height: 32),
          AppButton(
            label: saving.value ? l.weeklyGoalSaving : l.weeklyGoalSave,
            onPressed: saving.value ? null : save,
          ),
          if (saved.value) ...[
            const SizedBox(height: 12),
            Text(
              l.weeklyGoalSaved,
              textAlign: TextAlign.center,
              style: TextStyle(color: t.accent, fontSize: 14),
            ),
          ],
        ],
      ),
    );
  }
}

class _GoalField extends StatelessWidget {
  const _GoalField({
    required this.label,
    required this.unit,
    required this.controller,
    required this.keyboardType,
  });

  final String label;
  final String unit;
  final TextEditingController controller;
  final TextInputType keyboardType;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: TextStyle(fontFamily: 'Caveat', fontSize: 18)),
        const SizedBox(height: 6),
        Row(
          children: [
            SizedBox(
              width: 120,
              child: RoughTextField(
                controller: controller,
                keyboardType: keyboardType,
                hint: AppLocalizations.of(context)!.weeklyGoalNotSet,
              ),
            ),
            const SizedBox(width: 10),
            Flexible(
              child: Text(unit, style: TextStyle(fontSize: 13, color: t.muted)),
            ),
          ],
        ),
      ],
    );
  }
}
