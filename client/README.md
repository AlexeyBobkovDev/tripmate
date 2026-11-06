# TripMate Android client (Flutter + Dart only)

## Run it
```bash
cd tripmate_app
./setup.sh                 # once: creates android/, patches the manifest
flutter run                # emulator talks to your machine via 10.0.2.2:8080
```
Physical device / deployed API:
`flutter run --dart-define=API_BASE_URL=https://api.example.com/api/v1`

> Before a production release, drop `usesCleartextTraffic` from the manifest and serve the API over HTTPS.

## Screens
Welcome (animated) · Sign up · Sign in by ID · Trips (stats, swipe-to-delete, pull-to-refresh) ·
4-step trip wizard (create + edit, map picker) · Trip detail (info, mini-map, routes) ·
Route view (map with polyline + numbered stops, day tabs, timeline, legs) · Profile + edit.

## API coverage (all 14 operations)
| Operation | Where |
|---|---|
| POST /users | Sign up |
| GET /users/{id} | Sign in, app start |
| PATCH /users/{id} | Edit profile (changed fields + `version`) |
| DELETE /users/{id} | Profile → Delete account |
| GET /trips/{ownerID} | Trips tab |
| POST /trips/{ownerID} | Trip wizard |
| GET /trips/{ownerID}/{tripID} | Trip detail refresh |
| PATCH /trips/{ownerID}/{tripID} | Trip wizard in edit mode |
| DELETE /trips/{ownerID}/{tripID} | Swipe / detail |
| POST /routes/{ownerID}/{tripID} | "Generate a route" |
| GET /routes/{tripID} | Trip detail |
| GET /routes/{tripID}/{routeID} | Route screen |
| DELETE /routes/{tripID}/{routeID} | Route screen / detail |

## Assumptions to verify against your backend
1. **No auth endpoint** exists, so the session is the stored user ID (`AppState`).
2. `GET /trips/{ownerID}` is documented as returning one object; the client accepts an object or an array, and treats 404 as "no trips".
3. Route `distance` is treated as metres, `duration`/`visit_duration` as seconds.
4. Route timestamps are shown as wall-clock UTC (no time-zone database on the client).
5. 409 on PATCH = stale `version`; the client refreshes and asks the user to retry.

## Structure
```
lib/
  main.dart            app shell + session gate
  config.dart          base URL, map tiles
  theme.dart           design tokens (TM) + Material 3 dark theme
  models.dart          typed API models
  api_client.dart      REST client + error mapping
  app_state.dart       session + trips (provider)
  ui.dart              shared widgets (aurora bg, glass cards, buttons…)
  screens/             welcome, auth, home, profile, trip form/detail, route
```
