package sqlvet

import "github.com/MirrexOne/sqlvet/pkg/config"

// Settings represents the configuration for the sqlvet linter
// This is a type alias for golangci-lint integration
type Settings = config.SQLVetSettings

// DefaultSettings returns the default configuration for sqlvet
func DefaultSettings() Settings {
	return config.DefaultSettings()
}