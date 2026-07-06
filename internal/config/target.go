package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kizuren/pihole-sync/internal/pihole/model"
)

func (c *Config) loadTargets() error {
	primary, err := loadPrimary()
	if err != nil {
		return err
	}

	replicas, err := loadReplicas()
	if err != nil {
		return err
	}

	c.Primary = *primary
	c.Replicas = replicas
	return nil
}

func loadPrimary() (*model.PiHole, error) {
	env := "PRIMARY"
	var ph *model.PiHole
	var err error

	if fileValue := os.Getenv(fmt.Sprintf("%s_FILE", env)); len(fileValue) > 0 {
		bytes, readErr := os.ReadFile(filepath.Clean(fileValue))
		if readErr != nil {
			return nil, readErr
		}

		ph, err = parse(strings.TrimSpace(string(bytes)))
	} else if envValue := os.Getenv(env); len(envValue) > 0 {
		ph, err = parse(envValue)
	} else {
		return nil, fmt.Errorf("missing required env: %s/%s_FILE", env, env)
	}

	if err != nil {
		return nil, err
	}

	totpSecret, err := loadOptionalValue(fmt.Sprintf("%s_TOTP", env))
	if err != nil {
		return nil, err
	}
	ph.TOTPSecret = totpSecret

	return ph, nil
}

func loadReplicas() ([]model.PiHole, error) {
	env := "REPLICAS"
	var replicas []model.PiHole
	var err error

	if fileValue := os.Getenv(fmt.Sprintf("%s_FILE", env)); len(fileValue) > 0 {
		bytes, readErr := os.ReadFile(filepath.Clean(fileValue))
		if readErr != nil {
			return nil, readErr
		}

		replicas, err = parseMultiple(strings.Split(strings.TrimSpace(string(bytes)), ","))
	} else if envValue := os.Getenv(env); len(envValue) > 0 {
		replicas, err = parseMultiple(strings.Split(envValue, ","))
	} else {
		return nil, fmt.Errorf("missing required env: %s/%s_FILE", env, env)
	}

	if err != nil {
		return nil, err
	}

	totpValue, err := loadOptionalValue(fmt.Sprintf("%s_TOTP", env))
	if err != nil {
		return nil, err
	}
	if totpValue != "" {
		secrets := strings.Split(totpValue, ",")
		for i := range replicas {
			if i < len(secrets) {
				replicas[i].TOTPSecret = secrets[i]
			}
		}
	}

	return replicas, nil
}

// loadOptionalValue reads an optional env var, preferring the `_FILE` variant (matching the
// PRIMARY/REPLICAS convention) so secrets like TOTP seeds can be provided via Docker/file secrets too.
func loadOptionalValue(env string) (string, error) {
	if fileValue := os.Getenv(fmt.Sprintf("%s_FILE", env)); len(fileValue) > 0 {
		bytes, err := os.ReadFile(filepath.Clean(fileValue))
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(bytes)), nil
	}

	return os.Getenv(env), nil
}

func parse(value string) (*model.PiHole, error) {
	ph := model.PiHole{}
	if err := ph.Decode(value); err != nil {
		return nil, err
	}
	return &ph, nil
}

func parseMultiple(values []string) ([]model.PiHole, error) {
	replicas := []model.PiHole{}
	for _, value := range values {
		ph, err := parse(value)
		if err != nil {
			return nil, err
		}

		replicas = append(replicas, *ph)
	}
	return replicas, nil
}
