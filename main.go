package main

import (
	"os"

	"github.com/urfave/cli"
	log "unknwon.dev/clog/v2"

	"github.com/1gkx/salary/internal/cmd"
)

// func init() {
// 	conf.App.Version = "0.1.0+dev"
// }

func main() {
	app := cli.NewApp()
	app.Name = "Salary Project"
	app.Usage = "The service for approve salary card for banks clients"
	// app.Version = conf.App.Version
	app.Version = "0.1.0+dev"
	app.Commands = []cli.Command{
		cmd.Start,
		cmd.Install,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal("Failed to start application: %v", err)
	}
}
