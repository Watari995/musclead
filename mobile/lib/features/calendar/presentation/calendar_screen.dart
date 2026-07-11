import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:intl/intl.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/tab_page.dart';
import '../../../l10n/app_localizations.dart';
import '../../notifications/data/notification_repository.dart';
import '../../user/data/user_repository.dart';
import '../data/calendar_dtos.dart';
import '../data/calendar_repository.dart';

class CalendarScreen extends HookConsumerWidget {
  const CalendarScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final now = DateTime.now();
    final year = useState(now.year);
    final month = useState(now.month);
    final selectedDate = useState<DateTime?>(
      DateTime(now.year, now.month, now.day),
    );

    final monthlySummary = ref.watch(
      monthlySummaryProvider((year.value, month.value)),
    );
    final me = ref.watch(meProvider);
    final prefs = me.asData?.value.preferences;

    final trainingColor = _parseColor(prefs?.trainingColor ?? '#4A90E2');
    final mealColor = _parseColor(prefs?.mealColor ?? '#7ED321');
    final weightColor = _parseColor(prefs?.weightColor ?? '#FF6B6B');

    final unreadCount =
        ref.watch(notificationsProvider).asData?.value.unreadCount ?? 0;

    final l = AppLocalizations.of(context)!;
    return TabPage(
      title: l.calendarTitle,
      trailing: Stack(
        clipBehavior: Clip.none,
        children: [
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () => context.push('/notifications'),
          ),
          if (unreadCount > 0)
            Positioned(
              right: 6,
              top: 6,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
                decoration: BoxDecoration(
                  color: Colors.red,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  unreadCount > 99 ? '99+' : '$unreadCount',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 10,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
        ],
      ),
      onRefresh: () async {
        ref.invalidate(monthlySummaryProvider((year.value, month.value)));
        if (selectedDate.value != null) {
          ref.invalidate(dailySummaryProvider(selectedDate.value!));
        }
      },
      children: [
        AsyncValueView<GetMonthlySummaryResponse>(
          value: monthlySummary,
          onRetry: () =>
              ref.invalidate(monthlySummaryProvider((year.value, month.value))),
          data: (data) => Column(
            children: [
              _MonthHeader(
                year: year.value,
                month: month.value,
                onPrev: () {
                  if (month.value == 1) {
                    year.value--;
                    month.value = 12;
                  } else {
                    month.value--;
                  }
                  selectedDate.value = null;
                },
                onNext: () {
                  if (month.value == 12) {
                    year.value++;
                    month.value = 1;
                  } else {
                    month.value++;
                  }
                  selectedDate.value = null;
                },
              ),
              const SizedBox(height: 8),
              _CalendarGrid(
                year: year.value,
                month: month.value,
                days: data.days,
                selectedDate: selectedDate.value,
                trainingColor: trainingColor,
                mealColor: mealColor,
                weightColor: weightColor,
                onSelectDate: (date) => selectedDate.value = date,
              ),
              if (selectedDate.value != null) ...[
                const SizedBox(height: 16),
                _DailySummarySection(date: selectedDate.value!),
              ],
              const SizedBox(height: 12),
              _Legend(
                trainingColor: trainingColor,
                mealColor: mealColor,
                weightColor: weightColor,
              ),
            ],
          ),
        ),
      ],
    );
  }

  static Color _parseColor(String hex) {
    try {
      final h = hex.replaceAll('#', '');
      return Color(int.parse('FF$h', radix: 16));
    } catch (_) {
      return const Color(0xFF4A90E2);
    }
  }
}

class _MonthHeader extends StatelessWidget {
  const _MonthHeader({
    required this.year,
    required this.month,
    required this.onPrev,
    required this.onNext,
  });

  final int year;
  final int month;
  final VoidCallback onPrev;
  final VoidCallback onNext;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        IconButton(icon: const Icon(Icons.chevron_left), onPressed: onPrev),
        Text(
          AppLocalizations.of(context)!.calendarYearMonth(year, month),
          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
        ),
        IconButton(icon: const Icon(Icons.chevron_right), onPressed: onNext),
      ],
    );
  }
}

