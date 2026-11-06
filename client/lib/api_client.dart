import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:http/http.dart' as http;

import 'config.dart';
import 'models.dart';

class ApiException implements Exception {
  ApiException(this.status, this.message, [this.detail]);

  /// HTTP status, or 0 for network-level failures.
  final int status;
  final String message;
  final String? detail;

  bool get isNotFound => status == 404;
  bool get isConflict => status == 409;

  @override
  String toString() => message;
}

String friendlyError(Object e) {
  if (e is ApiException) {
    if (e.message.trim().isNotEmpty) return e.message;
    return 'Server error (${e.status}). Please try again.';
  }
  return 'Something went wrong. Please try again.';
}

/// Thin typed wrapper over the TripMate REST API (basePath /api/v1).
class ApiClient {
  ApiClient({http.Client? client, String? baseUrl})
      : _http = client ?? http.Client(),
        _base = baseUrl ?? AppConfig.apiBaseUrl;

  final http.Client _http;
  final String _base;

  Future<dynamic> _request(
    String method,
    String path, {
    Object? body,
    Duration timeout = const Duration(seconds: 15),
  }) async {
    final req = http.Request(method, Uri.parse('$_base$path'))
      ..headers['Accept'] = 'application/json';
    if (body != null) {
      req.headers['Content-Type'] = 'application/json';
      req.body = jsonEncode(body);
    }

    late http.Response res;
    try {
      res = await _http
          .send(req)
          .then(http.Response.fromStream)
          .timeout(timeout);
    } on TimeoutException {
      throw ApiException(0, 'The server took too long to respond.');
    } on SocketException {
      throw ApiException(
          0, "Can't reach TripMate servers. Check your connection.");
    } on http.ClientException {
      throw ApiException(
          0, "Can't reach TripMate servers. Check your connection.");
    }

    if (res.statusCode >= 400) throw _toError(res);
    if (res.statusCode == 204 || res.body.isEmpty) return null;
    return jsonDecode(utf8.decode(res.bodyBytes));
  }

  ApiException _toError(http.Response res) {
    String message = '';
    String? detail;
    try {
      final j = jsonDecode(utf8.decode(res.bodyBytes));
      if (j is Map<String, dynamic>) {
        message = (j['message'] as String?) ?? '';
        detail = j['error'] as String?;
        if (message.isEmpty) message = detail ?? '';
      }
    } catch (_) {}
    if (message.isEmpty) {
      message = switch (res.statusCode) {
        400 => 'Some of the data is invalid.',
        404 => 'Not found.',
        409 => 'This conflicts with the current state. Refresh and retry.',
        _ => 'Server error (${res.statusCode}).',
      };
    }
    return ApiException(res.statusCode, message, detail);
  }

  Map<String, dynamic> _obj(dynamic v) {
    if (v is Map<String, dynamic>) return v;
    throw ApiException(0, 'Unexpected response from server.');
  }

  // ── Users ──
  Future<AppUser> createUser(Map<String, dynamic> body) async =>
      AppUser.fromJson(_obj(await _request('POST', '/users', body: body)));

  Future<AppUser> getUser(int id) async =>
      AppUser.fromJson(_obj(await _request('GET', '/users/$id')));

  Future<AppUser> patchUser(int id, Map<String, dynamic> body) async =>
      AppUser.fromJson(
          _obj(await _request('PATCH', '/users/$id', body: body)));

  Future<void> deleteUser(int id) => _request('DELETE', '/users/$id');

  // ── Trips ──
  /// The swagger declares a single object for GET /trips/{ownerID}, but the
  /// endpoint is semantically a list, so both shapes are accepted.
  Future<List<Trip>> getTrips(int ownerId) async {
    try {
      final data = await _request('GET', '/trips/$ownerId');
      if (data is List) {
        return data
            .whereType<Map<String, dynamic>>()
            .map(Trip.fromJson)
            .toList();
      }
      if (data is Map<String, dynamic>) return [Trip.fromJson(data)];
      return [];
    } on ApiException catch (e) {
      if (e.isNotFound) return []; // no trips yet
      rethrow;
    }
  }

  Future<Trip> getTrip(int ownerId, int tripId) async =>
      Trip.fromJson(_obj(await _request('GET', '/trips/$ownerId/$tripId')));

  Future<Trip> createTrip(int ownerId, Map<String, dynamic> body) async =>
      Trip.fromJson(
          _obj(await _request('POST', '/trips/$ownerId', body: body)));

  Future<Trip> patchTrip(
          int ownerId, int tripId, Map<String, dynamic> body) async =>
      Trip.fromJson(_obj(
          await _request('PATCH', '/trips/$ownerId/$tripId', body: body)));

  Future<void> deleteTrip(int ownerId, int tripId) =>
      _request('DELETE', '/trips/$ownerId/$tripId');

  // ── Routes ──
  Future<List<RouteModel>> getRoutes(int tripId) async {
    try {
      final data = await _request('GET', '/routes/$tripId');
      if (data is List) {
        return data
            .whereType<Map<String, dynamic>>()
            .map(RouteModel.fromJson)
            .toList();
      }
      return [];
    } on ApiException catch (e) {
      if (e.isNotFound) return [];
      rethrow;
    }
  }

  Future<RouteModel> getRoute(int tripId, int routeId) async =>
      RouteModel.fromJson(
          _obj(await _request('GET', '/routes/$tripId/$routeId')));

  /// Route generation can take a while, so the timeout is generous.
  Future<RouteModel> createRoute(int ownerId, int tripId) async =>
      RouteModel.fromJson(_obj(await _request(
        'POST',
        '/routes/$ownerId/$tripId',
        timeout: const Duration(seconds: 90),
      )));

  Future<void> deleteRoute(int tripId, int routeId) =>
      _request('DELETE', '/routes/$tripId/$routeId');
}
