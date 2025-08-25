package sqlvet

import "github.com/MirrexOne/sqlvet/pkg/config"

type Settings = config.SQLVetSettings

func DefaultSettings() Settings {
	return config.DefaultSettings()
}
