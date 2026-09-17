package transform

import (
	"fmt"

	"github.com/google/uuid"
)

// Fixed OGSD namespace for deterministic UUID v5 derivation.
var ogsdNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8") // DNS namespace as base

func init() {
	// Derive a project-specific namespace from the DNS namespace + "equate-ogsd".
	ogsdNamespace = uuid.NewSHA1(ogsdNamespace, []byte("equate-ogsd"))
}

// SiteUUID returns a deterministic UUID for a collector site_id string.
func SiteUUID(siteID string) uuid.UUID {
	return uuid.NewSHA1(ogsdNamespace, []byte("site:"+siteID))
}

// DeviceUUID returns a deterministic UUID for a site+device pair.
func DeviceUUID(siteID, deviceID string) uuid.UUID {
	return uuid.NewSHA1(ogsdNamespace, []byte("device:"+siteID+"/"+deviceID))
}

// InterfaceUUID returns a deterministic UUID for a device+ifIndex pair.
func InterfaceUUID(deviceUUID uuid.UUID, ifIndex int) uuid.UUID {
	return uuid.NewSHA1(ogsdNamespace, []byte(fmt.Sprintf("interface:%s/%d", deviceUUID.String(), ifIndex)))
}
