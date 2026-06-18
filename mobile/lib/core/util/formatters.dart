/// 日付・時刻の軽量フォーマッタ。
/// API の日時は UTC で来るため、表示は常に端末ローカル(JST 等)へ変換する。
const _weekdaysJp = ['月', '火', '水', '木', '金', '土', '日'];

String weekdayJp(DateTime d) => _weekdaysJp[d.toLocal().weekday - 1];

/// 例: 6/17
String mdLabel(DateTime d) {
  final l = d.toLocal();
  return '${l.month}/${l.day}';
}

/// 例: 6月17日 (火)
String dateJpLong(DateTime d) {
  final l = d.toLocal();
  return '${l.month}月${l.day}日 (${weekdayJp(l)})';
}

/// 例: 6/17 (火)
String mdWeekday(DateTime d) {
  final l = d.toLocal();
  return '${l.month}/${l.day} (${weekdayJp(l)})';
}

/// 例: 07:30
String hhmm(DateTime d) {
  final l = d.toLocal();
  return '${l.hour.toString().padLeft(2, '0')}:${l.minute.toString().padLeft(2, '0')}';
}
