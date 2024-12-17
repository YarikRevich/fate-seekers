package main

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/runtime"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/logging"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	config.Init()
}

func main() {
	if err := ebiten.RunGame(runtime.NewRuntime()); err != nil {
		logging.GetInstance().Fatal(err.Error())
	}
}
