package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/equate/ogsd/services/snmp-collector/internal/tui/setup"
)

func resolveDeployDir() (string, error) {
	if v := os.Getenv("EQUATE_DEPLOY_DIR"); v != "" {
		return filepath.Abs(v)
	}
	const deployDirFile = "/etc/equate/deploy-dir"
	if data, err := os.ReadFile(deployDirFile); err == nil {
		if dir := strings.TrimSpace(string(data)); dir != "" && dir != "." {
			return filepath.Abs(dir)
		}
	}
	candidates := []string{
		"/opt/equate/current",
		"/opt/equate/releases/current",
	}
	for _, dir := range candidates {
		if _, err := os.Stat(dir); err == nil {
			return filepath.Abs(dir)
		}
	}
	return "", fmt.Errorf("deploy directory not found (set EQUATE_DEPLOY_DIR or /etc/equate/deploy-dir)")
}

func runDockerCompose(deployDir string, args ...string) error {
	files := dockerComposeFiles(deployDir)
	cmdArgs := []string{"compose"}
	const composeEnv = "/run/equate/rendered/compose.env"
	if _, err := os.Stat(composeEnv); err == nil {
		cmdArgs = append(cmdArgs, "--env-file", composeEnv)
	}
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			cmdArgs = append(cmdArgs, "-f", f)
		}
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Dir = deployDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func dockerComposeFiles(deployDir string) []string {
	return []string{
		filepath.Join(deployDir, "docker-compose.yml"),
		filepath.Join(deployDir, "docker-compose.sites.generated.yml"),
	}
}

func runSyncDBRoles(args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "usage: equate sync-db-roles")
		return 2
	}
	deployDir, err := resolveDeployDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync-db-roles: %v\n", err)
		return 1
	}
	if err := runSyncDBRolePasswords(deployDir); err != nil {
		fmt.Fprintf(os.Stderr, "sync-db-roles: %v\n", err)
		return 1
	}
	return 0
}

func runSyncDBRolePasswords(deployDir string) error {
	return setup.SyncDBRolePasswords(deployDir)
}
