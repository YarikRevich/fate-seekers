package selector

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/sessionid"
)

// ProcessChanges performs provided changes validation.
func ProcessChanges(sessionID string) bool {
	if !sessionid.Validate(sessionID) {
		notification.GetInstance().Push(
			translation.GetInstance().GetTranslation("client.selectormanager.invalid-session-id"),
			time.Second*3,
			common.NotificationErrorTextColor)

		return false
	}

	return true
}
