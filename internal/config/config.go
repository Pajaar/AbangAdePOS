package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	CookieName     string
	CookieSecure   bool
	CookieSameSite string
	CookieDomain   string

	CORSAllowedOrigins []string

	RateLimitLoginPerMinute float64
	RateLimitLoginBurst     int
}

// Load membaca .env (kalau ada) lalu variabel environment asli.
// Di Railway, .env tidak ada karena env var di-inject langsung ke proses —
// itu kondisi normal, bukan error.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		CookieName:     getEnv("COOKIE_NAME", "refresh_token"),
		CookieSecure:   getEnvBool("COOKIE_SECURE", false),
		CookieSameSite: getEnv("COOKIE_SAMESITE", "lax"),
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),

		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "")),
	}

	cfg.AccessTokenTTL = time.Duration(getEnvInt("ACCESS_TOKEN_TTL_MINUTES", 15)) * time.Minute
	cfg.RefreshTokenTTL = time.Duration(getEnvInt("REFRESH_TOKEN_TTL_DAYS", 7)) * 24 * time.Hour

	cfg.RateLimitLoginPerMinute = getEnvFloat("RATE_LIMIT_LOGIN_PER_MINUTE", 5)
	cfg.RateLimitLoginBurst = getEnvInt("RATE_LIMIT_LOGIN_BURST", 5)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string
	if c.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if c.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return fmt.Errorf("variabel environment wajib belum diisi: %s", strings.Join(missing, ", "))
	}
	return nil
}

// MySQLDSN membangun DSN untuk github.com/go-sql-driver/mysql.
// parseTime=true wajib agar DATETIME/TIMESTAMP ter-scan ke time.Time.
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC&timeout=5s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
