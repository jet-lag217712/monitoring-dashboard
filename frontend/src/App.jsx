import { Navigate, Outlet, Route, Routes, useMatch } from 'react-router-dom'
import AppShell from './layout/AppShell.jsx'
import DashboardPage from './dashboard/DashboardPage.jsx'
import SignInPage from './auth/SignInPage.jsx'
import WallDisplayPage from './wall/WallDisplayPage.jsx'
import WallEditorPage from './wall/WallEditorPage.jsx'
import { useAuth } from './hooks/useAuth.js'
import { useNetworkDashboard } from './hooks/useNetworkDashboard.js'
import './index.css'

export default function App() {
  const auth = useAuth()

  if (!auth.isAuthenticated) {
    return <SignInPage auth={auth} />
  }

  return (
    <Routes>
      <Route element={<AuthenticatedLayout auth={auth} />}>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/sites/:siteId" element={<DashboardPage />} />
        <Route path="/sites/:siteId/devices/:deviceKey" element={<DashboardPage />} />
        <Route path="/wall" element={<WallEditorPage />} />
        <Route path="/wall/display" element={<WallDisplayPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}

function AuthenticatedLayout({ auth }) {
  const deviceMatch = useMatch('/sites/:siteId/devices/:deviceKey')
  const siteMatch = useMatch('/sites/:siteId')
  const selectedSite = deviceMatch?.params.siteId ?? siteMatch?.params.siteId ?? null
  const selectedDevice = deviceMatch?.params.deviceKey ?? null
  const dashboard = useNetworkDashboard({
    enabled: true,
    selectedSite,
    selectedDevice,
    onUnauthorized: auth.signOut,
  })

  return (
    <AppShell dashboard={dashboard} auth={auth}>
      <Outlet context={dashboard} />
    </AppShell>
  )
}

