import 'dart:math' as math;

import 'package:flutter/material.dart';

import 'format.dart';
import 'models.dart';
import 'theme.dart';

// ───────────────────────── Navigation ─────────────────────────

Future<T?> pushPage<T>(BuildContext context, Widget page) {
  return Navigator.of(context).push<T>(PageRouteBuilder<T>(
    transitionDuration: const Duration(milliseconds: 380),
    reverseTransitionDuration: const Duration(milliseconds: 300),
    pageBuilder: (_, __, ___) => page,
    transitionsBuilder: (_, anim, __, child) {
      final curved = CurvedAnimation(parent: anim, curve: Curves.easeOutCubic);
      return FadeTransition(
        opacity: curved,
        child: SlideTransition(
          position: Tween<Offset>(begin: const Offset(0, .04), end: Offset.zero)
              .animate(curved),
          child: child,
        ),
      );
    },
  ));
}

// ───────────────────────── Feedback ─────────────────────────

void showSnack(BuildContext context, String message, {bool error = false}) {
  ScaffoldMessenger.of(context)
    ..hideCurrentSnackBar()
    ..showSnackBar(SnackBar(
      content: Row(children: [
        Icon(error ? Icons.error_outline_rounded : Icons.check_circle_outline,
            color: error ? TM.danger : TM.teal),
        const SizedBox(width: 12),
        Expanded(child: Text(message)),
      ]),
    ));
}

Future<bool> confirmDialog(
  BuildContext context, {
  required String title,
  required String message,
  String confirmLabel = 'Delete',
}) async {
  final r = await showDialog<bool>(
    context: context,
    builder: (ctx) => AlertDialog(
      title: Text(title),
      content: Text(message),
      actions: [
        TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('Cancel')),
        FilledButton(
          style: FilledButton.styleFrom(backgroundColor: TM.danger),
          onPressed: () => Navigator.pop(ctx, true),
          child: Text(confirmLabel),
        ),
      ],
    ),
  );
  return r ?? false;
}

// ───────────────────────── Backgrounds & surfaces ─────────────────────────

class AuroraBackground extends StatelessWidget {
  const AuroraBackground({super.key, required this.child});
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Stack(children: [
      const Positioned.fill(child: ColoredBox(color: TM.bg)),
      const Positioned(top: -140, left: -90, child: _Glow(TM.violet, 400)),
      const Positioned(top: 200, right: -160, child: _Glow(TM.teal, 340)),
      const Positioned(bottom: -160, left: -70, child: _Glow(TM.coral, 340)),
      Positioned.fill(child: child),
    ]);
  }
}

class _Glow extends StatelessWidget {
  const _Glow(this.color, this.size);
  final Color color;
  final double size;

  @override
  Widget build(BuildContext context) => IgnorePointer(
        child: Container(
          width: size,
          height: size,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            gradient: RadialGradient(
                colors: [color.fade(.30), color.fade(0)]),
          ),
        ),
      );
}

class GlassCard extends StatelessWidget {
  const GlassCard({
    super.key,
    required this.child,
    this.onTap,
    this.padding = const EdgeInsets.all(16),
    this.radius = 22,
    this.gradient,
  });

  final Widget child;
  final VoidCallback? onTap;
  final EdgeInsetsGeometry padding;
  final double radius;
  final Gradient? gradient;

  @override
  Widget build(BuildContext context) {
    final r = BorderRadius.circular(radius);
    return Material(
      color: Colors.transparent,
      child: Ink(
        decoration: BoxDecoration(
          color: gradient == null ? TM.surface.fade(.92) : null,
          gradient: gradient,
          borderRadius: r,
          border: Border.all(color: Colors.white.fade(.07)),
        ),
        child: InkWell(
          borderRadius: r,
          onTap: onTap,
          child: Padding(padding: padding, child: child),
        ),
      ),
    );
  }
}

class GradientButton extends StatelessWidget {
  const GradientButton({
    super.key,
    required this.label,
    required this.onPressed,
    this.icon,
    this.loading = false,
    this.gradient = TM.brand,
  });

  final String label;
  final VoidCallback? onPressed;
  final IconData? icon;
  final bool loading;
  final Gradient gradient;

