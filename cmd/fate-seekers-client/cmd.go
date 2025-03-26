package main

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/runtime"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/sound"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/db"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/hajimehoshi/ebiten/v2"
)

// init performs client internal components initialization.
func init() {
	config.SetupDefaultConfig()

	config.Init()
	db.Init()
	sound.GetInstance().InitSoundAmbientBatch()
}

func main() {
	if err := ebiten.RunGame(runtime.NewRuntime()); err != nil {
		logging.GetInstance().Fatal(err.Error())
	}
}
