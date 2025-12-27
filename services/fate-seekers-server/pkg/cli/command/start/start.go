package start

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/monitoring/manager"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/events"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/dashboards"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/sync"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/validator/encryptionkey"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/component/common"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/ui/manager/translation"
	"github.com/spf13/cobra"
)

var (
	ErrEncryptionKeyValidationFailed = errors.New("err happened failed to validate encryption key")
)

// Init performs initialization of start command.
func Init(root *cobra.Command) {
	command := &cobra.Command{
		Use:   "start",
		Short: "Starts FateSeekers server process",
		Long:  `Starts FateSeekers server process as a blocking operation.`,
		Run: func(cmd *cobra.Command, args []string) {
			db.Init()

			sync.Run()

			dashboards.Run()

			events.Run()

			if !encryptionkey.Validate(config.GetSettingsNetworkingEncryptionKey()) {
				logging.GetInstance().Fatal(ErrEncryptionKeyValidationFailed.Error())

				return
			}

			var (
				interruption     = make(chan os.Signal, 1)
				gracefulShutdown = make(chan bool, 1)
				done             = make(chan struct{})
			)

			signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)

			go func() {
				select {
				case <-interruption:
				case <-gracefulShutdown:
				}

				// 			"server.monitoring.start.title": {
				//     "one": "FateSeekers gaming server has been started!",
				//     "other": "FateSeekers gaming server has been started!"
				// },
				// "server.monitoring.start.monitoring.title": {
				//     "one": "FateSeekers gaming server has been started(including monitoring)!",
				//     "other": "FateSeekers gaming server has been started(including monitoring)!"
				// },
				// "server.monitoring.stop.title": {
				//     "one": "FateSeekers gaming server has been stoped!",
				//     "other": "FateSeekers gaming server has been stoped!"
				// },
				// "server.monitoring.stop.monitoring.title": {
				//     "one": "FateSeekers gaming server has been stoped(including monitoring)!",
				//     "other": "FateSeekers gaming server has been stoped(including monitoring)!"
				// }

				connector.GetInstance().Close(func(err error) {
					if err != nil {
						logging.GetInstance().Error(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("shared.networking.close-failure"),
								err.Error()))
					}

					if config.GetSettingsMonitoringEnabled() {
						manager.GetInstance().Remove(func(err error) {
							if err != nil {
								logging.GetInstance().Fatal(
									common.ComposeMessage(
										translation.GetInstance().GetTranslation("server.networking.start-failure"),
										err.Error()))
							} else {
								logging.GetInstance().Info(
									translation.GetInstance().GetTranslation("server.monitoring.stop.monitoring.title"))
							}

							close(done)
						})
					} else {
						logging.GetInstance().Info(
							translation.GetInstance().GetTranslation("server.monitoring.stop.title"))

						close(done)
					}
				})
			}()

			connector.GetInstance().Connect(func(err error) {
				if err != nil {
					logging.GetInstance().Error(
						common.ComposeMessage(
							translation.GetInstance().GetTranslation("server.networking.start-failure"),
							err.Error()))

					gracefulShutdown <- true

					return
				}

				if config.GetSettingsMonitoringEnabled() {
					manager.GetInstance().Deploy(func(err error) {
						if err != nil {
							logging.GetInstance().Error(
								common.ComposeMessage(
									translation.GetInstance().GetTranslation("server.networking.start-failure"),
									err.Error()))

							gracefulShutdown <- true

							return
						}

						logging.GetInstance().Info(
							translation.GetInstance().GetTranslation("server.monitoring.start.monitoring.title"))
					})
				} else {
					logging.GetInstance().Info(
						translation.GetInstance().GetTranslation("server.monitoring.start.title"))
				}
			})

			<-done
		},
	}

	root.AddCommand(command)
}
