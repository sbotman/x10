package main

import (
	"errors"
	"github.com/sbotman/x10/service"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)


// curl -H "Content-Type: application/json" -X POST -d '{"username":"xyz","password":"xyz"}' http://localhost:3000/api/login


func getConfig(c *cli.Context) (service.Config, error) {
	yamlPath := c.GlobalString("config")
	config := service.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, errors.New("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	return config, err
}

func main() {

	app := cli.NewApp()
	app.Name = "x10device"
	app.Usage = "work with the `x10Device` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config.yaml",
			Usage:  "config file to use",
			EnvVar: "APP_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.DeviceService{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "migratedb",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := service.DeviceService{}

				if err = svc.Migrate(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)

}
