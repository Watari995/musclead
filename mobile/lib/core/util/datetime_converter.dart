import 'package:json_annotation/json_annotation.dart';

/// リクエスト用 DateTime は常に UTC の RFC3339（末尾 Z）で送る。
/// バックエンドはタイムゾーン無しの ISO 文字列を受け付けないため。
class UtcDateTimeConverter implements JsonConverter<DateTime, String> {
  const UtcDateTimeConverter();

  @override
  DateTime fromJson(String json) => DateTime.parse(json).toLocal();

  @override
  String toJson(DateTime object) => object.toUtc().toIso8601String();
}

class NullableUtcDateTimeConverter
    implements JsonConverter<DateTime?, String?> {
  const NullableUtcDateTimeConverter();

  @override
  DateTime? fromJson(String? json) =>
      (json == null || json.isEmpty) ? null : DateTime.parse(json).toLocal();

  @override
  String? toJson(DateTime? object) => object?.toUtc().toIso8601String();
}
