package core

import (
	"context"
	"testing"

	"github.com/gosnmp/gosnmp"
)

type fakeWalker struct {
	columns map[string][]gosnmp.SnmpPDU
}

func (f *fakeWalker) Walk(_ context.Context, rootOID string, walkFn gosnmp.WalkFunc) error {
	for _, pdu := range f.columns[rootOID] {
		if err := walkFn(pdu); err != nil {
			return err
		}
	}
	return nil
}

func TestParseInterfaceWalk(t *testing.T) {
	t.Parallel()

	pdus := []gosnmp.SnmpPDU{
		{Name: OIDIfHCInOctets + ".1", Type: gosnmp.Counter64, Value: uint64(100)},
		{Name: OIDIfHCInOctets + ".2", Type: gosnmp.Counter64, Value: uint64(200)},
	}
	got, err := ParseInterfaceWalk(OIDIfHCInOctets, pdus)
	if err != nil {
		t.Fatalf("ParseInterfaceWalk: %v", err)
	}
	if got[1] != 100 || got[2] != 200 {
		t.Fatalf("unexpected map: %#v", got)
	}
}

func TestPollInterfacesHC(t *testing.T) {
	t.Parallel()

	w := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex: {
			{Name: OIDIfIndex + ".1", Type: gosnmp.Integer, Value: 1},
			{Name: OIDIfIndex + ".2", Type: gosnmp.Integer, Value: 2},
		},
		OIDIfHCInOctets: {
			{Name: OIDIfHCInOctets + ".1", Type: gosnmp.Counter64, Value: uint64(10)},
			{Name: OIDIfHCInOctets + ".2", Type: gosnmp.Counter64, Value: uint64(20)},
		},
		OIDIfHCOutOctets: {
			{Name: OIDIfHCOutOctets + ".1", Type: gosnmp.Counter64, Value: uint64(11)},
			{Name: OIDIfHCOutOctets + ".2", Type: gosnmp.Counter64, Value: uint64(21)},
		},
		OIDIfInErrors: {
			{Name: OIDIfInErrors + ".1", Type: gosnmp.Counter32, Value: uint32(1)},
			{Name: OIDIfInErrors + ".2", Type: gosnmp.Counter32, Value: uint32(2)},
		},
		OIDIfOutErrors: {
			{Name: OIDIfOutErrors + ".1", Type: gosnmp.Counter32, Value: uint32(3)},
			{Name: OIDIfOutErrors + ".2", Type: gosnmp.Counter32, Value: uint32(4)},
		},
	}}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 2 {
		t.Fatalf("got %d readings, want 2", len(readings))
	}
	if readings[0].IfIndex != 1 || readings[0].InOctets != 10 || readings[0].OutOctets != 11 {
		t.Fatalf("reading[0]=%#v", readings[0])
	}
	if readings[1].IfIndex != 2 || readings[1].InErrors != 2 || readings[1].OutErrors != 4 {
		t.Fatalf("reading[1]=%#v", readings[1])
	}
}

func TestPollInterfacesFallback32Bit(t *testing.T) {
	t.Parallel()

	w := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex: {
			{Name: OIDIfIndex + ".5", Type: gosnmp.Integer, Value: 5},
		},
		OIDIfHCInOctets:  {},
		OIDIfHCOutOctets: {},
		OIDIfInOctets: {
			{Name: OIDIfInOctets + ".5", Type: gosnmp.Counter32, Value: uint32(99)},
		},
		OIDIfOutOctets: {
			{Name: OIDIfOutOctets + ".5", Type: gosnmp.Counter32, Value: uint32(88)},
		},
		OIDIfInErrors: {
			{Name: OIDIfInErrors + ".5", Type: gosnmp.Counter32, Value: uint32(0)},
		},
		OIDIfOutErrors: {
			{Name: OIDIfOutErrors + ".5", Type: gosnmp.Counter32, Value: uint32(0)},
		},
	}}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("got %d readings, want 1", len(readings))
	}
	if readings[0].InOctets != 99 || readings[0].OutOctets != 88 {
		t.Fatalf("unexpected fallback values: %#v", readings[0])
	}
}

