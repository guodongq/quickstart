package config

const (
	GitHubAppClientIDEnvKey   = "GITHUB_APP_CLIENT_ID"
	GitHubAppPrivateKeyEnvKey = "GITHUB_APP_PRIVATE_KEY"
)

const (
	LoggerLevelEnvKey  = "LOGGER_LEVEL"
	LoggerFormatEnvKey = "LOGGER_FORMAT"
	LoggerOutputEnvKey = "LOGGER_OUTPUT"

	ServerBaseURLEnvKey = "SERVER_BASE_URL"
	ServerNameEnvKey    = "SERVER_NAME"

	ProbesPortEnvKey              = "PROBES_PORT"
	ProbesEnabledEnvKey           = "PROBES_ENABLED"
	ProbesLivenessEndpointEnvKey  = "PROBES_LIVENESS_ENDPOINT"
	ProbesReadinessEndpointEnvKey = "PROBES_READINESS_ENDPOINT"
)
