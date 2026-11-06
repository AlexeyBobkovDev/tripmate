import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../api_client.dart';
import '../app_state.dart';
import '../theme.dart';
import '../ui.dart';

class _AuthScaffold extends StatelessWidget {
  const _AuthScaffold(
      {required this.title, required this.subtitle, required this.child});
  final String title;
  final String subtitle;
  final Widget child;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: AuroraBackground(
        child: SafeArea(
          child: Column(children: [
            Align(
              alignment: Alignment.centerLeft,
              child: IconButton(
                icon: const Icon(Icons.arrow_back_rounded),
                onPressed: () => Navigator.of(context).maybePop(),
              ),
            ),
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(title,
                        style: const TextStyle(
                            fontSize: 32,
                            fontWeight: FontWeight.w900,
                            letterSpacing: -.5)),
                    const SizedBox(height: 8),
                    Text(subtitle,
                        style: const TextStyle(color: TM.textDim, height: 1.4)),
                    const SizedBox(height: 28),
                    child,
                  ],
                ),
              ),
            ),
          ]),
        ),
      ),
    );
  }
}

// ───────────────────────── Sign up ─────────────────────────

class SignUpScreen extends StatefulWidget {
  const SignUpScreen({super.key});

  @override
  State<SignUpScreen> createState() => _SignUpScreenState();
}

class _SignUpScreenState extends State<SignUpScreen> {
  final _form = GlobalKey<FormState>();
  final _name = TextEditingController();
  final _surname = TextEditingController();
  final _username = TextEditingController();
  final _email = TextEditingController();
  final _phone = TextEditingController();
  final _birth = TextEditingController();
  final _password = TextEditingController();
  final _about = TextEditingController();
  bool _hide = true;
  bool _busy = false;

  @override
  void dispose() {
    for (final c in [_name, _surname, _username, _email, _phone, _birth, _password, _about]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_form.currentState!.validate()) return;
    setState(() => _busy = true);
    final state = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      await state.signUp({
        'name': _name.text.trim(),
        'surname': _surname.text.trim(),
        'username': _username.text.trim(),
        'email': _email.text.trim(),
        'phone_number': _phone.text.trim(),
        'birth_date': _birth.text.trim(),
        'password': _password.text,
        if (_about.text.trim().isNotEmpty) 'description': _about.text.trim(),
      });
      nav.popUntil((r) => r.isFirst);
    } catch (e) {
      if (!mounted) return;
      final conflict = e is ApiException && e.isConflict;
      showSnack(
        context,
        conflict
            ? 'That email, username or phone number is already registered.'
            : friendlyError(e),
        error: true,
      );
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return _AuthScaffold(
      title: 'Create your\nTripMate account',
      subtitle: 'A minute now, unforgettable routes later.',
      child: Form(
        key: _form,
        autovalidateMode: AutovalidateMode.onUserInteraction,
        child: Column(children: [
          Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Expanded(
              child: TextFormField(
                controller: _name,
                textCapitalization: TextCapitalization.words,
                decoration: const InputDecoration(labelText: 'Name'),
                validator: (v) => requiredText(v, max: 50, label: 'Name'),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: TextFormField(
                controller: _surname,
                textCapitalization: TextCapitalization.words,
                decoration: const InputDecoration(labelText: 'Surname'),
                validator: (v) => requiredText(v, max: 50, label: 'Surname'),
              ),
            ),
          ]),
          const SizedBox(height: 14),
          TextFormField(
            controller: _username,
            decoration: const InputDecoration(
                labelText: 'Username', prefixIcon: Icon(Icons.alternate_email_rounded)),
            validator: (v) => requiredText(v, max: 100, label: 'Username'),
          ),
          const SizedBox(height: 14),
          TextFormField(
            controller: _email,
            keyboardType: TextInputType.emailAddress,
            decoration: const InputDecoration(
                labelText: 'Email', prefixIcon: Icon(Icons.mail_outline_rounded)),
            validator: (v) {
              final t = (v ?? '').trim();
              if (!RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$').hasMatch(t)) {
                return 'Enter a valid email';
              }
              return null;
            },
          ),
          const SizedBox(height: 14),
          TextFormField(
            controller: _phone,
            keyboardType: TextInputType.phone,
            inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[0-9+]'))],
            decoration: const InputDecoration(
                labelText: 'Phone (e.g. +79990746978)',
                prefixIcon: Icon(Icons.phone_outlined)),
            validator: (v) {
              final t = (v ?? '').trim();
              if (!RegExp(r'^\+?[0-9]{8,15}$').hasMatch(t)) {
                return 'Enter a valid phone number';
              }
              return null;
            },
          ),
          const SizedBox(height: 14),
          TextFormField(
            controller: _birth,
            readOnly: true,
            onTap: () => pickBirthDate(context, _birth),
            decoration: const InputDecoration(
                labelText: 'Date of birth',
                prefixIcon: Icon(Icons.cake_outlined)),
            validator: (v) =>
                (v ?? '').isEmpty ? 'Pick your date of birth' : null,
          ),
          const SizedBox(height: 14),
          TextFormField(
            controller: _password,
            obscureText: _hide,
            decoration: InputDecoration(
              labelText: 'Password (8+ characters)',
              prefixIcon: const Icon(Icons.lock_outline_rounded),
              suffixIcon: IconButton(
                icon: Icon(_hide
                    ? Icons.visibility_outlined
                    : Icons.visibility_off_outlined),
                onPressed: () => setState(() => _hide = !_hide),
              ),
            ),
            validator: (v) {
              final t = v ?? '';
              if (t.length < 8) return 'At least 8 characters';
              if (t.length > 100) return 'At most 100 characters';
              return null;
            },
          ),
          const SizedBox(height: 14),
          TextFormField(
            controller: _about,
            maxLines: 3,
            maxLength: 1000,
            decoration: const InputDecoration(
                labelText: 'About you (optional)', alignLabelWithHint: true),
          ),
          const SizedBox(height: 18),
          SizedBox(
            width: double.infinity,
            child: GradientButton(
                label: 'Create account', loading: _busy, onPressed: _submit),
          ),
        ]),
      ),
    );
  }
}

