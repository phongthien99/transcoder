package main

import (
	src "example/src"
	"example/src/config"

	_ "github.com/sigmaott/gest/package/technique/version"
)

func main() {
	//rebuild
	app := src.NewApp(config.GetConfiguration())
	app.Run()

}
