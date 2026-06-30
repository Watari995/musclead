import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
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
        await ref.read(notificationRepositoryProvider).upsertWeeklyGoal(
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

    return Scaffold(
      appBar: AppBar(
        title: const Text('週次目標'),
        backgroundColor: t.surface,
        foregroundColor: t.ink,
        elevation: 0,
      ),
      backgroundColor: t.bg,
      body: ListView(
        padding: const EdgeInsets.all(18),
        children: [
          Text(
            '毎週日曜に達成チェックの通知が届きます。\n設定しない項目は空欄のままにしてください。',
            style: TextStyle(fontSize: 13, color: t.muted),
          ),
          const SizedBox(height: 24),
          _GoalField(
            label: 'トレーニング回数',
            unit: '回 / 週',
            controller: trainingCtrl,
            keyboardType: TextInputType.number,
          ),
          const SizedBox(height: 16),
          _GoalField(
            label: '平均カロリー目標',
            unit: 'kcal 以内 / 日',
            controller: calorieCtrl,
            keyboardType: TextInputType.number,
          ),
          const SizedBox(height: 16),
          _GoalField(
            label: '体重変化目標',
            unit: 'kg（例: -0.5 で減量）',
            controller: weightCtrl,
            keyboardType: const TextInputType.numberWithOptions(
              signed: true,
              decimal: true,
            ),
          ),
          const SizedBox(height: 32),
          AppButton(
            label: saving.value ? '保存中…' : '保存',
            onPressed: saving.value ? null : save,
          ),
          if (saved.value) ...[
            const SizedBox(height: 12),
            Text(
              '保存しました',
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
        Text(
          label,
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: t.ink,
          ),
        ),
        const SizedBox(height: 6),
        Row(
          children: [
            SizedBox(
              width: 120,
              child: TextField(
                controller: controller,
                keyboardType: keyboardType,
                decoration: InputDecoration(
                  hintText: '未設定',
                  hintStyle: TextStyle(color: t.muted),
                  filled: true,
                  fillColor: t.surface,
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 10,
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: t.line),
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                    borderSide: BorderSide(color: t.line),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 10),
            Flexible(
              child: Text(
                unit,
                style: TextStyle(fontSize: 13, color: t.muted),
              ),
            ),
          ],
        ),
      ],
    );
  }
}
