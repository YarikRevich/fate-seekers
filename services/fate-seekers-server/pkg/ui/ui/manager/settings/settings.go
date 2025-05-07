package settings

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/state/value"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/validator/encryptionkey"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/validator/port"
)

// ProcessChanges performs provided changes application.
func ProcessChanges(networkingPort, networkingEncryptionKey, language string) bool {
	var applied, demandRestart bool

	if config.GetSettingsNetworkingServerHost() != networkingPort {
		if !port.Validate(networkingPort) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("settingsmanager.invalid-networking-host"),
				time.Second*3,
				common.NotificationErrorTextColor)

			return false
		}
	}

	if config.GetSettingsNetworkingEncryptionKey() != networkingEncryptionKey {
		if !encryptionkey.Validate(networkingEncryptionKey) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("settingsmanager.invalid-networking-encryption-key"),
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

	if config.GetSettingsNetworkingServerHost() != networkingHost {
		config.SetSettingsNetworkingHost(networkingHost)

		applied = true
	}

	if config.GetSettingsNetworkingEncryptionKey() != networkingEncryptionKey {
		config.SetSettingsNetworkingEncryptionKey(networkingEncryptionKey)

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
func AnyProvidedChanges(soundMusic, soundFX int, networkingHost, networkingEncryptionKey, language string) bool {
	return config.GetSettingsSoundMusic() != soundMusic ||
		config.GetSettingsSoundFX() != soundFX ||
		config.GetSettingsNetworkingServerHost() != networkingHost ||
		config.GetSettingsNetworkingEncryptionKey() != networkingEncryptionKey ||
		config.GetSettingsLanguage() != language
}
