package validate_test

import (
	"testing"

	"github.com/equate/ogsd/services/ingestion-service/internal/validate"
)

func TestValidate_RejectsV1Topics(t *testing.T) {
	payload := []byte(`{"timestamp":"2026-06-01T18:00:00Z","metric":"uptime_seconds","value":1}`)
	topics := []string{
		"site/site-001/device/dev-001/metric/device",
		"site/site-001/device/dev-001/metric/interface",
		"bad",
	}
	for _, topic := range topics {
		if _, err := validate.Validate(topic, payload); err == nil {
			t.Fatalf("expected error for %q", topic)
		}
	}
}
