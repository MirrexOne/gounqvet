package gounqvet

import "github.com/MirrexOne/gounqvet/pkg/config"

// Settings is a type alias for GounqvetSettings from the config package.
type Settings = config.GounqvetSettings

// DefaultSettings returns the default configuration for Gounqvet.
func DefaultSettings() Settings {
	return config.DefaultSettings()
}
