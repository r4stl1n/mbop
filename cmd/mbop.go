package main

import (
	"github.com/r4stl1n/mbop/internal/mbop/client"
	"github.com/r4stl1n/mbop/pkg/util"
	"go.uber.org/zap"
)

func init() {
	new(util.Logger).Init("mbop")
}

func main() {

	new(util.Utils).Init().PrintBanner()

	runError := new(client.CLI).Init().Run()
	if runError != nil {
		zap.L().Error("run error occurred", zap.Error(runError))
	}
}
