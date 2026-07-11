import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../../../l10n/app_localizations.dart';
import '../data/food_product_dtos.dart';
import '../data/food_product_repository.dart';

Future<FoodProductDto?> showFoodRegisterSheet(
  BuildContext context, {
  String? initialBarcode,
}) {
  return showModalBottomSheet<FoodProductDto>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => _FoodRegisterSheet(initialBarcode: initialBarcode),
  );
}

class _FoodRegisterSheet extends HookConsumerWidget {
  const _FoodRegisterSheet({this.initialBarcode});
  final String? initialBarcode;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final barcodeCtrl = useTextEditingController(text: initialBarcode ?? '');
    final nameCtrl = useTextEditingController();
    final caloriesCtrl = useTextEditingController();
    final proteinCtrl = useTextEditingController();
    final fatCtrl = useTextEditingController();
    final carbCtrl = useTextEditingController();
    final loading = useState(false);
    final error = useState<String?>(null);
    final t = context.tokens;
    const numeric = TextInputType.numberWithOptions(decimal: true);
    final l = AppLocalizations.of(context)!;

    Future<void> submit() async {
      final kcal = int.tryParse(caloriesCtrl.text.trim());
      if (nameCtrl.text.trim().isEmpty) {
        error.value = l.foodNameRequired;
        return;
      }
      if (kcal == null) {
        error.value = l.commonCaloriesRequired;
        return;
      }
      loading.value = true;
      error.value = null;
      try {
        final req = CreateFoodProductRequest(
          barcode: barcodeCtrl.text.trim().isEmpty
              ? null
              : barcodeCtrl.text.trim(),
          name: nameCtrl.text.trim(),
          calories: kcal,
          proteinG: proteinCtrl.text.trim().isEmpty
              ? null
              : proteinCtrl.text.trim(),
          fatG: fatCtrl.text.trim().isEmpty ? null : fatCtrl.text.trim(),
          carbohydrateG: carbCtrl.text.trim().isEmpty
              ? null
              : carbCtrl.text.trim(),
        );
        final id = await ref.read(foodProductRepositoryProvider).create(req);
        if (context.mounted) {
          Navigator.of(context).pop(
            FoodProductDto(
              id: id,
              barcode: req.barcode,
              name: req.name,
              calories: req.calories,
              proteinG: req.proteinG,
              fatG: req.fatG,
              carbohydrateG: req.carbohydrateG,
              registerSource: 'user',
            ),
          );
        }
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = l.foodRegisterFailed;
      } finally {
        if (context.mounted) loading.value = false;
      }
    }

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
                l.foodRegisterTitle,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(height: 16),
              if (initialBarcode != null) ...[
                AppTextField(
                  label: l.foodBarcodeLabel2,
                  controller: barcodeCtrl,
                  hint: '4901085615881',
                  keyboardType: TextInputType.number,
                ),
                const SizedBox(height: 14),
              ],
              AppTextField(
                label: l.foodNameField,
                controller: nameCtrl,
                hint: 'Snickers',
              ),
              const SizedBox(height: 14),
              AppTextField(
                label: l.foodCaloriesField,
                controller: caloriesCtrl,
                hint: '250',
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 14),
              Row(
                children: [
                  Expanded(
                    child: AppTextField(
                      label: 'P (g)',
                      controller: proteinCtrl,
                      hint: '4',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'F (g)',
                      controller: fatCtrl,
                      hint: '11',
                      keyboardType: numeric,
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(
                    child: AppTextField(
                      label: 'C (g)',
                      controller: carbCtrl,
                      hint: '33',
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
                label: l.foodRegisterBtn,
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
