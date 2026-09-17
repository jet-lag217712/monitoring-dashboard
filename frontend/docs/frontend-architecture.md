# Frontend Architecture

This frontend is a Vite + React single-page dashboard. `react-router-dom` owns the current screen. The URL names which view is visible; search, selected interface, and wall-slot config stay out of the address bar.

## Where to Start

Start with `src/main.jsx`. It wraps the tree in `BrowserRouter`, then `src/App.jsx`:

- It calls `useAuth()` for appliance-local PAM sessions.
- Unauthenticated operators see `SignInPage` while the current path stays in the address bar, so `/wall/display` survives login.
- Authenticated operators render `AppShell` plus `Routes` for All Sites, site detail, device detail, the wall editor, and the wall display.

After that, read `src/hooks/useNetworkDashboard.js`. That hook loads live site data, tracks search text, and polls for updates. Site and device selection are passed in from the matched URL. API failures surface as an error state; there is no mock fallback.

## Folder Structure

### `src/config`

- `api.js` contains the API base URL, polling interval, auth mode, and URL helper.
- `paths.js` builds every in-app route. Components must not concatenate path strings.

### `src/auth`

- `SignInPage.jsx` is the username/password form for appliance-local sessions.

### `src/services`

API functions live here. `sitesApi.js` contains small functions for backend requests using cookie sessions and CSRF:

- `fetchSitesFromApi()`
- `fetchSiteDetailFromApi(siteId)`

Future API calls should be added to this folder first. Components should not call `fetch()` directly.

### `src/hooks`

- `useAuth.js` owns the PAM cookie session.
- `useNetworkDashboard.js` owns site list, alerts, search, last updated time, live/error data mode, polling, and navigation helpers. `selectedSite` and `selectedDevice` are inputs from the router.

### `src/utils`

Pure helper functions live here. `siteData.js` transforms or summarizes API site responses.

### `src/layout`

- `AppShell.jsx` renders the fixed navigation, alert banner, and main content area.
- `Nav.jsx` renders the top navigation bar, including the Wall link.

### `src/dashboard`

- `DashboardPage.jsx` renders All Sites, site detail, or device detail from the matched route.

### `src/wall`

- `WallEditorPage.jsx` is the kiosk control plane (`/wall`).
- `WallDisplayPage.jsx` is the TV presentation route (`/wall/display`).

### `src/sites`

Site-focused UI: `SitesGrid.jsx`, `SiteCard.jsx`, `SiteDetail.jsx`.

### `src/devices`

Device-focused UI: `DeviceRow.jsx`.

### `src/tables`

- `DevicesTable.jsx` owns the table structure for device detail rows.

### `src/charts`

Chart components for utilization history and interface traffic.

### `src/alerts`

- `AlertBanner.jsx` renders the fixed banner shown when one or more sites need attention.

### `src/common`

Small reusable UI components: `BackButton.jsx`, `LoadingSkeleton.jsx`, `SearchBar.jsx`, `StatusBadge.jsx`, `UtilizationBar.jsx`.

## Data Flow

1. `main.jsx` mounts `BrowserRouter` around `App`.
2. `App.jsx` calls `useAuth()`.
3. Unauthenticated operators see `SignInPage` for the current URL.
4. Authenticated layout reads the matched site/device params and calls `useNetworkDashboard()`.
5. `useNetworkDashboard()` loads live data through `src/services/sitesApi.js`.
6. If live requests fail, the hook sets `dataMode` to `error` and a load error message.
7. The hook normalizes and summarizes data with helpers from `src/utils/siteData.js`.
8. `AppShell` receives dashboard state for nav, search, and alerts.
9. `Routes` render `DashboardPage`, `WallEditorPage`, or `WallDisplayPage`.
10. Clicks that change screens call `navigate()` through the path helpers in `src/config/paths.js`.

## Routing

`src/config/paths.js` is the source of truth for path strings.

| Path | View |
|---|---|
| `/` | All Sites (`SitesGrid`) |
| `/sites/:siteId` | Site detail |
| `/sites/:siteId/devices/:deviceKey` | Device detail |
| `/wall` | Wall editor (control plane) |
| `/wall/display` | Wall display (kiosk presentation) |
| anything else | Redirect to `/` |

`siteId` is the existing `site_id`. `deviceKey` is the current map key (IP, else hostname) and is encoded with `encodeURIComponent`.

The URL names the screen, not wall-slot contents. Search query, selected interface, and the four circuit assignments stay out of the query string. `/wall` is the only place kiosk management belongs; `/wall/display` presents the same wall.

Refresh and browser back/forward restore the same site or device. Nginx already serves `index.html` for unknown paths (`try_files`).

## Authentication

Appliance-local PAM sessions live in:

- `src/auth` for the sign-in page
- `src/hooks/useAuth.js` for session state
- cookie + CSRF header logic in `src/services/sitesApi.js`

`BrowserRouter` sits above `SignInPage`, so a TV bookmarked at `/wall/display` returns to that path after login.
