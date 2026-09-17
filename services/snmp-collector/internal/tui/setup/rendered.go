package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// RunRenderedDir is the ephemeral rendered-secret directory (tmpfs).
	RunRenderedDir = "/run/equate/rendered"
	// DurableRenderedDir is the reboot-surviving copy of rendered secrets.
	DurableRenderedDir = "/var/lib/equate/rendered"
)

// PersistRenderedSecrets copies /run/equate/rendered to /var/lib/equate/rendered.
func PersistRenderedSecrets() error {
	return PersistRenderedDir(RunRenderedDir, DurableRenderedDir)
}

// PersistRenderedDir copies src rendered secrets over dest, replacing dest.
func PersistRenderedDir(src, dest string) error {
	if _, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("persist rendered secrets: %w", err)
	}
	if err := copyDirReplace(src, dest); err != nil {
		return fmt.Errorf("persist rendered secrets: %w", err)
	}
	return nil
}

// RestoreRenderedSecrets copies the durable backup onto /run/equate/rendered
// when compose.env is missing.
func RestoreRenderedSecrets() error {
	return RestoreRenderedDir(DurableRenderedDir, RunRenderedDir)
}

// RestoreRenderedDir copies durableDir onto runDir when runDir/compose.env is absent.
func RestoreRenderedDir(durableDir, runDir string) error {
	if _, err := os.Stat(filepath.Join(runDir, "compose.env")); err == nil {
		return nil
	}
	if _, err := os.Stat(filepath.Join(durableDir, "compose.env")); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no rendered secret backup; run equate configure once, then restore will persist it")
		}
		return fmt.Errorf("rendered secret backup: %w", err)
	}
	if err := copyDirReplace(durableDir, runDir); err != nil {
		return fmt.Errorf("restore rendered secrets: %w", err)
	}
	return nil
}

// OverlaySNMPFromEnvFile copies SNMP_COMMUNITY and SNMP_DISCOVERY_COMMUNITY from
// deploy .env into compose.env and collector.env under runDir.
func OverlaySNMPFromEnvFile(dotEnvPath, runDir string) error {
	src, err := loadEnvMap(dotEnvPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("overlay SNMP: read %s: %w", dotEnvPath, err)
	}
	updates := snmpUpdatesFromMap(src)
	if len(updates) == 0 {
		return nil
	}
	for _, name := range []string{"compose.env", "collector.env"} {
		path := filepath.Join(runDir, name)
		if err := PatchEnvFileKeys(path, updates); err != nil {
			return err
		}
	}
	return nil
}

func snmpUpdatesFromMap(src map[string]string) map[string]string {
	updates := make(map[string]string, 2)
	if v := strings.TrimSpace(src["SNMP_COMMUNITY"]); v != "" && v != "CHANGE_ME" {
		updates["SNMP_COMMUNITY"] = v
	}
	if v := strings.TrimSpace(src["SNMP_DISCOVERY_COMMUNITY"]); v != "" && v != "CHANGE_ME" {
		updates["SNMP_DISCOVERY_COMMUNITY"] = v
	}
	if _, ok := updates["SNMP_DISCOVERY_COMMUNITY"]; !ok {
		if v := updates["SNMP_COMMUNITY"]; v != "" {
			updates["SNMP_DISCOVERY_COMMUNITY"] = v
		}
	}
	return updates
}

// PatchEnvFileKeys updates or appends keys in an env file. Missing files are skipped.
func PatchEnvFileKeys(path string, updates map[string]string) error {
	if len(updates) == 0 {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("patch env %s: %w", path, err)
	}
	seen := make(map[string]bool, len(updates))
	var b strings.Builder
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			key, _, ok := strings.Cut(trimmed, "=")
			key = strings.TrimSpace(key)
			if ok {
				if val, hit := updates[key]; hit && strings.TrimSpace(val) != "" {
					fmt.Fprintf(&b, "%s=%s\n", key, val)
					seen[key] = true
					continue
				}
			}
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	for key, val := range updates {
		if seen[key] || strings.TrimSpace(val) == "" {
			continue
		}
		fmt.Fprintf(&b, "%s=%s\n", key, val)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return fmt.Errorf("patch env %s: %w", path, err)
	}
	return nil
}

// RequireConfigured checks that site artifacts exist so restore is not first boot.
func RequireConfigured(deployDir string) error {
	if _, err := os.Stat(generatedComposePath(deployDir)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("appliance is not configured (missing %s)", generatedComposeFile)
		}
		return err
	}
	if _, err := LoadManifest(deployDir); err != nil {
		return fmt.Errorf("appliance is not configured: %w", err)
	}
	return nil
}

// WaitForSiteCollectors waits until each site collector answers status.summary or /healthz.
func WaitForSiteCollectors(deployDir string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 3 * time.Minute
	}
	manifest, err := LoadManifest(deployDir)
	if err != nil {
		return err
	}
	for _, spec := range manifest.Sites {
		client := newDeployControl(deployDir, spec)
		if err := waitForCollector(spec.AdminURL(), client, timeout); err != nil {
			return fmt.Errorf("%s: %w", spec.SiteID, err)
		}
	}
	return nil
}

func persistRenderedSNMPFromDotEnv(dotEnvPath string) error {
	if err := OverlaySNMPFromEnvFile(dotEnvPath, RunRenderedDir); err != nil {
		return err
	}
	return PersistRenderedSecrets()
}

func copyDirReplace(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp := dest + ".copying"
	if err := os.RemoveAll(tmp); err != nil {
		return err
	}
	cmd := exec.Command("cp", "-a", src, tmp)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(tmp)
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("copy %s: %s", src, msg)
	}
	if err := os.RemoveAll(dest); err != nil {
		_ = os.RemoveAll(tmp)
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.RemoveAll(tmp)
		return err
	}
	return nil
}
