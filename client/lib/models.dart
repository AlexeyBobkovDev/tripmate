import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:latlong2/latlong.dart';

int _int(dynamic v, [int d = 0]) => v is num ? v.toInt() : d;
double _dbl(dynamic v) => v is num ? v.toDouble() : 0.0;
String _str(dynamic v) => v is String ? v : '';
DateTime? _dt(dynamic v) =>
    v is String && v.isNotEmpty ? DateTime.tryParse(v) : null;
Map<String, dynamic> _map(dynamic v) =>
    v is Map<String, dynamic> ? v : <String, dynamic>{};
List<dynamic> _list(dynamic v) => v is List ? v : const <dynamic>[];

// ───────────────────────── Enums ─────────────────────────

enum Interest {
  history('History', Icons.account_balance_rounded, Color(0xFFFFB84D)),
  nature('Nature', Icons.forest_rounded, Color(0xFF3DDC97)),
  food('Food', Icons.ramen_dining_rounded, Color(0xFFFF7A59)),
  architecture('Architecture', Icons.location_city_rounded, Color(0xFF5AA9FF));

  const Interest(this.label, this.icon, this.color);
  final String label;
  final IconData icon;
  final Color color;

  static Interest? tryParse(dynamic v) {
    for (final i in Interest.values) {
      if (i.name == v) return i;
    }
    return null;
  }
}

enum TravelMode {
  walking('Walk', Icons.directions_walk_rounded),
  cycling('Cycle', Icons.directions_bike_rounded),
  driving('Drive', Icons.directions_car_rounded),
  transit('Transit', Icons.directions_transit_rounded);

  const TravelMode(this.label, this.icon);
  final String label;
  final IconData icon;

  static TravelMode tryParse(dynamic v) {
    for (final m in TravelMode.values) {
      if (m.name == v) return m;
    }
    return TravelMode.walking;
  }
}

enum TripPhase { upcoming, active, past }

// ───────────────────────── Value objects ─────────────────────────

class Coordinates {
  const Coordinates(this.latitude, this.longitude);
  final double latitude;
  final double longitude;

  factory Coordinates.fromJson(Map<String, dynamic> j) =>
      Coordinates(_dbl(j['latitude']), _dbl(j['longitude']));

  Map<String, dynamic> toJson() =>
      {'latitude': latitude, 'longitude': longitude};

  LatLng get latLng => LatLng(latitude, longitude);
}

class Money {
  const Money(this.amountMinorUnits, this.currency);
  final int amountMinorUnits;
  final String currency;

  /// Minor units per major unit (cents per dollar, kopecks per rouble…).
  static int factorFor(String currency) => currency == 'JPY' ? 1 : 100;

  double get major => amountMinorUnits / factorFor(currency);

  String get pretty => '${NumberFormat('#,##0.##').format(major)} $currency';

  factory Money.fromJson(Map<String, dynamic> j) =>
      Money(_int(j['amount_minor_units']), _str(j['currency']));

  Map<String, dynamic> toJson() =>
      {'amount_minor_units': amountMinorUnits, 'currency': currency};
}

// ───────────────────────── Users ─────────────────────────

class AppUser {
  const AppUser({
    required this.id,
    required this.version,
    required this.name,
    required this.surname,
    required this.username,
    required this.email,
    required this.phoneNumber,
    required this.birthDate,
    required this.description,
  });

  final int id;
  final int version;
  final String name;
  final String surname;
  final String username;
  final String email;
  final String phoneNumber;
  final String birthDate; // yyyy-MM-dd
  final String description;

  String get fullName => '$name $surname'.trim();
  String get initials {
    final a = name.isNotEmpty ? name[0] : '';
    final b = surname.isNotEmpty ? surname[0] : '';
    final r = (a + b).toUpperCase();
    return r.isEmpty ? '?' : r;
  }

  factory AppUser.fromJson(Map<String, dynamic> j) {
    var bd = _str(j['birth_date']);
    if (bd.length > 10) bd = bd.substring(0, 10);
    return AppUser(
      id: _int(j['id']),
      version: _int(j['version'], 1),
      name: _str(j['name']),
      surname: _str(j['surname']),
      username: _str(j['username']),
      email: _str(j['email']),
      phoneNumber: _str(j['phone_number']),
      birthDate: bd,
      description: _str(j['description']),
    );
  }
}

// ───────────────────────── Trips ─────────────────────────

class Trip {
  const Trip({
    required this.id,
    required this.ownerId,
    required this.version,
    required this.name,
    required this.planningArea,
    required this.timeZone,
    required this.startDate,
    required this.endDate,
    required this.startLocation,
    required this.budget,
    required this.interests,
  });

  final int id;
  final int ownerId;
  final int version;
  final String name;
  final String planningArea;
  final String timeZone;
  final DateTime startDate;
  final DateTime endDate;
  final Coordinates startLocation;
  final Money budget;
  final List<Interest> interests;

  int get dayCount {
    final s = startDate.toUtc();
    final e = endDate.toUtc();
    return DateTime.utc(e.year, e.month, e.day)
            .difference(DateTime.utc(s.year, s.month, s.day))
            .inDays +
        1;
  }

