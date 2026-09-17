package derive

import (
	"strings"
)

// SiteLabelInput carries site identity for upstream label resolution.
type SiteLabelInput struct {
	Name     string
	Location string
}

// UpstreamSiteLabel returns the human upstream site label for a device at deviceSiteName.
// IDF sites show their campus prefix; MDF sites show the district/core site display name.
func UpstreamSiteLabel(deviceSiteName string, upstreamSiteIDs []string, sitesByName map[string]SiteLabelInput) string {
	lower := strings.ToLower(strings.TrimSpace(deviceSiteName))
	if lower == "" || isCoreSite(lower) {
		return ""
	}

	if campus := campusFromIDFSiteName(deviceSiteName); campus != "" {
		return campus
	}
	for _, upstream := range upstreamSiteIDs {
		if campus := campusFromMDFSiteName(upstream); campus != "" {
			return campus
		}
	}

	if strings.Contains(lower, "-mdf") || hasMDFUpstream(upstreamSiteIDs) {
		return districtLabel(upstreamSiteIDs, sitesByName)
	}

	return ""
}

func campusFromIDFSiteName(siteName string) string {
	lower := strings.ToLower(siteName)
	i := strings.Index(lower, "-idf")
	if i <= 0 {
		return ""
	}
	prefix := siteName[:i]
	// Require a multi-part campus id (e.g. site-a-idf1). Generic names like
	// remote-idf fall back to the upstream MDF site name instead.
	if !strings.Contains(prefix, "-") {
		return ""
	}
	return prettyPrintSiteID(prefix)
}

func campusFromMDFSiteName(mdfSiteName string) string {
	lower := strings.ToLower(strings.TrimSpace(mdfSiteName))
	if !strings.Contains(lower, "-mdf") {
		return ""
	}
	i := strings.Index(lower, "-mdf")
	if i <= 0 {
		return ""
	}
	return prettyPrintSiteID(mdfSiteName[:i])
}

func hasMDFUpstream(upstreamSiteIDs []string) bool {
	for _, upstream := range upstreamSiteIDs {
		if strings.Contains(strings.ToLower(upstream), "-mdf") {
			return true
		}
	}
	return false
}

func districtLabel(upstreamSiteIDs []string, sitesByName map[string]SiteLabelInput) string {
	coreSiteID := resolveCoreSiteID(upstreamSiteIDs, sitesByName)
	if coreSiteID == "" {
		return ""
	}
	if site, ok := sitesByName[coreSiteID]; ok {
		return coreDisplayName(site)
	}
	return coreDisplayName(SiteLabelInput{Name: coreSiteID})
}

func isCoreSite(lower string) bool {
	switch lower {
	case "do-core", "district-office", "district":
		return true
	}
	return strings.Contains(lower, "core") &&
		!strings.Contains(lower, "-mdf") &&
		!strings.Contains(lower, "-idf")
}

func resolveCoreSiteID(upstreamSiteIDs []string, sitesByName map[string]SiteLabelInput) string {
	for _, upstream := range upstreamSiteIDs {
		if isCoreSite(strings.ToLower(upstream)) {
			return upstream
		}
	}
	if len(upstreamSiteIDs) > 0 {
		return upstreamSiteIDs[0]
	}
	for name := range sitesByName {
		if isCoreSite(strings.ToLower(name)) {
			return name
		}
	}
	return ""
}

func coreDisplayName(site SiteLabelInput) string {
	if site.Location != "" {
		return site.Location
	}
	lower := strings.ToLower(site.Name)
	if isCoreSite(lower) {
		return "District-Office"
	}
	return prettyPrintSiteID(site.Name)
}

func prettyPrintSiteID(siteID string) string {
	parts := strings.Split(siteID, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, "-")
}
