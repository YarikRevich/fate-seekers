package settings

import (
	"fmt"
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
)

// ProcessChanges performs provided changes application.
func ProcessChanges(soundMusic, soundFX int, networkingHost, language string) {
	var applied, demandRestart bool

	if config.GetSettingsSoundMusic() != soundMusic {
		config.SetSettingsSoundMusic(soundMusic)

		applied = true
	}

	if config.GetSettingsSoundFX() != soundFX {
		config.SetSettingsSoundFX(soundFX)

		applied = true
	}

	if config.GetSettingsNetworkingHost() != networkingHost {
		config.SetSettingsNetworkingHost(networkingHost)

		applied = true
		demandRestart = true
	}

	if config.GetSettingsLanguage() != language {
		config.SetSettingsLanguage(language)

		applied = true
		demandRestart = true
	}

	if applied {
		fmt.Println(demandRestart)

		notification.GetInstance().Push(
			translation.GetInstance().GetTranslation("settingsmanager.success"),
			time.Second*3,
			common.NotificationInfoTextColor)
	}
}

// AnyProvidedChanges checks if there are any new provided changes.
func AnyProvidedChanges(soundMusic, soundFX int, networkingHost, language string) bool {
	return config.GetSettingsSoundMusic() != soundMusic ||
		config.GetSettingsSoundFX() != soundFX ||
		config.GetSettingsNetworkingHost() != networkingHost ||
		config.GetSettingsLanguage() != language
}
