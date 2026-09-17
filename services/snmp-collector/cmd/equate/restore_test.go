package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/equate/ogsd/services/snmp-collector/internal/tui/setup"
)

func TestRunRestoreRejectsArgs(t *testing.T) {
	if code := runRestore([]string{"--scan"}); code != 2 {
		t.Fatalf("code=%d", code)
	}
}

func TestRestoreApplianceRefusesUnconfigured(t *testing.T) {
	dir := t.TempDir()
	err := restoreApplianceWithDirs(dir, filepath.Join(dir, "durable"), filepath.Join(dir, "run"))
	if err == nil {
		t.Fatal("expected not configured")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("err=%v", err)
	}
}

func TestRestoreApplianceRefusesMissingBackup(t *testing.T) {
	deployDir := configuredDeployDir(t)
	err := restoreApplianceWithDirs(deployDir, filepath.Join(t.TempDir(), "durable"), filepath.Join(t.TempDir(), "run"))
	if err == nil {
		t.Fatal("expected missing backup")
	}
	if !strings.Contains(err.Error(), "no rendered secret backup") {
		t.Fatalf("err=%v", err)
	}
}

func TestEnsureApplianceRenderedSecretsUsesBackup(t *testing.T) {
	deployDir := t.TempDir()
	durable := filepath.Join(t.TempDir(), "durable")
	runDir := filepath.Join(t.TempDir(), "run")
	if err := os.MkdirAll(durable, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(durable, "compose.env"), []byte("SNMP_COMMUNITY=CHANGE_ME\nPOSTGRES_PASSWORD=kept\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(deployDir, ".env"), []byte("SNMP_COMMUNITY=live-comm\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ensureApplianceRenderedSecretsAt(deployDir, durable, runDir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(runDir, "compose.env"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "POSTGRES_PASSWORD=kept") {
		t.Fatalf("did not restore backup: %s", got)
	}
	if !strings.Contains(got, "SNMP_COMMUNITY=live-comm") {
		t.Fatalf("did not overlay SNMP: %s", got)
	}
	if _, err := os.Stat(filepath.Join(deployDir, "scripts", "configure-vm.sh")); err == nil {
		t.Fatal("test deploy should not contain configure-vm.sh")
	}
}

func TestEnsureApplianceRenderedSecretsSkipsWhenPresent(t *testing.T) {
	deployDir := t.TempDir()
	runDir := filepath.Join(t.TempDir(), "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, "compose.env"), []byte("already=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ensureApplianceRenderedSecretsAt(deployDir, filepath.Join(t.TempDir(), "missing"), runDir); err != nil {
		t.Fatal(err)
	}
}

func configuredDeployDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	specs, err := setup.BuildSiteSpecs(setup.ProfileAppliance, 1, []string{"campus"}, []string{"10.0.0.0/24"})
	if err != nil {
		t.Fatal(err)
	}
	if err := setup.WriteManifest(dir, setup.Manifest{SiteCount: 1, Sites: specs}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.sites.generated.yml"), []byte("services: {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}
