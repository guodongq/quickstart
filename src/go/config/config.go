package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type (
	Config struct {
		Logger *Logger
		Server *Server
		Probes *Probes

		GitHubApp *GitHubAppConfig
		Cache     *Cache
	}

	Logger struct {
		Level     string
		Formatter string
		Output    string
	}

	Server struct {
		Name     string
		BasePath string
	}

	Probes struct {
		Port              int
		Enabled           bool
		LivenessEndpoint  string
		ReadinessEndpoint string
	}

	Cache struct {
		LRU *LRUConfig
		TTL *TTLConfig
	}

	LRUConfig struct {
		DefaultExpiration time.Duration
		EvictionInterval  time.Duration
		MaxEntries        int32
	}

	TTLConfig struct {
		DefaultExpiration time.Duration
		EvictionInterval  time.Duration
	}

	GitHubAppConfig struct {
		ClientID   string
		PrivateKey []byte
	}
)

func NewConfig() *Config {
	v := viper.New()
	v.SetDefault(LoggerLevelEnvKey, "debug")
	v.SetDefault(LoggerFormatEnvKey, "text")
	v.SetDefault(LoggerOutputEnvKey, "stderr")

	binaryPath := strings.Split(os.Args[0], "/")
	v.SetDefault(ServerNameEnvKey, binaryPath[len(binaryPath)-1])
	v.SetDefault(ServerBaseURLEnvKey, "/api/academy")

	v.SetDefault(ProbesPortEnvKey, 8000)
	v.SetDefault(ProbesEnabledEnvKey, true)
	v.SetDefault(ProbesLivenessEndpointEnvKey, "/healthz")
	v.SetDefault(ProbesReadinessEndpointEnvKey, "/readyz")

	return &Config{
		Logger: &Logger{
			Level:     v.GetString(LoggerLevelEnvKey),
			Formatter: v.GetString(LoggerFormatEnvKey),
			Output:    v.GetString(LoggerOutputEnvKey),
		},
		Server: &Server{
			Name:     v.GetString(ServerNameEnvKey),
			BasePath: v.GetString(ServerBaseURLEnvKey),
		},
		Probes: &Probes{
			Port:              v.GetInt(ProbesPortEnvKey),
			Enabled:           v.GetBool(ProbesEnabledEnvKey),
			LivenessEndpoint:  v.GetString(ProbesLivenessEndpointEnvKey),
			ReadinessEndpoint: v.GetString(ProbesReadinessEndpointEnvKey),
		},

		GitHubApp: &GitHubAppConfig{
			ClientID:   v.GetString(GitHubAppClientIDEnvKey),
			PrivateKey: []byte(v.GetString(GitHubAppPrivateKeyEnvKey)),
		},
	}
}
