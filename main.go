package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-amplience/internal"
)

func main() {
	p := internal.NewAmpliencePlugin()
	plugin.ServePlugin(p)
}
