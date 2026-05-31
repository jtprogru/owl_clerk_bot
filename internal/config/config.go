package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken  string
	OwnerID   int64
	Debug     bool
	ServeHost string
	ServePort int
	DBPath    string
	WebUser   string
	WebPass   string
}

func (c Config) ServeAddr() string {
	return fmt.Sprintf("%s:%d", c.ServeHost, c.ServePort)
}

func Load() (*Config, error) {
	c := &Config{
		ServeHost: getenv("OWL_SERVE_HOST", "127.0.0.1"),
		DBPath:    getenv("OWL_DB_PATH", "./data/owl.db"),
		WebUser:   getenv("OWL_WEB_USER", "admin"),
		WebPass:   os.Getenv("OWL_WEB_PASS"),
	}

	c.BotToken = strings.TrimSpace(os.Getenv("OWL_BOT_TOKEN"))
	if c.BotToken == "" {
		return nil, errors.New("OWL_BOT_TOKEN is required")
	}

	rawOwner := strings.TrimSpace(os.Getenv("OWL_OWNER_ID"))
	if rawOwner == "" {
		return nil, errors.New("OWL_OWNER_ID is required")
	}
	owner, err := strconv.ParseInt(rawOwner, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("OWL_OWNER_ID: %w", err)
	}
	c.OwnerID = owner

	port, err := strconv.Atoi(getenv("OWL_SERVE_PORT", "8000"))
	if err != nil {
		return nil, fmt.Errorf("OWL_SERVE_PORT: %w", err)
	}
	c.ServePort = port

	c.Debug = parseBool(os.Getenv("OWL_DEBUG"))

	if c.WebPass == "" {
		return nil, errors.New("OWL_WEB_PASS is required")
	}

	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}