class _CalendarGrid extends StatelessWidget {
  const _CalendarGrid({
    required this.year,
    required this.month,
    required this.days,
    required this.selectedDate,
    required this.trainingColor,
    required this.mealColor,
    required this.weightColor,
    required this.onSelectDate,
  });

  final int year;
  final int month;
  final List<MonthlySummaryDayDto> days;
  final DateTime? selectedDate;
  final Color trainingColor;
  final Color mealColor;
  final Color weightColor;
  final ValueChanged<DateTime> onSelectDate;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final weekdays = [l.calendarSun, l.calendarMon, l.calendarTue, l.calendarWed, l.calendarThu, l.calendarFri, l.calendarSat];
    final t = context.tokens;
    final dayMap = {for (final d in days) d.date: d};
    final firstWeekday = DateTime(year, month, 1).weekday % 7; // Sun=0
    final daysInMonth = DateTime(year, month + 1, 0).day;
    final now = DateTime.now();
    final todayStr = DateFormat('yyyy-MM-dd').format(now);

    return Column(
      children: [
        Row(
          children: weekdays.asMap().entries.map((e) {
            final i = e.key;
            final d = e.value;
            final color = i == 0
                ? Theme.of(context).colorScheme.error
                : i == 6
                ? Colors.blue
                : t.muted;
            return Expanded(
              child: Center(
                child: Text(
                  d,
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: FontWeight.w600,
                    color: color,
                  ),
                ),
              ),
            );
          }).toList(),
        ),
        const SizedBox(height: 4),
        GridView.builder(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: 7,
            childAspectRatio: 0.85,
          ),
          itemCount: firstWeekday + daysInMonth,
          itemBuilder: (context, index) {
            if (index < firstWeekday) return const SizedBox.shrink();
            final day = index - firstWeekday + 1;
            final dateStr = DateFormat(
              'yyyy-MM-dd',
            ).format(DateTime(year, month, day));
            final info = dayMap[dateStr];
            final isToday = dateStr == todayStr;
            final isSelected =
                selectedDate != null &&
                DateFormat('yyyy-MM-dd').format(selectedDate!) == dateStr;
            final dow = (firstWeekday + day - 1) % 7;

            return GestureDetector(
              onTap: () => onSelectDate(DateTime(year, month, day)),
              child: Container(
                margin: const EdgeInsets.all(1),
                decoration: BoxDecoration(
                  color: isSelected
                      ? Theme.of(context).colorScheme.onSurface
                      : Colors.transparent,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Container(
                      width: 28,
                      height: 28,
                      decoration: isToday && !isSelected
                          ? BoxDecoration(
                              shape: BoxShape.circle,
                              border: Border.all(
                                color: Theme.of(context).colorScheme.onSurface,
                                width: 1.5,
                              ),
                            )
                          : null,
                      child: Center(
                        child: Text(
                          '$day',
                          style: TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.w500,
                            color: isSelected
                                ? Theme.of(context).colorScheme.surface
                                : dow == 0
                                ? Theme.of(context).colorScheme.error
                                : dow == 6
                                ? Colors.blue
                                : null,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(height: 2),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        if (info?.hasTraining == true)
                          _Dot(color: trainingColor),
                        if (info?.hasMeal == true) _Dot(color: mealColor),
                        if (info?.hasWeight == true) _Dot(color: weightColor),
                      ],
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      ],
    );
  }
}

class _Dot extends StatelessWidget {
  const _Dot({required this.color});
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 5,
      height: 5,
      margin: const EdgeInsets.symmetric(horizontal: 1),
      decoration: BoxDecoration(color: color, shape: BoxShape.circle),
    );
  }
}

class _DailySummarySection extends ConsumerWidget {
  const _DailySummarySection({required this.date});
  final DateTime date;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final summary = ref.watch(dailySummaryProvider(date));
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    final label = l.calendarDayLabel(date.month, date.day);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Divider(color: t.hairline),
        const SizedBox(height: 8),
        Text(
          label,
          style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w700),
        ),
        const SizedBox(height: 8),
        summary.when(
          loading: () =>
              const Center(child: CircularProgressIndicator(strokeWidth: 2.4)),
          error: (e, _) => Text(
            l.calendarError(e),
            style: TextStyle(
              color: Theme.of(context).colorScheme.error,
              fontSize: 12,
            ),
          ),
          data: (data) {
            if (data.trainings.isEmpty &&
                data.meals.isEmpty &&
                data.weights.isEmpty) {
              return Text(
                l.calendarNoRecords,
                style: TextStyle(fontSize: 13, color: t.muted),
              );
            }
            final totalMealCalories = data.totalCalories;
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (data.trainings.isNotEmpty) ...[
                  _SummarySection(
                    title: l.calendarTraining,
                    children: data.trainings
                        .map(
                          (tr) => _SummaryRow(
                            left: l.calendarExerciseSets(tr.exerciseCount, tr.setCount),
                            right: DateFormat(
                              'HH:mm',
                            ).format(tr.startedAt.toLocal()),
                            onTap: () => context.go('/trainings'),
                          ),
                        )
                        .toList(),
                  ),
                  const SizedBox(height: 8),
                ],
                if (data.meals.isNotEmpty) ...[
                  _SummarySection(
                    title: l.calendarMeal,
                    children: [
                      _SummaryRow(
                        left: l.calendarTotal,
                        right: '${totalMealCalories}kcal',
                        onTap: () => context.go('/meals'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                ],
                if (data.weights.isNotEmpty)
                  _SummarySection(
                    title: l.calendarWeight,
                    children: data.weights
                        .map(
                          (w) => _SummaryRow(
                            left: '${w.weightKg}kg',
                            right: w.bodyFatPercentage != null
                                ? l.weightBodyFat(w.bodyFatPercentage!)
                                : '',
                            onTap: () => context.go('/weights'),
                          ),
                        )
                        .toList(),
                  ),
              ],
            );
          },
        ),
      ],
    );
  }
}

