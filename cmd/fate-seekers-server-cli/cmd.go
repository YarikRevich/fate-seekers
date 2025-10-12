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

// TODO: make monitoring start based on config properties and during a laucn
