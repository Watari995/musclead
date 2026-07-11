import 'package:decimal/decimal.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/number_stepper.dart';
import '../data/weight_dtos.dart';
import '../data/weight_repository.dart';
import '../../../l10n/app_localizations.dart';

/// 体重記録のモーダルボトムシートを開く。
Future<void> showWeightRecordSheet(BuildContext context) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => const _WeightRecordSheet(),
  );
}

class _WeightRecordSheet extends HookConsumerWidget {
  const _WeightRecordSheet();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final list = ref.read(weightsProvider).asData?.value ?? const <WeightDto>[];
    final last = list.isNotEmpty ? list.first : null;

    final weight = useTextEditingController(
      text: last?.weightKg.toString() ?? '',
    );
    final bodyFat = useTextEditingController(
      text: last?.bodyFatPercentage?.toString() ?? '',
    );
    final muscle = useTextEditingController(
      text: last?.skeletalMuscleKg?.toString() ?? '',
    );
    final loading = useState(false);
    final error = useState<String?>(null);

    Future<void> submit() async {
      final w = Decimal.tryParse(weight.text.trim());
      if (w == null) {
        error.value = l.weightInputRequired;
        return;
      }
      loading.value = true;
      error.value = null;
      try {
        await ref
            .read(weightRepositoryProvider)
            .upsert(
              UpsertWeightRequest(
                weightKg: w,
                measuredAt: DateTime.now(),
                bodyFatPercentage: Decimal.tryParse(bodyFat.text.trim()),
                skeletalMuscleKg: Decimal.tryParse(muscle.text.trim()),
              ),
            );
        ref.invalidate(weightsProvider);
        if (context.mounted) Navigator.of(context).pop();
      } on Failure catch (f) {
        error.value = f.message;
      } catch (_) {
        error.value = l.weightRecordFailed;
      } finally {
        if (context.mounted) loading.value = false;
      }
    }

    return Padding(
      padding: EdgeInsets.only(
        bottom: MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(20, 4, 20, 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Text(
                l.weightRecordTitle,
                style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
              ),
              const SizedBox(height: 18),
              NumberStepper(label: l.weightKgLabel, controller: weight, hint: '72.5'),
              const SizedBox(height: 14),
              NumberStepper(
                label: l.weightBodyFatLabel,
                controller: bodyFat,
                max: 100,
                hint: '18.2',
              ),
              const SizedBox(height: 14),
              NumberStepper(
                label: l.weightMuscleLabel,
                controller: muscle,
                hint: '33.1',
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
                label: l.commonRecord,
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
