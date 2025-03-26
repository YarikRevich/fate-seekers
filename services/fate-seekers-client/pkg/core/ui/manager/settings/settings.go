package settings

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/host"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
)

// ProcessChanges performs provided changes application.
func ProcessChanges(soundMusic, soundFX int, networkingHost, language string) bool {
	var applied, demandRestart bool

	if config.GetSettingsNetworkingHost() != networkingHost {
		if !host.Validate(networkingHost) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("settingsmanager.invalid-networking-host"),
				time.Second*3,
				common.NotificationErrorTextColor)

			return false
		}
	}

	if config.GetSettingsSoundMusic() != soundMusic {
		config.SetSettingsSoundMusic(soundMusic)

		dispatcher.GetInstance().Dispatch(
			action.NewSetSoundMusicUpdated(value.SOUND_MUSIC_UPDATED_FALSE_VALUE))

		applied = true
	}

	if config.GetSettingsSoundFX() != soundFX {
		config.SetSettingsSoundFX(soundFX)

		dispatcher.GetInstance().Dispatch(
			action.NewSetSoundFXUpdated(value.SOUND_FX_UPDATED_FALSE_VALUE))

		applied = true
	}

	if config.GetSettingsNetworkingHost() != networkingHost {
		config.SetSettingsNetworkingHost(networkingHost)

		applied = true
	}

	if config.GetSettingsLanguage() != language {
		config.SetSettingsLanguage(language)

		applied = true
		demandRestart = true
	}

	if applied {
		if demandRestart {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("settingsmanager.restart-demand"),
				time.Second*3,
				common.NotificationInfoTextColor)
		} else {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("settingsmanager.success"),
				time.Second*3,
				common.NotificationInfoTextColor)
		}
	}

	return true
}

// AnyProvidedChanges checks if there are any new provided changes.
func AnyProvidedChanges(soundMusic, soundFX int, networkingHost, language string) bool {
	return config.GetSettingsSoundMusic() != soundMusic ||
		config.GetSettingsSoundFX() != soundFX ||
		config.GetSettingsNetworkingHost() != networkingHost ||
		config.GetSettingsLanguage() != language
}