  @override
  Widget build(BuildContext context) {
    final enabled = onPressed != null && !loading;
    return Opacity(
      opacity: onPressed == null && !loading ? .5 : 1,
      child: Material(
        color: Colors.transparent,
        child: Ink(
          decoration: BoxDecoration(
            gradient: gradient,
            borderRadius: BorderRadius.circular(18),
            boxShadow: [
              BoxShadow(
                  color: TM.violet.fade(.35),
                  blurRadius: 22,
                  offset: const Offset(0, 10)),
            ],
          ),
          child: InkWell(
            borderRadius: BorderRadius.circular(18),
            onTap: enabled ? onPressed : null,
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 17, horizontal: 24),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                mainAxisSize: MainAxisSize.min,
                children: [
                  if (loading)
                    const SizedBox(
                        width: 22,
                        height: 22,
                        child: CircularProgressIndicator(
                            strokeWidth: 2.6, color: Colors.white))
                  else ...[
                    if (icon != null) ...[
                      Icon(icon, color: Colors.white, size: 20),
                      const SizedBox(width: 10),
                    ],
                    Flexible(
                      child: Text(label,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.w700,
                              fontSize: 16,
                              letterSpacing: .2)),
                    ),
                  ],
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class FadeSlideIn extends StatelessWidget {
  const FadeSlideIn({super.key, required this.child, this.index = 0});
  final Widget child;
  final int index;

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0, end: 1),
      duration: Duration(milliseconds: 450 + math.min(index, 8) * 80),
      curve: Curves.easeOutCubic,
      builder: (context, v, child) => Opacity(
        opacity: v,
        child: Transform.translate(offset: Offset(0, 26 * (1 - v)), child: child),
      ),
      child: child,
    );
  }
}

class Skeleton extends StatefulWidget {
  const Skeleton({super.key, this.height = 150, this.radius = 24});
  final double height;
  final double radius;

  @override
  State<Skeleton> createState() => _SkeletonState();
}

class _SkeletonState extends State<Skeleton>
    with SingleTickerProviderStateMixin {
  late final AnimationController _c = AnimationController(
      vsync: this, duration: const Duration(milliseconds: 1100))
    ..repeat(reverse: true);

  @override
  void dispose() {
    _c.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => AnimatedBuilder(
        animation: _c,
        builder: (_, __) => Container(
          height: widget.height,
          decoration: BoxDecoration(
            color: Color.lerp(TM.surface, TM.surfaceHigh, _c.value),
            borderRadius: BorderRadius.circular(widget.radius),
          ),
        ),
      );
}

// ───────────────────────── States ─────────────────────────

class EmptyState extends StatelessWidget {
  const EmptyState({
    super.key,
    required this.icon,
    required this.title,
    required this.subtitle,
    this.action,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final Widget? action;

  @override
  Widget build(BuildContext context) => Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 36),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Container(
            width: 96,
            height: 96,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: TM.brand,
              boxShadow: [
                BoxShadow(color: TM.violet.fade(.4), blurRadius: 36),
              ],
            ),
            child: Icon(icon, size: 46, color: Colors.white),
          ),
          const SizedBox(height: 22),
          Text(title,
              textAlign: TextAlign.center,
              style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
          const SizedBox(height: 8),
          Text(subtitle,
              textAlign: TextAlign.center,
              style: const TextStyle(color: TM.textDim, height: 1.4)),
          if (action != null) ...[const SizedBox(height: 22), action!],
        ]),
      );
}

class ErrorState extends StatelessWidget {
  const ErrorState({super.key, required this.message, required this.onRetry});
  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) => EmptyState(
        icon: Icons.cloud_off_rounded,
        title: 'Connection hiccup',
        subtitle: message,
        action: OutlinedButton.icon(
          onPressed: onRetry,
          icon: const Icon(Icons.refresh_rounded),
          label: const Text('Try again'),
        ),
      );
}

class SectionTitle extends StatelessWidget {
  const SectionTitle(this.text, {super.key, this.trailing});
  final String text;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) => Padding(
        padding: const EdgeInsets.only(top: 8, bottom: 12),
        child: Row(children: [
          Expanded(
            child: Text(text,
                style: const TextStyle(
                    fontSize: 18, fontWeight: FontWeight.w800, letterSpacing: .2)),
          ),
          if (trailing != null) trailing!,
        ]),
      );
}

class Pill extends StatelessWidget {
  const Pill({super.key, required this.text, this.icon, this.color});
  final String text;
  final IconData? icon;
  final Color? color;

  @override
  Widget build(BuildContext context) {
    final c = color ?? Colors.white;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: Colors.black.fade(.22),
        borderRadius: BorderRadius.circular(30),
        border: Border.all(color: Colors.white.fade(.12)),
      ),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        if (icon != null) ...[
          Icon(icon, size: 14, color: c),
          const SizedBox(width: 5),
        ],
        Flexible(
          child: Text(text,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                  color: c, fontSize: 12, fontWeight: FontWeight.w600)),
        ),
      ]),
    );
  }
}

// ───────────────────────── Helpers ─────────────────────────

LinearGradient tripGradient(Trip t) {
  if (t.interests.isEmpty) return TM.brand;
  final c = t.interests.first.color;
  return LinearGradient(
    colors: [c, Color.lerp(c, const Color(0xFF2B1F7A), .68)!],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}

Future<void> pickBirthDate(
    BuildContext context, TextEditingController controller) async {
  final now = DateTime.now();
  final initial = DateTime.tryParse(controller.text) ??
      DateTime(now.year - 25, 1, 1);
  final d = await showDatePicker(
    context: context,
    initialDate: initial,
    firstDate: DateTime(1900),
    lastDate: now,
    helpText: 'Date of birth',
  );
  if (d != null) controller.text = apiDate(d);
}

String? requiredText(String? v, {int max = 100, String label = 'This field'}) {
  final t = (v ?? '').trim();
  if (t.isEmpty) return '$label is required';
  if (t.length > max) return '$label must be at most $max characters';
  return null;
}
