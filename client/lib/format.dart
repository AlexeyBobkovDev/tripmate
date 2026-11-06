import 'package:intl/intl.dart';

String fmtDate(DateTime d) => DateFormat('d MMM yyyy').format(d.toUtc());
String fmtDayMonth(DateTime d) => DateFormat('d MMM').format(d.toUtc());

String fmtRange(DateTime a, DateTime b) {
  final x = a.toUtc(), y = b.toUtc();
  if (x.year == y.year && x.month == y.month && x.day == y.day) {
    return fmtDate(a);
  }
  return '${fmtDayMonth(a)} – ${fmtDate(b)}';
}

String fmtTime(DateTime d) => DateFormat('HH:mm').format(d.toUtc());

/// Wire format for date-only fields, e.g. 2026-11-14.
String apiDate(DateTime d) => DateFormat('yyyy-MM-dd').format(d);

String fmtDistance(int metres) => metres >= 1000
    ? '${(metres / 1000).toStringAsFixed(metres >= 10000 ? 0 : 1)} km'
    : '$metres m';

String fmtDuration(int seconds) {
  if (seconds < 60) return '$seconds s';
  final m = (seconds / 60).round();
  if (m < 60) return '$m min';
  final h = m ~/ 60;
  final r = m % 60;
  return r == 0 ? '$h h' : '$h h $r min';
}

String greeting() {
  final h = DateTime.now().hour;
  if (h < 5) return 'Good night';
  if (h < 12) return 'Good morning';
  if (h < 18) return 'Good afternoon';
  return 'Good evening';
}
