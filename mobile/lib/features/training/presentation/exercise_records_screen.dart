import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../../core/theme/app_tokens.dart';
import '../../../core/theme/sketchy.dart';
import '../../../core/util/formatters.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/async_value_view.dart';
import '../../../core/widgets/section_title.dart';
import '../data/exercise_dtos.dart';
import '../data/training_repository.dart';
import '../../../l10n/app_localizations.dart';

class ExerciseRecordsScreen extends ConsumerStatefulWidget {
  const ExerciseRecordsScreen({super.key});

  @override
  ConsumerState<ExerciseRecordsScreen> createState() =>
      _ExerciseRecordsScreenState();
}

class _ExerciseRecordsScreenState extends ConsumerState<ExerciseRecordsScreen> {
  String? _selectedExerciseId;
  String _period = '1month';

  @override
  Widget build(BuildContext context) {
    final l = AppLocalizations.of(context)!;
    final periods = [
      ('1week', l.exerciseRecordsPeriod1week),
      ('1month', l.exerciseRecordsPeriod1month),
      ('3months', l.exerciseRecordsPeriod3months),
      ('halfyear', l.exerciseRecordsPeriodHalfYear),
      ('1year', l.exerciseRecordsPeriod1year),
    ];
    final exercisesAsync = ref.watch(exercisesProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(l.exerciseRecordsTitle),
        backgroundColor: Colors.transparent,
      ),
      body: SafeArea(
        child: AsyncValueView<List<ExerciseDto>>(
          value: exercisesAsync,
          onRetry: () => ref.invalidate(exercisesProvider),
          data: (exercises) {
            if (exercises.isEmpty) {
              return Center(child: Text(l.exerciseRecordsNoExercises));
            }
            _selectedExerciseId ??= exercises.first.id;
            return Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                _ExerciseSelector(
                  exercises: exercises,
                  selectedId: _selectedExerciseId!,
                  onChanged: (id) => setState(() => _selectedExerciseId = id),
                ),
                _PeriodSelector(
                  periods: periods,
                  selected: _period,
                  onChanged: (p) => setState(() => _period = p),
                ),
                Expanded(
                  child: _TimeseriesBody(
                    exerciseId: _selectedExerciseId!,
                    period: _period,
                  ),
                ),
              ],
            );
          },
        ),
      ),
    );
  }
}

class _ExerciseSelector extends StatelessWidget {
  const _ExerciseSelector({
    required this.exercises,
    required this.selectedId,
    required this.onChanged,
  });

  final List<ExerciseDto> exercises;
  final String selectedId;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
      child: RoughBox(
        padding: const EdgeInsets.symmetric(horizontal: 12),
        radius: BorderRadius.circular(12),
        child: DropdownButtonHideUnderline(
          child: DropdownButton<String>(
            value: selectedId,
            isExpanded: true,
            onChanged: (v) {
              if (v != null) onChanged(v);
            },
            items: exercises
                .map(
                  (ex) => DropdownMenuItem(
                    value: ex.id,
                    child: Text(
                      ex.name,
                      style: const TextStyle(fontWeight: FontWeight.w600),
                    ),
                  ),
                )
                .toList(),
          ),
        ),
      ),
    );
  }
}

class _PeriodSelector extends StatelessWidget {
  const _PeriodSelector({
    required this.periods,
    required this.selected,
    required this.onChanged,
  });

  final List<(String, String)> periods;
  final String selected;
  final ValueChanged<String> onChanged;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 4),
      child: Row(
        children: periods.map((pair) {
          final (value, label) = pair;
          final isSelected = value == selected;
          return Expanded(
            child: GestureDetector(
              onTap: () => onChanged(value),
              child: Container(
                margin: const EdgeInsets.symmetric(horizontal: 2),
                padding: const EdgeInsets.symmetric(vertical: 6),
                decoration: BoxDecoration(
                  color: isSelected ? t.accent : Colors.transparent,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  label,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 11,
                    fontWeight: isSelected
                        ? FontWeight.w700
                        : FontWeight.normal,
                    color: isSelected ? Colors.white : t.muted,
                  ),
                ),
              ),
            ),
          );
        }).toList(),
      ),
    );
  }
}

class _TimeseriesBody extends ConsumerWidget {
  const _TimeseriesBody({required this.exerciseId, required this.period});

  final String exerciseId;
  final String period;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final l = AppLocalizations.of(context)!;
    final async = ref.watch(
      exerciseBestSetTimeseriesProvider((exerciseId, period)),
    );

    return AsyncValueView<BestSetTimeseriesResponseDto>(
      value: async,
      onRetry: () => ref.invalidate(
        exerciseBestSetTimeseriesProvider((exerciseId, period)),
      ),
      data: (res) {
        final pts = res.dataPoints;
        if (pts.isEmpty) {
          return Center(child: Text(l.exerciseRecordsNoData));
        }
        return SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              SectionTitle(l.exerciseRecordsWeightChart),
              _WeightChart(dataPoints: pts),
              const SizedBox(height: 16),
              SectionTitle(l.exerciseRecordsRepsChart),
              _RepsChart(dataPoints: pts),
            ],
          ),
        );
      },
    );
  }
}

class _WeightChart extends StatelessWidget {
  const _WeightChart({required this.dataPoints});

