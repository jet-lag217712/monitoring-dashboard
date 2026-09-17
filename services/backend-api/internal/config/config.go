package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the backend API runtime configuration.
type Config struct {
	API             APIConfig      `yaml:"api"`
	Admin           AdminConfig    `yaml:"admin"`
	Auth            AuthConfig     `yaml:"auth"`
	Database        DatabaseConfig `yaml:"database"`
	OnlineThreshold time.Duration  `yaml:"online_threshold"`
}

// APIConfig controls the public REST server.
type APIConfig struct {
	Listen      string `yaml:"listen"`
	CORSOrigins string `yaml:"cors_origins"`
}

// Supported authentication modes.
const (
	AuthModeApplianceLocal = "appliance_local"
)

// AuthConfig controls authentication for the public API.
type AuthConfig struct {
	// Mode must be appliance_local.
	Mode string `yaml:"mode"`
	// BrokerSocket is the host authentication broker Unix socket.
	BrokerSocket string `yaml:"broker_socket"`
	// BrokerTimeout bounds each broker request.
	BrokerTimeout time.Duration `yaml:"broker_timeout"`
	// SessionTTL controls appliance-local web session lifetime.
	SessionTTL time.Duration `yaml:"session_ttl"`
	// LoginRateLimit controls failed login attempts per key and window.
	LoginRateLimit int `yaml:"login_rate_limit"`
	LoginRateWindow time.Duration `yaml:"login_rate_window"`
	// LoginRateEntries bounds in-memory rate-limiter cardinality.
	LoginRateEntries int `yaml:"login_rate_entries"`
}

// AdminConfig controls the admin HTTP server (metrics/health).
type AdminConfig struct {
	Listen string `yaml:"listen"`
}

// DatabaseConfig controls the PostgreSQL connection pool.
type DatabaseConfig struct {
	URLEnv      string        `yaml:"url_env"`
	MaxConns    int32         `yaml:"max_conns"`
	MinConns    int32         `yaml:"min_conns"`
	MaxLifetime time.Duration `yaml:"max_lifetime"`
}

// Load reads and validates an API config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.applyDefaults()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.API.Listen == "" {
		c.API.Listen = ":8000"
	}
	if c.Admin.Listen == "" {
		c.Admin.Listen = ":9092"
	}
	if c.Database.URLEnv == "" {
		c.Database.URLEnv = "DATABASE_URL"
	}
	if c.Database.MaxConns == 0 {
		c.Database.MaxConns = 10
	}
	if c.Database.MinConns == 0 {
		c.Database.MinConns = 1
	}
	if c.Database.MaxLifetime == 0 {
		c.Database.MaxLifetime = time.Hour
	}
	if c.OnlineThreshold == 0 {
		c.OnlineThreshold = 5 * time.Minute
	}
	if strings.TrimSpace(c.Auth.Mode) == "" {
		c.Auth.Mode = AuthModeApplianceLocal
	}
	if strings.TrimSpace(c.Auth.BrokerSocket) == "" {
		c.Auth.BrokerSocket = "/run/equate/auth.sock"
	}
	if c.Auth.BrokerTimeout == 0 {
		c.Auth.BrokerTimeout = 3 * time.Second
	}
	if c.Auth.SessionTTL == 0 {
		c.Auth.SessionTTL = 12 * time.Hour
	}
	if c.Auth.LoginRateLimit == 0 {
		c.Auth.LoginRateLimit = 5
	}
	if c.Auth.LoginRateWindow == 0 {
		c.Auth.LoginRateWindow = 5 * time.Minute
	}
	if c.Auth.LoginRateEntries == 0 {
		c.Auth.LoginRateEntries = 2048
	}
}

// AuthMode returns the selected authentication mode.
func (c *Config) AuthMode() string {
	if mode := strings.TrimSpace(c.Auth.Mode); mode != "" {
		return mode
	}
	return AuthModeApplianceLocal
}

// DatabaseURL returns the PostgreSQL URL from the configured environment variable.
func (c *Config) DatabaseURL() string {
	if c.Database.URLEnv == "" {
		return ""
	}
	return os.Getenv(c.Database.URLEnv)
}

// CORSOriginList returns allowed CORS origins (trimmed, non-empty).
func (c *Config) CORSOriginList() []string {
	if strings.TrimSpace(c.API.CORSOrigins) == "" {
		return nil
	}
	parts := strings.Split(c.API.CORSOrigins, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Validate checks required fields.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.API.Listen) == "" {
		return fmt.Errorf("api.listen is required")
	}
	if strings.TrimSpace(c.Admin.Listen) == "" {
		return fmt.Errorf("admin.listen is required")
	}
	if strings.TrimSpace(c.Database.URLEnv) == "" {
		return fmt.Errorf("database.url_env is required")
	}
	if c.DatabaseURL() == "" {
		return fmt.Errorf("environment variable %q is required", c.Database.URLEnv)
	}
	if c.Database.MaxConns <= 0 {
		return fmt.Errorf("database.max_conns must be positive")
	}
	if c.Database.MinConns < 0 {
		return fmt.Errorf("database.min_conns must be >= 0")
	}
	if c.OnlineThreshold <= 0 {
		return fmt.Errorf("online_threshold must be positive")
	}
	if c.AuthMode() != AuthModeApplianceLocal {
		return fmt.Errorf("auth.mode must be %s", AuthModeApplianceLocal)
	}
	if !strings.HasPrefix(c.Auth.BrokerSocket, "/") {
		return fmt.Errorf("auth.broker_socket must be an absolute path")
	}
	if c.Auth.BrokerTimeout <= 0 {
		return fmt.Errorf("auth.broker_timeout must be positive")
	}
	if c.Auth.SessionTTL <= 0 {
		return fmt.Errorf("auth.session_ttl must be positive")
	}
	if c.Auth.LoginRateLimit <= 0 {
		return fmt.Errorf("auth.login_rate_limit must be positive")
	}
	if c.Auth.LoginRateWindow <= 0 {
		return fmt.Errorf("auth.login_rate_window must be positive")
	}
	if c.Auth.LoginRateEntries <= 0 {
		return fmt.Errorf("auth.login_rate_entries must be positive")
	}
	return nil
}
