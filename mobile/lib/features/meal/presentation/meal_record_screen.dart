import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../../food/data/food_product_dtos.dart';
import '../../food/presentation/food_register_sheet.dart';
import '../../food/presentation/food_search_section.dart';
import '../../../l10n/app_localizations.dart';
import '../data/meal_dtos.dart';
import '../data/meal_repository.dart';
import '../data/meal_template_dtos.dart';

/// 食事記録・編集ページ。[existing] を渡すと編集モード、[fromTemplate] でプリフィル。
class MealRecordScreen extends HookConsumerWidget {
  const MealRecordScreen({super.key, this.existing, this.fromTemplate});

  final MealDto? existing;
  final MealTemplateDto? fromTemplate;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final edit = existing;
    final tpl = fromTemplate;
    final isEdit = edit != null;

    final l = AppLocalizations.of(context)!;
    final mealTypes = [l.mealBreakfast, l.mealLunch, l.mealDinner, l.mealSnack];
    final mealType = useState(edit?.mealType ?? tpl?.mealType ?? l.mealBreakfast);
    final eatenAt = useState(edit?.eatenAt ?? DateTime.now());
    final caloriesCtrl = useTextEditingController(
      text: edit != null
          ? '${edit.calories}'
          : tpl != null
          ? '${tpl.calories}'
          : '',
    );
    final proteinCtrl = useTextEditingController(
      text: edit?.proteinG?.toString() ?? tpl?.proteinG ?? '',
    );
    final fatCtrl = useTextEditingController(
      text: edit?.fatG?.toString() ?? tpl?.fatG ?? '',
    );
    final carbCtrl = useTextEditingController(
      text: edit?.carbohydrateG?.toString() ?? tpl?.carbohydrateG ?? '',
    );
    final memoCtrl = useTextEditingController(text: edit?.memo ?? '');
    final loading = useState(false);
    final error = useState<String?>(null);
    final selectedFoodId = useState<String?>(edit?.foodProductId);
    final baseFood = useState<FoodProductDto?>(null);
    final servingCount = useState<double>(
      double.tryParse(edit?.servingCount ?? '') ?? 1.0,
    );
    final t = context.tokens;
    const numeric = TextInputType.numberWithOptions(decimal: true);

    void updateServing(double s) {
      final base = baseFood.value;
      if (base == null) return;
      servingCount.value = s;
      caloriesCtrl.text = '${(base.calories * s).round()}';
      proteinCtrl.text = base.proteinG != null
          ? ((double.tryParse(base.proteinG!) ?? 0) * s).toStringAsFixed(1)
          : '';
      fatCtrl.text = base.fatG != null
          ? ((double.tryParse(base.fatG!) ?? 0) * s).toStringAsFixed(1)
          : '';
      carbCtrl.text = base.carbohydrateG != null
          ? ((double.tryParse(base.carbohydrateG!) ?? 0) * s).toStringAsFixed(1)
          : '';
    }

    void applyFood(FoodProductDto food) {
      selectedFoodId.value = food.id;
      baseFood.value = food;
      servingCount.value = 1.0;
      caloriesCtrl.text = '${food.calories}';
      proteinCtrl.text = food.proteinG ?? '';
      fatCtrl.text = food.fatG ?? '';
      carbCtrl.text = food.carbohydrateG ?? '';
      if (memoCtrl.text.trim().isEmpty) {
        memoCtrl.text = food.name;
      }
    }

    Future<void> pickDateTime() async {
      final date = await showDatePicker(
        context: context,
        initialDate: eatenAt.value,
        firstDate: DateTime(2020),
        lastDate: DateTime.now().add(const Duration(days: 1)),
      );
      if (date == null || !context.mounted) return;
      final time = await showTimePicker(
        context: context,
        initialTime: TimeOfDay.fromDateTime(eatenAt.value),
      );
      if (time == null) return;
      eatenAt.value = DateTime(
        date.year,
        date.month,
        date.day,
        time.hour,
        time.minute,
      );
    }

    Future<void> submit() async {
      final kcal = int.tryParse(caloriesCtrl.text.trim());
      if (kcal == null) {
        error.value = l.commonCaloriesRequired;
        return;
      }
      loading.value = true;
      error.value = null;

      // 編集時は既存の写真パスを保持して渡す
      final photos =
          edit?.photos
              .where((p) => p.imagePath != null)
              .map(
                (p) => MealPhotoInput(
                  displayOrder: p.displayOrder,
                  imagePath: p.imagePath!,
                ),
              )
              .toList() ??
          [];

      final req = RecordMealRequest(
        eatenAt: eatenAt.value,
        mealType: mealType.value,
        calories: kcal,
        proteinG: double.tryParse(proteinCtrl.text.trim()),
        fatG: double.tryParse(fatCtrl.text.trim()),
        carbohydrateG: double.tryParse(carbCtrl.text.trim()),
        memo: memoCtrl.text.trim().isEmpty ? null : memoCtrl.text.trim(),
        photos: photos,
        foodProductId: selectedFoodId.value,
        servingCount: selectedFoodId.value != null ? servingCount.value : null,
      );
      try {
        final repo = ref.read(mealRepositoryProvider);
        if (isEdit) {
          await repo.update(edit.id, req);
        } else {
          await repo.record(req);
        }
        ref.invalidate(mealsProvider);
        if (context.mounted) context.pop();
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = '保存に失敗しました';
      } finally {
        if (context.mounted) loading.value = false;
      }
    }

