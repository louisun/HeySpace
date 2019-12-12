package main

type ModeConfig struct {
	MarkdownMode bool // markdown 模式
	ServerMode   bool // server监听模式
	PDFMode      bool // pdf 模式
}

var Config *ModeConfig = &ModeConfig{
	MarkdownMode: true,
	ServerMode:   false,
	PDFMode:      false,
}
