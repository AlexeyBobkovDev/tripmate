import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:latlong2/latlong.dart';
import 'package:provider/provider.dart';

import '../api_client.dart';
import '../app_state.dart';
import '../config.dart';
import '../format.dart';
import '../models.dart';
import '../theme.dart';
import '../ui.dart';

const _timeZones = [
  'Europe/Moscow', 'Europe/London', 'Europe/Paris', 'Europe/Berlin',
  'Europe/Istanbul', 'Asia/Dubai', 'Asia/Almaty', 'Asia/Kolkata',
  'Asia/Singapore', 'Asia/Tokyo', 'Australia/Sydney', 'Africa/Cairo',
  'America/New_York', 'America/Chicago', 'America/Los_Angeles',
  'America/Sao_Paulo', 'UTC',
];
const _currencies = ['RUB', 'USD', 'EUR', 'GBP', 'KZT', 'TRY', 'AED', 'JPY'];

/// 4-step wizard used both to create a trip and (with [initial]) to edit one.
class TripFormScreen extends StatefulWidget {
  const TripFormScreen({super.key, this.initial});
  final Trip? initial;

  @override
  State<TripFormScreen> createState() => _TripFormScreenState();
}

class _TripFormScreenState extends State<TripFormScreen> {
  static const _steps = 4;

  final _page = PageController();
  final _k0 = GlobalKey<FormState>();
  final _k2 = GlobalKey<FormState>();
  final _name = TextEditingController();
  final _area = TextEditingController();
  final _budget = TextEditingController();

  int _step = 0;
  bool _busy = false;
  DateTimeRange? _range;
  String _tz = 'Europe/Moscow';
  String _currency = 'RUB';
  final Set<Interest> _interests = {};
  LatLng _start = const LatLng(55.7539, 37.6208);

  bool get _edit => widget.initial != null;

  DateTime _dateOnly(DateTime d) {
    final u = d.toUtc();
    return DateTime(u.year, u.month, u.day);
  }

  @override
  void initState() {
    super.initState();
    final t = widget.initial;
    if (t != null) {
      _name.text = t.name;
      _area.text = t.planningArea;
      _range = DateTimeRange(start: _dateOnly(t.startDate), end: _dateOnly(t.endDate));
      _tz = t.timeZone;
      _currency = t.budget.currency;
      final v = t.budget.major;
      _budget.text = v == v.roundToDouble() ? v.toStringAsFixed(0) : v.toStringAsFixed(2);
      _interests.addAll(t.interests);
      _start = t.startLocation.latLng;
    }
  }

  @override
  void dispose() {
    _page.dispose();
    _name.dispose();
    _area.dispose();
    _budget.dispose();
    super.dispose();
  }

  // ── Navigation between steps ──
  void _go(int step) {
    setState(() => _step = step);
    _page.animateToPage(step,
        duration: const Duration(milliseconds: 380), curve: Curves.easeOutCubic);
  }

  void _next() {
    FocusScope.of(context).unfocus();
    switch (_step) {
      case 0:
        if (!_k0.currentState!.validate()) return;
      case 1:
        if (_range == null) {
          showSnack(context, 'Choose your travel dates', error: true);
          return;
        }
      case 2:
        if (_interests.isEmpty) {
          showSnack(context, 'Pick at least one interest', error: true);
          return;
        }
        if (!_k2.currentState!.validate()) return;
      case 3:
        _submit();
        return;
    }
    _go(_step + 1);
  }

  int get _minor {
    final raw = double.parse(_budget.text.trim().replaceAll(',', '.').replaceAll(' ', ''));
    return (raw * Money.factorFor(_currency)).round();
  }

