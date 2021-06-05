package space

import (
	"bytes"
	"regexp"
	"strings"
)

type TextHandler interface {
	HandleText() string
	handleLine(string) string
}

type MarkdownHandler struct {
	Text             *string
	CodeFlag         bool
	EmptyFlag        bool
	HeaderRegex      *regexp.Regexp
	OrderedListRegex *regexp.Regexp
}

type PlainTextHandler struct {
	Text *string
}

func NewMarkdownHandler(text *string) *MarkdownHandler {
	return &MarkdownHandler{
		Text:             text,
		CodeFlag:         false,
		EmptyFlag:        false,
		HeaderRegex:      regexp.MustCompile(`^(#+?\s)(.*)`),
		OrderedListRegex: regexp.MustCompile(`^(\d+\.\s)(.*)`),
	}
}

func NewPlainTextHander(text *string) *PlainTextHandler {
	return &PlainTextHandler{
		Text: text,
	}
}

// Markdown 处理文本
func (content *MarkdownHandler) HandleText() string {
	originalLines := strings.Split(*content.Text, "\n")
	outputLines := make([]string, 0, len(originalLines))
	var newLine string
	for _, line := range originalLines {
		newLine = content.handleLine(line)
		// 非代码块状态下，连续多个空行，只保留一个
		if strings.TrimSpace(newLine) == "" && !content.CodeFlag {
			if content.EmptyFlag {
				continue
			}
			newLine = ""
			content.EmptyFlag = true
		} else {
			content.EmptyFlag = false
		}
		outputLines = append(outputLines, newLine)
	}
	return strings.Join(outputLines, "\n")

}

// Markdown 处理一行文本
func (content *MarkdownHandler) handleLine(line string) string {
	if !content.CodeFlag {
		if strings.HasPrefix(line, "```") {
			// 代码块开头
			content.CodeFlag = true
			return line
		} else {
			// 非代码块内容的处理
			uLine := []rune(line)
			if strings.HasPrefix(line, "> ") {
				// 引用
				return "> " + content.handleBlock([]rune(strings.TrimLeft(string(uLine[2:]), " ")))
			} else if match := content.HeaderRegex.FindStringSubmatch(line); match != nil {
				// 标题
				return match[1] + content.handleBlock([]rune(strings.TrimLeft(string(match[2]), " ")))
			} else if match := content.OrderedListRegex.FindStringSubmatch(line); match != nil {
				// 有序列表
				return match[1] + content.handleBlock([]rune(strings.TrimLeft(string(match[2]), " ")))
			} else if strings.HasPrefix(line, "- ") {
				// 无序列表
				return "- " + content.handleBlock([]rune(strings.TrimLeft(string(uLine[2:]), " ")))
			} else {
				// 正常文本
				return content.handleBlock([]rune(line))
			}
		}
	} else {
		if strings.HasPrefix(line, "```") {
			// 代码块结尾
			content.CodeFlag = false
		}
		// 处于代码块之间的内容直接返回
		return line
	}
}

func (content *MarkdownHandler) handleBlock(line []rune) string {
	var (
		buffer       bytes.Buffer
		outputString string
		length       = len(line)

		italicCnt   = 0  // 斜体 * 的计数
		boldCnt     = 0  // 粗体 ** 的计数
		backtickCnt = 0  // 反引号 ` 的计数
		preRune     rune // 前一个字符
	)

	// buffer 写入方式：先写字符，后判断是否写入空格
	for idx, currentRune := range line {
		buffer.WriteRune(currentRune)

		// 1: 计数
		if currentRune == '*' {
			if preRune == '*' {
				boldCnt++
				italicCnt--
			} else {
				italicCnt++
			}
		}

		if currentRune == '`' {
			backtickCnt++
		}

		// 2. 判断是否要加空格
		if idx < length-1 {
			nextRune := line[idx+1]

			if isZh(currentRune) && isGeneralEn(nextRune) {
				// 中文 + 泛用英文 -> 加空格
				buffer.WriteString(" ")
			} else if isGeneralEn(currentRune) && isZh(nextRune) {
				// 泛用英文 + 中文 -> 加空格
				buffer.WriteString(" ")
			} else if (isZh(currentRune) && isEnLeftBracket(nextRune)) || (isEnRightBracket(currentRune) && isZh(nextRune)) {
				// 由于链接、图片等格式，只用于这样的情况 )中文(
				buffer.WriteString(" ")
			}

			// 有几种情况要特殊处理，核心是要分清在标点内部还是外部

			if isZh(currentRune) {
				// 粗体中文**abc**
				// 斜体中文*abc*
				// 点中文`abc`
				if nextRune == '*' {
					if idx < length-2 {
						cn2 := line[idx+2]
						if cn2 == '*' {
							// 粗体要看后面第三个字符是否是英文
							if idx < length-3 {
								cn3 := line[idx+3]
								// 粗体中文**a**
								if boldCnt%2 == 0 && isGeneralEn(cn3) {
									// 一个新粗体的开始
									buffer.WriteString(" ")
								}
							}
						} else {
							// 斜体要看后面第二个字符是否是英文
							// 斜体中文*a*
							if italicCnt%2 == 0 && isGeneralEn(cn2) {
								// 一个新斜体的开始
								buffer.WriteString(" ")
							}
						}
					}
				} else if nextRune == '`' {
					if idx < length-2 {
						cn2 := line[idx+2]
						// 小代码块要看后面第二个字符是否是英文
						// 点中文`a`
						if backtickCnt%2 == 0 && isGeneralEn(cn2) {
							// 一个新代码块的开始
							buffer.WriteString(" ")
						}
					}
				}

			} else if currentRune == '*' && isZh(nextRune) {
				// *abc*斜体中文

				// 需要判断 * 之前的是否是英文，中文不需要加空格
				if preRune == '*' {
					if boldCnt%2 == 0 {
						// **abc**粗体中文
						if idx-2 > 0 && line[idx-2] != '*' && !isZh(line[idx-2]) {
							buffer.WriteString(" ")
						}

						// ***abc***粗体中文
						if line[idx-2] == '*' && idx-3 > 0 && !isZh(line[idx-3]) {
							buffer.WriteString(" ")
						}
					}
				} else {
					if italicCnt%2 == 0 {
						// *abc*粗体中文
						if idx-1 > 0 && !isZh(line[idx-1]) {
							buffer.WriteString(" ")
						}
					}
				}
			} else if currentRune == '`' && isZh(nextRune) {
				if backtickCnt%2 == 0 {
					if idx-1 > 0 && !isZh(line[idx-1]) {
						// `abc`点中文
						buffer.WriteString(" ")
					}
				}
			}

			preRune = currentRune
		}
	}
	outputString = buffer.String()
	return outputString
}

// Markdown 处理文本
func (content *PlainTextHandler) HandleText() string {
	return ""
}

func (content *PlainTextHandler) handleLine(line string) string {
	return ""
}
