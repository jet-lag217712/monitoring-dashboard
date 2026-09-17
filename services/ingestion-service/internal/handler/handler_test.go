package handler_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/equate/ogsd/services/ingestion-service/internal/handler"
	"github.com/equate/ogsd/services/ingestion-service/internal/metrics"
	"github.com/equate/ogsd/services/ingestion-service/internal/store"
	"github.com/equate/ogsd/services/ingestion-service/internal/transform"
	"github.com/prometheus/client_golang/prometheus"
)

type fakeStore struct {
	deviceV2Result    store.Result
	deviceV2Err       error
	ifaceV2Result     store.Result
	ifaceV2Err        error
	healthResult      store.Result
	healthErr         error
	heartbeatResult   store.Result
	heartbeatErr      error
}

func (f *fakeStore) PersistDeviceTelemetry(context.Context, transform.DeviceTelemetrySample) (store.Result, error) {
	return f.deviceV2Result, f.deviceV2Err
}

func (f *fakeStore) PersistInterfaceTelemetry(context.Context, transform.InterfaceTelemetrySample) (store.Result, error) {
	return f.ifaceV2Result, f.ifaceV2Err
}

func (f *fakeStore) PersistHealth(context.Context, transform.HealthSample) (store.Result, error) {
	return f.healthResult, f.healthErr
}

func (f *fakeStore) PersistHeartbeat(context.Context, transform.HeartbeatSample) (store.Result, error) {
	return f.heartbeatResult, f.heartbeatErr
}

func newHandler(t *testing.T, s *fakeStore) *handler.Handler {
	t.Helper()
	m := metrics.NewWithRegisterer(prometheus.NewRegistry())
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return handler.New(s, m, log)
}

func TestHandle_RejectACKs(t *testing.T) {
	h := newHandler(t, &fakeStore{})
	ack := h.Handle(context.Background(), "site/a/device/b/metric/device", []byte(`{`))
	if !ack {
		t.Fatal("invalid payload should ACK")
	}
}

func TestHandle_V2HeartbeatACKs(t *testing.T) {
	h := newHandler(t, &fakeStore{heartbeatResult: store.ResultInserted})
	payload := []byte(`{
		"schema_version":"2.0",
		"event_id":"018f3e2c-7a9d-7b20-8f63-1e2d3c4b5a63",
		"event_type":"collector_heartbeat",
		"site_id":"site-001",
		"collector_id":"collector-west-01",
		"observed_at":"2026-07-16T18:00:00Z",
		"emitted_at":"2026-07-16T18:00:01Z",
		"config_revision":"rev-1",
		"payload":{
			"hostname":"collector-host",
			"version":"unknown",
			"git_commit":"unknown",
			"build_time":"unknown",
			"uptime_seconds":10,
			"sqlite_queue_depth":0,
			"memory_usage_bytes":1000,
			"goroutine_count":5
		}
	}`)
	ack := h.Handle(context.Background(), "site/site-001/collector/collector-west-01/telemetry/v2/heartbeat", payload)
	if !ack {
		t.Fatal("v2 heartbeat insert should ACK")
	}
}

func TestHandle_V2RejectACKs(t *testing.T) {
	h := newHandler(t, &fakeStore{})
	ack := h.Handle(context.Background(), "site/site-001/device/dev-001/telemetry/v2/device", []byte(`{"schema_version":"1.0"}`))
	if !ack {
		t.Fatal("invalid v2 payload should ACK")
	}
}

func TestHandle_V2DBErrorNoACK(t *testing.T) {
	h := newHandler(t, &fakeStore{heartbeatErr: errors.New("db down")})
	payload := []byte(`{
		"schema_version":"2.0",
		"event_id":"018f3e2c-7a9d-7b20-8f63-1e2d3c4b5a63",
		"event_type":"collector_heartbeat",
		"site_id":"site-001",
		"collector_id":"collector-west-01",
		"observed_at":"2026-07-16T18:00:00Z",
		"emitted_at":"2026-07-16T18:00:01Z",
		"config_revision":"rev-1",
		"payload":{
			"hostname":"collector-host",
			"version":"unknown",
			"git_commit":"unknown",
			"build_time":"unknown",
			"uptime_seconds":10,
			"sqlite_queue_depth":0,
			"memory_usage_bytes":1000,
			"goroutine_count":5
		}
	}`)
	ack := h.Handle(context.Background(), "site/site-001/collector/collector-west-01/telemetry/v2/heartbeat", payload)
	if ack {
		t.Fatal("v2 db error must not ACK")
	}
}