// ───────────────────────── Sign in ─────────────────────────

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignInScreenState();
}

class _SignInScreenState extends State<SignInScreen> {
  final _form = GlobalKey<FormState>();
  final _id = TextEditingController();
  bool _busy = false;

  @override
  void dispose() {
    _id.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_form.currentState!.validate()) return;
    setState(() => _busy = true);
    final state = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      await state.signIn(int.parse(_id.text.trim()));
      nav.popUntil((r) => r.isFirst);
    } catch (e) {
      if (!mounted) return;
      showSnack(
        context,
        e is ApiException && e.isNotFound
            ? 'No account with that ID.'
            : friendlyError(e),
        error: true,
      );
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return _AuthScaffold(
      title: 'Welcome back',
      subtitle: 'Enter your TripMate ID to pick up where you left off.',
      child: Form(
        key: _form,
        child: Column(children: [
          TextFormField(
            controller: _id,
            keyboardType: TextInputType.number,
            inputFormatters: [FilteringTextInputFormatter.digitsOnly],
            decoration: const InputDecoration(
                labelText: 'TripMate ID', prefixIcon: Icon(Icons.badge_outlined)),
            validator: (v) =>
                int.tryParse((v ?? '').trim()) == null ? 'Enter your numeric ID' : null,
          ),
          const SizedBox(height: 14),
          GlassCard(
            padding: const EdgeInsets.all(14),
            child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
              const Icon(Icons.info_outline_rounded, color: TM.teal, size: 20),
              const SizedBox(width: 10),
              Expanded(
                child: Text(
                  'Password sign-in is not available in the API yet, so your ID '
                  'is your key. You can find it in Profile after signing up.',
                  style: TextStyle(color: TM.text.fade(.75), height: 1.4, fontSize: 13),
                ),
              ),
            ]),
          ),
          const SizedBox(height: 22),
          SizedBox(
            width: double.infinity,
            child: GradientButton(
                label: 'Continue', loading: _busy, onPressed: _submit),
          ),
        ]),
      ),
    );
  }
}
