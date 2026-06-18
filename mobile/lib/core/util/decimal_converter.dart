import 'package:decimal/decimal.dart';
import 'package:json_annotation/json_annotation.dart';

/// バックエンドは精度保持のため小数を文字列で返す（例: weight_kg:"72.5"）。
/// DTO 側では [Decimal] として扱うための変換。
class DecimalStringConverter implements JsonConverter<Decimal, String> {
  const DecimalStringConverter();

  @override
  Decimal fromJson(String json) => Decimal.parse(json);

  @override
  String toJson(Decimal object) => object.toString();
}

class NullableDecimalStringConverter
    implements JsonConverter<Decimal?, String?> {
  const NullableDecimalStringConverter();

  @override
  Decimal? fromJson(String? json) =>
      (json == null || json.isEmpty) ? null : Decimal.parse(json);

  @override
  String? toJson(Decimal? object) => object?.toString();
}
