import 'package:flutter/material.dart';
import 'package:flutter_map/flutter_map.dart';
import 'package:provider/provider.dart';

import '../api_client.dart';
import '../app_state.dart';
import '../config.dart';
import '../format.dart';
import '../models.dart';
import '../theme.dart';
import '../ui.dart';
import 'route_screen.dart';
import 'trip_form_screen.dart';

class TripDetailScreen extends StatefulWidget {
  const TripDetailScreen({super.key, required this.trip});
  final Trip trip;

  @override
  State<TripDetailScreen> createState() => _TripDetailScreenState();
}

class _TripDetailScreenState extends State<TripDetailScreen> {
  late Trip _trip = widget.trip;
  List<RouteModel> _routes = [];
  bool _loading = true;
  bool _generating = false;
  String? _error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final s = context.read<AppState>();
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      // Refresh the trip first so edits always use the latest version number.
      try {
        final fresh = await s.api.getTrip(s.ownerId, _trip.id);
        if (mounted) setState(() => _trip = fresh);
        s.replaceTripLocally(fresh);
      } catch (_) {}
      final routes = await s.api.getRoutes(_trip.id);
      routes.sort((a, b) => b.id.compareTo(a.id));
      if (mounted) setState(() => _routes = routes);
    } catch (e) {
      if (mounted) setState(() => _error = friendlyError(e));
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _generate() async {
    final s = context.read<AppState>();
    setState(() => _generating = true);
    try {
      final r = await s.api.createRoute(s.ownerId, _trip.id);
      if (!mounted) return;
      setState(() => _routes.insert(0, r));
      _open(r);
    } catch (e) {
      if (mounted) showSnack(context, friendlyError(e), error: true);
    } finally {
      if (mounted) setState(() => _generating = false);
    }
  }

  Future<void> _open(RouteModel r) async {
    final deleted = await pushPage<bool>(context, RouteScreen(trip: _trip, route: r));
    if (deleted == true && mounted) _load();
  }

  Future<void> _deleteRoute(RouteModel r) async {
    final ok = await confirmDialog(context,
        title: 'Delete this route?', message: 'The generated route will be removed.');
    if (!ok || !mounted) return;
    final s = context.read<AppState>();
    try {
      await s.api.deleteRoute(_trip.id, r.id);
      if (mounted) setState(() => _routes.removeWhere((x) => x.id == r.id));
    } catch (e) {
      if (mounted) showSnack(context, friendlyError(e), error: true);
    }
  }

  Future<void> _edit() async {
    final updated = await pushPage<Trip>(context, TripFormScreen(initial: _trip));
    if (updated != null && mounted) setState(() => _trip = updated);
  }

  Future<void> _deleteTrip() async {
    final ok = await confirmDialog(context,
        title: 'Delete "${_trip.name}"?',
        message: 'The trip and its routes will be removed.');
    if (!ok || !mounted) return;
    final s = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      await s.deleteTrip(_trip.id);
      nav.pop();
    } catch (e) {
      if (mounted) showSnack(context, friendlyError(e), error: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final t = _trip;
    final main = t.interests.isNotEmpty ? t.interests.first : null;
    final tileW = (MediaQuery.of(context).size.width - 40 - 12) / 2;

    return Scaffold(
      body: AuroraBackground(
        child: CustomScrollView(
          slivers: [
            SliverAppBar(
              pinned: true,
              expandedHeight: 230,
              backgroundColor: TM.bg,
              surfaceTintColor: Colors.transparent,
              actions: [
                IconButton(icon: const Icon(Icons.edit_rounded), onPressed: _edit),
                IconButton(
                    icon: const Icon(Icons.delete_outline_rounded),
                    onPressed: _deleteTrip),
              ],
              flexibleSpace: FlexibleSpaceBar(
                titlePadding: const EdgeInsetsDirectional.only(start: 56, bottom: 16, end: 110),
                title: Text(t.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(fontWeight: FontWeight.w900, color: Colors.white)),
                background: Stack(fit: StackFit.expand, children: [
                  Hero(
                    tag: 'trip-bg-${t.id}',
                    child: DecoratedBox(
                        decoration: BoxDecoration(gradient: tripGradient(t))),
                  ),
                  if (main != null)
                    Positioned(
                      right: -20,
                      bottom: -10,
                      child: Icon(main.icon, size: 190, color: Colors.white.fade(.12)),
                    ),
                  Positioned(
                    left: 20,
                    bottom: 66,
                    right: 20,
                    child: Wrap(spacing: 8, runSpacing: 8, children: [
                      Pill(text: t.planningArea, icon: Icons.place_rounded),
                      Pill(text: fmtRange(t.startDate, t.endDate), icon: Icons.calendar_month_rounded),
                    ]),
                  ),
                ]),
              ),
            ),
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(20, 20, 20, 48),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  FadeSlideIn(
                    child: Wrap(spacing: 12, runSpacing: 12, children: [
                      _Info(width: tileW, icon: Icons.timelapse_rounded, label: 'Duration', value: '${t.dayCount} day${t.dayCount == 1 ? '' : 's'}', color: TM.violet),
                      _Info(width: tileW, icon: Icons.account_balance_wallet_rounded, label: 'Budget', value: t.budget.pretty, color: TM.teal),
                      _Info(width: tileW, icon: Icons.schedule_rounded, label: 'Time zone', value: t.timeZone, color: TM.amber),
                      _Info(width: tileW, icon: Icons.tag_rounded, label: 'Version', value: 'v${t.version}', color: TM.coral),
                    ]),
                  ),
                  const SizedBox(height: 22),
                  const SectionTitle('Interests'),
                  Wrap(spacing: 10, runSpacing: 10, children: [
                    for (final i in t.interests)
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                        decoration: BoxDecoration(
                          color: i.color.fade(.16),
                          borderRadius: BorderRadius.circular(30),
                          border: Border.all(color: i.color.fade(.5)),
                        ),
                        child: Row(mainAxisSize: MainAxisSize.min, children: [
                          Icon(i.icon, size: 18, color: i.color),
                          const SizedBox(width: 8),
                          Text(i.label, style: const TextStyle(fontWeight: FontWeight.w700)),
                        ]),
                      ),
                  ]),
                  const SizedBox(height: 22),
                  const SectionTitle('Starting point'),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(22),
                    child: SizedBox(
                      height: 160,
                      child: FlutterMap(
                        options: MapOptions(
                          initialCenter: t.startLocation.latLng,
                          initialZoom: 13,
                          interactionOptions:
                              const InteractionOptions(flags: InteractiveFlag.none),
                        ),
                        children: [
                          TileLayer(
                            urlTemplate: AppConfig.mapTiles,
                            subdomains: AppConfig.mapSubdomains,
                            userAgentPackageName: AppConfig.userAgentPackage,
                          ),
                          MarkerLayer(markers: [
                            Marker(
                              point: t.startLocation.latLng,
                              width: 44,
                              height: 44,
                              alignment: Alignment.topCenter,
                              child: const Icon(Icons.location_on_rounded, color: TM.coral, size: 42),
                            ),
                          ]),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 26),
                  SectionTitle('Routes', trailing: _loading ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2)) : null),
                  SizedBox(
                    width: double.infinity,
                    child: GradientButton(
                      label: _generating ? 'Crafting your route…' : 'Generate a route',
                      icon: Icons.auto_awesome_rounded,
                      loading: _generating,
                      onPressed: _generate,
                    ),
                  ),
                  const SizedBox(height: 14),
                  if (_error != null)
                    ErrorState(message: _error!, onRetry: _load)
                  else if (!_loading && _routes.isEmpty)
                    const Padding(
                      padding: EdgeInsets.symmetric(vertical: 18),
                      child: Text(
                        'No routes yet. Generate one and TripMate will plan every day for you.',
                        textAlign: TextAlign.center,
                        style: TextStyle(color: TM.textDim, height: 1.4),
                      ),
                    )
                  else
                    for (var i = 0; i < _routes.length; i++)
                      FadeSlideIn(
                        index: i,
                        child: Padding(
                          padding: const EdgeInsets.only(bottom: 12),
                          child: _RouteTile(
                            route: _routes[i],
                            outdated: _routes[i].tripVersion < t.version,
                            onTap: () => _open(_routes[i]),
                            onDelete: () => _deleteRoute(_routes[i]),
                          ),
                        ),
                      ),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Info extends StatelessWidget {
  const _Info({required this.width, required this.icon, required this.label, required this.value, required this.color});
  final double width;
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  @override
  Widget build(BuildContext context) => SizedBox(
        width: width,
        child: GlassCard(
          padding: const EdgeInsets.all(14),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Icon(icon, color: color, size: 22),
            const SizedBox(height: 10),
            Text(label, style: const TextStyle(color: TM.textDim, fontSize: 12)),
            const SizedBox(height: 2),
            Text(value,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15)),
          ]),
        ),
      );
}

class _RouteTile extends StatelessWidget {
  const _RouteTile({required this.route, required this.outdated, required this.onTap, required this.onDelete});
  final RouteModel route;
  final bool outdated;
  final VoidCallback onTap;
  final VoidCallback onDelete;

  @override
  Widget build(BuildContext context) {
    return GlassCard(
      onTap: onTap,
      child: Row(children: [
        Container(
          width: 52,
          height: 52,
          decoration: BoxDecoration(gradient: TM.brand, borderRadius: BorderRadius.circular(16)),
          child: const Icon(Icons.alt_route_rounded, color: Colors.white),
        ),
        const SizedBox(width: 14),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(children: [
              Text('Route #${route.id}',
                  style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
              const SizedBox(width: 8),
              if (outdated)
                const Pill(text: 'Outdated', icon: Icons.update_rounded, color: TM.amber),
            ]),
            const SizedBox(height: 4),
            Text(
              '${route.days.length} day${route.days.length == 1 ? '' : 's'} · ${route.stopCount} stops${route.status.isEmpty ? '' : ' · ${route.status}'}',
              style: const TextStyle(color: TM.textDim, fontSize: 13),
            ),
          ]),
        ),
        IconButton(
          icon: const Icon(Icons.delete_outline_rounded, color: TM.textDim),
          onPressed: onDelete,
        ),
      ]),
    );
  }
}
