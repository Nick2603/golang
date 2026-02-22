package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
	JWTExpiry time.Duration
}

func Load() *Config {
	jwtExpiry := 24 * time.Hour
	if hours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24")); err == nil {
		jwtExpiry = time.Duration(hours) * time.Hour
	}

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		MongoURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGODB_DATABASE", "backend"),
		JWTSecret: getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiry: jwtExpiry,
	}

	slog.Info("Config loaded",
		"port", cfg.Port,
		"mongoURI", cfg.MongoURI,
		"database", cfg.MongoDB,
		"jwtExpiry", cfg.JWTExpiry,
	)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
