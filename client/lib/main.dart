import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import 'api_client.dart';
import 'app_state.dart';
import 'screens/home_screen.dart';
import 'screens/welcome_screen.dart';
import 'theme.dart';
import 'ui.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setSystemUIOverlayStyle(const SystemUiOverlayStyle(
    statusBarColor: Colors.transparent,
    statusBarIconBrightness: Brightness.light,
    systemNavigationBarColor: TM.surface,
    systemNavigationBarIconBrightness: Brightness.light,
  ));
  runApp(
    ChangeNotifierProvider(
      create: (_) => AppState(ApiClient())..init(),
      child: const TripMateApp(),
    ),
  );
}

class TripMateApp extends StatelessWidget {
  const TripMateApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'TripMate',
      debugShowCheckedModeBanner: false,
      theme: buildTheme(),
      home: const _Gate(),
    );
  }
}

class _Gate extends StatelessWidget {
  const _Gate();

  @override
  Widget build(BuildContext context) {
    final s = context.watch<AppState>();
    final Widget child = switch (s.status) {
      AuthStatus.booting => const _Splash(key: ValueKey('splash')),
      AuthStatus.bootFailed => Scaffold(
          key: const ValueKey('boot-failed'),
          body: AuroraBackground(
            child: Center(
              child: ErrorState(
                  message: s.bootError ?? 'Unable to start.', onRetry: s.init),
            ),
          ),
        ),
      AuthStatus.signedOut => const WelcomeScreen(key: ValueKey('welcome')),
      AuthStatus.signedIn => const HomeScreen(key: ValueKey('home')),
    };
    return AnimatedSwitcher(
        duration: const Duration(milliseconds: 420), child: child);
  }
}

class _Splash extends StatelessWidget {
  const _Splash({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: AuroraBackground(
        child: Center(
          child: Container(
            width: 104,
            height: 104,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: TM.brand,
              boxShadow: [BoxShadow(color: TM.violet.fade(.55), blurRadius: 48)],
            ),
            child: const Icon(Icons.explore_rounded, size: 58, color: Colors.white),
          ),
        ),
      ),
    );
  }
}
