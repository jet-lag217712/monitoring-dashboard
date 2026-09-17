import { Link } from 'react-router-dom'
import PageHeader from '../dashboard/PageHeader.jsx'
import { paths } from '../config/paths.js'

export default function WallDisplayPage() {
  return (
    <div className="wall-page">
      <PageHeader
        eyebrow="Wall Display"
        title="Wall display"
        subtitle="Presentation view of the same wall. Management stays on the editor page."
      />
      <p className="page-sub wall-page-copy">
        This route is what a TV will load. Exit display to change slots, name, or ranges from the
        webpage.
      </p>
      <Link className="wall-action-link" to={paths.wall()}>
        Exit display
      </Link>
    </div>
  )
}
