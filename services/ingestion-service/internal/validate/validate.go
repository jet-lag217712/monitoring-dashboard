package validate

import (
	"fmt"
	"strings"
	"time"
)

// Kind identifies the metric message type from the topic path.
type Kind string

// Message is the result of validating a topic + payload.
type Message struct {
	Kind        Kind
	DeviceV2    *DeviceTelemetryV2
	InterfaceV2 *InterfaceTelemetryV2
	HealthV2    *HealthStateV2
	HeartbeatV2 *HeartbeatV2
}

// Validate parses and validates a v2 MQTT topic and JSON payload.
func Validate(topic string, payload []byte) (Message, error) {
	return ValidateV2(topic, payload)
}

func parseTimestamp(raw string) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, fmt.Errorf("timestamp is required")
	}
	ts, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		// Also accept RFC3339Nano from collectors.
		ts, err = time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid timestamp: %w", err)
		}
	}
	return ts.UTC(), nil
}
