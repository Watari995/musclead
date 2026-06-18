import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../data/meal_dtos.dart';
import '../data/meal_repository.dart';

const _mealTypes = ['朝食', '昼食', '夕食', '間食'];

/// 食事記録のモーダルボトムシート。[existing] を渡すと編集モード。
Future<void> showMealRecordSheet(BuildContext context, {MealDto? existing}) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => _MealRecordSheet(existing: existing),
  );
}

class _MealRecordSheet extends HookConsumerWidget {
  const _MealRecordSheet({this.existing});

  final MealDto? existing;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final edit = existing;
    final isEdit = edit != null;
    final type = useState(edit?.mealType ?? '朝食');
    final calories = useTextEditingController(
      text: edit != null ? '${edit.calories}' : '',
    );
    final protein = useTextEditingController(
      text: edit?.proteinG?.toString() ?? '',
    );
    final fat = useTextEditingController(text: edit?.fatG?.toString() ?? '');
    final carb = useTextEditingController(
      text: edit?.carbohydrateG?.toString() ?? '',
    );
    final memo = useTextEditingController(text: edit?.memo ?? '');
    final loading = useState(false);
    final error = useState<String?>(null);
    final t = context.tokens;

    Future<void> submit() async {
      final kcal = int.tryParse(calories.text.trim());
      if (kcal == null) {
        error.value = 'カロリーを入力してください';
        return;
      }
      loading.value = true;
      error.value = null;
      final req = RecordMealRequest(
        eatenAt: edit?.eatenAt ?? DateTime.now(),
        mealType: type.value,
        calories: kcal,
        proteinG: double.tryParse(protein.text.trim()),
        fatG: double.tryParse(fat.text.trim()),
        carbohydrateG: double.tryParse(carb.text.trim()),
        memo: memo.text.trim().isEmpty ? null : memo.text.trim(),
      );
      try {
        final repo = ref.read(mealRepositoryProvider);
        if (isEdit) {
          await repo.update(edit.id, req);
        } else {
          await repo.record(req);
        }
        ref.invalidate(mealsProvider);
        if (context.mounted) Navigator.of(context).pop();
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = '保存に失敗しました';
      } finally {
        if (context.mounted) loading.value = false;
      }
    }

    Future<void> deleteMeal() async {
      loading.value = true;
      try {
        await ref.read(mealRepositoryProvider).delete(edit!.id);
        ref.invalidate(mealsProvider);
        if (context.mounted) Navigator.of(context).pop();
      } catch (_) {
        error.value = '削除に失敗しました';
        loading.value = false;
      }
    }

    const numeric = TextInputType.numberWithOptions(decimal: true);

    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(20, 4, 20, 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                isEdit ? '食事を編集' : '食事を記録',
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(height: 16),
              Wrap(
                spacing: 8,
                children: [
                  for (final m in _mealTypes)
                    ChoiceChip(
                      label: Text(m),
                      selected: type.value == m,
                      showCheckmark: false,
                      selectedColor: t.accentWeak,
                      onSelected: (_) => type.value = m,
                    ),
                ],
              ),
              const SizedBox(height: 14),
              AppTextField(
                label: 'カロリー (kcal)',
                controller: calories,
                hint: '420',
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 14),
              Row(
                children: [
                  Expanded(
                    child: AppTextField(
                      label: 'P (g)',
                      controller: protein,
                      hint: '28',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'F (g)',
                      controller: fat,
                      hint: '9',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'C (g)',
                      controller: carb,
                      hint: '58',
                      keyboardType: numeric,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 14),
              AppTextField(
                label: 'メモ ・任意',
                controller: memo,
                hint: 'オートミール・バナナ・卵',
              ),
              if (error.value != null) ...[
                const SizedBox(height: 12),
                Text(
                  error.value!,
                  style: TextStyle(color: t.accent, fontSize: 13),
                ),
              ],
              const SizedBox(height: 20),
              AppButton(
                label: isEdit ? '更新する' : '記録する',
                loading: loading.value,
                onPressed: submit,
              ),
              if (isEdit) ...[
                const SizedBox(height: 8),
                AppButton(
                  label: 'この記録を削除',
                  variant: AppButtonVariant.text,
                  onPressed: loading.value ? null : deleteMeal,
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
