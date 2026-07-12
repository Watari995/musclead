import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:mobile_scanner/mobile_scanner.dart';

import '../../../core/error/failure.dart';
import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/widgets/app_button.dart';
import '../../../core/widgets/app_text_field.dart';
import '../../../l10n/app_localizations.dart';
import '../data/food_product_dtos.dart';
import '../data/food_product_repository.dart';

enum _SearchMode { name, barcode }

typedef FoodSelectCallback = void Function(FoodProductDto food);
typedef FoodNotFoundCallback = void Function(String? barcode);

class FoodSearchSection extends HookConsumerWidget {
  const FoodSearchSection({
    super.key,
    required this.onSelect,
    required this.onNotFound,
  });

  final FoodSelectCallback onSelect;
  final FoodNotFoundCallback onNotFound;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final mode = useState(_SearchMode.name);
    final nameCtrl = useTextEditingController();
    final barcodeCtrl = useTextEditingController();
    final results = useState<List<FoodProductDto>>([]);
    final isLoading = useState(false);
    final error = useState<String?>(null);
    final isFetched = useState(false);
    final debounce = useRef<Timer?>(null);
    final t = context.tokens;

    Future<void> searchByName(String query) async {
      if (query.isEmpty) {
        results.value = [];
        isFetched.value = false;
        return;
      }
      isLoading.value = true;
      error.value = null;
      try {
        results.value = await ref
            .read(foodProductRepositoryProvider)
            .searchByName(query);
        isFetched.value = true;
      } on Failure catch (f) {
        error.value = f.message;
      } finally {
        isLoading.value = false;
      }
    }

    Future<void> searchByBarcode(String barcode) async {
      if (barcode.isEmpty) return;
      isLoading.value = true;
      error.value = null;
      isFetched.value = false;
      try {
        results.value = await ref
            .read(foodProductRepositoryProvider)
            .searchByBarcode(barcode);
        isFetched.value = true;
      } on Failure catch (f) {
        error.value = f.message;
      } finally {
        isLoading.value = false;
      }
    }

    useEffect(() {
      void listener() {
        debounce.value?.cancel();
        debounce.value = Timer(const Duration(milliseconds: 400), () {
          searchByName(nameCtrl.text.trim());
        });
      }

      if (mode.value == _SearchMode.name) {
        nameCtrl.addListener(listener);
        return () {
          nameCtrl.removeListener(listener);
          debounce.value?.cancel();
        };
      }
      return null;
    }, [mode.value]);

    void switchMode(_SearchMode m) {
      mode.value = m;
      results.value = [];
      isFetched.value = false;
      error.value = null;
      nameCtrl.clear();
      barcodeCtrl.clear();
    }

    void handleSelect(FoodProductDto food) {
      results.value = [];
      isFetched.value = false;
      nameCtrl.clear();
      barcodeCtrl.clear();
      onSelect(food);
    }

    Future<void> openScanner() async {
      final barcode = await Navigator.of(context).push<String>(
        MaterialPageRoute(
          fullscreenDialog: true,
          builder: (_) => const _BarcodeScannerPage(),
        ),
      );
      if (barcode != null && context.mounted) {
        barcodeCtrl.text = barcode;
        await searchByBarcode(barcode);
      }
    }

    final showEmpty =
        isFetched.value &&
        !isLoading.value &&
        results.value.isEmpty &&
        (mode.value == _SearchMode.name
            ? nameCtrl.text.trim().isNotEmpty
            : barcodeCtrl.text.trim().isNotEmpty);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // mode tabs
        Row(
          children: [
            _ModeTab(
              label: l.foodSearchByName,
              selected: mode.value == _SearchMode.name,
              onTap: () => switchMode(_SearchMode.name),
            ),
            const SizedBox(width: 8),
            _ModeTab(
              label: l.foodSearchByBarcode,
              selected: mode.value == _SearchMode.barcode,
              onTap: () => switchMode(_SearchMode.barcode),
            ),
          ],
        ),
        const SizedBox(height: 12),

