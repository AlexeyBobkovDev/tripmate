import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../models.dart';
import '../theme.dart';
import '../ui.dart';
import 'auth_screens.dart';

class WelcomeScreen extends StatefulWidget {
  const WelcomeScreen({super.key});

  @override
  State<WelcomeScreen> createState() => _WelcomeScreenState();
}

class _WelcomeScreenState extends State<WelcomeScreen>
    with SingleTickerProviderStateMixin {
  late final AnimationController _c =
      AnimationController(vsync: this, duration: const Duration(seconds: 28))
        ..repeat();

  @override
  void dispose() {
    _c.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: AuroraBackground(
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Column(children: [
              const Spacer(flex: 2),
              _orbit(),
              const SizedBox(height: 28),
              FadeSlideIn(
                child: ShaderMask(
                  shaderCallback: (r) => TM.brand.createShader(r),
                  child: const Text('TripMate',
                      style: TextStyle(
                          fontSize: 46,
                          fontWeight: FontWeight.w900,
                          letterSpacing: -1,
                          color: Colors.white)),
                ),
              ),
              const SizedBox(height: 10),
              const FadeSlideIn(
                index: 2,
                child: Text(
                  'Tell us what you love.\nWe\'ll plan the route. You just travel.',
                  textAlign: TextAlign.center,
                  style: TextStyle(color: TM.textDim, fontSize: 16, height: 1.45),
                ),
              ),
              const Spacer(flex: 3),
              FadeSlideIn(
                index: 4,
                child: SizedBox(
                  width: double.infinity,
                  child: GradientButton(
                    label: 'Create account',
                    icon: Icons.rocket_launch_rounded,
                    onPressed: () => pushPage(context, const SignUpScreen()),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              FadeSlideIn(
                index: 5,
                child: SizedBox(
                  width: double.infinity,
                  child: OutlinedButton(
                    style: OutlinedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 17),
                      side: BorderSide(color: Colors.white.fade(.18)),
                      shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(18)),
                    ),
                    onPressed: () => pushPage(context, const SignInScreen()),
                    child: const Text('I already have an account',
                        style: TextStyle(
                            color: TM.text, fontWeight: FontWeight.w600)),
                  ),
                ),
              ),
              const SizedBox(height: 20),
            ]),
          ),
        ),
      ),
    );
  }

  Widget _orbit() {
    const items = Interest.values;
    return SizedBox(
      height: 280,
      child: AnimatedBuilder(
        animation: _c,
        builder: (context, _) => Stack(alignment: Alignment.center, children: [
          _ring(250, .07),
          _ring(170, .05),
          Container(
            width: 104,
            height: 104,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: TM.brand,
              boxShadow: [BoxShadow(color: TM.violet.fade(.55), blurRadius: 44)],
            ),
            child: const Icon(Icons.explore_rounded, size: 58, color: Colors.white),
          ),
          for (var i = 0; i < items.length; i++)
            Builder(builder: (context) {
              final a = _c.value * 2 * math.pi + i * math.pi / 2;
              return Transform.translate(
                offset: Offset(math.cos(a) * 125, math.sin(a) * 125),
                child: Container(
                  width: 52,
                  height: 52,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: TM.surface,
                    border: Border.all(color: items[i].color.fade(.6), width: 1.5),
                    boxShadow: [
                      BoxShadow(color: items[i].color.fade(.35), blurRadius: 18),
                    ],
                  ),
                  child: Icon(items[i].icon, color: items[i].color, size: 26),
                ),
              );
            }),
        ]),
      ),
    );
  }

  Widget _ring(double size, double opacity) => Container(
        width: size,
        height: size,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          border: Border.all(color: Colors.white.fade(opacity)),
        ),
      );
}