func TestPollInterfacesInventoryAndLastChange(t *testing.T) {
	t.Parallel()

	w := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex:       {{Name: OIDIfIndex + ".7", Type: gosnmp.Integer, Value: 7}},
		OIDIfDescr:       {{Name: OIDIfDescr + ".7", Type: gosnmp.OctetString, Value: []byte("Ethernet7")}},
		OIDIfName:        {{Name: OIDIfName + ".7", Type: gosnmp.OctetString, Value: "Et7"}},
		OIDIfAlias:       {{Name: OIDIfAlias + ".7", Type: gosnmp.OctetString, Value: "uplink"}},
		OIDIfType:        {{Name: OIDIfType + ".7", Type: gosnmp.Integer, Value: 6}},
		OIDIfSpeed:       {{Name: OIDIfSpeed + ".7", Type: gosnmp.Gauge32, Value: uint32(1_000_000_000)}},
		OIDIfHighSpeed:   {{Name: OIDIfHighSpeed + ".7", Type: gosnmp.Gauge32, Value: uint32(10_000)}},
		OIDIfAdminStatus: {{Name: OIDIfAdminStatus + ".7", Type: gosnmp.Integer, Value: 1}},
		OIDIfOperStatus:  {{Name: OIDIfOperStatus + ".7", Type: gosnmp.Integer, Value: 2}},
		OIDIfLastChange:  {{Name: OIDIfLastChange + ".7", Type: gosnmp.TimeTicks, Value: uint32(1234)}},
		OIDIfHCInOctets:  {{Name: OIDIfHCInOctets + ".7", Type: gosnmp.Counter64, Value: uint64(10)}},
		OIDIfHCOutOctets: {{Name: OIDIfHCOutOctets + ".7", Type: gosnmp.Counter64, Value: uint64(20)}},
		OIDIfInErrors:    {{Name: OIDIfInErrors + ".7", Type: gosnmp.Counter32, Value: uint32(1)}},
		OIDIfOutErrors:   {{Name: OIDIfOutErrors + ".7", Type: gosnmp.Counter32, Value: uint32(2)}},
	}}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("got %d readings, want 1", len(readings))
	}
	got := readings[0]
	if got.IfDescr != "Ethernet7" || got.IfName != "Et7" || got.IfAlias != "uplink" {
		t.Fatalf("identity fields=%#v", got)
	}
	if got.IfTypeName != "ethernetcsmacd" || got.SpeedBPS != 10_000_000_000 {
		t.Fatalf("type/speed=%#v", got)
	}
	if got.AdminStatus != "up" || got.OperStatus != "down" || got.LastChangeSeconds != 12.34 {
		t.Fatalf("status/last change=%#v", got)
	}
	if !got.HasCounters {
		t.Fatal("expected counters to be present")
	}
}

func TestDuplexStatusName(t *testing.T) {
	t.Parallel()
	if DuplexStatusName(2) != "half" || DuplexStatusName(3) != "full" || DuplexStatusName(1) != "unknown" {
		t.Fatalf("unexpected duplex names")
	}
}

func TestPollInterfacesDuplexAndPackets(t *testing.T) {
	t.Parallel()

	w := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex: {{Name: OIDIfIndex + ".1", Type: gosnmp.Integer, Value: 1}},
		OIDDot3StatsDuplexStatus: {
			{Name: OIDDot3StatsDuplexStatus + ".1", Type: gosnmp.Integer, Value: 3},
		},
		OIDIfHCInOctets:  {{Name: OIDIfHCInOctets + ".1", Type: gosnmp.Counter64, Value: uint64(10)}},
		OIDIfHCOutOctets: {{Name: OIDIfHCOutOctets + ".1", Type: gosnmp.Counter64, Value: uint64(20)}},
		OIDIfHCInUcastPkts: {
			{Name: OIDIfHCInUcastPkts + ".1", Type: gosnmp.Counter64, Value: uint64(100)},
		},
		OIDIfHCOutUcastPkts: {
			{Name: OIDIfHCOutUcastPkts + ".1", Type: gosnmp.Counter64, Value: uint64(200)},
		},
		OIDIfInErrors:  {{Name: OIDIfInErrors + ".1", Type: gosnmp.Counter32, Value: uint32(0)}},
		OIDIfOutErrors: {{Name: OIDIfOutErrors + ".1", Type: gosnmp.Counter32, Value: uint32(0)}},
	}}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("got %d readings, want 1", len(readings))
	}
	got := readings[0]
	if got.Duplex == nil || *got.Duplex != "full" {
		t.Fatalf("duplex=%v want full", got.Duplex)
	}
	if !got.HasPackets || got.InPackets != 100 || got.OutPackets != 200 {
		t.Fatalf("packets=%#v", got)
	}
}