        if (mode.value == _SearchMode.name)
          AppTextField(
            label: l.foodNameLabel,
            controller: nameCtrl,
            hint: l.foodNameHint,
          )
        else
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Expanded(
                child: AppTextField(
                  label: l.foodBarcodeLabel,
                  controller: barcodeCtrl,
                  hint: '4901085615881',
                  keyboardType: TextInputType.number,
                ),
              ),
              const SizedBox(width: 8),
              SizedBox(
                height: 48,
                child: IconButton(
                  icon: Icon(Icons.qr_code_scanner, color: t.accent),
                  tooltip: l.foodCameraScan,
                  onPressed: openScanner,
                ),
              ),
              const SizedBox(width: 4),
              AppButton(
                label: l.foodSearch,
                expand: false,
                onPressed: () => searchByBarcode(barcodeCtrl.text.trim()),
              ),
            ],
          ),

        if (isLoading.value) ...[
          const SizedBox(height: 10),
          const LinearProgressIndicator(),
        ],

        if (error.value != null) ...[
          const SizedBox(height: 8),
          Text(error.value!, style: TextStyle(color: t.accent, fontSize: 13)),
        ],

        if (results.value.isNotEmpty) ...[
          const SizedBox(height: 10),
          RoughBox(
            radius: BorderRadius.circular(13),
            clipBehavior: Clip.antiAlias,
            child: Column(
              children: [
                for (int i = 0; i < results.value.length; i++) ...[
                  if (i > 0) Divider(height: 1, color: t.hairline),
                  _ResultRow(
                    food: results.value[i],
                    onTap: () => handleSelect(results.value[i]),
                  ),
                ],
              ],
            ),
          ),
        ],

        if (showEmpty) ...[
          const SizedBox(height: 10),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                l.foodNotFound,
                style: TextStyle(color: t.muted, fontSize: 13),
              ),
              GestureDetector(
                onTap: () => onNotFound(
                  mode.value == _SearchMode.barcode
                      ? barcodeCtrl.text.trim()
                      : null,
                ),
                child: Text(
                  l.foodRegisterLink,
                  style: TextStyle(
                    color: t.accent,
                    fontSize: 13,
                    decoration: TextDecoration.underline,
                  ),
                ),
              ),
            ],
          ),
        ],
      ],
    );
  }
}

class _ModeTab extends StatelessWidget {
  const _ModeTab({
    required this.label,
    required this.selected,
    required this.onTap,
  });

  final String label;
  final bool selected;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 150),
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
        decoration: BoxDecoration(
          color: selected ? t.accent : Colors.transparent,
          border: Border.all(color: selected ? t.accent : t.border),
          borderRadius: BorderRadius.circular(20),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontSize: 12.5,
            fontWeight: FontWeight.w600,
            color: selected ? context.colors.onPrimary : t.muted,
          ),
        ),
      ),
    );
  }
}

class _ResultRow extends StatelessWidget {
  const _ResultRow({required this.food, required this.onTap});

  final FoodProductDto food;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(13),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    food.name,
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 14,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    [
                      '${food.calories} kcal',
                      if (food.proteinG != null) 'P ${food.proteinG}g',
                      if (food.fatG != null) 'F ${food.fatG}g',
                      if (food.carbohydrateG != null)
                        'C ${food.carbohydrateG}g',
                    ].join(' · '),
                    style: TextStyle(color: t.muted, fontSize: 12),
                  ),
                ],
              ),
            ),
            Icon(Icons.chevron_right, color: t.subtle, size: 18),
          ],
        ),
      ),
    );
  }
}

class _BarcodeScannerPage extends StatelessWidget {
  const _BarcodeScannerPage();

  @override
  Widget build(BuildContext context) {
    bool detected = false;
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(
        backgroundColor: Colors.black,
        foregroundColor: Colors.white,
        title: Text(AppLocalizations.of(context)!.foodScanTitle),
      ),
      body: Stack(
        children: [
          MobileScanner(
            onDetect: (capture) {
              if (detected) return;
              final barcode = capture.barcodes.firstOrNull?.rawValue;
              if (barcode != null && context.mounted) {
                detected = true;
                Navigator.of(context).pop(barcode);
              }
            },
          ),
          Center(
            child: Container(
              width: 240,
              height: 160,
              decoration: BoxDecoration(
                border: Border.all(color: Colors.white, width: 2),
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
          Positioned(
            bottom: 40,
            left: 0,
            right: 0,
            child: Text(
              AppLocalizations.of(context)!.foodScanHint,
              textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.white, fontSize: 14),
            ),
          ),
        ],
      ),
    );
  }
}
