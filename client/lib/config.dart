/// Central configuration.
///
/// The Android emulator reaches your computer's localhost via 10.0.2.2.
/// Override at build time:
///   flutter run --dart-define=API_BASE_URL=http://192.168.1.50:8080/api/v1
class AppConfig {
  AppConfig._();

  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://10.0.2.2:8080/api/v1',
  );

  /// Dark basemap tiles (no API key needed). Check the provider's terms
  /// before shipping to production and switch to your own tile provider.
  static const String mapTiles =
      'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}.png';
  static const List<String> mapSubdomains = ['a', 'b', 'c', 'd'];
  static const String userAgentPackage = 'app.tripmate.client';
}
