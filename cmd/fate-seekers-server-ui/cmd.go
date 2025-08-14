package main

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/logging"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/repository/sync"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/ui/runtime"
	"github.com/hajimehoshi/ebiten/v2"
)

// init performs client internal components initialization.
func init() {
	config.SetupDefaultConfig()
	config.Init()

	db.Init()

	sync.Run()
}

func main() {
	if err := ebiten.RunGame(runtime.NewRuntime()); err != nil {
		logging.GetInstance().Fatal(err.Error())
	}
}
