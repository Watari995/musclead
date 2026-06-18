import 'package:decimal/decimal.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:musclead/features/meal/data/meal_dtos.dart';
import 'package:musclead/features/training/data/training_dtos.dart';
import 'package:musclead/features/weight/data/weight_dtos.dart';

void main() {
  test('MealDto: snake_case + 文字列小数を正しく解釈', () {
    final m = MealDto.fromJson({
      'id': 'm1',
      'eaten_at': '2026-06-17T07:30:00Z',
      'meal_type': '朝食',
      'calories': 420,
      'protein_g': '28.5',
      'fat_g': '9',
      'carbohydrate_g': '58',
      'memo': null,
      'photos': [
        {'display_order': 0, 'image_url': 'https://x/y.jpg'},
      ],
      'user_id': 'u1',
    });
    expect(m.id, 'm1');
    expect(m.mealType, '朝食');
    expect(m.calories, 420);
    expect(m.proteinG, Decimal.parse('28.5'));
    expect(m.photos.single.imageUrl, 'https://x/y.jpg');
  });

  test('WeightDto: 必須 Decimal + nullable Decimal', () {
    final w = WeightDto.fromJson({
      'id': 'w1',
      'weight_kg': '72.5',
      'measured_at': '2026-06-17T00:00:00Z',
      'body_fat_percentage': '18.2',
      'skeletal_muscle_kg': null,
    });
    expect(w.weightKg, Decimal.parse('72.5'));
    expect(w.bodyFatPercentage, Decimal.parse('18.2'));
    expect(w.skeletalMuscleKg, isNull);
  });

  test('RecordMealRequest.toJson: null 省略・数値送出・UTC 日時', () {
    final json = const RecordMealRequestBuilder().build().toJson();
    expect(json['meal_type'], '朝食');
    expect(json['calories'], 420);
    expect(json['protein_g'], 28.5);
    expect(json.containsKey('memo'), isFalse); // include_if_null:false
    expect(json['eaten_at'], '2026-06-17T07:30:00.000Z');
  });

  test('TrainingDto: 入れ子(種目>セット)と weight_kg を解釈', () {
    final t = TrainingDto.fromJson({
      'id': 't1',
      'started_at': '2026-06-16T10:00:00Z',
      'exercises': [
        {
          'id': 'e1',
          'exercise_id': 'x1',
          'display_order': 0,
          'sets': [
            {'id': 's1', 'set_number': 1, 'weight_kg': '60', 'reps': 10},
          ],
        },
      ],
    });
    final set = t.exercises.single.sets.single;
    expect(set.weightKg, Decimal.parse('60'));
    expect(set.reps, 10);
  });
}

class RecordMealRequestBuilder {
  const RecordMealRequestBuilder();
  RecordMealRequest build() => RecordMealRequest(
    eatenAt: DateTime.utc(2026, 6, 17, 7, 30),
    mealType: '朝食',
    calories: 420,
    proteinG: 28.5,
  );
}
