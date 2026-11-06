import 'package:flutter/material.dart';

/// TripMate design tokens.
class TM {
  TM._();

  static const bg = Color(0xFF0A0E1A);
  static const surface = Color(0xFF131A2C);
  static const surfaceHigh = Color(0xFF1B2440);
  static const violet = Color(0xFF7C6CFF);
  static const teal = Color(0xFF22D3C5);
  static const coral = Color(0xFFFF7A59);
  static const amber = Color(0xFFFFB84D);
  static const danger = Color(0xFFFF6B6B);
  static const text = Color(0xFFEAF0FF);
  static const textDim = Color(0xFF93A0C4);

  static const brand = LinearGradient(
    colors: [violet, teal],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
  static const sunset = LinearGradient(
    colors: [coral, violet],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}

extension ColorFade on Color {
  /// Same colour with the given opacity (0..1).
  Color fade(double opacity) =>
      withAlpha((opacity.clamp(0.0, 1.0) * 255).round());
}

OutlineInputBorder _border(Color c, [double w = 1]) => OutlineInputBorder(
      borderRadius: BorderRadius.circular(16),
      borderSide: BorderSide(color: c, width: w),
    );

ThemeData buildTheme() {
  final scheme = ColorScheme.fromSeed(
    seedColor: TM.violet,
    brightness: Brightness.dark,
  ).copyWith(
    primary: TM.violet,
    secondary: TM.teal,
    tertiary: TM.coral,
    surface: TM.surface,
    surfaceContainerHigh: TM.surfaceHigh,
    error: TM.danger,
  );
  final base = ThemeData(
    useMaterial3: true,
    colorScheme: scheme,
    brightness: Brightness.dark,
  );
  return base.copyWith(
    scaffoldBackgroundColor: TM.bg,
    textTheme: base.textTheme.apply(bodyColor: TM.text, displayColor: TM.text),
    appBarTheme: const AppBarTheme(
      backgroundColor: Colors.transparent,
      elevation: 0,
      scrolledUnderElevation: 0,
      centerTitle: false,
      foregroundColor: TM.text,
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: TM.surface,
      contentPadding: const EdgeInsets.symmetric(horizontal: 18, vertical: 18),
      border: _border(Colors.transparent),
      enabledBorder: _border(TM.surfaceHigh),
      focusedBorder: _border(TM.violet, 1.6),
      errorBorder: _border(TM.danger),
      focusedErrorBorder: _border(TM.danger, 1.6),
      labelStyle: const TextStyle(color: TM.textDim),
      hintStyle: const TextStyle(color: TM.textDim),
      prefixIconColor: TM.textDim,
    ),
    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: TM.surface,
      surfaceTintColor: Colors.transparent,
      indicatorColor: TM.violet.fade(.28),
      height: 70,
    ),
    snackBarTheme: SnackBarThemeData(
      behavior: SnackBarBehavior.floating,
      backgroundColor: TM.surfaceHigh,
      contentTextStyle: const TextStyle(color: TM.text),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: TM.violet,
      foregroundColor: Colors.white,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
    ),
  );
}
