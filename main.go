package main

import (
	"HeySpace/clipboard"
	"HeySpace/config"
	"HeySpace/file"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "HeySpace",
		Usage:    "在中英文之间添加空格",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Renzo",
				Email: "luyang.sun@outlook.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "输入文件路径",
				DefaultText: "默认剪贴板输入",
			},
			&cli.StringFlag{
				Name:        "out",
				Aliases:     []string{"o"},
				Usage:       "输出文件路径",
				DefaultText: "默认剪贴板输出",
			},
			&cli.StringFlag{
				Name:    "backup",
				Aliases: []string{"b"},
				Usage:   "备份目录路径",
			},
			&cli.BoolFlag{
				Name:        "server",
				Aliases:     []string{"s"},
				Usage:       "服务器监听模式",
				Value:       false,
				DefaultText: "关闭",
			},
			&cli.BoolFlag{
				Name:        "markdown",
				Aliases:     []string{"m"},
				Usage:       "Markdown 模式",
				Value:       true,
				DefaultText: "开启",
			},
			&cli.BoolFlag{
				Name:        "pdf",
				Aliases:     []string{"p"},
				Usage:       "PDF 模式",
				Value:       true,
				DefaultText: "开启",
			},
		},
		Action: func(c *cli.Context) error {
			config.GlobalConfig.MarkdownMode = c.Bool("markdown")
			config.GlobalConfig.ServerMode = c.Bool("server")
			config.GlobalConfig.PDFMode = c.Bool("pdf") // TODO

			if c.String("file") == "" {
				if !config.GlobalConfig.ServerMode {
					if err := clipboard.HandleClipboardInput(); err != nil {
						log.Fatal(err)
					}
					log.Println("已成功处理并发送至剪贴板，请复制查看")
				} else {
					startServe()
				}

			} else {
				// 若从文件输入，输出也会被设定为文件输出
				if err := file.HandlePathInput(c.String("file"), c.String("out"), c.String("backup")); err != nil {
					log.Fatal(err)
				}
				log.Println("已成功处理所有文件，请查看")
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// 服务模式：定时监听剪贴板，实时输出到剪贴板 TODO
func startServe() {

}
