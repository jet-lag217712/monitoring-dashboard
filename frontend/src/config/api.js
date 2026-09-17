export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8000'
export const POLL_INTERVAL_MS = 5000
export const AUTH_MODE = import.meta.env.VITE_AUTH_MODE ?? 'appliance_local'

export function isApplianceAuth() {
  return AUTH_MODE === 'appliance_local'
}

export function apiUrl(path) {
  const normalized = path.startsWith('/') ? path : `/${path}`
  if (API_BASE_URL.startsWith('http://') || API_BASE_URL.startsWith('https://')) {
    return new URL(normalized, API_BASE_URL).toString()
  }
  // Same-origin relative URL (appliance: VITE_API_BASE_URL=/api).
  return normalized
}
