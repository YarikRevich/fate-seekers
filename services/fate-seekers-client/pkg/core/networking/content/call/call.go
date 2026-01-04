package call

import (
	contentv1 "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/api"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/networking/content/handler"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/state/store"
	"google.golang.org/protobuf/proto"
)

// PerformHitPlayerWithFist performs hit player with fist operation request.
func PerformHitPlayerWithFist(sessionID int64, callback func(err error)) {
	go func() {
		message, err := proto.Marshal(&contentv1.HitPlayerWithFistRequest{
			Issuer:    store.GetRepositoryUUID(),
			SessionId: sessionID,
		})
		if err != nil {
			callback(err)

			return
		}

		err = handler.GetInstance().Send(contentv1.HIT_PLAYER_WITH_FIST_REQUEST, message)
		if err != nil {
			callback(err)

			return
		}

		callback(nil)
	}()
}
