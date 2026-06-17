import 'package:decimal/decimal.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:musclead/core/util/decimal_converter.dart';

void main() {
  group('DecimalStringConverter', () {
    const c = DecimalStringConverter();
    test('文字列 → Decimal（精度保持）', () {
      expect(c.fromJson('72.5'), Decimal.parse('72.5'));
      expect(c.fromJson('0.1'), Decimal.parse('0.1'));
    });
    test('Decimal → 文字列', () {
      expect(c.toJson(Decimal.parse('72.5')), '72.5');
      expect(c.toJson(Decimal.fromInt(60)), '60');
    });
  });

  group('NullableDecimalStringConverter', () {
    const c = NullableDecimalStringConverter();
    test('null / 空文字 → null', () {
      expect(c.fromJson(null), isNull);
      expect(c.fromJson(''), isNull);
      expect(c.toJson(null), isNull);
    });
    test('値あり', () {
      expect(c.fromJson('18.2'), Decimal.parse('18.2'));
    });
  });
}
