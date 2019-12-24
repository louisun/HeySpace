package clipboard

import (
	"github.com/louisun/heyspace/config"
	"github.com/louisun/heyspace/space"

	"github.com/atotto/clipboard"
)

// 处理剪贴板输入
func HandleClipboardInput() error {
	inContent, err := clipboard.ReadAll()
	if err != nil {
		return err
	}
	var handler space.TextHandler
	if config.GlobalConfig.MarkdownMode {
		handler = space.NewMarkdownHandler(&inContent)
	} else {
		handler = space.NewPlainTextHander(&inContent)
	}
	return clipboard.WriteAll(handler.HandleText())
}
