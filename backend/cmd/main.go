package main

import (
	"github.com/wozhdeleniye/redclass-app/internal/config"
	"github.com/wozhdeleniye/redclass-app/internal/modules"
	"go.uber.org/fx"
)

func init() {

}

func main() {
	app := fx.New(
		fx.Provide(config.Load),

		modules.GormModule,
	)

	app.Run()
}
