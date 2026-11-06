import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../api_client.dart';
import '../app_state.dart';
import '../config.dart';
import '../models.dart';
import '../theme.dart';
import '../ui.dart';

class ProfileTab extends StatelessWidget {
  const ProfileTab({super.key});

  @override
  Widget build(BuildContext context) {
    final s = context.watch<AppState>();
    final u = s.user;
    if (u == null) return const SizedBox.shrink();

    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 24, 20, 120),
      children: [
        FadeSlideIn(
          child: Column(children: [
            Container(
              width: 104,
              height: 104,
              alignment: Alignment.center,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: TM.sunset,
                boxShadow: [BoxShadow(color: TM.coral.fade(.4), blurRadius: 36)],
              ),
              child: Text(u.initials,
                  style: const TextStyle(
                      fontSize: 38,
                      fontWeight: FontWeight.w900,
                      color: Colors.white)),
            ),
            const SizedBox(height: 16),
            Text(u.fullName,
                style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w900)),
            Text('@${u.username}', style: const TextStyle(color: TM.textDim)),
            if (u.description.isNotEmpty) ...[
              const SizedBox(height: 12),
              Text(u.description,
                  textAlign: TextAlign.center,
                  style: TextStyle(color: TM.text.fade(.8), height: 1.4)),
            ],
          ]),
        ),
        const SizedBox(height: 26),
        FadeSlideIn(
          index: 1,
          child: GlassCard(
            child: Column(children: [
              _Row(Icons.mail_outline_rounded, 'Email', u.email),
              _Row(Icons.phone_outlined, 'Phone', u.phoneNumber),
              _Row(Icons.cake_outlined, 'Born', u.birthDate),
              _Row(Icons.badge_outlined, 'TripMate ID', '${u.id}', copy: true),
            ]),
          ),
        ),
        const SizedBox(height: 18),
        FadeSlideIn(
          index: 2,
          child: GradientButton(
            label: 'Edit profile',
            icon: Icons.edit_rounded,
            onPressed: () => pushPage(context, EditProfileScreen(user: u)),
          ),
        ),
        const SizedBox(height: 12),
        FadeSlideIn(
          index: 3,
          child: OutlinedButton.icon(
            style: OutlinedButton.styleFrom(
              padding: const EdgeInsets.symmetric(vertical: 16),
              side: BorderSide(color: Colors.white.fade(.18)),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
            ),
            onPressed: () => s.signOut(),
            icon: const Icon(Icons.logout_rounded, color: TM.text),
            label: const Text('Sign out', style: TextStyle(color: TM.text)),
          ),
        ),
        const SizedBox(height: 6),
        TextButton(
          onPressed: () async {
            final ok = await confirmDialog(
              context,
              title: 'Delete your account?',
              message: 'This permanently removes your profile. This cannot be undone.',
              confirmLabel: 'Delete account',
            );
            if (!ok) return;
            try {
              await s.deleteAccount();
            } catch (e) {
              if (context.mounted) showSnack(context, friendlyError(e), error: true);
            }
          },
          child: const Text('Delete account', style: TextStyle(color: TM.danger)),
        ),
        const SizedBox(height: 10),
        Center(
          child: Text('API · ${AppConfig.apiBaseUrl}',
              style: TextStyle(color: TM.textDim.fade(.6), fontSize: 11)),
        ),
      ],
    );
  }
}

class _Row extends StatelessWidget {
  const _Row(this.icon, this.label, this.value, {this.copy = false});
  final IconData icon;
  final String label;
  final String value;
  final bool copy;

  @override
  Widget build(BuildContext context) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 10),
        child: Row(children: [
          Icon(icon, size: 20, color: TM.teal),
          const SizedBox(width: 14),
          Text(label, style: const TextStyle(color: TM.textDim)),
          const SizedBox(width: 16),
          Expanded(
            child: Text(value.isEmpty ? '—' : value,
                textAlign: TextAlign.end,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(fontWeight: FontWeight.w600)),
          ),
          if (copy)
            InkWell(
              onTap: () {
                Clipboard.setData(ClipboardData(text: value));
                showSnack(context, 'ID copied');
              },
              child: const Padding(
                padding: EdgeInsets.only(left: 10),
                child: Icon(Icons.copy_rounded, size: 18, color: TM.textDim),
              ),
            ),
        ]),
      );
}

// ───────────────────────── Edit profile ─────────────────────────

class EditProfileScreen extends StatefulWidget {
  const EditProfileScreen({super.key, required this.user});
  final AppUser user;

  @override
  State<EditProfileScreen> createState() => _EditProfileScreenState();
}

