import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:intl/intl.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'calendar_dtos.dart';

class CalendarRepository {
  CalendarRepository(this._dio);

  final Dio _dio;

  Future<GetMonthlySummaryResponse> getMonthlySummary(int year, int month) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/calendar/monthly-summary',
          queryParameters: {'year': year, 'month': month},
        );
        return GetMonthlySummaryResponse.fromJson(res.data!);
      });

  Future<GetDailySummaryResponse> getDailySummary(DateTime date) =>
      guardApi(() async {
        final dateStr = DateFormat('yyyy-MM-dd').format(date);
        final res = await _dio.get<Map<String, dynamic>>(
          '/calendar/daily-summary',
          queryParameters: {'date': dateStr},
        );
        return GetDailySummaryResponse.fromJson(res.data!);
      });
}

final calendarRepositoryProvider = Provider<CalendarRepository>(
  (ref) => CalendarRepository(ref.watch(dioProvider)),
);

final monthlySummaryProvider =
    FutureProvider.family<GetMonthlySummaryResponse, (int, int)>((ref, args) {
      final (year, month) = args;
      return ref
          .watch(calendarRepositoryProvider)
          .getMonthlySummary(year, month);
    });

final dailySummaryProvider =
    FutureProvider.family<GetDailySummaryResponse, DateTime>(
      (ref, date) =>
          ref.watch(calendarRepositoryProvider).getDailySummary(date),
    );
