//nolint:tagliatelle // exclude this linter until it supports UPPER_SNAKE_CASE
package config

import (
	"log/slog"
	"time"
)

type LogLeveler string

func (l LogLeveler) Level() slog.Level {
	var level slog.Level

	_ = level.UnmarshalText([]byte(l))

	return level
}

// Config holds the server configuration.
type Config struct {
	LogLevel             LogLeveler         `mapstructure:"LOG_LEVEL"`
	ServiceTokens        string             `mapstructure:"SERVICE_TOKENS"`
	TracingEnabled       bool               `mapstructure:"TRACING_ENABLED"`
	ProfilingEnabled     bool               `mapstructure:"PROFILING_ENABLED"`
	RequestTimeThreshold time.Duration      `mapstructure:"REQUEST_TIME_THRESHOLD"`
	DB                   DB                 `mapstructure:",squash"`
	HTTP                 HTTP               `mapstructure:",squash"`
	HTTPCaller           HTTPCaller         `mapstructure:",squash"`
	Locales              Locales            `mapstructure:",squash"`
	UserService          UserService        `mapstructure:",squash"`
	ListingService       ListingService     `mapstructure:",squash"`
	ListingViewService   ListingViewService `mapstructure:",squash"`
}

type UserService struct {
	URL      string        `mapstructure:"USER_SERVICE_URL"`
	MaxRetry int           `mapstructure:"USER_SERVICE_MAX_RETRY"`
	Timeout  time.Duration `mapstructure:"USER_SERVICE_TIMEOUT"`
}

type ListingService struct {
	URL      string        `mapstructure:"LISTING_SERVICE_URL"`
	MaxRetry int           `mapstructure:"LISTING_SERVICE_MAX_RETRY"`
	Timeout  time.Duration `mapstructure:"LISTING_SERVICE_TIMEOUT"`
}

type ListingViewService struct {
	URL      string        `mapstructure:"LISTING_VIEW_SERVICE_URL"`
	MaxRetry int           `mapstructure:"LISTING_VIEW_SERVICE_MAX_RETRY"`
	Timeout  time.Duration `mapstructure:"LISTING_VIEW_SERVICE_TIMEOUT"`
}

type DB struct {
	DSN                   string        `mapstructure:"DB_DSN"`
	MaxOpenConnections    int           `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections    int           `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	MaxConnectionLifetime time.Duration `mapstructure:"DB_MAX_CONNECTIONS_LIFETIME"`
	MaxIdleConnectionTime time.Duration `mapstructure:"DB_MAX_IDLE_CONNECTIONS_TIME"`
}

type HTTP struct {
	Port          int           `mapstructure:"HTTP_PORT"`
	Timeout       time.Duration `mapstructure:"HTTP_TIMEOUT"`
	PprofEnabled  bool          `mapstructure:"PPROF_ENABLED"`
	PprofPort     int           `mapstructure:"PPROF_PORT"`
	AllowedOrigin []string      `mapstructure:"ALLOWED_ORIGIN"`
}

type HTTPCaller struct {
	Timeout time.Duration `mapstructure:"HTTP_CALLER_TIMEOUT"`
}

type Locales struct {
	BasePath           string `mapstructure:"LOCALES_BASE_PATH"`
	SupportedLanguages string `mapstructure:"LOCALES_SUPPORTED_LANGUAGES"`
}
