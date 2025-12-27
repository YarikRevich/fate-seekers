package dashboards

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/services"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/cache"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository"
)

// Run starts dashboards data synchronization.
func Run() {
	go func() {
		cache.GetInstance().BeginSessionsTransaction()

		sessionsCount, err := repository.GetSessionsRepository().Count()
		if err != nil {
			logging.GetInstance().Fatal(err.Error())
		}

		services.SetAvailableSession(sessionsCount)

		cache.GetInstance().CommitSessionsTransaction()

		cache.GetInstance().BeginLobbySetTransaction()

		lobbiesCount, err := repository.GetLobbiesRepository().Count()
		if err != nil {
			logging.GetInstance().Fatal(err.Error())
		}

		services.SetAvailableLobby(lobbiesCount)

		cache.GetInstance().CommitLobbySetTransaction()
	}()
}
