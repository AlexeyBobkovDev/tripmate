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

class RouteScreen extends StatefulWidget {
  const RouteScreen({super.key, required this.trip, required this.route});
  final Trip trip;
  final RouteModel route;

  @override
  State<RouteScreen> createState() => _RouteScreenState();
}

class _RouteScreenState extends State<RouteScreen> {
  late RouteModel _route = widget.route;
  int _day = 0;

  @override
  void initState() {
    super.initState();
    _refresh();
  }

  /// Silently pulls the freshest copy of the route.
  Future<void> _refresh() async {
    final s = context.read<AppState>();
    try {
      final r = await s.api.getRoute(widget.trip.id, _route.id);
      if (mounted) setState(() => _route = r);
    } catch (_) {}
  }

  Future<void> _delete() async {
    final ok = await confirmDialog(context,
        title: 'Delete this route?', message: 'This cannot be undone.');
    if (!ok || !mounted) return;
    final s = context.read<AppState>();
    final nav = Navigator.of(context);
    try {
      await s.api.deleteRoute(widget.trip.id, _route.id);
      nav.pop(true);
    } catch (e) {
      if (mounted) showSnack(context, friendlyError(e), error: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final days = _route.sortedDays;
    final idx = days.isEmpty ? 0 : _day.clamp(0, days.length - 1);
    final day = days.isEmpty ? null : days[idx];
    final stops = day?.sortedStops ?? <RouteStop>[];
    final points = <LatLng>[
      for (final s in stops)
        if (s.coordinates != null) s.coordinates!.latLng,
    ];
    final totalDist = day?.legs.fold<int>(0, (n, l) => n + l.distance) ?? 0;
    final totalTime = day?.legs.fold<int>(0, (n, l) => n + l.duration) ?? 0;
    final top = MediaQuery.of(context).padding.top;

    return Scaffold(
      body: AuroraBackground(
        child: Column(children: [
          SizedBox(
            height: 330,
            child: Stack(children: [
              Positioned.fill(child: _buildMap(idx, stops, points)),
              Positioned(
                top: top + 8,
                left: 12,
                child: _RoundBtn(Icons.arrow_back_rounded, () => Navigator.of(context).pop()),
              ),
              Positioned(
                top: top + 8,
                right: 12,
                child: _RoundBtn(Icons.delete_outline_rounded, _delete),
              ),
            ]),
          ),
          Expanded(
            child: Container(
              width: double.infinity,
              decoration: const BoxDecoration(
                color: TM.bg,
                borderRadius: BorderRadius.vertical(top: Radius.circular(30)),
              ),
              child: ClipRRect(
                borderRadius: const BorderRadius.vertical(top: Radius.circular(30)),
                child: days.isEmpty
                    ? const EmptyState(
                        icon: Icons.map_outlined,
                        title: 'Nothing planned yet',
                        subtitle: 'This route has no days. Try generating a new one.')
                    : Column(children: [
                        const SizedBox(height: 16),
                        SizedBox(
                          height: 46,
                          child: ListView.separated(
                            scrollDirection: Axis.horizontal,
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            itemCount: days.length,
                            separatorBuilder: (_, __) => const SizedBox(width: 10),
                            itemBuilder: (_, i) {
                              final d = days[i];
                              final on = i == idx;
                              return GestureDetector(
                                onTap: () => setState(() => _day = i),
                                child: AnimatedContainer(
                                  duration: const Duration(milliseconds: 250),
                                  padding: const EdgeInsets.symmetric(horizontal: 18),
                                  alignment: Alignment.center,
                                  decoration: BoxDecoration(
                                    gradient: on ? TM.brand : null,
                                    color: on ? null : TM.surface,
                                    borderRadius: BorderRadius.circular(30),
                                    border: Border.all(color: Colors.white.fade(on ? 0 : .08)),
                                  ),
                                  child: Text(
                                    'Day ${d.dayNumber}${d.date == null ? '' : ' · ${fmtDayMonth(d.date!)}'}',
                                    style: TextStyle(
                                        fontWeight: FontWeight.w700,
                                        color: on ? Colors.white : TM.textDim),
                                  ),
                                ),
                              );
                            },
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.fromLTRB(20, 14, 20, 4),
                          child: Wrap(spacing: 8, runSpacing: 8, children: [
                            Pill(text: '${stops.length} stops', icon: Icons.flag_rounded, color: TM.teal),
                            if (totalDist > 0)
                              Pill(text: fmtDistance(totalDist), icon: Icons.straighten_rounded, color: TM.amber),
                            if (totalTime > 0)
                              Pill(text: '${fmtDuration(totalTime)} travel', icon: Icons.timer_outlined, color: TM.coral),
                          ]),
                        ),
                        Expanded(
                          child: stops.isEmpty
                              ? const Center(
                                  child: Text('No stops on this day.',
                                      style: TextStyle(color: TM.textDim)))
                              : ListView.builder(
                                  padding: const EdgeInsets.fromLTRB(20, 14, 20, 40),
                                  itemCount: stops.length,
                                  itemBuilder: (_, i) => FadeSlideIn(
                                    index: i,
                                    child: _StopTile(
                                      number: i + 1,
                                      stop: stops[i],
                                      leg: _legFrom(day!, stops[i]),
                                      isLast: i == stops.length - 1,
                                    ),
                                  ),
                                ),
                        ),
                      ]),
              ),
            ),
          ),
        ]),
      ),
    );
  }

  RouteLeg? _legFrom(RouteDay day, RouteStop s) {
    for (final l in day.legs) {
      if (l.fromStopId == s.id) return l;
    }
    return null;
  }

  Widget _buildMap(int dayIdx, List<RouteStop> stops, List<LatLng> points) {
    final fallback = widget.trip.startLocation.latLng;
    return FlutterMap(
      // New key per day so the camera re-fits to that day's stops.
      key: ValueKey('map-${_route.id}-$dayIdx-${points.length}'),
      options: MapOptions(
        initialCenter: points.isNotEmpty ? points.first : fallback,
        initialZoom: 13,
        initialCameraFit: points.length > 1
            ? CameraFit.bounds(
                bounds: LatLngBounds.fromPoints(points),
                padding: const EdgeInsets.fromLTRB(56, 96, 56, 56),
              )
            : null,
      ),
      children: [
        TileLayer(
          urlTemplate: AppConfig.mapTiles,
          subdomains: AppConfig.mapSubdomains,
          userAgentPackageName: AppConfig.userAgentPackage,
        ),
        if (points.length > 1)
          PolylineLayer(polylines: [
            Polyline(points: points, strokeWidth: 5, color: TM.teal),
          ]),
        MarkerLayer(markers: [
          for (var i = 0; i < stops.length; i++)
            if (stops[i].coordinates != null)
              Marker(
                point: stops[i].coordinates!.latLng,
                width: 36,
                height: 36,
                child: Container(
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    gradient: TM.sunset,
                    shape: BoxShape.circle,
                    border: Border.all(color: Colors.white, width: 2),
                    boxShadow: [BoxShadow(color: TM.coral.fade(.5), blurRadius: 10)],
                  ),
                  child: Text('${i + 1}',
                      style: const TextStyle(
                          color: Colors.white, fontWeight: FontWeight.w900)),
                ),
              ),
        ]),
      ],
    );
  }
}

class _RoundBtn extends StatelessWidget {
  const _RoundBtn(this.icon, this.onTap);
  final IconData icon;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) => Material(
        color: TM.bg.fade(.75),
        shape: const CircleBorder(),
        child: InkWell(
          customBorder: const CircleBorder(),
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.all(11),
            child: Icon(icon, color: TM.text, size: 22),
          ),
        ),
      );
}