  Future<void> _submit() async {
    setState(() => _busy = true);
    final s = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      final Trip result;
      if (_edit) {
        final i = widget.initial!;
        // PATCH: send only the changed fields and the version we edited.
        final body = <String, dynamic>{'version': i.version};
        if (_name.text.trim() != i.name) body['name'] = _name.text.trim();
        if (_area.text.trim() != i.planningArea) body['planning_area'] = _area.text.trim();
        if (apiDate(_range!.start) != apiDate(_dateOnly(i.startDate))) {
          body['start_date'] = apiDate(_range!.start);
        }
        if (apiDate(_range!.end) != apiDate(_dateOnly(i.endDate))) {
          body['end_date'] = apiDate(_range!.end);
        }
        if (_tz != i.timeZone) body['time_zone'] = _tz;
        final sameInterests = _interests.length == i.interests.length &&
            _interests.containsAll(i.interests);
        if (!sameInterests) body['interests'] = _interests.map((e) => e.name).toList();
        if (_minor != i.budget.amountMinorUnits || _currency != i.budget.currency) {
          body['budget'] = Money(_minor, _currency).toJson();
        }
        if (_start.latitude != i.startLocation.latitude ||
            _start.longitude != i.startLocation.longitude) {
          body['start_location'] =
              Coordinates(_start.latitude, _start.longitude).toJson();
        }
        if (body.length == 1) {
          nav.pop();
          return;
        }
        result = await s.patchTrip(i.id, body);
      } else {
        result = await s.createTrip({
          'name': _name.text.trim(),
          'planning_area': _area.text.trim(),
          'start_date': apiDate(_range!.start),
          'end_date': apiDate(_range!.end),
          'time_zone': _tz,
          'interests': _interests.map((e) => e.name).toList(),
          'budget': Money(_minor, _currency).toJson(),
          'start_location': Coordinates(_start.latitude, _start.longitude).toJson(),
        });
      }
      nav.pop(result);
    } catch (e) {
      if (!mounted) return;
      if (e is ApiException && e.isConflict && _edit) {
        s.loadTrips();
        showSnack(context,
            'This trip was changed elsewhere. It has been refreshed, please review and retry.',
            error: true);
        nav.pop();
      } else {
        showSnack(context, friendlyError(e), error: true);
      }
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _pickDates() async {
    final now = DateTime.now();
    final r = await showDateRangePicker(
      context: context,
      firstDate: DateTime(now.year - 1),
      lastDate: DateTime(now.year + 5),
      initialDateRange: _range,
      helpText: 'When are you travelling?',
    );
    if (r != null) setState(() => _range = r);
  }

  // ── Steps ──
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_edit ? 'Edit trip' : 'New trip'),
        leading: IconButton(
          icon: Icon(_step == 0 ? Icons.close_rounded : Icons.arrow_back_rounded),
          onPressed: _busy
              ? null
              : () => _step == 0 ? Navigator.of(context).pop() : _go(_step - 1),
        ),
      ),
      body: AuroraBackground(
        child: SafeArea(
          child: Column(children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 4, 24, 12),
              child: Row(children: [
                for (var i = 0; i < _steps; i++)
                  Expanded(
                    child: AnimatedContainer(
                      duration: const Duration(milliseconds: 300),
                      height: 5,
                      margin: EdgeInsets.only(right: i == _steps - 1 ? 0 : 6),
                      decoration: BoxDecoration(
                        gradient: i <= _step ? TM.brand : null,
                        color: i <= _step ? null : Colors.white.fade(.1),
                        borderRadius: BorderRadius.circular(4),
                      ),
                    ),
                  ),
              ]),
            ),
            Expanded(
              child: PageView(
                controller: _page,
                physics: const NeverScrollableScrollPhysics(),
                children: [_stepBasics(), _stepDates(), _stepTaste(), _stepStart()],
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 8, 24, 16),
              child: SizedBox(
                width: double.infinity,
                child: GradientButton(
                  label: _step == _steps - 1
                      ? (_edit ? 'Save changes' : 'Create trip')
                      : 'Continue',
                  icon: _step == _steps - 1 ? Icons.check_rounded : null,
                  loading: _busy,
                  onPressed: _next,
                ),
              ),
            ),
          ]),
        ),
      ),
    );
  }

  Widget _title(String t, String sub) => Padding(
        padding: const EdgeInsets.only(bottom: 24),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(t, style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w900, letterSpacing: -.4)),
          const SizedBox(height: 6),
          Text(sub, style: const TextStyle(color: TM.textDim, height: 1.4)),
        ]),
      );

  Widget _stepBasics() => SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(24, 12, 24, 24),
        child: Form(
          key: _k0,
          autovalidateMode: AutovalidateMode.onUserInteraction,
          child: Column(children: [
            _title('Name your adventure', 'Give it a title and tell us where you are headed.'),
            TextFormField(
              controller: _name,
              maxLength: 100,
              textCapitalization: TextCapitalization.sentences,
              decoration: const InputDecoration(
                  labelText: 'Trip name',
                  hintText: 'Walk around Moscow',
                  prefixIcon: Icon(Icons.edit_note_rounded)),
              validator: (v) => requiredText(v, max: 100, label: 'Trip name'),
            ),
            const SizedBox(height: 8),
            TextFormField(
              controller: _area,
              maxLength: 100,
              textCapitalization: TextCapitalization.words,
              decoration: const InputDecoration(
                  labelText: 'Planning area',
                  hintText: 'Moscow',
                  prefixIcon: Icon(Icons.place_outlined)),
              validator: (v) => requiredText(v, max: 100, label: 'Planning area'),
            ),
          ]),
        ),
      );

  Widget _stepDates() {
    final r = _range;
    return SingleChildScrollView(
      padding: const EdgeInsets.fromLTRB(24, 12, 24, 24),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        _title('When are you going?', 'Pick your dates and the time zone of your destination.'),
        GlassCard(
          onTap: _pickDates,
          padding: const EdgeInsets.all(20),
          child: Row(children: [
            Container(
              width: 52,
              height: 52,
              decoration: BoxDecoration(gradient: TM.brand, borderRadius: BorderRadius.circular(16)),
              child: const Icon(Icons.date_range_rounded, color: Colors.white),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(r == null ? 'Choose dates' : fmtRange(r.start, r.end),
                    style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
                const SizedBox(height: 2),
                Text(
                    r == null
                        ? 'Tap to open the calendar'
                        : '${r.end.difference(r.start).inDays + 1} day(s)',
                    style: const TextStyle(color: TM.textDim)),
              ]),
            ),
            const Icon(Icons.chevron_right_rounded, color: TM.textDim),
          ]),
        ),
        const SizedBox(height: 18),
        DropdownButtonFormField<String>(
          initialValue: _timeZones.contains(_tz) ? _tz : null,
          isExpanded: true,
          dropdownColor: TM.surfaceHigh,
          decoration: InputDecoration(
              labelText: 'Time zone',
              helperText: _timeZones.contains(_tz) ? null : 'Current: $_tz',
              prefixIcon: const Icon(Icons.schedule_rounded)),
          items: [for (final z in _timeZones) DropdownMenuItem(value: z, child: Text(z))],
          onChanged: (v) => setState(() => _tz = v ?? _tz),
        ),
      ]),
    );
  }

  Widget _stepTaste() => SingleChildScrollView(
        padding: const EdgeInsets.fromLTRB(24, 12, 24, 24),
        child: Form(
          key: _k2,
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            _title('What do you love?', 'Choose your interests and set a budget. We shape the route around both.'),
            GridView.count(
              crossAxisCount: 2,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
              childAspectRatio: 1.35,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              children: [for (final i in Interest.values) _interestTile(i)],
            ),
            const SizedBox(height: 22),
            const SectionTitle('Budget'),
            Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Expanded(
                flex: 3,
                child: TextFormField(
                  controller: _budget,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                      labelText: 'Amount', prefixIcon: Icon(Icons.payments_outlined)),
                  validator: (v) {
                    final x = double.tryParse((v ?? '').trim().replaceAll(',', '.').replaceAll(' ', ''));
                    if (x == null || (x * Money.factorFor(_currency)).round() < 1) {
                      return 'Enter an amount above 0';
                    }
                    return null;
                  },
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                flex: 2,
                child: DropdownButtonFormField<String>(
                  initialValue: _currencies.contains(_currency) ? _currency : null,
                  dropdownColor: TM.surfaceHigh,
                  decoration: const InputDecoration(labelText: 'Currency'),
                  items: [for (final c in _currencies) DropdownMenuItem(value: c, child: Text(c))],
                  onChanged: (v) => setState(() => _currency = v ?? _currency),
                ),
              ),
            ]),
          ]),
        ),
      );

  Widget _interestTile(Interest i) {
    final on = _interests.contains(i);
    return GestureDetector(
      onTap: () => setState(() => on ? _interests.remove(i) : _interests.add(i)),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 260),
        curve: Curves.easeOutCubic,
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          gradient: on
              ? LinearGradient(colors: [
                  i.color,
                  Color.lerp(i.color, const Color(0xFF2B1F7A), .6)!,
                ], begin: Alignment.topLeft, end: Alignment.bottomRight)
              : null,
          color: on ? null : TM.surface,
          borderRadius: BorderRadius.circular(22),
          border: Border.all(color: on ? Colors.white.fade(.3) : Colors.white.fade(.07)),
          boxShadow: on ? [BoxShadow(color: i.color.fade(.35), blurRadius: 20, offset: const Offset(0, 8))] : null,
        ),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(children: [
            Icon(i.icon, size: 30, color: on ? Colors.white : i.color),
            const Spacer(),
            AnimatedOpacity(
              opacity: on ? 1 : 0,
              duration: const Duration(milliseconds: 200),
              child: const Icon(Icons.check_circle_rounded, color: Colors.white),
            ),
          ]),
          const Spacer(),
          Text(i.label,
              style: TextStyle(
                  fontSize: 17,
                  fontWeight: FontWeight.w800,
                  color: on ? Colors.white : TM.text)),
        ]),
      ),
    );
  }

  Widget _stepStart() => Padding(
        padding: const EdgeInsets.fromLTRB(24, 12, 24, 8),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          _title('Where do we start?', 'Tap the map to drop your starting point.'),
          Expanded(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(24),
              child: FlutterMap(
                options: MapOptions(
                  initialCenter: _start,
                  initialZoom: 11,
                  onTap: (_, p) => setState(() => _start = p),
                ),
                children: [
                  TileLayer(
                    urlTemplate: AppConfig.mapTiles,
                    subdomains: AppConfig.mapSubdomains,
                    userAgentPackageName: AppConfig.userAgentPackage,
                  ),
                  MarkerLayer(markers: [
                    Marker(
                      point: _start,
                      width: 52,
                      height: 52,
                      alignment: Alignment.topCenter,
                      child: const Icon(Icons.location_on_rounded, color: TM.coral, size: 50),
                    ),
                  ]),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),
          Row(children: [
            const Icon(Icons.my_location_rounded, size: 16, color: TM.teal),
            const SizedBox(width: 8),
            Text('${_start.latitude.toStringAsFixed(5)}, ${_start.longitude.toStringAsFixed(5)}',
                style: const TextStyle(color: TM.textDim)),
          ]),
        ]),
      );
}
