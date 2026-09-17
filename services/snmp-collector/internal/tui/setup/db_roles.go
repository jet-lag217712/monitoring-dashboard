package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SyncDBRolePasswords applies ogsd_ingestion / ogsd_api passwords from rendered secrets.
func SyncDBRolePasswords(deployDir string) error {
	if _, err := os.Stat(applianceComposeEnvPath); err != nil {
		return nil
	}
	if _, err := os.Stat(appliancePostgresEnvPath); err != nil {
		return fmt.Errorf("sync database role passwords: missing %s", appliancePostgresEnvPath)
	}
	if _, err := os.Stat(filepath.Join(deployDir, "docker-compose.yml")); err != nil {
		return fmt.Errorf("sync database role passwords: missing docker-compose.yml")
	}

	composeEnv, err := loadEnvMap(applianceComposeEnvPath)
	if err != nil {
		return fmt.Errorf("sync database role passwords: %w", err)
	}
	postgresEnv, err := loadEnvMap(appliancePostgresEnvPath)
	if err != nil {
		return fmt.Errorf("sync database role passwords: %w", err)
	}

	user := firstNonEmpty(composeEnv["POSTGRES_USER"], postgresEnv["POSTGRES_USER"], "ogsd")
	db := firstNonEmpty(composeEnv["POSTGRES_DB"], postgresEnv["POSTGRES_DB"], "ogsd")
	ingestion := firstNonEmpty(composeEnv["OGSD_INGESTION_PASSWORD"], postgresEnv["OGSD_INGESTION_PASSWORD"])
	api := firstNonEmpty(composeEnv["OGSD_API_PASSWORD"], postgresEnv["OGSD_API_PASSWORD"])
	if ingestion == "" || api == "" {
		return fmt.Errorf("sync database role passwords: OGSD_INGESTION_PASSWORD and OGSD_API_PASSWORD are required")
	}

	fmt.Fprintln(os.Stderr, "sync-db-role-passwords: ensuring postgres is running...")
	if err := runApplianceCompose(deployDir, "up", "-d", "postgres"); err != nil {
		return fmt.Errorf("sync database role passwords: start postgres: %w", err)
	}
	fmt.Fprintln(os.Stderr, "sync-db-role-passwords: waiting for postgres...")
	if err := waitForPostgres(deployDir, user, db); err != nil {
		return fmt.Errorf("sync database role passwords: %w", err)
	}

	sql := fmt.Sprintf(
		"ALTER ROLE ogsd_ingestion WITH PASSWORD '%s';\nALTER ROLE ogsd_api WITH PASSWORD '%s';\n",
		escapeSQLLiteral(ingestion),
		escapeSQLLiteral(api),
	)
	fmt.Fprintln(os.Stderr, "sync-db-role-passwords: applying ogsd_ingestion and ogsd_api passwords...")
	if err := ExecPostgresSQL(deployDir, sql); err != nil {
		return fmt.Errorf("sync database role passwords: %w", err)
	}
	fmt.Fprintln(os.Stderr, "sync-db-role-passwords: done")
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
