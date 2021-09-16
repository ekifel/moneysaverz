package main

import (
	"github.com/ekifel/moneysaverz/internal/app"
)

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
