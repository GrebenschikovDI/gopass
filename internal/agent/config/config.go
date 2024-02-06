package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ClientConfig struct {
	Server string
	Ping   time.Duration
	OS     string
}

const (
	defaultServer = "localhost:8000"
	defaultPing   = 3 * time.Second
)

func LoadConfig() (*ClientConfig, error) {
	config := &ClientConfig{}
	err := config.parseFlags()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *ClientConfig) parseFlags() error {
	server := flag.String("s", defaultServer, "server to connect")
	ping := flag.String("p", defaultPing.String(), "check rate")
	c.Server = fmt.Sprintf("http://%s/", *server)
	duration, err := parseDuration(*ping, defaultPing)
	if err != nil {
		return err
	}
	c.Ping = duration
	c.OS, err = checkOS()
	if err != nil {
		return err
	}
	return nil
}

// parseDuration разбирает строку и возвращает длительность времени.
func parseDuration(value string, defaultValue time.Duration) (time.Duration, error) {
	if _, err := strconv.Atoi(value); err == nil {
		value = fmt.Sprintf("%ss", value)
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue, fmt.Errorf("cannot parse interval to Duration: %w", err)
	}
	return duration, nil
}

func checkOS() (string, error) {
	info, err := os.Stat("/")
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	if info.IsDir() {
		return "unix", nil
	} else {
		return "windows", nil
	}
}
