import { useCallback, useEffect, useRef, useState } from 'react'
import { apiUrl } from '../config/api.js'
import { setCsrfTokenProvider } from '../services/sitesApi.js'

function profileFromApplianceUser(username) {
  return {
    email: '',
    name: username,
    picture: null,
    sub: username,
  }
}

async function fetchApplianceSession() {
  const sessionUrl = apiUrl('/auth/me')
  const res = await fetch(sessionUrl, {
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  if (res.status === 401) {
    return null
  }
  if (!res.ok) {
    throw new Error('Failed to load session')
  }
  const data = await res.json()
  return {
    user: profileFromApplianceUser(data.username),
    csrfToken: data.csrf_token,
  }
}

/**
 * Authentication hook for appliance-local cookie sessions.
 */
export function useAuth() {
  const [status, setStatus] = useState('loading')
  const [user, setUser] = useState(null)
  const [error, setError] = useState(null)
  const csrfRef = useRef(null)

  const clearSession = useCallback(() => {
    csrfRef.current = null
    setCsrfTokenProvider(() => null)
    setUser(null)
    setStatus('signed_out')
  }, [])

  const applyApplianceSession = useCallback(session => {
    if (!session?.user) {
      clearSession()
      return false
    }
    csrfRef.current = session.csrfToken
    setCsrfTokenProvider(() => csrfRef.current)
    setUser(session.user)
    setStatus('signed_in')
    setError(null)
    return true
  }, [clearSession])

  useEffect(() => {
    let cancelled = false

    async function initAppliance() {
      try {
        const session = await fetchApplianceSession()
        if (cancelled) return
        if (session) {
          applyApplianceSession(session)
        } else {
          setStatus('signed_out')
        }
      } catch (err) {
        if (!cancelled) {
          setError(err.message ?? 'Failed to initialize sign-in')
          setStatus('error')
        }
      }
    }

    initAppliance()
    return () => {
      cancelled = true
      setCsrfTokenProvider(() => null)
    }
  }, [applyApplianceSession])

  const signIn = useCallback(async (username, password) => {
    setError(null)
    const res = await fetch(apiUrl('/auth/login'), {
      method: 'POST',
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
    if (res.status === 401 || res.status === 429) {
      const payload = await res.json().catch(() => null)
      const message = payload?.error?.message ?? 'Invalid username or password'
      setError(message)
      setStatus('signed_out')
      return false
    }
    if (!res.ok) {
      setError('Sign-in is temporarily unavailable')
      setStatus('error')
      return false
    }
    const data = await res.json()
    return applyApplianceSession({
      user: profileFromApplianceUser(data.username),
      csrfToken: data.csrf_token,
    })
  }, [applyApplianceSession])

  const signOut = useCallback(async () => {
    const csrf = csrfRef.current
    try {
      await fetch(apiUrl('/auth/logout'), {
        method: 'POST',
        credentials: 'include',
        headers: {
          Accept: 'application/json',
          ...(csrf ? { 'X-CSRF-Token': csrf } : {}),
        },
      })
    } catch {
      // Best-effort logout; clear local state regardless.
    }
    clearSession()
  }, [clearSession])

  return {
    error,
    signIn,
    signOut,
    status,
    user,
    isAuthenticated: status === 'signed_in',
    isConfigured: true,
    isAppliance: true,
  }
}
