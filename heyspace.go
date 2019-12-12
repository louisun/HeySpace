package main

import (
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"

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
			Config.MarkdownMode = c.Bool("markdown")
			Config.ServerMode = c.Bool("server")
			Config.PDFMode = c.Bool("pdf") // TODO

			if c.String("file") == "" {
				if !Config.ServerMode {
					if err := handleClipboardInput(); err != nil {
						log.Fatal(err)
					}
					log.Println("已成功处理并发送至剪贴板，请复制查看")
				} else {
					startServe()
				}

			} else {
				// 若从文件输入，输出也会被设定为文件输出 TODO
				if err := handleFileInput(c.String("file"), c.String("out"), c.String("backup")); err != nil {
					log.Fatal(err)
				}
				log.Println("已成功处理并发送至剪贴板，请复制查看")
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// 处理剪贴板输入
func handleClipboardInput() error {
	inContent, err := clipboard.ReadAll()
	if err != nil {
		return err
	}
	var handler TextHandler
	if Config.MarkdownMode {
		handler = NewMarkdownHandler(&inContent)
	} else {
		handler = NewPlainTextHander(&inContent)

	}
	return clipboard.WriteAll(handler.HandleText())
}

// 处理文件输入 TODO
func handleFileInput(inPath string, outPath string, backupPath string) error {
	// 输入为文件路径

	// 输入为目录路径（递归遍历处理）

	// 输出到指定目录（默认覆盖当前文件）

	// 原文件备份到备份内目录
	return nil
}

// 服务模式：定时监听剪贴板，实时输出到剪贴板 TODO
func startServe() {

}
