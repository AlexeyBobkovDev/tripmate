import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../api_client.dart';
import '../app_state.dart';
import '../format.dart';
import '../models.dart';
import '../theme.dart';
import '../ui.dart';
import 'profile_screens.dart';
import 'trip_detail_screen.dart';
import 'trip_form_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _tab = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      extendBody: true,
      body: AuroraBackground(
        child: SafeArea(
          bottom: false,
          child: IndexedStack(
            index: _tab,
            children: const [TripsTab(), ProfileTab()],
          ),
        ),
      ),
      floatingActionButton: _tab == 0
          ? FloatingActionButton.extended(
              onPressed: () => pushPage(context, const TripFormScreen()),
              icon: const Icon(Icons.add_rounded),
              label: const Text('New trip',
                  style: TextStyle(fontWeight: FontWeight.w700)),
            )
          : null,
      bottomNavigationBar: NavigationBar(
        selectedIndex: _tab,
        onDestinationSelected: (i) => setState(() => _tab = i),
        destinations: const [
          NavigationDestination(
              icon: Icon(Icons.luggage_outlined),
              selectedIcon: Icon(Icons.luggage_rounded),
              label: 'Trips'),
          NavigationDestination(
              icon: Icon(Icons.person_outline_rounded),
              selectedIcon: Icon(Icons.person_rounded),
              label: 'Profile'),
        ],
      ),
    );
  }
}

// ───────────────────────── Trips tab ─────────────────────────

class TripsTab extends StatelessWidget {
  const TripsTab({super.key});

  Future<bool> _confirmDelete(BuildContext context, AppState s, Trip t) async {
    final ok = await confirmDialog(
      context,
      title: 'Delete "${t.name}"?',
      message: 'The trip and its generated routes will be removed.',
    );
    if (!ok) return false;
    try {
      await s.api.deleteTrip(s.ownerId, t.id);
      return true;
    } catch (e) {
      if (context.mounted) showSnack(context, friendlyError(e), error: true);
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    final s = context.watch<AppState>();
    final trips = [...s.trips]
      ..sort((a, b) => a.startDate.compareTo(b.startDate));
    final current = trips.where((t) => t.phase != TripPhase.past).toList();
    final past =
        trips.where((t) => t.phase == TripPhase.past).toList().reversed.toList();
    final upcoming = trips.where((t) => t.phase == TripPhase.upcoming).length;

    final items = <Widget>[];
    var n = 0;
    Widget card(Trip t) => FadeSlideIn(
          index: n++,
          child: Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: Dismissible(
              key: ValueKey('trip-${t.id}'),
              direction: DismissDirection.endToStart,
              confirmDismiss: (_) => _confirmDelete(context, s, t),
              onDismissed: (_) => s.dropTripLocally(t.id),
              background: Container(
                alignment: Alignment.centerRight,
                padding: const EdgeInsets.only(right: 28),
                decoration: BoxDecoration(
                    color: TM.danger.fade(.85),
                    borderRadius: BorderRadius.circular(26)),
                child: const Icon(Icons.delete_outline_rounded,
                    color: Colors.white, size: 30),
              ),
              child: TripCard(
                trip: t,
                onTap: () => pushPage(context, TripDetailScreen(trip: t)),
              ),
            ),
          ),
        );

    if (s.tripsLoading && s.trips.isEmpty) {
      items.addAll(const [
        Padding(padding: EdgeInsets.only(bottom: 16), child: Skeleton(height: 176, radius: 26)),
        Padding(padding: EdgeInsets.only(bottom: 16), child: Skeleton(height: 176, radius: 26)),
        Padding(padding: EdgeInsets.only(bottom: 16), child: Skeleton(height: 176, radius: 26)),
      ]);
    } else if (s.tripsError != null && s.trips.isEmpty) {
      items.add(ErrorState(message: s.tripsError!, onRetry: s.loadTrips));
    } else if (s.trips.isEmpty) {
      items.add(EmptyState(
        icon: Icons.travel_explore_rounded,
        title: 'Your adventures start here',
        subtitle:
            'Pick a city, your interests and a budget. TripMate builds the route.',
        action: GradientButton(
          label: 'Plan your first trip',
          icon: Icons.add_rounded,
          onPressed: () => pushPage(context, const TripFormScreen()),
        ),
      ));
    } else {
      if (current.isNotEmpty) {
        items.add(const SectionTitle('Upcoming & active'));
        items.addAll(current.map(card));
      }
      if (past.isNotEmpty) {
        items.add(const SectionTitle('Past trips'));
        items.addAll(past.map(card));
      }
    }

    return RefreshIndicator(
      color: TM.teal,
      backgroundColor: TM.surfaceHigh,
      onRefresh: s.loadTrips,
      child: CustomScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        slivers: [
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(20, 18, 20, 20),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('${greeting()},',
                    style: const TextStyle(color: TM.textDim, fontSize: 15)),
                const SizedBox(height: 2),
                ShaderMask(
                  shaderCallback: (r) => TM.brand.createShader(r),
                  child: Text(s.user?.name ?? '',
                      style: const TextStyle(
                          fontSize: 34,
                          fontWeight: FontWeight.w900,
                          letterSpacing: -.5,
                          color: Colors.white)),
                ),
                const SizedBox(height: 4),
                const Text('Where to next?',
                    style: TextStyle(fontSize: 17, fontWeight: FontWeight.w600)),
                if (s.trips.isNotEmpty) ...[
                  const SizedBox(height: 18),
                  Row(children: [
                    Expanded(
                        child: _Stat(
                            icon: Icons.luggage_rounded,
                            value: '${s.trips.length}',
                            label: 'Trips',
                            color: TM.violet)),
                    const SizedBox(width: 12),
                    Expanded(
                        child: _Stat(
                            icon: Icons.flight_takeoff_rounded,
                            value: '$upcoming',
                            label: 'Upcoming',
                            color: TM.teal)),
                  ]),
                ],
              ]),
            ),
          ),
          SliverPadding(
            padding: const EdgeInsets.fromLTRB(20, 0, 20, 120),
            sliver: SliverList(delegate: SliverChildListDelegate(items)),
          ),
        ],
      ),
    );
  }
}

