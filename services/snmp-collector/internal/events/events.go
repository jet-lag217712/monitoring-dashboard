package events

// Event is a telemetry payload that can be published.
type Event interface {
	Topic() string
}
