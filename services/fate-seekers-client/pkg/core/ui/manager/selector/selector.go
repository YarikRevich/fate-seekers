package selector

import (
	"time"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/notification"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/manager/translation"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/ui/validator/sessionname"
)

// ProcessChanges performs provided changes validation.
func ProcessChanges(sessionName string) bool {
	if !sessionname.Validate(sessionName) {
		notification.GetInstance().Push(
			translation.GetInstance().GetTranslation("client.selectormanager.invalid-session-name"),
			time.Second*3,
			common.NotificationErrorTextColor)

		return false
	}

	return true
}
