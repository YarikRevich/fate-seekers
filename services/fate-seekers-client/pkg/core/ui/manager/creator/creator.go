package creator

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/sessionname"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/sessionseed"
)

// ProcessChanges performs provided changes validation.
func ProcessChanges(name, seed string) bool {
	if !sessionname.Validate(name) {
		notification.GetInstance().Push(
			translation.GetInstance().GetTranslation("client.creatormanager.invalid-session-name"),
			time.Second*3,
			common.NotificationErrorTextColor)

		return false
	}

	if !sessionseed.Validate(seed) {
		notification.GetInstance().Push(
			translation.GetInstance().GetTranslation("client.creatormanager.invalid-session-seed"),
			time.Second*3,
			common.NotificationErrorTextColor)

		return false
	}

	return true
}
