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

  // --- Routines ---
  Future<List<RoutineDto>> listRoutines({int limit = 100, int offset = 0}) =>
      guardApi(() async {
        final res = await _dio.get<Map<String, dynamic>>(
          '/routines',
          queryParameters: {'limit': limit, 'offset': offset},
        );
        return ListRoutinesResponse.fromJson(res.data!).routines;
      });
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
