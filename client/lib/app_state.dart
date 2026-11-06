import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'api_client.dart';
import 'models.dart';

enum AuthStatus { booting, signedOut, signedIn, bootFailed }

/// App-wide state: session + the signed-in user's trips.
///
/// The API has no authentication endpoint yet, so the session is simply the
/// user ID persisted on the device. When auth lands, swap it for a token here.
class AppState extends ChangeNotifier {
  AppState(this.api);

  final ApiClient api;

  static const _kUserId = 'tripmate.user_id';

  AuthStatus status = AuthStatus.booting;
  String? bootError;
  AppUser? user;

  List<Trip> trips = [];
  bool tripsLoading = false;
  String? tripsError;

  int get ownerId => user!.id;

  // ── Session ──
  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    final id = prefs.getInt(_kUserId);
    if (id == null) {
      status = AuthStatus.signedOut;
      notifyListeners();
      return;
    }
    status = AuthStatus.booting;
    notifyListeners();
    try {
      user = await api.getUser(id);
      status = AuthStatus.signedIn;
      notifyListeners();
      loadTrips();
    } on ApiException catch (e) {
      if (e.isNotFound) {
        await prefs.remove(_kUserId);
        status = AuthStatus.signedOut;
      } else {
        bootError = e.message;
        status = AuthStatus.bootFailed;
      }
      notifyListeners();
    }
  }

  Future<void> _startSession(AppUser u) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt(_kUserId, u.id);
    user = u;
    trips = [];
    tripsError = null;
    status = AuthStatus.signedIn;
    notifyListeners();
    loadTrips();
  }

  Future<void> signUp(Map<String, dynamic> body) async =>
      _startSession(await api.createUser(body));

  Future<void> signIn(int id) async => _startSession(await api.getUser(id));

  Future<void> signOut() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_kUserId);
    user = null;
    trips = [];
    status = AuthStatus.signedOut;
    notifyListeners();
  }

  Future<void> updateProfile(Map<String, dynamic> body) async {
    user = await api.patchUser(user!.id, body);
    notifyListeners();
  }

  /// Re-fetches the profile (e.g. after a version conflict). Errors are ignored.
  Future<void> refreshUser() async {
    if (user == null) return;
    try {
      user = await api.getUser(user!.id);
      notifyListeners();
    } catch (_) {}
  }

  Future<void> deleteAccount() async {
    await api.deleteUser(user!.id);
    await signOut();
  }

  // ── Trips ──
  Future<void> loadTrips() async {
    if (user == null) return;
    tripsLoading = true;
    tripsError = null;
    notifyListeners();
    try {
      trips = await api.getTrips(user!.id);
    } catch (e) {
      tripsError = friendlyError(e);
    } finally {
      tripsLoading = false;
      notifyListeners();
    }
  }

  Future<Trip> createTrip(Map<String, dynamic> body) async {
    final t = await api.createTrip(ownerId, body);
    trips = [t, ...trips];
    notifyListeners();
    return t;
  }

  Future<Trip> patchTrip(int tripId, Map<String, dynamic> body) async {
    final t = await api.patchTrip(ownerId, tripId, body);
    trips = [for (final x in trips) x.id == tripId ? t : x];
    notifyListeners();
    return t;
  }

  /// Deletes on the server and removes locally.
  Future<void> deleteTrip(int tripId) async {
    await api.deleteTrip(ownerId, tripId);
    dropTripLocally(tripId);
  }

  void replaceTripLocally(Trip t) {
    trips = [for (final x in trips) x.id == t.id ? t : x];
    notifyListeners();
  }

  void dropTripLocally(int tripId) {
    trips = trips.where((t) => t.id != tripId).toList();
    notifyListeners();
  }
}