class _SummarySection extends StatelessWidget {
  const _SummarySection({required this.title, required this.children});
  final String title;
  final List<Widget> children;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: TextStyle(
            fontSize: 11,
            fontWeight: FontWeight.w600,
            color: t.muted,
          ),
        ),
        const SizedBox(height: 4),
        ...children,
      ],
    );
  }
}

class _SummaryRow extends StatelessWidget {
  const _SummaryRow({required this.left, required this.right, this.onTap});
  final String left;
  final String right;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: 4),
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: BoxDecoration(
          color: context.colors.surface,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: t.border),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(left, style: const TextStyle(fontSize: 13)),
            if (right.isNotEmpty)
              Row(
                children: [
                  Text(right, style: TextStyle(fontSize: 12, color: t.muted)),
                  if (onTap != null) ...[
                    const SizedBox(width: 4),
                    Icon(Icons.chevron_right, size: 14, color: t.muted),
                  ],
                ],
              ),
          ],
        ),
      ),
    );
  }
}

class _Legend extends StatelessWidget {
  const _Legend({
    required this.trainingColor,
    required this.mealColor,
    required this.weightColor,
  });
  final Color trainingColor;
  final Color mealColor;
  final Color weightColor;

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final t = context.tokens;
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _LegendItem(color: trainingColor, label: l.calendarTraining, t: t),
        const SizedBox(width: 12),
        _LegendItem(color: mealColor, label: l.calendarMeal, t: t),
        const SizedBox(width: 12),
        _LegendItem(color: weightColor, label: l.calendarWeight, t: t),
      ],
    );
  }
}

class _LegendItem extends StatelessWidget {
  const _LegendItem({
    required this.color,
    required this.label,
    required this.t,
  });
  final Color color;
  final String label;
  final AppTokens t;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(color: color, shape: BoxShape.circle),
        ),
        const SizedBox(width: 4),
        Text(label, style: TextStyle(fontSize: 11, color: t.muted)),
      ],
    );
  }
}
