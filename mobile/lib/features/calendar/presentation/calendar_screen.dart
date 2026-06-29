import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:intl/intl.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/tab_page.dart';
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

    final monthlySummary =
        ref.watch(monthlySummaryProvider((year.value, month.value)));
    final me = ref.watch(meProvider);
    final prefs = me.asData?.value.preferences;

    final trainingColor = _parseColor(prefs?.trainingColor ?? '#4A90E2');
    final mealColor = _parseColor(prefs?.mealColor ?? '#7ED321');
    final weightColor = _parseColor(prefs?.weightColor ?? '#FF6B6B');

    return TabPage(
      title: 'カレンダー',
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
          '$year年$month月',
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

  static const _weekdays = ['日', '月', '火', '水', '木', '金', '土'];

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final dayMap = {for (final d in days) d.date: d};
    final firstWeekday = DateTime(year, month, 1).weekday % 7; // Sun=0
    final daysInMonth = DateTime(year, month + 1, 0).day;
    final now = DateTime.now();
    final todayStr = DateFormat('yyyy-MM-dd').format(now);

    return Column(
      children: [
        Row(
          children: _weekdays.asMap().entries.map((e) {
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
            final dateStr = DateFormat('yyyy-MM-dd').format(
              DateTime(year, month, day),
            );
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
                                color:
                                    Theme.of(context).colorScheme.onSurface,
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
                        if (info?.hasTraining == true) _Dot(color: trainingColor),
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
    final t = context.tokens;
    final label = '${date.month}月${date.day}日';

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
            'エラー: $e',
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
                'この日の記録はありません',
                style: TextStyle(fontSize: 13, color: t.muted),
              );
            }
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (data.trainings.isNotEmpty) ...[
                  _SummarySection(
                    title: 'トレーニング',
                    children: data.trainings
                        .map(
                          (tr) => _SummaryRow(
                            left: '${tr.exerciseCount}種目 / ${tr.setCount}セット',
                            right: DateFormat('HH:mm').format(
                              tr.startedAt.toLocal(),
                            ),
                          ),
                        )
                        .toList(),
                  ),
                  const SizedBox(height: 8),
                ],
                if (data.meals.isNotEmpty) ...[
                  _SummarySection(
                    title: '食事',
                    children: data.meals
                        .map(
                          (m) => _SummaryRow(
                            left: m.mealType,
                            right: '${m.calories}kcal',
                          ),
                        )
                        .toList(),
                  ),
                  const SizedBox(height: 8),
                ],
                if (data.weights.isNotEmpty)
                  _SummarySection(
                    title: '体重',
                    children: data.weights
                        .map(
                          (w) => _SummaryRow(
                            left: '${w.weightKg}kg',
                            right: w.bodyFatPercentage != null
                                ? '体脂肪 ${w.bodyFatPercentage}%'
                                : '',
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
  const _SummaryRow({required this.left, required this.right});
  final String left;
  final String right;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Container(
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
            Text(right, style: TextStyle(fontSize: 12, color: t.muted)),
        ],
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
    final t = context.tokens;
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _LegendItem(color: trainingColor, label: 'トレーニング', t: t),
        const SizedBox(width: 12),
        _LegendItem(color: mealColor, label: '食事', t: t),
        const SizedBox(width: 12),
        _LegendItem(color: weightColor, label: '体重', t: t),
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
