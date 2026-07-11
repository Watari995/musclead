import 'package:decimal/decimal.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_badge.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../data/meal_dtos.dart';
import '../../../l10n/app_localizations.dart';
import '../data/meal_repository.dart';
import 'meal_template_sheet.dart';

class MealsScreen extends ConsumerWidget {
  const MealsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final meals = ref.watch(mealsProvider);
    return TabPage(
      title: l.mealTitle,
      subtitle: dateJpLong(DateTime.now()),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          IconButton(
            icon: Icon(Icons.bookmark_outline, color: context.tokens.accent),
            tooltip: l.mealTemplateFromRecord,
            onPressed: () async {
              final template = await showMealTemplateSheet(context);
              if (template != null && context.mounted) {
                // テンプレートから記録画面を開く（extra で渡す）
                context.push('/meals/new', extra: template);
              }
            },
          ),
          IconButton(
            icon: Icon(Icons.add, color: context.tokens.accent),
            onPressed: () => context.push('/meals/new'),
          ),
        ],
      ),
      onRefresh: () => ref.refresh(mealsProvider.future),
      children: [
        AsyncValueView<List<MealDto>>(
          value: meals,
          onRetry: () => ref.invalidate(mealsProvider),
          data: (list) => _MealsBody(meals: list),
        ),
      ],
    );
  }
}

class _MealsBody extends StatelessWidget {
  const _MealsBody({required this.meals});

  final List<MealDto> meals;

  bool _isToday(DateTime d) {
    final now = DateTime.now();
    return d.year == now.year && d.month == now.month && d.day == now.day;
  }

  @override
  Widget build(BuildContext context) {
    final today = meals.where((m) => _isToday(m.eatenAt)).toList();
    final kcal = today.fold<int>(0, (s, m) => s + m.calories);
    Decimal sum(Decimal? Function(MealDto) f) =>
        today.fold(Decimal.zero, (s, m) => s + (f(m) ?? Decimal.zero));
    final p = sum((m) => m.proteinG);
    final f = sum((m) => m.fatG);
    final c = sum((m) => m.carbohydrateG);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _SummaryCard(kcal: kcal, protein: p, fat: f, carb: c),
        SectionTitle(AppLocalizations.of(context)!.mealToday),
        if (today.isEmpty)
          _Empty(text: AppLocalizations.of(context)!.mealEmpty)
        else
          AppListBox(children: [for (final m in today) _MealRow(meal: m)]),
      ],
    );
  }
}

class _SummaryCard extends StatelessWidget {
  const _SummaryCard({
    required this.kcal,
    required this.protein,
    required this.fat,
    required this.carb,
  });

  final int kcal;
  final Decimal protein;
  final Decimal fat;
  final Decimal carb;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            AppLocalizations.of(context)!.mealTotal,
            style: TextStyle(fontSize: 12, color: t.muted),
          ),
          const SizedBox(height: 2),
          Row(
            crossAxisAlignment: CrossAxisAlignment.baseline,
            textBaseline: TextBaseline.alphabetic,
            children: [
              Text(
                '$kcal',
                style: const TextStyle(
                  fontSize: 30,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(width: 4),
              Text('kcal', style: TextStyle(fontSize: 14, color: t.muted)),
            ],
          ),
          const SizedBox(height: 14),
          _MacroBar(
            label: AppLocalizations.of(context)!.mealProtein,
            grams: protein,
            color: t.macroP,
          ),
          const SizedBox(height: 10),
          _MacroBar(
            label: AppLocalizations.of(context)!.mealFat,
            grams: fat,
            color: t.macroF,
          ),
          const SizedBox(height: 10),
          _MacroBar(
            label: AppLocalizations.of(context)!.mealCarb,
            grams: carb,
            color: t.macroC,
          ),
        ],
      ),
    );
  }
}

class _MacroBar extends StatelessWidget {
  const _MacroBar({
    required this.label,
    required this.grams,
    required this.color,
  });

  final String label;
  final Decimal grams;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Row(
      children: [
        SizedBox(
          width: 64,
          child: Text(label, style: const TextStyle(fontSize: 12.5)),
        ),
        const SizedBox(width: 10),
        Expanded(
          child: Container(
            height: 7,
            decoration: BoxDecoration(
              color: t.accentWeak.withValues(alpha: 0),
              borderRadius: BorderRadius.circular(4),
            ),
            child: ColoredBox(color: color.withValues(alpha: 0.18)),
          ),
        ),
        const SizedBox(width: 10),
        Text(
          '${grams.toStringAsFixed(0)}g',
          style: const TextStyle(fontSize: 12.5, fontWeight: FontWeight.w700),
        ),
      ],
    );
  }
}

class _MealRow extends StatelessWidget {
  const _MealRow({required this.meal});

  final MealDto meal;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return AppListRow(
      onTap: () => context.push('/meals/${meal.id}/edit', extra: meal),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  meal.memo?.isNotEmpty == true ? meal.memo! : meal.mealType,
                  style: const TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 2),
                Row(
                  children: [
                    Text(
                      hhmm(meal.eatenAt),
                      style: TextStyle(fontSize: 12, color: t.muted),
                    ),
                    const SizedBox(width: 8),
                    AppBadge(meal.mealType),
                  ],
                ),
              ],
            ),
          ),
          Text(
            '${meal.calories}',
            style: const TextStyle(fontWeight: FontWeight.w700),
          ),
          Text(' kcal', style: TextStyle(fontSize: 11, color: t.muted)),
        ],
      ),
    );
  }
}

class _Empty extends StatelessWidget {
  const _Empty({required this.text});
  final String text;
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 28),
      child: Center(
        child: Text(text, style: TextStyle(color: context.tokens.muted)),
      ),
    );
  }
}
