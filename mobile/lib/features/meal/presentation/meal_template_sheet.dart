import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/app_text_field.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../data/meal_template_dtos.dart';
import '../data/meal_template_repository.dart';

const _mealTypes = ['朝食', '昼食', '夕食', '間食'];

/// テンプレート一覧シート。[onSelect] でテンプレートを選択して閉じる。
Future<MealTemplateDto?> showMealTemplateSheet(BuildContext context) {
  return showModalBottomSheet<MealTemplateDto>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => const _MealTemplateSheet(),
  );
}

class _MealTemplateSheet extends ConsumerWidget {
  const _MealTemplateSheet();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final templates = ref.watch(mealTemplatesProvider);
    return DraggableScrollableSheet(
      initialChildSize: 0.6,
      minChildSize: 0.4,
      maxChildSize: 0.9,
      expand: false,
      builder: (context, scrollController) => Padding(
        padding: const EdgeInsets.fromLTRB(20, 4, 20, 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text(
                  '食事テンプレート',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
                ),
                TextButton(
                  onPressed: () => showCreateTemplateSheet(context, ref),
                  child: const Text('+ 新規'),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Expanded(
              child: AsyncValueView<List<MealTemplateDto>>(
                value: templates,
                onRetry: () => ref.invalidate(mealTemplatesProvider),
                data: (list) {
                  if (list.isEmpty) {
                    return Center(
                      child: Text(
                        'テンプレートはまだありません',
                        style: TextStyle(color: context.tokens.muted),
                      ),
                    );
                  }
                  return ListView.separated(
                    controller: scrollController,
                    itemCount: list.length,
                    separatorBuilder: (_, _) => const SizedBox(height: 8),
                    itemBuilder: (context, i) => _TemplateRow(
                      template: list[i],
                      onTap: () => Navigator.of(context).pop(list[i]),
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TemplateRow extends ConsumerWidget {
  const _TemplateRow({required this.template, required this.onTap});

  final MealTemplateDto template;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final t = context.tokens;
    return AppListRow(
      onTap: onTap,
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  template.name,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '${template.calories} kcal  P ${template.proteinG ?? "0"}g  F ${template.fatG ?? "0"}g  C ${template.carbohydrateG ?? "0"}g',
                  style: TextStyle(fontSize: 12, color: t.muted),
                ),
              ],
            ),
          ),
          IconButton(
            icon: Icon(Icons.delete_outline, size: 20, color: t.muted),
            onPressed: () async {
              final confirmed = await showDialog<bool>(
                context: context,
                builder: (_) => AlertDialog(
                  title: const Text('削除'),
                  content: Text('「${template.name}」を削除しますか?'),
                  actions: [
                    TextButton(
                      onPressed: () => Navigator.pop(context, false),
                      child: const Text('キャンセル'),
                    ),
                    TextButton(
                      onPressed: () => Navigator.pop(context, true),
                      child: Text('削除', style: TextStyle(color: t.accent)),
                    ),
                  ],
                ),
              );
              if (confirmed == true) {
                await ref
                    .read(mealTemplateRepositoryProvider)
                    .delete(template.id);
                ref.invalidate(mealTemplatesProvider);
              }
            },
          ),
        ],
      ),
    );
  }
}

/// テンプレート作成シート
Future<void> showCreateTemplateSheet(BuildContext context, WidgetRef outerRef) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => _CreateTemplateSheet(outerRef: outerRef),
  );
}

class _CreateTemplateSheet extends HookConsumerWidget {
  const _CreateTemplateSheet({required this.outerRef});

  final WidgetRef outerRef;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final type = useState('朝食');
    final name = useTextEditingController();
    final calories = useTextEditingController();
    final protein = useTextEditingController();
    final fat = useTextEditingController();
    final carb = useTextEditingController();
    final loading = useState(false);
    final error = useState<String?>(null);
    final t = context.tokens;

    Future<void> submit() async {
      final kcal = int.tryParse(calories.text.trim());
      if (name.text.trim().isEmpty) {
        error.value = '名前を入力してください';
        return;
      }
      if (kcal == null) {
        error.value = 'カロリーを入力してください';
        return;
      }
      loading.value = true;
      error.value = null;
      try {
        await ref
            .read(mealTemplateRepositoryProvider)
            .create(
              UpsertMealTemplateRequest(
                name: name.text.trim(),
                mealType: type.value,
                calories: kcal,
                proteinG: double.tryParse(protein.text.trim()),
                fatG: double.tryParse(fat.text.trim()),
                carbohydrateG: double.tryParse(carb.text.trim()),
              ),
            );
        outerRef.invalidate(mealTemplatesProvider);
        if (context.mounted) Navigator.of(context).pop();
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = '保存に失敗しました';
      } finally {
        if (context.mounted) loading.value = false;
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
              const Text(
                'テンプレートを作成',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 16),
              AppTextField(
                label: 'テンプレート名',
                controller: name,
                hint: 'プロテインシェイク',
              ),
              const SizedBox(height: 14),
              const SectionTitle('種類'),
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
              if (error.value != null) ...[
                const SizedBox(height: 12),
                Text(
                  error.value!,
                  style: TextStyle(color: t.accent, fontSize: 13),
                ),
              ],
              const SizedBox(height: 20),
              AppButton(
                label: loading.value ? '保存中…' : '保存する',
                loading: loading.value,
                onPressed: submit,
              ),
            ],
          ),
        ),
      ),
    );
  }
}
