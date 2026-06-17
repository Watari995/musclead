import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../../../core/widgets/tab_page.dart';
import '../data/weight_dtos.dart';
import '../data/weight_repository.dart';

class WeightsScreen extends ConsumerWidget {
  const WeightsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final weights = ref.watch(weightsProvider);
    return TabPage(
      title: '体重',
      onRefresh: () => ref.refresh(weightsProvider.future),
      children: [
        AsyncValueView<List<WeightDto>>(
          value: weights,
          onRetry: () => ref.invalidate(weightsProvider),
          data: (list) {
            if (list.isEmpty) {
              return const Padding(
                padding: EdgeInsets.symmetric(vertical: 40),
                child: Center(child: Text('まだ記録がありません')),
              );
            }
            // API は新しい順。グラフは古い順、リストは新しい順で使う。
            final asc = list.reversed.toList();
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                _ChartCard(
                  ascending: asc,
                  latest: list.first,
                  previous: list.length > 1 ? list[1] : null,
                ),
                const SectionTitle('最近の記録'),
                AppListBox(
                  children: [
                    for (final w in list.take(10)) _WeightRow(weight: w),
                  ],
                ),
              ],
            );
          },
        ),
      ],
    );
  }
}

class _ChartCard extends StatelessWidget {
  const _ChartCard({
    required this.ascending,
    required this.latest,
    this.previous,
  });

  final List<WeightDto> ascending;
  final WeightDto latest;
  final WeightDto? previous;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final spots = <FlSpot>[
      for (var i = 0; i < ascending.length; i++)
        FlSpot(i.toDouble(), ascending[i].weightKg.toDouble()),
    ];
    final delta = previous == null
        ? null
        : latest.weightKg - previous!.weightKg;

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            crossAxisAlignment: CrossAxisAlignment.baseline,
            textBaseline: TextBaseline.alphabetic,
            children: [
              Text(
                latest.weightKg.toString(),
                style: const TextStyle(
                  fontSize: 34,
                  fontWeight: FontWeight.w800,
                ),
              ),
              const SizedBox(width: 4),
              Text('kg', style: TextStyle(fontSize: 16, color: t.muted)),
              const Spacer(),
              if (delta != null)
                Text(
                  '${delta.toDouble() >= 0 ? '▲' : '▼'} ${delta.abs().toStringAsFixed(1)} kg',
                  style: TextStyle(fontWeight: FontWeight.w700, color: t.muted),
                ),
            ],
          ),
          if (latest.bodyFatPercentage != null)
            Padding(
              padding: const EdgeInsets.only(top: 2),
              child: Text(
                '体脂肪 ${latest.bodyFatPercentage}%',
                style: TextStyle(fontSize: 13, color: t.muted),
              ),
            ),
          const SizedBox(height: 14),
          SizedBox(
            height: 150,
            child: spots.length < 2
                ? Center(
                    child: Text(
                      '記録が増えるとグラフが表示されます',
                      style: TextStyle(color: t.muted, fontSize: 12),
                    ),
                  )
                : LineChart(
                    LineChartData(
                      gridData: const FlGridData(show: false),
                      titlesData: const FlTitlesData(show: false),
                      borderData: FlBorderData(show: false),
                      lineTouchData: const LineTouchData(enabled: false),
                      lineBarsData: [
                        LineChartBarData(
                          spots: spots,
                          isCurved: true,
                          color: t.accent,
                          barWidth: 2.5,
                          dotData: const FlDotData(show: false),
                          belowBarData: BarAreaData(
                            show: true,
                            color: t.accent.withValues(alpha: 0.10),
                          ),
                        ),
                      ],
                    ),
                  ),
          ),
        ],
      ),
    );
  }
}

class _WeightRow extends StatelessWidget {
  const _WeightRow({required this.weight});

  final WeightDto weight;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return AppListRow(
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  mdWeekday(weight.measuredAt),
                  style: const TextStyle(fontWeight: FontWeight.w600),
                ),
                if (weight.bodyFatPercentage != null)
                  Text(
                    '体脂肪 ${weight.bodyFatPercentage}%',
                    style: TextStyle(fontSize: 11, color: t.muted),
                  ),
              ],
            ),
          ),
          Text(
            '${weight.weightKg} kg',
            style: const TextStyle(fontWeight: FontWeight.w700),
          ),
        ],
      ),
    );
  }
}