    Future<void> deleteMeal() async {
      final confirmed = await showDialog<bool>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('削除の確認'),
          content: const Text('この記録を削除しますか？'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(ctx).pop(false),
              child: const Text('キャンセル'),
            ),
            TextButton(
              onPressed: () => Navigator.of(ctx).pop(true),
              child: Text('削除', style: TextStyle(color: t.accent)),
            ),
          ],
        ),
      );
      if (confirmed != true || !context.mounted) return;
      loading.value = true;
      try {
        await ref.read(mealRepositoryProvider).delete(edit!.id);
        ref.invalidate(mealsProvider);
        if (context.mounted) context.pop();
      } catch (_) {
        error.value = '削除に失敗しました';
        loading.value = false;
      }
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(isEdit ? '食事を編集' : '食事を記録'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(20, 16, 20, 32),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              FoodSearchSection(
                onSelect: applyFood,
                onNotFound: (barcode) async {
                  final food = await showFoodRegisterSheet(
                    context,
                    initialBarcode: barcode,
                  );
                  if (food != null && context.mounted) applyFood(food);
                },
              ),

              if (baseFood.value != null) ...[
                const SizedBox(height: 12),
                Row(
                  children: [
                    Text(
                      'サービング数',
                      style: TextStyle(fontSize: 13, color: t.muted),
                    ),
                    const Spacer(),
                    IconButton(
                      icon: const Icon(Icons.remove, size: 18),
                      onPressed: () => updateServing(
                        (servingCount.value - 0.5).clamp(0.5, 99),
                      ),
                    ),
                    SizedBox(
                      width: 40,
                      child: Text(
                        servingCount.value % 1 == 0
                            ? '${servingCount.value.toInt()}'
                            : servingCount.value.toStringAsFixed(1),
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.add, size: 18),
                      onPressed: () => updateServing(servingCount.value + 0.5),
                    ),
                  ],
                ),
              ],

              const SizedBox(height: 24),
              const Divider(),
              const SizedBox(height: 16),

              // 種類
              Wrap(
                spacing: 8,
                children: [
                  for (final m in _mealTypes)
                    ChoiceChip(
                      label: Text(m),
                      selected: mealType.value == m,
                      showCheckmark: false,
                      selectedColor: t.accentWeak,
                      onSelected: (_) => mealType.value = m,
                    ),
                ],
              ),

              const SizedBox(height: 16),

              // 日時
              GestureDetector(
                onTap: pickDateTime,
                child: Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 13,
                    vertical: 14,
                  ),
                  decoration: BoxDecoration(
                    border: Border.all(color: t.border),
                    borderRadius: BorderRadius.circular(13),
                    color: Theme.of(context).colorScheme.surface,
                  ),
                  child: Row(
                    children: [
                      Icon(Icons.calendar_today, size: 16, color: t.muted),
                      const SizedBox(width: 8),
                      Text(
                        _formatDateTime(eatenAt.value),
                        style: const TextStyle(fontSize: 16),
                      ),
                    ],
                  ),
                ),
              ),

              const SizedBox(height: 14),
              AppTextField(
                label: 'カロリー (kcal)',
                controller: caloriesCtrl,
                hint: '420',
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 14),
              Row(
                children: [
                  Expanded(
                    child: AppTextField(
                      label: 'P (g)',
                      controller: proteinCtrl,
                      hint: '28',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'F (g)',
                      controller: fatCtrl,
                      hint: '9',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'C (g)',
                      controller: carbCtrl,
                      hint: '58',
                      keyboardType: numeric,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 14),
              AppTextField(
                label: 'メモ · 任意',
                controller: memoCtrl,
                hint: 'オートミール・バナナ・卵',
              ),

              if (error.value != null) ...[
                const SizedBox(height: 12),
                Text(
                  error.value!,
                  style: TextStyle(color: t.accent, fontSize: 13),
                ),
              ],

              const SizedBox(height: 24),
              AppButton(
                label: isEdit ? '保存する' : '記録する',
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

String _formatDateTime(DateTime dt) {
  String pad(int n) => n.toString().padLeft(2, '0');
  return '${dt.year}/${pad(dt.month)}/${pad(dt.day)} ${pad(dt.hour)}:${pad(dt.minute)}';
}
