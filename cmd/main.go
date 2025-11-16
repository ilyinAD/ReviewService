package main

import (
	"avitostazhko/internal/utils"

	"go.uber.org/fx"
)

func main() {
	utils.LoadEnv()
	fx.New(BuildApp()).Run()
}
