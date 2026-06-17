/// 日付・時刻の軽量フォーマッタ（locale 初期化不要な範囲）。
const _weekdaysJp = ['月', '火', '水', '木', '金', '土', '日'];

String weekdayJp(DateTime d) => _weekdaysJp[d.weekday - 1];

/// 例: 6/17
String mdLabel(DateTime d) => '${d.month}/${d.day}';

/// 例: 6月17日 (火)
String dateJpLong(DateTime d) => '${d.month}月${d.day}日 (${weekdayJp(d)})';

/// 例: 6/17 (火)
String mdWeekday(DateTime d) => '${d.month}/${d.day} (${weekdayJp(d)})';

/// 例: 07:30
String hhmm(DateTime d) =>
    '${d.hour.toString().padLeft(2, '0')}:${d.minute.toString().padLeft(2, '0')}';
