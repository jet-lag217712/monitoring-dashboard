package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func sqlTextArray(values []string) string {
	if len(values) == 0 {
		return "ARRAY[]::text[]"
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, "'"+escapeSQLLiteral(value)+"'")
	}
	return "ARRAY[" + strings.Join(parts, ",") + "]::text[]"
}

// SiteUpsertSQL returns INSERT ... ON CONFLICT statements for configured sites.
func SiteUpsertSQL(specs []SiteSpec) string {
	var b strings.Builder
	for _, spec := range specs {
		siteID := strings.TrimSpace(spec.SiteID)
		if siteID == "" {
			continue
		}
		sid := SiteUUID(siteID).String()
		fmt.Fprintf(&b,
			"INSERT INTO sites (id, name, upstream_site_ids, hub_device_ids) VALUES ('%s', '%s', %s, %s) ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, upstream_site_ids = EXCLUDED.upstream_site_ids, hub_device_ids = EXCLUDED.hub_device_ids;\n",
			sid, escapeSQLLiteral(siteID), sqlTextArray(spec.UpstreamSiteIDs), sqlTextArray(spec.HubDeviceIDs),
		)
	}
	return b.String()
}

// SyncSiteTopology upserts sites from the appliance manifest and prunes orphans.
func SyncSiteTopology(deployDir string) error {
	manifestPath := filepath.Join(deployDir, manifestFile)
	if _, err := os.Stat(manifestPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	manifest, err := LoadManifest(deployDir)
	if err != nil {
		return err
	}
	if len(manifest.Sites) == 0 {
		return fmt.Errorf("manifest has no sites")
	}
	keep := make([]string, 0, len(manifest.Sites))
	for _, spec := range manifest.Sites {
		keep = append(keep, spec.SiteID)
	}
	sql := SiteUpsertSQL(manifest.Sites) + OrphanSitesPruneSQL(keep)
	if err := ExecPostgresSQL(deployDir, sql); err != nil {
		return fmt.Errorf("site topology sync: %w", err)
	}
	return nil
}
