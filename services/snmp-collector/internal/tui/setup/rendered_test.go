package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPatchEnvFileKeysOverlaysSNMP(t *testing.T) {
	dir := t.TempDir()
	compose := filepath.Join(dir, "compose.env")
	if err := os.WriteFile(compose, []byte("POSTGRES_PASSWORD=keep\nSNMP_COMMUNITY=CHANGE_ME\nMQTT_BROKER=tls://mosquitto:8883\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := PatchEnvFileKeys(compose, map[string]string{
		"SNMP_COMMUNITY":           "lab-read",
		"SNMP_DISCOVERY_COMMUNITY": "lab-read",
	}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(compose)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "POSTGRES_PASSWORD=keep") {
		t.Fatalf("lost postgres password: %s", got)
	}
	if !strings.Contains(got, "SNMP_COMMUNITY=lab-read") {
		t.Fatalf("community not patched: %s", got)
	}
	if strings.Contains(got, "CHANGE_ME") {
		t.Fatalf("CHANGE_ME remains: %s", got)
	}
}

func TestOverlaySNMPFromEnvFile(t *testing.T) {
	dir := t.TempDir()
	runDir := filepath.Join(dir, "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dotEnv := filepath.Join(dir, ".env")
	if err := os.WriteFile(dotEnv, []byte("SNMP_COMMUNITY=site-comm\nSNMP_DISCOVERY_COMMUNITY=disc-comm\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "compose.env"), []byte("SNMP_COMMUNITY=CHANGE_ME\nMQTT_PASSWORD=secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "collector.env"), []byte("SNMP_COMMUNITY=CHANGE_ME\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := OverlaySNMPFromEnvFile(dotEnv, runDir); err != nil {
		t.Fatal(err)
	}
	compose, err := os.ReadFile(filepath.Join(runDir, "compose.env"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(compose), "SNMP_COMMUNITY=site-comm") || !strings.Contains(string(compose), "SNMP_DISCOVERY_COMMUNITY=disc-comm") {
		t.Fatalf("compose.env=%s", compose)
	}
	if !strings.Contains(string(compose), "MQTT_PASSWORD=secret") {
		t.Fatalf("lost mqtt password: %s", compose)
	}
	collector, err := os.ReadFile(filepath.Join(runDir, "collector.env"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(collector), "SNMP_COMMUNITY=site-comm") {
		t.Fatalf("collector.env=%s", collector)
	}
}

func TestRestoreRenderedDirMissingBackup(t *testing.T) {
	dir := t.TempDir()
	err := RestoreRenderedDir(filepath.Join(dir, "durable"), filepath.Join(dir, "run"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "no rendered secret backup") {
		t.Fatalf("err=%v", err)
	}
}

func TestRestoreRenderedDirCopiesBackup(t *testing.T) {
	dir := t.TempDir()
	durable := filepath.Join(dir, "durable")
	runDir := filepath.Join(dir, "run")
	if err := os.MkdirAll(durable, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(durable, "compose.env"), []byte("POSTGRES_PASSWORD=kept\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := RestoreRenderedDir(durable, runDir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(runDir, "compose.env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "POSTGRES_PASSWORD=kept\n" {
		t.Fatalf("got %q", data)
	}
}

func TestRequireConfigured(t *testing.T) {
	dir := t.TempDir()
	if err := RequireConfigured(dir); err == nil {
		t.Fatal("expected not configured")
	}
	specs, err := BuildSiteSpecs(ProfileAppliance, 1, []string{"campus"}, []string{"10.0.0.0/24"})
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteManifest(dir, Manifest{SiteCount: 1, Sites: specs}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(generatedComposePath(dir), []byte("services: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := RequireConfigured(dir); err != nil {
		t.Fatal(err)
	}
}