  final List<BestSetTimeseriesDataPointDto> dataPoints;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final spots = <FlSpot>[
      for (var i = 0; i < dataPoints.length; i++)
        FlSpot(i.toDouble(), dataPoints[i].weightKg.toDouble()),
    ];

    return AppCard(
      child: SizedBox(
        height: 160,
        child: spots.length < 2
            ? Center(
                child: Text(
                  AppLocalizations.of(context)!.commonGraphHint,
                  style: TextStyle(color: t.muted, fontSize: 12),
                ),
              )
            : LineChart(
                LineChartData(
                  gridData: const FlGridData(show: false),
                  borderData: FlBorderData(show: false),
                  lineTouchData: LineTouchData(
                    touchTooltipData: LineTouchTooltipData(
                      getTooltipItems: (spots) => spots
                          .map(
                            (s) => LineTooltipItem(
                              '${s.y.toStringAsFixed(1)} kg\n${mdLabel(dataPoints[s.x.toInt()].performedAt)}',
                              const TextStyle(
                                color: Colors.white,
                                fontSize: 12,
                              ),
                            ),
                          )
                          .toList(),
                    ),
                  ),
                  titlesData: FlTitlesData(
                    leftTitles: AxisTitles(
                      sideTitles: SideTitles(
                        showTitles: true,
                        reservedSize: 40,
                        getTitlesWidget: (v, _) => Text(
                          v.toStringAsFixed(0),
                          style: TextStyle(fontSize: 10, color: t.muted),
                        ),
                      ),
                    ),
                    bottomTitles: AxisTitles(
                      sideTitles: SideTitles(
                        showTitles: true,
                        interval: (dataPoints.length / 4).ceilToDouble(),
                        getTitlesWidget: (v, _) {
                          final i = v.toInt();
                          if (i < 0 || i >= dataPoints.length) {
                            return const SizedBox.shrink();
                          }
                          return Text(
                            mdLabel(dataPoints[i].performedAt),
                            style: TextStyle(fontSize: 10, color: t.muted),
                          );
                        },
                      ),
                    ),
                    topTitles: const AxisTitles(
                      sideTitles: SideTitles(showTitles: false),
                    ),
                    rightTitles: const AxisTitles(
                      sideTitles: SideTitles(showTitles: false),
                    ),
                  ),
                  lineBarsData: [
                    LineChartBarData(
                      spots: spots,
                      isCurved: true,
                      color: t.accent,
                      barWidth: 2.5,
                      dotData: FlDotData(
                        show: dataPoints.length <= 12,
                        getDotPainter: (spot, xPercentage, bar, index) =>
                            FlDotCirclePainter(
                              radius: 3,
                              color: t.accent,
                              strokeWidth: 0,
                            ),
                      ),
                      belowBarData: BarAreaData(
                        show: true,
                        color: t.accent.withValues(alpha: 0.10),
                      ),
                    ),
                  ],
                ),
              ),
      ),
    );
  }
}

class _RepsChart extends StatelessWidget {
  const _RepsChart({required this.dataPoints});

  final List<BestSetTimeseriesDataPointDto> dataPoints;

  @override
  Widget build(BuildContext context) {
    final t = context.tokens;
    final groups = <BarChartGroupData>[
      for (var i = 0; i < dataPoints.length; i++)
        BarChartGroupData(
          x: i,
          barRods: [
            BarChartRodData(
              toY: dataPoints[i].reps.toDouble(),
              color: t.accent,
              width: dataPoints.length > 20 ? 4 : 12,
              borderRadius: const BorderRadius.vertical(
                top: Radius.circular(3),
              ),
            ),
          ],
        ),
    ];

    return AppCard(
      child: SizedBox(
        height: 160,
        child: BarChart(
          BarChartData(
            gridData: const FlGridData(show: false),
            borderData: FlBorderData(show: false),
            barTouchData: BarTouchData(
              touchTooltipData: BarTouchTooltipData(
                getTooltipItem: (group, groupIndex, rod, rodIndex) =>
                    BarTooltipItem(
                      '${rod.toY.toInt()} reps\n${mdLabel(dataPoints[group.x].performedAt)}',
                      const TextStyle(color: Colors.white, fontSize: 12),
                    ),
              ),
            ),
            titlesData: FlTitlesData(
              leftTitles: AxisTitles(
                sideTitles: SideTitles(
                  showTitles: true,
                  reservedSize: 32,
                  getTitlesWidget: (v, _) => Text(
                    v.toInt().toString(),
                    style: TextStyle(fontSize: 10, color: t.muted),
                  ),
                ),
              ),
              bottomTitles: AxisTitles(
                sideTitles: SideTitles(
                  showTitles: true,
                  interval: (dataPoints.length / 4).ceilToDouble(),
                  getTitlesWidget: (v, _) {
                    final i = v.toInt();
                    if (i < 0 || i >= dataPoints.length) {
                      return const SizedBox.shrink();
                    }
                    return Text(
                      mdLabel(dataPoints[i].performedAt),
                      style: TextStyle(fontSize: 10, color: t.muted),
                    );
                  },
                ),
              ),
              topTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false),
              ),
              rightTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false),
              ),
            ),
            barGroups: groups,
          ),
        ),
      ),
    );
  }
}
