import { Link } from 'react-router-dom'
import PageHeader from '../dashboard/PageHeader.jsx'
import { paths } from '../config/paths.js'

export default function WallEditorPage() {
  return (
    <div className="wall-page">
      <PageHeader
        eyebrow="Wall"
        title="Wall"
        subtitle="Assign the four circuits here. This is the control plane for the TV display."
      />
      <p className="page-sub wall-page-copy">
        Charts and slot pickers land in a later slice. Open display to see the kiosk route the TV
        will bookmark.
      </p>
      <Link className="wall-action-link" to={paths.wallDisplay()}>
        Open display
      </Link>
    </div>
  )
}
