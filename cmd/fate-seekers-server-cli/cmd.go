package main

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/cli/command"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-server/pkg/shared/config"
)

// init performs client internal components initialization.
func init() {
	config.SetupDefaultConfig()
	config.Init()
}

func main() {
	command.Init()
}
