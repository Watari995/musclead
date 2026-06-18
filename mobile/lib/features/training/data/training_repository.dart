import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/data/api_guard.dart';
import '../../../core/providers/core_providers.dart';
import 'exercise_dtos.dart';
import 'routine_dtos.dart';
import 'training_dtos.dart';

class TrainingRepository {
  TrainingRepository(this._dio);

  final Dio _dio;

  // --- Trainings ---
  Future<List<TrainingDto>> listTrainings({int limit = 50, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/trainings',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListTrainingsResponse.fromJson(res.data!).trainings;
      });

  Future<void> recordTraining(RecordTrainingRequest request) =>
      guardApi(() => _dio.post<void>('/trainings', data: request.toJson()));

  Future<void> updateTraining(String id, RecordTrainingRequest request) =>
      guardApi(
        () => _dio.put<void>('/trainings/$id', data: request.toJson()),
      );

  Future<void> deleteTraining(String id) =>
      guardApi(() => _dio.delete<void>('/trainings/$id'));

  // --- Exercises ---
  Future<List<ExerciseDto>> listExercises({int limit = 100, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/exercises',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListExercisesResponse.fromJson(res.data!).exercises;
      });

  Future<void> createExercise(String name) => guardApi(
    () => _dio.post<void>(
      '/exercises',
      data: UpsertExerciseRequest(name: name).toJson(),
    ),
  );

  Future<void> deleteExercise(String id) =>
      guardApi(() => _dio.delete<void>('/exercises/$id'));

  /// 複数種目の自己ベスト。exercise_ids は repeat 形式で渡す。
  Future<List<BestSetDto>> bestSets(List<String> exerciseIds) =>
      guardApi(() async {
        if (exerciseIds.isEmpty) return <BestSetDto>[];
        final res = await _dio.get<Map<String, dynamic>>(
          '/exercises/best-sets',
          queryParameters: {'exercise_ids': exerciseIds},
          options: Options(listFormat: ListFormat.multi),
        );
        return ListBestSetsResponse.fromJson(res.data!).bestSets;
      });

  /// 種目マスタを渡した順に並び替える。
  Future<void> reorderExercises(List<String> exerciseIds) => guardApi(
    () => _dio.post<void>(
      '/exercises/reorder',
      data: {'exercise_ids': exerciseIds},
    ),
  );

  // --- Routines ---
  Future<List<RoutineDto>> listRoutines({int limit = 100, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/routines',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListRoutinesResponse.fromJson(res.data!).routines;
      });

  Future<void> createRoutine(UpsertRoutineRequest request) =>
      guardApi(() => _dio.post<void>('/routines', data: request.toJson()));

  Future<void> deleteRoutine(String id) =>
      guardApi(() => _dio.delete<void>('/routines/$id'));

  /// ルーティンを渡した順に並び替える。
  Future<void> reorderRoutines(List<String> routineIds) => guardApi(
    () =>
        _dio.post<void>('/routines/reorder', data: {'routine_ids': routineIds}),
  );
}

final trainingRepositoryProvider = Provider<TrainingRepository>(
  (ref) => TrainingRepository(ref.watch(dioProvider)),
);

final trainingsProvider = FutureProvider<List<TrainingDto>>(
  (ref) => ref.watch(trainingRepositoryProvider).listTrainings(),
);

final exercisesProvider = FutureProvider<List<ExerciseDto>>(
  (ref) => ref.watch(trainingRepositoryProvider).listExercises(),
);

final routinesProvider = FutureProvider<List<RoutineDto>>(
  (ref) => ref.watch(trainingRepositoryProvider).listRoutines(),
);
