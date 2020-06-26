package main

import (
	"fmt"
	"os"

	"github.com/boriska70/gorming/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Hello")

	app := app.NewApp()
	defer app.Close(false)

	if err := app.Initialize(); err != nil {
		log.Panicf("exiting - cannot initialize application: %s", err.Error())
		os.Exit(2)
	}

	app.Run()

}
