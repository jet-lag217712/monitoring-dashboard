package setup

import (
	"strings"
	"testing"
)

func TestSiteUpsertSQL(t *testing.T) {
	sql := SiteUpsertSQL([]SiteSpec{
		{SiteID: "district", UpstreamSiteIDs: []string{"core"}, HubDeviceIDs: []string{"do-core"}},
		{SiteID: "and"},
	})
	if !strings.Contains(sql, "INSERT INTO sites") {
		t.Fatalf("sql=%s", sql)
	}
	if !strings.Contains(sql, SiteUUID("district").String()) {
		t.Fatalf("missing district uuid in %s", sql)
	}
	if !strings.Contains(sql, "ARRAY['core']::text[]") {
		t.Fatalf("missing upstream array in %s", sql)
	}
	if !strings.Contains(sql, "ARRAY[]::text[]") {
		t.Fatalf("missing empty array in %s", sql)
	}
}