class _StopTile extends StatelessWidget {
  const _StopTile({required this.number, required this.stop, required this.leg, required this.isLast});
  final int number;
  final RouteStop stop;
  final RouteLeg? leg;
  final bool isLast;

  @override
  Widget build(BuildContext context) {
    final a = stop.arrivalTime;
    final d = stop.departureTime;
    final time = a == null
        ? null
        : d == null
            ? fmtTime(a)
            : '${fmtTime(a)} – ${fmtTime(d)}';

    return IntrinsicHeight(
      child: Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        SizedBox(
          width: 44,
          child: Column(children: [
            Container(
              width: 34,
              height: 34,
              alignment: Alignment.center,
              decoration: const BoxDecoration(gradient: TM.sunset, shape: BoxShape.circle),
              child: Text('$number',
                  style: const TextStyle(fontWeight: FontWeight.w900, color: Colors.white)),
            ),
            if (!isLast)
              Expanded(
                child: Container(
                  width: 2,
                  margin: const EdgeInsets.symmetric(vertical: 4),
                  color: Colors.white.fade(.12),
                ),
              ),
          ]),
        ),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.only(bottom: 14),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              GlassCard(
                padding: const EdgeInsets.all(14),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text(stop.name.isEmpty ? 'Stop $number' : stop.name,
                      style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800)),
                  if (time != null || stop.visitDuration > 0) ...[
                    const SizedBox(height: 8),
                    Wrap(spacing: 8, runSpacing: 6, children: [
                      if (time != null)
                        Pill(text: time, icon: Icons.schedule_rounded, color: TM.teal),
                      if (stop.visitDuration > 0)
                        Pill(text: 'Stay ${fmtDuration(stop.visitDuration)}', icon: Icons.hourglass_bottom_rounded, color: TM.amber),
                    ]),
                  ],
                  if (stop.description.isNotEmpty) ...[
                    const SizedBox(height: 10),
                    Text(stop.description,
                        style: TextStyle(color: TM.text.fade(.78), height: 1.4, fontSize: 13.5)),
                  ],
                ]),
              ),
              if (leg != null)
                Padding(
                  padding: const EdgeInsets.only(top: 10, left: 4),
                  child: Row(children: [
                    Icon(leg!.travelMode.icon, size: 18, color: TM.violet),
                    const SizedBox(width: 8),
                    Text(
                      '${leg!.travelMode.label} · ${fmtDuration(leg!.duration)} · ${fmtDistance(leg!.distance)}',
                      style: const TextStyle(color: TM.textDim, fontSize: 12.5),
                    ),
                  ]),
                ),
            ]),
          ),
        ),
      ]),
    );
  }
}
