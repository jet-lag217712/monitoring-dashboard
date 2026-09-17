# Frontend Architecture

This frontend is a Vite + React single-page dashboard. It does not use a routing library yet. Instead, the current screen is controlled by React state: the app shows the sites overview until a site is selected, then it shows that site's detail page.

## Where to Start

Start with `src/App.jsx`. It is intentionally small:

- It calls `useAuth()` for appliance-local PAM sessions.
- It calls `useNetworkDashboard()` to get all dashboard state and actions.
- It renders `SignInPage` until the operator is authenticated.
- It renders `AppShell` for the page frame and `DashboardPage` for the current dashboard screen.

After that, read `src/hooks/useNetworkDashboard.js`. That hook loads live site data, tracks search text and the selected site, and polls for updates. API failures surface as an error state; there is no mock fallback.

## Folder Structure

### `src/config`

- `api.js` contains the API base URL, polling interval, auth mode, and URL helper.

### `src/auth`

- `SignInPage.jsx` is the username/password form for appliance-local sessions.

### `src/services`

API functions live here. `sitesApi.js` contains small functions for backend requests using cookie sessions and CSRF:

- `fetchSitesFromApi()`
- `fetchSiteDetailFromApi(siteId)`

Future API calls should be added to this folder first. Components should not call `fetch()` directly.

### `src/hooks`

- `useAuth.js` owns the PAM cookie session.
- `useNetworkDashboard.js` owns site list, alerts, selection, search, last updated time, live/error data mode, and polling.

### `src/utils`

Pure helper functions live here. `siteData.js` transforms or summarizes API site responses.

### `src/layout`

- `AppShell.jsx` renders the fixed navigation, alert banner, and main content area.
- `Nav.jsx` renders the top navigation bar.

### `src/dashboard`

- `DashboardPage.jsx` decides whether to show the sites overview or the selected site detail.

### `src/sites`

Site-focused UI: `SitesGrid.jsx`, `SiteCard.jsx`, `SiteDetail.jsx`.

### `src/devices`

Device-focused UI: `DeviceRow.jsx`.

### `src/tables`

- `DevicesTable.jsx` owns the table structure for device detail rows.

### `src/charts`

Intended home for future chart components.

### `src/alerts`

- `AlertBanner.jsx` renders the fixed banner shown when one or more sites need attention.

### `src/common`

Small reusable UI components: `BackButton.jsx`, `LoadingSkeleton.jsx`, `SearchBar.jsx`, `StatusBadge.jsx`, `UtilizationBar.jsx`.

## Data Flow

1. `App.jsx` calls `useAuth()` and `useNetworkDashboard()`.
2. Unauthenticated operators see `SignInPage`.
3. `useNetworkDashboard()` loads live data through `src/services/sitesApi.js`.
4. If live requests fail, the hook sets `dataMode` to `error` and a load error message.
5. The hook normalizes and summarizes data with helpers from `src/utils/siteData.js`.
6. `App.jsx` passes the dashboard state into `AppShell` and `DashboardPage`.
7. `DashboardPage` renders either `SitesGrid` or `SiteDetail`.
8. User actions, such as searching or clicking a site, call handlers from `useNetworkDashboard()`.

## Routing

There is no URL-based routing yet. The app uses `selectedSite` state:

- `selectedSite === null` means the all-sites overview is visible.
- `selectedSite !== null` means the site detail screen is visible.

## Authentication

Appliance-local PAM sessions live in:

- `src/auth` for the sign-in page
- `src/hooks/useAuth.js` for session state
- cookie + CSRF header logic in `src/services/sitesApi.js`