func TestPollInterfacesPacketFallback32Bit(t *testing.T) {
	t.Parallel()

	w := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex:       {{Name: OIDIfIndex + ".5", Type: gosnmp.Integer, Value: 5}},
		OIDIfHCInOctets:  {{Name: OIDIfHCInOctets + ".5", Type: gosnmp.Counter64, Value: uint64(1)}},
		OIDIfHCOutOctets: {{Name: OIDIfHCOutOctets + ".5", Type: gosnmp.Counter64, Value: uint64(2)}},
		OIDIfHCInUcastPkts:  {},
		OIDIfHCOutUcastPkts: {},
		OIDIfInUcastPkts: {
			{Name: OIDIfInUcastPkts + ".5", Type: gosnmp.Counter32, Value: uint32(55)},
		},
		OIDIfOutUcastPkts: {
			{Name: OIDIfOutUcastPkts + ".5", Type: gosnmp.Counter32, Value: uint32(66)},
		},
		OIDIfInErrors:  {{Name: OIDIfInErrors + ".5", Type: gosnmp.Counter32, Value: uint32(0)}},
		OIDIfOutErrors: {{Name: OIDIfOutErrors + ".5", Type: gosnmp.Counter32, Value: uint32(0)}},
	}}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("got %d readings, want 1", len(readings))
	}
	got := readings[0]
	if !got.HasPackets || got.InPackets != 55 || got.OutPackets != 66 {
		t.Fatalf("packet fallback=%#v", got)
	}
}

type failingDuplexWalker struct {
	fakeWalker
	failDuplex bool
}

func (f *failingDuplexWalker) Walk(ctx context.Context, rootOID string, walkFn gosnmp.WalkFunc) error {
	if f.failDuplex && rootOID == OIDDot3StatsDuplexStatus {
		return context.Canceled
	}
	return f.fakeWalker.Walk(ctx, rootOID, walkFn)
}

func TestPollInterfacesDuplexWalkFailureIsBestEffort(t *testing.T) {
	t.Parallel()

	base := &fakeWalker{columns: map[string][]gosnmp.SnmpPDU{
		OIDIfIndex:       {{Name: OIDIfIndex + ".1", Type: gosnmp.Integer, Value: 1}},
		OIDIfHCInOctets:  {{Name: OIDIfHCInOctets + ".1", Type: gosnmp.Counter64, Value: uint64(10)}},
		OIDIfHCOutOctets: {{Name: OIDIfHCOutOctets + ".1", Type: gosnmp.Counter64, Value: uint64(20)}},
		OIDIfHCInUcastPkts: {
			{Name: OIDIfHCInUcastPkts + ".1", Type: gosnmp.Counter64, Value: uint64(1)},
		},
		OIDIfHCOutUcastPkts: {
			{Name: OIDIfHCOutUcastPkts + ".1", Type: gosnmp.Counter64, Value: uint64(2)},
		},
		OIDIfInErrors:  {{Name: OIDIfInErrors + ".1", Type: gosnmp.Counter32, Value: uint32(0)}},
		OIDIfOutErrors: {{Name: OIDIfOutErrors + ".1", Type: gosnmp.Counter32, Value: uint32(0)}},
	}}
	w := &failingDuplexWalker{fakeWalker: *base, failDuplex: true}

	readings, err := PollInterfaces(context.Background(), w)
	if err != nil {
		t.Fatalf("PollInterfaces: %v", err)
	}
	if len(readings) != 1 {
		t.Fatalf("got %d readings, want 1", len(readings))
	}
	if readings[0].Duplex != nil {
		t.Fatalf("duplex=%v want nil on walk failure", readings[0].Duplex)
	}
}
