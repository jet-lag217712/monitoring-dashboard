package transform_test

import (
	"testing"

	"github.com/equate/ogsd/services/ingestion-service/internal/transform"
)

func TestUUIDV5_Stable(t *testing.T) {
	a := transform.SiteUUID("site-001")
	b := transform.SiteUUID("site-001")
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
	devA := transform.DeviceUUID("site-001", "dev-001")
	devB := transform.DeviceUUID("site-001", "dev-001")
	if devA != devB {
		t.Fatalf("%s != %s", devA, devB)
	}
	ifaceA := transform.InterfaceUUID(devA, 2)
	ifaceB := transform.InterfaceUUID(devA, 2)
	if ifaceA != ifaceB {
		t.Fatalf("%s != %s", ifaceA, ifaceB)
	}
}

func TestUUIDV5_DifferentIDsDiffer(t *testing.T) {
	if transform.SiteUUID("a") == transform.SiteUUID("b") {
		t.Fatal("site UUIDs should differ")
	}
	if transform.DeviceUUID("s", "d1") == transform.DeviceUUID("s", "d2") {
		t.Fatal("device UUIDs should differ")
	}
	dev := transform.DeviceUUID("s", "d")
	if transform.InterfaceUUID(dev, 1) == transform.InterfaceUUID(dev, 2) {
		t.Fatal("interface UUIDs should differ")
	}
}
