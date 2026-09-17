import { Link, useLocation } from 'react-router-dom'
import logoUrl from '../../assets/logo.svg'
import SearchBar from '../common/SearchBar.jsx'
import { isWallPath, paths } from '../config/paths.js'

export default function Nav({
  onLogoClick,
  onWallClick,
  user,
  onSignOut,
  searchQuery,
  onSearchQueryChange,
  onSearchFocus,
  onSearchClear,
}) {
  const location = useLocation()
  const wallActive = isWallPath(location.pathname)

  return (
    <nav className="app-nav">
      <Link to={paths.home()} className="nav-logo" onClick={onLogoClick}>
        <span className="logo-mark">
          <img src={logoUrl} alt="Equate Logo" />
        </span>
        Equate
      </Link>

      <SearchBar
        variant="nav"
        value={searchQuery}
        onChange={onSearchQueryChange}
        onFocus={onSearchFocus}
        onClear={onSearchClear}
        placeholder="Search sites, hostname, or IP…"
        id="nav-search"
      />

      <div className="nav-right">
        <Link
          to={paths.wall()}
          className={`nav-wall${wallActive ? ' is-active' : ''}`}
          onClick={onWallClick}
        >
          Wall
        </Link>
        {user && (
          <>
            {(user.name || user.email) && (
              <span className="nav-user" title={user.email || user.name}>
                {user.name || user.email}
              </span>
            )}
            <button type="button" className="nav-sign-out" onClick={onSignOut}>
              Log out
            </button>
          </>
        )}
      </div>
    </nav>
  )
}
