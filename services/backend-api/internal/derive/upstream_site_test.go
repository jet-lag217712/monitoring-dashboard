package derive

import "testing"

func TestUpstreamSiteLabel_IDF(t *testing.T) {
	got := UpstreamSiteLabel("site-a-idf1", nil, nil)
	if got != "Site-A" {
		t.Fatalf("UpstreamSiteLabel=%q want Site-A", got)
	}
}

func TestUpstreamSiteLabel_IDFUppercase(t *testing.T) {
	got := UpstreamSiteLabel("SITE-A-IDF1", nil, nil)
	if got != "Site-A" {
		t.Fatalf("UpstreamSiteLabel=%q want Site-A", got)
	}
}

func TestUpstreamSiteLabel_IDFFromUpstreamMDF(t *testing.T) {
	got := UpstreamSiteLabel("remote-idf", []string{"site-a-mdf"}, nil)
	if got != "Site-A" {
		t.Fatalf("UpstreamSiteLabel=%q want Site-A", got)
	}
}

func TestUpstreamSiteLabel_MDFWithCoreLocation(t *testing.T) {
	sites := map[string]SiteLabelInput{
		"do-core": {Name: "do-core", Location: "District Office"},
	}
	got := UpstreamSiteLabel("site-a-mdf", []string{"do-core"}, sites)
	if got != "District Office" {
		t.Fatalf("UpstreamSiteLabel=%q want District Office", got)
	}
}

func TestUpstreamSiteLabel_MDFWithoutCoreLocation(t *testing.T) {
	sites := map[string]SiteLabelInput{
		"do-core": {Name: "do-core"},
	}
	got := UpstreamSiteLabel("site-a-mdf", []string{"do-core"}, sites)
	if got != "District-Office" {
		t.Fatalf("UpstreamSiteLabel=%q want District-Office", got)
	}
}

func TestUpstreamSiteLabel_MDFWithDistrictOfficeUpstream(t *testing.T) {
	sites := map[string]SiteLabelInput{
		"district-office": {Name: "district-office", Location: "District Office"},
	}
	got := UpstreamSiteLabel("site-a-mdf", []string{"district-office"}, sites)
	if got != "District Office" {
		t.Fatalf("UpstreamSiteLabel=%q want District Office", got)
	}
}

func TestUpstreamSiteLabel_Core(t *testing.T) {
	got := UpstreamSiteLabel("do-core", nil, nil)
	if got != "" {
		t.Fatalf("UpstreamSiteLabel=%q want empty", got)
	}
}