class _EditProfileScreenState extends State<EditProfileScreen> {
  final _form = GlobalKey<FormState>();
  late final _name = TextEditingController(text: widget.user.name);
  late final _surname = TextEditingController(text: widget.user.surname);
  late final _username = TextEditingController(text: widget.user.username);
  late final _email = TextEditingController(text: widget.user.email);
  late final _phone = TextEditingController(text: widget.user.phoneNumber);
  late final _birth = TextEditingController(text: widget.user.birthDate);
  late final _about = TextEditingController(text: widget.user.description);
  bool _busy = false;

  @override
  void dispose() {
    for (final c in [_name, _surname, _username, _email, _phone, _birth, _about]) {
      c.dispose();
    }
    super.dispose();
  }

  Future<void> _save() async {
    if (!_form.currentState!.validate()) return;
    final u = widget.user;
    // PATCH semantics: send only what changed, plus the version we edited.
    final body = <String, dynamic>{'version': u.version};
    void diff(String key, String now, String before) {
      if (now != before) body[key] = now;
    }

    diff('name', _name.text.trim(), u.name);
    diff('surname', _surname.text.trim(), u.surname);
    diff('username', _username.text.trim(), u.username);
    diff('email', _email.text.trim(), u.email);
    diff('phone_number', _phone.text.trim(), u.phoneNumber);
    diff('birth_date', _birth.text.trim(), u.birthDate);
    diff('description', _about.text.trim(), u.description);

    if (body.length == 1) {
      Navigator.of(context).pop();
      return;
    }

    setState(() => _busy = true);
    final state = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      await state.updateProfile(body);
      nav.pop();
    } catch (e) {
      if (!mounted) return;
      if (e is ApiException && e.isConflict) {
        // Someone (or another device) changed the profile first: resync.
        await state.refreshUser();
        if (!mounted) return;
        showSnack(context,
            'Your profile changed elsewhere (or a value is already taken). Review and try again.',
            error: true);
        nav.pop();
      } else {
        showSnack(context, friendlyError(e), error: true);
      }
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Edit profile')),
      body: AuroraBackground(
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.fromLTRB(24, 8, 24, 32),
            child: Form(
              key: _form,
              autovalidateMode: AutovalidateMode.onUserInteraction,
              child: Column(children: [
                Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Expanded(
                    child: TextFormField(
                      controller: _name,
                      decoration: const InputDecoration(labelText: 'Name'),
                      validator: (v) => requiredText(v, max: 50, label: 'Name'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextFormField(
                      controller: _surname,
                      decoration: const InputDecoration(labelText: 'Surname'),
                      validator: (v) => requiredText(v, max: 50, label: 'Surname'),
                    ),
                  ),
                ]),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _username,
                  decoration: const InputDecoration(
                      labelText: 'Username',
                      prefixIcon: Icon(Icons.alternate_email_rounded)),
                  validator: (v) => requiredText(v, max: 100, label: 'Username'),
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _email,
                  keyboardType: TextInputType.emailAddress,
                  decoration: const InputDecoration(
                      labelText: 'Email', prefixIcon: Icon(Icons.mail_outline_rounded)),
                  validator: (v) =>
                      RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$').hasMatch((v ?? '').trim())
                          ? null
                          : 'Enter a valid email',
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _phone,
                  keyboardType: TextInputType.phone,
                  inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[0-9+]'))],
                  decoration: const InputDecoration(
                      labelText: 'Phone', prefixIcon: Icon(Icons.phone_outlined)),
                  validator: (v) => RegExp(r'^\+?[0-9]{8,15}$').hasMatch((v ?? '').trim())
                      ? null
                      : 'Enter a valid phone number',
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _birth,
                  readOnly: true,
                  onTap: () => pickBirthDate(context, _birth),
                  decoration: const InputDecoration(
                      labelText: 'Date of birth', prefixIcon: Icon(Icons.cake_outlined)),
                  // The API rejects a null birth date, so it must stay set.
                  validator: (v) => (v ?? '').isEmpty ? 'Date of birth is required' : null,
                ),
                const SizedBox(height: 14),
                TextFormField(
                  controller: _about,
                  maxLines: 3,
                  maxLength: 1000,
                  decoration: const InputDecoration(
                      labelText: 'About you', alignLabelWithHint: true),
                ),
                const SizedBox(height: 14),
                SizedBox(
                  width: double.infinity,
                  child: GradientButton(
                      label: 'Save changes', loading: _busy, onPressed: _save),
                ),
              ]),
            ),
          ),
        ),
      ),
    );
  }
}
