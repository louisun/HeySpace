package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	paramIn     = "in"
	paramOut    = "out"
	paramBackup = "backup"
	paramQuiet  = "quiet"
)

type mode struct {
	MarkdownMode bool // 兼容 Markdown 格式
	QuietMode    bool // 不输出具体文件日志
}

var globalConfig = &mode{
	QuietMode: false,
}

func main() {
	app := &cli.App{
		Name:     "HeySpace",
		Usage:    "在中英文之间添加空格",
		Version:  "v0.0.2",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Renzo",
				Email: "luyang.sun@outlook.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        paramIn,
				Aliases:     []string{"i"},
				Usage:       "输入文件路径",
				DefaultText: "默认剪贴板输入",
			},
			&cli.StringFlag{
				Name:        paramOut,
				Aliases:     []string{"o"},
				Usage:       "输出文件路径",
				DefaultText: "默认剪贴板输出",
			},
			&cli.StringFlag{
				Name:    paramBackup,
				Aliases: []string{"b"},
				Usage:   "备份目录路径",
			},
			&cli.BoolFlag{
				Name:        paramQuiet,
				Aliases:     []string{"q"},
				Usage:       "不输出具体文件日志",
				Value:       false,
				DefaultText: "关闭",
			},
		},
		Action: runApp,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp(c *cli.Context) error {
	globalConfig.QuietMode = c.Bool(paramQuiet)

	if len(c.String("in")) == 0 {
		if err := handleClipboard(); err != nil {
			log.Fatal(err)
		}

		return nil
	}

	if err := handlePathInput(c.String(paramIn), c.String(paramOut), c.String(paramBackup)); err != nil {
		log.Fatal(err)
	}

	return nil
}
