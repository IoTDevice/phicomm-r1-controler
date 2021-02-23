//go:generate chmod +x ./scripts/generate/generate-bindata.sh
//go:generate ./scripts/generate/generate-bindata.sh
package main

import (
	"fmt"
	"github.com/IoTDevice/phicomm-r1-controler/config"
	"github.com/IoTDevice/phicomm-r1-controler/services"
	//_ "github.com/go-bindata/go-bindata"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "phicomm-r1-controler"
	myApp.Usage = "-c [config file path]"
	myApp.Version = buildVersion(version, commit, date, builtBy)
	myApp.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       config.ConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"ConfigFilePath"},
					Destination: &config.ConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				config.LoadSnapcraftConfigPath()
				config.InitConfigFile()
				return nil
			},
		},
		{
			Name:    "r1-ip",
			Aliases: []string{"r"},
			Usage:   "斐讯R1的ip直接免配置文件运行",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ip",
					Aliases:     []string{"i"},
					Value:       "",
					Usage:       "斐讯R1的ip地址",
					EnvVars:     []string{"R1IP"},
					Destination: &config.SingleIpPort,
				},
			},
			Action: func(c *cli.Context) error {
				config.LoadSnapcraftConfigPath()
				return services.Run(c)
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test this command",
			Action: func(c *cli.Context) error {
				fmt.Println("ok")
				return nil
			},
		},
	}
	myApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       config.ConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"ConfigFilePath"},
			Destination: &config.ConfigFilePath,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config.LoadSnapcraftConfigPath()
		_, err := os.Stat(config.ConfigFilePath)
		if err != nil {
			config.InitConfigFile()
		}
		config.UseConfigFile()
		return services.Run(c)
	}
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
