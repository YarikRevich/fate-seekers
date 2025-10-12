package start

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/connector"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/networking/metadata/events"
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

			events.Run()

			if !encryptionkey.Validate(config.GetSettingsNetworkingEncryptionKey()) {
				logging.GetInstance().Fatal(ErrEncryptionKeyValidationFailed.Error())

				return
			}

			connector.GetInstance().Connect(func(err error) {
				if err != nil {
					logging.GetInstance().Fatal(
						common.ComposeMessage(
							translation.GetInstance().GetTranslation("server.networking.start-failure"),
							err.Error()))

					return
				}

				logging.GetInstance().Info("FateSeekers gaming server has been started!")
			})

			interruption := make(chan os.Signal, 1)
			done := make(chan struct{})

			signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-interruption

				connector.GetInstance().Close(func(err error) {
					if err != nil {
						logging.GetInstance().Error(
							common.ComposeMessage(
								translation.GetInstance().GetTranslation("shared.networking.close-failure"),
								err.Error()))

						return
					}

					logging.GetInstance().Info("FateSeekers gaming server has been stoped!")

					close(done)
				})
			}()

			<-done
		},
	}

	root.AddCommand(command)
}
