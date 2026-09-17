# Layout & Page Structure

> **For AI agents:** This file documents how the dashboard's pages are structured, how the layout shell works, how navigation integrates with content, and how the routed views (All Sites, Site Detail, Device Detail, Wall) are composed. Use this as the reference when building new pages or extending existing ones.

---

## 1 · Overall Application Shell

**File:** `src/layout/AppShell.jsx`

The `AppShell` wraps every screen in the application. It renders three things in order:

1. `<Nav>` — fixed top navigation bar (always visible)
2. `<AlertBanner>` — fixed banner below nav (renders nothing when no alerts)
3. `<main className="page-content">` — scrollable content area

```
┌──────────────────────────────────────────────────────────────┐  ← fixed, z-index: 100
│  Nav (60px tall)                                             │
├──────────────────────────────────────────────────────────────┤  ← fixed, z-index: 99 (conditional)
│  Alert Banner (~40px tall, only when alerts exist)           │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  page-content (scrollable)                                   │
│  max-width: 1280px, centered, padded                         │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### `.app-layout`
```css
display: flex;
flex-direction: column;
min-height: 100vh;
```

### `.page-content`
```css
flex: 1;
padding: 88px 48px 48px;   /* top = nav height (60px) + spacing (28px) */
max-width: 1280px;
margin: 0 auto;
width: 100%;
```

When the alert banner is visible, `AppShell` overrides `padding-top` via inline style:
- No alert banner: `padding-top: 88px`
- Alert banner present: `padding-top: 124px` (`60px nav + ~40px banner + 24px gap`)

### Z-index Stack
| Layer | z-index |
|---|---|
| Nav | 100 |
| Alert Banner | 99 |
| Page content | default (auto) |

---

## 2 · Navigation

**File:** `src/layout/Nav.jsx`

### Positioning Model
The nav is `position: fixed`, spanning the full viewport width (`left: 0; right: 0`) at `top: 0`. Content must compensate with top padding — the `AppShell` handles this automatically via `.page-content`'s `padding-top`.

### Content Areas

```
┌────────────────────────────────────────────────────────┐
│  [logo mark]  Equate          ⌕ search    Wall  user   │
│   ←─ .nav-logo ──────────────────────────────────────→ │
│                                     ←─── .nav-right ──→│
└────────────────────────────────────────────────────────┘
```

**Left section (`.nav-logo`):**
- Logo SVG (`30×30px`) + wordmark text "Equate"
- Entire `.nav-logo` is a `<Link>` to `/` (All Sites)
- Clicking also closes search

**Wall link (`.nav-wall`):**
- Lives in `.nav-right`, before the user name
- `<Link>` to `/wall`
- Active on `/wall` and `/wall/display`
- This is how operators reach kiosk management from the webpage

**Right section (`.nav-right`):**
- Wall link, signed-in user name, and Log out
- The `·` separator uses `.nav-sep` class for color differentiation when counts are shown

### Logo Mark SVG
The `//` double-slash SVG is imported as an image asset from `assets/logo.svg`. It is rendered as `<img>` inside `.logo-mark` (not inline SVG), so it cannot be styled via CSS color properties.

---

## 3 · Page Structure: All Sites View

This is the default/home view rendered by `src/sites/SitesGrid.jsx` inside `DashboardPage`.

### Full Page Composition (top to bottom)

```
─ Nav (fixed, 60px) ──────────────────────────────────────────
─ Alert Banner (fixed, conditional) ──────────────────────────
┌─────────────────────────────────────────────────────────────┐
│  Page Header                              Updated 4:43 PM   │  ← 32px margin-bottom
│  • NETWORK DASHBOARD                                        │
│  All Sites                                                  │
├─────────────────────────────────────────────────────────────┤
│  [ Total Sites: 4 ]  [ Devices: 12 ]  [ Critical: 0 ]  ... │  ← Stat Strip, 32px mb
├─────────────────────────────────────────────────────────────┤
│  ⌕  Search by site name, type, status, or ID               │  ← SearchBar, 24px mb
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Site Card    │  │ Site Card    │  │ Site Card    │      │  ← Sites Grid
│  │              │  │              │  │              │      │    auto-fill columns
│  └──────────────┘  └──────────────┘  └──────────────┘      │    min 280px each
│                                                             │    16px gap
└─────────────────────────────────────────────────────────────┘
```

### Grid Behavior
- **Desktop (>900px):** As many columns as fit at ≥280px minimum. With 4 sites and a 1280px max-width, this typically renders as 4 columns.
- **Tablet (≤900px):** Forced to 1 column (`grid-template-columns: 1fr` via media query override)
- **Mobile (≤540px):** 1 column (inherits from tablet breakpoint)

### Empty State
If search returns no matches, the grid is replaced with a `.empty-state` block:
```
┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─┐
│ No matching sites                                      │  ← dashed border
│ Try a different search term...                         │
└ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─┘
```

---

## 4 · Page Structure: Site Detail View

This view is rendered by `src/sites/SiteDetail.jsx` when a site is selected.

### Full Page Composition (top to bottom)