  TripPhase get phase {
    final now = DateTime.now().toUtc();
    if (now.isAfter(endDate.toUtc())) return TripPhase.past;
    if (now.isBefore(startDate.toUtc())) return TripPhase.upcoming;
    return TripPhase.active;
  }

  int get daysUntilStart {
    final n = DateTime.now().toUtc();
    final s = startDate.toUtc();
    return DateTime.utc(s.year, s.month, s.day)
        .difference(DateTime.utc(n.year, n.month, n.day))
        .inDays;
  }

  factory Trip.fromJson(Map<String, dynamic> j) {
    final now = DateTime.now().toUtc();
    return Trip(
      id: _int(j['id']),
      ownerId: _int(j['owner_id']),
      version: _int(j['version'], 1),
      name: _str(j['name']),
      planningArea: _str(j['planning_area']),
      timeZone: _str(j['time_zone']),
      startDate: _dt(j['start_date']) ?? now,
      endDate: _dt(j['end_date']) ?? now,
      startLocation: Coordinates.fromJson(_map(j['start_location'])),
      budget: Money.fromJson(_map(j['budget'])),
      interests: _list(j['interests'])
          .map(Interest.tryParse)
          .whereType<Interest>()
          .toList(),
    );
  }
}

// ───────────────────────── Routes ─────────────────────────

class RouteStop {
  const RouteStop({
    required this.id,
    required this.name,
    required this.description,
    required this.position,
    required this.visitDuration,
    this.placeId,
    this.coordinates,
    this.arrivalTime,
    this.departureTime,
  });

  final int id;
  final String name;
  final String description;
  final int position;
  final int visitDuration; // seconds
  final int? placeId;
  final Coordinates? coordinates;
  final DateTime? arrivalTime;
  final DateTime? departureTime;

  factory RouteStop.fromJson(Map<String, dynamic> j) => RouteStop(
        id: _int(j['id']),
        name: _str(j['name']),
        description: _str(j['description']),
        position: _int(j['position']),
        visitDuration: _int(j['visit_duration']),
        placeId: j['place_id'] is num ? _int(j['place_id']) : null,
        coordinates: j['coordinates'] is Map
            ? Coordinates.fromJson(_map(j['coordinates']))
            : null,
        arrivalTime: _dt(j['arrival_time']),
        departureTime: _dt(j['departure_time']),
      );
}

class RouteLeg {
  const RouteLeg({
    required this.id,
    required this.fromStopId,
    required this.toStopId,
    required this.distance,
    required this.duration,
    required this.travelMode,
  });

  final int id;
  final int fromStopId;
  final int toStopId;
  final int distance; // metres (assumed)
  final int duration; // seconds (assumed)
  final TravelMode travelMode;

  factory RouteLeg.fromJson(Map<String, dynamic> j) => RouteLeg(
        id: _int(j['id']),
        fromStopId: _int(j['from_stop_id']),
        toStopId: _int(j['to_stop_id']),
        distance: _int(j['distance']),
        duration: _int(j['duration']),
        travelMode: TravelMode.tryParse(j['travel_mode']),
      );
}

class RouteDay {
  const RouteDay({
    required this.id,
    required this.dayNumber,
    required this.stops,
    required this.legs,
    this.date,
  });

  final int id;
  final int dayNumber;
  final DateTime? date;
  final List<RouteStop> stops;
  final List<RouteLeg> legs;

  List<RouteStop> get sortedStops =>
      [...stops]..sort((a, b) => a.position.compareTo(b.position));

  factory RouteDay.fromJson(Map<String, dynamic> j) => RouteDay(
        id: _int(j['id']),
        dayNumber: _int(j['day_number']),
        date: _dt(j['date']),
        stops: _list(j['stops'])
            .whereType<Map<String, dynamic>>()
            .map(RouteStop.fromJson)
            .toList(),
        legs: _list(j['legs'])
            .whereType<Map<String, dynamic>>()
            .map(RouteLeg.fromJson)
            .toList(),
      );
}

class RouteModel {
  const RouteModel({
    required this.id,
    required this.ownerId,
    required this.status,
    required this.tripVersion,
    required this.version,
    required this.days,
    this.createdAt,
  });

  final int id;
  final int ownerId;
  final String status;
  final int tripVersion;
  final int version;
  final List<RouteDay> days;
  final DateTime? createdAt;

  List<RouteDay> get sortedDays =>
      [...days]..sort((a, b) => a.dayNumber.compareTo(b.dayNumber));

  int get stopCount => days.fold(0, (n, d) => n + d.stops.length);

  factory RouteModel.fromJson(Map<String, dynamic> j) => RouteModel(
        id: _int(j['id']),
        ownerId: _int(j['owner_id']),
        status: _str(j['status']),
        tripVersion: _int(j['trip_version'], 1),
        version: _int(j['version'], 1),
        createdAt: _dt(j['created_at']),
        days: _list(j['days'])
            .whereType<Map<String, dynamic>>()
            .map(RouteDay.fromJson)
            .toList(),
      );
}
