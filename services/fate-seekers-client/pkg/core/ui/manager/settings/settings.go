package settings

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/encryptionkey"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/host"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/action"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/dispatcher"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/value"
)

// ProcessChanges performs provided changes application.
func ProcessChanges(soundMusic, soundFX int, networkingServerHost, networkingEncryptionKey, language string) bool {
	var applied, demandRestart bool

	if config.GetSettingsNetworkingServerHost() != networkingServerHost {
		if !host.Validate(networkingServerHost) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("client.settingsmanager.invalid-networking-server-host"),
				time.Second*3,
				common.NotificationErrorTextColor)

			return false
		}
	}

	if config.GetSettingsNetworkingEncryptionKey() != networkingEncryptionKey {
		if !encryptionkey.Validate(networkingEncryptionKey) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("shared.settingsmanager.invalid-networking-encryption-key"),
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

		applied = true
	}

	if config.GetSettingsNetworkingServerHost() != networkingServerHost {
		config.SetSettingsNetworkingServerHost(networkingServerHost)

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
				translation.GetInstance().GetTranslation("shared.settingsmanager.restart-demand"),
				time.Second*3,
				common.NotificationInfoTextColor)
		} else {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("shared.settingsmanager.success"),
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
