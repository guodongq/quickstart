package env

import (
	"os"
	"strconv"
	"strings"
)

func SetFromEnvVal(dst *string, keys []string) {
	for _, k := range keys {
		if v := os.Getenv(k); len(v) != 0 {
			*dst = v
			break
		}
	}
}

func SetBoolPtrFromEnvVal(dst **bool, keys []string) {
	for _, k := range keys {
		value := os.Getenv(k)
		if len(value) == 0 {
			continue
		}

		switch {
		case strings.EqualFold(value, "false"):
			*dst = new(bool)
			**dst = false
		case strings.EqualFold(value, "true"):
			*dst = new(bool)
			**dst = true
		}
	}
}

func GetEnvAsStringOrFallback(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func GetEnvAsIntOrFallback(key string, defaultValue int) (int, error) {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

func GetEnvAsFloat64OrFallback(key string, defaultValue float64) (float64, error) {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}
