package sqlvet

import "github.com/MirrexOne/sqlvet/pkg/config"

// Settings is a type alias for SQLVetSettings from the config package.
type Settings = config.SQLVetSettings

// DefaultSettings returns the default configuration for SQLVet.
func DefaultSettings() Settings {
	return config.DefaultSettings()
}
