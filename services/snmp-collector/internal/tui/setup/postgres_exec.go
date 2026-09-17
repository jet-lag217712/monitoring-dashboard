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
	applianceComposeEnvPath  = "/run/equate/rendered/compose.env"
	appliancePostgresEnvPath = "/run/equate/rendered/postgres.env"
)

func lookupEnvFile(path, key string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	prefix := key + "="
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			return strings.Trim(strings.TrimSpace(strings.TrimPrefix(line, prefix)), `"'`)
		}
	}
	return ""
}

func loadEnvMap(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if key != "" {
			out[key] = val
		}
	}
	return out, nil
}

func applianceComposeArgs(deployDir string) []string {
	args := []string{"compose"}
	if _, err := os.Stat(applianceComposeEnvPath); err == nil {
		args = append(args, "--env-file", applianceComposeEnvPath)
	}
	for _, name := range []string{"docker-compose.yml", generatedComposeFile} {
		path := filepath.Join(deployDir, name)
		if _, err := os.Stat(path); err == nil {
			args = append(args, "-f", path)
		}
	}
	return args
}

func runApplianceCompose(deployDir string, args ...string) error {
	cmd := exec.Command("docker", append(applianceComposeArgs(deployDir), args...)...)
	cmd.Dir = deployDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func postgresIdent(composeEnv string) (user, db string) {
	user = lookupEnvFile(composeEnv, "POSTGRES_USER")
	if user == "" {
		user = "ogsd"
	}
	db = lookupEnvFile(composeEnv, "POSTGRES_DB")
	if db == "" {
		db = "ogsd"
	}
	return user, db
}

// ExecPostgresSQL runs SQL against appliance Postgres via docker compose, or
// psql DATABASE_URL when that is set and compose env is missing.
func ExecPostgresSQL(deployDir, sql string) error {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil
	}
	if _, err := os.Stat(applianceComposeEnvPath); err == nil {
		user, db := postgresIdent(applianceComposeEnvPath)
		cmd := exec.Command("docker", append(applianceComposeArgs(deployDir), "exec", "-T", "postgres",
			"psql", "-U", user, "-d", db, "-v", "ON_ERROR_STOP=1")...)
		cmd.Dir = deployDir
		cmd.Stdin = strings.NewReader(sql + "\n")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cmd := exec.Command("psql", databaseURL, "-v", "ON_ERROR_STOP=1")
		cmd.Stdin = strings.NewReader(sql + "\n")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return fmt.Errorf("DATABASE_URL or compose env (%s) is required", applianceComposeEnvPath)
}

func waitForPostgres(deployDir, user, db string) error {
	var last error
	for i := 0; i < 60; i++ {
		cmd := exec.Command("docker", append(applianceComposeArgs(deployDir), "exec", "-T", "postgres",
			"pg_isready", "-U", user, "-d", db)...)
		cmd.Dir = deployDir
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			last = err
		}
		time.Sleep(2 * time.Second)
	}
	if last == nil {
		last = fmt.Errorf("postgres is not ready")
	}
	return fmt.Errorf("postgres is not ready: %w", last)
}
