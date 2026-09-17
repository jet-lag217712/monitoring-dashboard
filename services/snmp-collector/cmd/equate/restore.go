package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/equate/ogsd/services/snmp-collector/internal/tui/setup"
)

func runRestore(args []string) int {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "usage: equate restore")
		return 2
	}
	deployDir, err := resolveDeployDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "restore: %v\n", err)
		return 1
	}
	if err := restoreAppliance(deployDir); err != nil {
		fmt.Fprintf(os.Stderr, "restore: %v\n", err)
		return 1
	}
	fmt.Fprintln(os.Stdout, "restore complete")
	return 0
}

func restoreAppliance(deployDir string) error {
	return restoreApplianceWithDirs(deployDir, setup.DurableRenderedDir, setup.RunRenderedDir)
}

func restoreApplianceWithDirs(deployDir, durableDir, runDir string) error {
	if err := setup.RequireConfigured(deployDir); err != nil {
		return err
	}
	if err := setup.RestoreRenderedDir(durableDir, runDir); err != nil {
		return err
	}
	if err := setup.OverlaySNMPFromEnvFile(filepath.Join(deployDir, ".env"), runDir); err != nil {
		return err
	}
	if err := setup.PersistRenderedDir(runDir, durableDir); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "restore: starting stack...")
	if err := runDockerCompose(deployDir, "up", "-d", "--remove-orphans"); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}
	if err := setup.SyncDBRolePasswords(deployDir); err != nil {
		return err
	}
	if err := setup.AppliancePostConfigure(deployDir); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "restore: waiting for collectors...")
	if err := setup.WaitForSiteCollectors(deployDir, 3*time.Minute); err != nil {
		return err
	}
	return nil
}
