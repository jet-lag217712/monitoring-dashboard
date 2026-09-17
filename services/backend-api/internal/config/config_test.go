package config

import "testing"

func TestAuthModeDefaultsToApplianceLocal(t *testing.T) {
	cfg := Config{}
	if got := cfg.AuthMode(); got != AuthModeApplianceLocal {
		t.Fatalf("AuthMode=%q, want %s", got, AuthModeApplianceLocal)
	}
}
