package config

type ModeConfig struct {
	MarkdownMode bool // markdown 模式
	ServerMode   bool // server监听模式
	PDFMode      bool // pdf 模式
	QuietMode    bool // 不输出具体文件日志
}

var GlobalConfig *ModeConfig = &ModeConfig{
	MarkdownMode: true,
	ServerMode:   false,
	PDFMode:      false,
	QuietMode:    false,
}