class _Stat extends StatelessWidget {
  const _Stat(
      {required this.icon,
      required this.value,
      required this.label,
      required this.color});
  final IconData icon;
  final String value;
  final String label;
  final Color color;

  @override
  Widget build(BuildContext context) => GlassCard(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        child: Row(children: [
          Container(
            width: 42,
            height: 42,
            decoration: BoxDecoration(
                color: color.fade(.18), borderRadius: BorderRadius.circular(14)),
            child: Icon(icon, color: color),
          ),
          const SizedBox(width: 12),
          Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(value,
                style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w900)),
            Text(label, style: const TextStyle(color: TM.textDim, fontSize: 12)),
          ]),
        ]),
      );
}

// ───────────────────────── Trip card ─────────────────────────

class TripCard extends StatelessWidget {
  const TripCard({super.key, required this.trip, required this.onTap});
  final Trip trip;
  final VoidCallback onTap;

  String get _phaseText => switch (trip.phase) {
        TripPhase.active => 'Happening now',
        TripPhase.past => 'Completed',
        TripPhase.upcoming => trip.daysUntilStart <= 0
            ? 'Starts today'
            : 'In ${trip.daysUntilStart} d',
      };

  @override
  Widget build(BuildContext context) {
    final main = trip.interests.isNotEmpty ? trip.interests.first : null;
    final glow = main?.color ?? TM.violet;
    return SizedBox(
      height: 176,
      child: Stack(children: [
        Positioned.fill(
          child: Hero(
            tag: 'trip-bg-${trip.id}',
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: tripGradient(trip),
                borderRadius: BorderRadius.circular(26),
                boxShadow: [
                  BoxShadow(
                      color: glow.fade(.28),
                      blurRadius: 24,
                      offset: const Offset(0, 12)),
                ],
              ),
            ),
          ),
        ),
        if (main != null)
          Positioned.fill(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(26),
              child: Stack(children: [
                Positioned(
                  right: -16,
                  bottom: -22,
                  child: Icon(main.icon,
                      size: 140, color: Colors.white.fade(.13)),
                ),
              ]),
            ),
          ),
        Positioned.fill(
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              borderRadius: BorderRadius.circular(26),
              onTap: onTap,
              child: Padding(
                padding: const EdgeInsets.all(18),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Row(children: [
                    Flexible(
                        child: Pill(
                            text: trip.planningArea, icon: Icons.place_rounded)),
                    const Spacer(),
                    Pill(
                      text: _phaseText,
                      color: trip.phase == TripPhase.active
                          ? TM.teal
                          : Colors.white,
                    ),
                  ]),
                  const Spacer(),
                  Text(trip.name,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                          fontSize: 23,
                          fontWeight: FontWeight.w900,
                          height: 1.1,
                          color: Colors.white)),
                  const SizedBox(height: 8),
                  Row(children: [
                    const Icon(Icons.calendar_month_rounded,
                        size: 15, color: Colors.white70),
                    const SizedBox(width: 6),
                    Text(fmtRange(trip.startDate, trip.endDate),
                        style: const TextStyle(color: Colors.white70, fontSize: 13)),
                    const Spacer(),
                    Text(trip.budget.pretty,
                        style: const TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.w700,
                            fontSize: 13)),
                  ]),
                ]),
              ),
            ),
          ),
        ),
      ]),
    );
  }
}