```
─ Nav (fixed, 60px) ──────────────────────────────────────────
─ Alert Banner (fixed, conditional) ──────────────────────────
┌─────────────────────────────────────────────────────────────┐
│  ← All Sites                                                │  ← BackButton, 28px mb
├─────────────────────────────────────────────────────────────┤
│  • SITE DETAIL                             [Alert badge?]   │  ← 28px mb
│  District Office                                            │
│  3 devices · 3 online                                       │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐   │
│  │ IP Address │ Hostname │ Role │ Status │ CPU │ Mem ... │   │  ← Devices Table
│  ├────────────┼──────────┼──────┼────────┼─────┼──────  │   │
│  │ 10.10.0.1  │ dist-... │ Core │ ● HLTY │ ██  │ 44%    │   │
│  │ 10.10.0.2  │ dist-... │ Fire │ ● HLTY │ █░  │ 52%    │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Loading State
While `siteDetail` data is `null`, the detail page shows:
```
← All Sites

[████████████████████████ shimmer skeleton ████████████████]
```
The skeleton fills `200px` height with `border-radius: 16px`.

### Alert Badge in Detail Header
If `summary.active_alerts > 0`, a `.status-badge.alert` renders in the right column of the detail header with the text "Critical Alerts Active". It is positioned with `height: fit-content` to align at the top.

---

## 5 · URL Routing

The application uses `react-router-dom`. Path helpers live in `src/config/paths.js`. Do not concatenate route strings in components.

| Path | View rendered |
|---|---|
| `/` | `<SitesGrid>` (all-sites overview) |
| `/sites/:siteId` | `<SiteDetail>` for that site |
| `/sites/:siteId/devices/:deviceKey` | `<DeviceDetail>` for that device |
| `/wall` | `<WallEditorPage>` (kiosk control plane) |
| `/wall/display` | `<WallDisplayPage>` (TV presentation) |
| anything else | `<Navigate>` to `/` |

`deviceKey` is the existing map key (IP, else hostname). Encode it with `encodeURIComponent`.

Clicking a `<SiteCard>` navigates to `/sites/:siteId`. Clicking a device row navigates to `/sites/:siteId/devices/:deviceKey`. The Back button and Nav logo navigate to `/`. Device back navigates to that site's path.

Search remains an overlay and is not reflected in the URL. Selected interface on device detail is in-memory only. Wall slot assignments, time ranges, and the wall title must not be stored in the query string — `/wall` manages them on the page; `/wall/display` only presents.

**When adding new pages:** Add a path helper in `src/config/paths.js`, a `<Route>` in `App.jsx`, and update this table. Do not go back to in-memory screen state for navigable views.

---

## 6 · Data Flow Through the Layout

```
BrowserRouter
       │
       ▼
    App.jsx
    ├── unauthenticated → SignInPage (URL unchanged)
    └── authenticated
         ├── useMatch → selectedSite, selectedDevice
         ├── useNetworkDashboard()
         │      sites, alerts, search → AppShell → Nav, AlertBanner
         └── Routes
              ├── / → SitesGrid
              ├── /sites/:siteId → SiteDetail
              ├── /sites/:siteId/devices/:deviceKey → DeviceDetail
              ├── /wall → WallEditorPage
              └── /wall/display → WallDisplayPage
```

Dashboard data originates in `useNetworkDashboard.js`. Screen identity originates in the URL. Display components are presentational and receive data through props or outlet context. No component fetches data directly.

---

## 7 · Responsive Layout Reference

### Desktop (>900px)
```css
nav { padding: 0 48px; height: 60px; }
.page-content { padding: 88px 48px 48px; }
.alert-banner { padding: 10px 48px; }
.stat-strip { grid-template-columns: repeat(4, 1fr); }
.sites-grid { grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); }
```

### Tablet (≤900px)
```css
nav { padding: 0 20px; }
.page-content { padding: 88px 20px 40px; }
.alert-banner { padding: 10px 20px; }
.stat-strip { grid-template-columns: 1fr 1fr; }
.sites-grid { grid-template-columns: 1fr; }
.nav-right { gap: 8px; }
```

### Mobile (≤540px)
```css
nav {
  height: auto;
  min-height: 60px;
  padding-top: 10px;
  padding-bottom: 10px;
  flex-wrap: wrap;
  gap: 8px 12px;
}
.nav-right {
  width: 100%;
  margin-left: 0;
  flex-wrap: wrap;
}
.stat-strip { grid-template-columns: 1fr; }
```

---

## 8 · Maximum Content Width

All page content is constrained to `max-width: 1280px; margin: 0 auto` via `.page-content`. On very wide screens the content centers with equal space on both sides. The nav and alert banner always span the full viewport (`left: 0; right: 0`).

---

## 9 · Font Smoothing

Applied globally to the `body`:
```css
-webkit-font-smoothing: antialiased;
```

This ensures Epilogue at heavy weights renders cleanly on macOS/iOS without sub-pixel rendering artifacts.

---

## 10 · Box Sizing Reset

Applied globally:
```css
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
```

All sizing calculations include padding and border. No margin/padding defaults exist — every spacing value is intentional.

---

## 11 · Scroll Behavior

```css
html { scroll-behavior: smooth; }
```

Applied globally. The dashboard doesn't currently use anchor links, but this is in place for future use.
