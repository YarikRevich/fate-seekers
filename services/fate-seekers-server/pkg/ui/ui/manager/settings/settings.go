package settings

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/validator/encryptionkey"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/validator/port"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
)

// ProcessChanges performs provided changes application.
func ProcessChanges(networkingServerPort string, networkingEncryptionKey, language string) bool {
	var applied, demandRestart bool

	if config.GetSettingsNetworkingServerPort() != networkingServerPort {
		if !port.Validate(networkingServerPort) {
			notification.GetInstance().Push(
				translation.GetInstance().GetTranslation("shared.settingsmanager.invalid-networking-server-port"),
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

	if config.GetSettingsNetworkingServerPort() != networkingServerPort {
		config.SetSettingsNetworkingServerPort(networkingServerPort)

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
func AnyProvidedChanges(networkingServerPort, networkingEncryptionKey, language string) bool {
	return config.GetSettingsNetworkingServerPort() != networkingServerPort ||
		config.GetSettingsNetworkingEncryptionKey() != networkingEncryptionKey ||
		config.GetSettingsLanguage() != language
}
