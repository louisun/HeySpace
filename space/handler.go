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
				return match[1] + content.handleBlock([]rune(match[2]))
			} else if match := content.OrderedListRegex.FindStringSubmatch(line); match != nil {
				// 有序列表
				return match[1] + content.handleBlock([]rune(match[2]))
			} else if strings.HasPrefix(line, "- ") {
				// 无序列表
				return "- " + content.handleBlock(uLine[2:])
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
	var buffer bytes.Buffer
	var outputString string
	var length = len(line)
	var italicCnt = 0
	var boldCnt = 0
	var pointCnt = 0
	var preRune rune
	for i, c := range line {
		// 先写入字符
		buffer.WriteRune(c)
		if c == '*' {
			if preRune == '*' {
				boldCnt++
				italicCnt--
			} else {
				italicCnt++
			}
		}
		if c == '`' {
			pointCnt++
		}
		if i < length-1 {
			cn := line[i+1]
			// 每次看看后面是什么字符，以判断是否需要加空格
			//log.Println("当前字符：", string(c), "  下个一字符：", string(cn))
			if isZh(c) && isGeneralEn(cn) {
				// 中文 + 泛用英文 -> 加空格
				buffer.WriteString(" ")
			} else if isGeneralEn(c) && isZh(cn) {
				// 泛用英文 + 中文 -> 加空格
				buffer.WriteString(" ")
			} else if (isZh(c) && isEnLeftBracket(cn)) || (isEnRightBracket(c) && isZh(cn)) {
				// 由于链接、图片等格式，只用于这样的情况 )中文(
				buffer.WriteString(" ")
			}
			// 有几种情况要特殊处理，核心是要分清在标点内部还是外部

			if isZh(c) {
				// 粗体中文**abc**
				// 斜体中文*abc*
				// 点中文`abc`
				if cn == '*' {
					if i < length-2 {
						cn2 := line[i+2]
						if cn2 == '*' {
							// 粗体要看后面第三个字符是否是英文
							if i < length-3 {
								cn3 := line[i+3]
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
				} else if cn == '`' {
					if i < length-2 {
						cn2 := line[i+2]
						// 小代码块要看后面第二个字符是否是英文
						// 点中文`a`
						if pointCnt%2 == 0 && isGeneralEn(cn2) {
							// 一个新代码块的开始
							buffer.WriteString(" ")
						}
					}
				}

			} else if c == '*' && isZh(cn) {
				// **abc**粗体中文
				// *abc*斜体中文
				if preRune == '*' {
					if boldCnt%2 == 0 {
						buffer.WriteString(" ")
					}
				} else {
					if italicCnt%2 == 0 {
						buffer.WriteString(" ")
					}
				}
			} else if c == '`' && isZh(cn) {
				// `abc`点中文
				buffer.WriteString(" ")
			}

			preRune = c
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
