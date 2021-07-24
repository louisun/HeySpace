package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"

	"github.com/otiai10/copy"
)

const (
	NoBackupFlag        = "nobackup"
	DefaultBackupSuffix = "_bk"
	DefaultOutputSuffix = "_new"
	MarkdownSuffix      = ".md"
)

// markdown regex pre compiled
var (
	headerRegex      = regexp.MustCompile(`^(#+?\s)(.*)`)               // # header
	orderedListRegex = regexp.MustCompile(`^(\d+\.\s)(.*)`)             // 1. ordered list
	linkRegex        = regexp.MustCompile(`.*?([!]?\[.*?\]\(.*?\)).*?`) // ![]() OR []()
	linkDetailRegex  = regexp.MustCompile(`([!]?)\[(.*?)\]\((.*?)\)`)   // extract link text
)

var (
	codeFlag  bool
	emptyFlag bool
)

type matchIdxPair struct {
	startIdx int
	endIdx   int
}

func handleClipboard() error {
	// read content from clipboard
	content, err := clipboard.ReadAll()
	if err != nil {
		return err
	}

	return clipboard.WriteAll(FormatMarkdown(content))
}

func FormatMarkdown(text string) string {
	var (
		lines       = strings.Split(text, "\n")
		outputLines = make([]string, 0, len(lines))

		newLine string
	)

	for _, line := range lines {
		newLine = formatLine(line)

		// 位于代码块中
		if codeFlag {
			emptyFlag = false
			outputLines = append(outputLines, newLine)
			continue
		}

		// 位于非代码块中，连续多个空行只保留一个
		if strings.TrimSpace(newLine) == "" {
			if emptyFlag {
				continue
			}

			emptyFlag = true
		} else {
			emptyFlag = false
			outputLines = append(outputLines, newLine)
		}
	}

	return strings.Join(outputLines, "\n")
}

func formatLine(line string) string {
	// 位于代码块中
	if codeFlag {
		if strings.HasPrefix(line, "```") {
			codeFlag = false
		}
		return line
	}

	// 位于非代码块中
	if strings.HasPrefix(line, "```") {
		codeFlag = true
		return line
	}

	runes := []rune(line)

	// 引用
	if strings.HasPrefix(line, "> ") {
		return "> " + doFormat([]rune(strings.TrimLeft(string(runes[2:]), " ")))
	}
	// 标题
	if match := headerRegex.FindStringSubmatch(line); match != nil {
		return match[1] + doFormat([]rune(strings.TrimLeft(string(match[2]), " ")))
	}
	// 有序列表
	if match := orderedListRegex.FindStringSubmatch(line); match != nil {
		return match[1] + doFormat([]rune(strings.TrimLeft(string(match[2]), " ")))
	}
	// 无序列表
	if strings.HasPrefix(line, "- ") {
		return "- " + doFormat([]rune(strings.TrimLeft(string(runes[2:]), " ")))
	}
	// 包含链接的文本
	linkMatchIdx := linkRegex.FindAllStringSubmatchIndex(line, -1)
	if len(linkMatchIdx) != 0 {
		return doFormatWithLink(line, linkMatchIdx)
	}

	// 正常文本
	return doFormat([]rune(line))
}

func doFormatWithLink(line string, linkMatchIdx [][]int) string {
	pairs := make([]matchIdxPair, 0)
	for _, idxList := range linkMatchIdx {
		if len(idxList) == 0 {
			return doFormat([]rune(line))
		}

		if len(idxList)%2 != 0 {
			log.Println("idxList not in pairs")
			return doFormat([]rune(line))
		}

		start := 0
		end := 0

		// get (start, end) pairs
		for i, idx := range idxList {
			// skip the first and second index
			if i < 2 {
				continue
			}

			if i%2 == 0 {
				start = idx
			} else {
				end = idx
				pairs = append(pairs, matchIdxPair{
					startIdx: start,
					endIdx:   end,
				})
			}
		}
	}

	// like 0 .... (10 ... 20) ... (30 ... 40) ...
	resultBuf := bytes.Buffer{}
	buf := bytes.Buffer{}
	prevEndIdx := 0
	resultBuf.Grow(len(line)>>1 + len(line))

	for i, pair := range pairs {
		buf.Reset()

		// 处理 link 与文本之间的数据
		buf.WriteString(doFormat([]rune(line[prevEndIdx:pair.startIdx])))
		prevEndIdx = pair.endIdx

		// 处理 link 数据
		buf.WriteString(handleLinks(line[pair.startIdx:pair.endIdx]))

		// 处理最后的 link 与文本直接的数据
		if i == len(pairs)-1 {
			buf.WriteString(doFormat([]rune(line[pair.endIdx:])))
			prevEndIdx = pair.endIdx
		}

		resultBuf.WriteString(buf.String())
	}

	return resultBuf.String()
}

func doFormat(line []rune) string {
	var (
		preRune rune         // 前一个字符
		length  = len(line)  // 行字符数
		buffer  bytes.Buffer // 字符串缓冲区
	)

	var (
		italicCnt    = 0 // 斜体 * 计数
		boldCnt      = 0 // 粗体 ** 计数
		backQuoteCnt = 0 // 反引号 ` 计数
	)

	// buffer 写入方式：先写字符，后判断是否写入空格
	for idx, currentRune := range line {
		buffer.WriteRune(currentRune)

		// 相关符号数量统计
		switch currentRune {
		case '*':
			if preRune == '*' {
				boldCnt++
				italicCnt--
			} else {
				italicCnt++
			}
		case '`':
			backQuoteCnt++
		}

		// 判断是否要加空格
		if idx < length-1 {
			nextRune := line[idx+1]

			// 注：泛用英文不包括 Markdown 中的特殊符号 * ` [ ] ( )
			if isZh(currentRune) && isGeneralEn(nextRune) {
				// 中文 + 泛用英文 -> 加空格
				buffer.WriteString(" ")
			} else if isGeneralEn(currentRune) && isZh(nextRune) {
				// 泛用英文 + 中文 -> 加空格
				buffer.WriteString(" ")
			} else if (isZh(currentRune) && isEnLeftBracket(nextRune)) || (isEnRightBracket(currentRune) && isZh(nextRune)) {
				// 只用于这样的情况 “中文(” 或者 “)中文”，主要针对链接、图片等格式
				buffer.WriteString(" ")
			}

			// 有几种情况要特殊处理，核心是要分清在标点内部还是外部

			// 粗体中文**abc**
			// 斜体中文*abc*
			// 点中文`abc`
			if isZh(currentRune) {
				switch nextRune {
				case '*':
					doZhStar(&buffer, line, idx, boldCnt, italicCnt)
				case '`':
					doZhBackQuote(&buffer, line, idx, backQuoteCnt)
				}

				preRune = nextRune
				continue
			}

			// *abc*中文
			if currentRune == '*' && isZh(nextRune) {
				// * 之前的字符是英文则需要加空格
				// 区分 bold 和 italic
				switch preRune {
				case '*':
					doBoldStarZh(&buffer, line, idx, boldCnt)
				default:
					doSingleStarZh(&buffer, line, idx, italicCnt)
				}

				preRune = nextRune
				continue
			}

			if currentRune == '`' && isZh(nextRune) {
				doBackQuoteZh(&buffer, line, idx, backQuoteCnt)
				preRune = currentRune
				continue
			}
		}
	}

	return buffer.String()
}

func doZhStar(buffer *bytes.Buffer, line []rune, idx, boldCnt, italicCnt int) {
	length := len(line)
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
}

func doZhBackQuote(buffer *bytes.Buffer, line []rune, idx, backQuoteCnt int) {
	if idx < len(line)-2 {
		cn2 := line[idx+2]
		// 小代码块要看后面第二个字符是否是英文
		// 点中文`a`
		if backQuoteCnt%2 == 0 && isGeneralEn(cn2) {
			// 一个新代码块的开始
			buffer.WriteString(" ")
		}
	}
}

func doBoldStarZh(buffer *bytes.Buffer, line []rune, idx, boldCnt int) {
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
}

func doSingleStarZh(buffer *bytes.Buffer, line []rune, idx, italicCnt int) {
	if italicCnt%2 == 0 {
		// *abc*粗体中文
		if idx-1 > 0 && !isZh(line[idx-1]) {
			buffer.WriteString(" ")
		}
	}
}

func doBackQuoteZh(buffer *bytes.Buffer, line []rune, idx, backQuoteCnt int) {
	if backQuoteCnt%2 == 0 {
		if idx-1 > 0 && !isZh(line[idx-1]) {
			// `abc`点中文
			buffer.WriteString(" ")
		}
	}
}

func handleLinks(text string) string {
	if match := linkDetailRegex.FindStringSubmatch(text); len(match) > 3 {
		linkText := doFormat([]rune(match[2]))
		return fmt.Sprintf("%s[%s](%s)", match[1], linkText, match[3])
	}

	return text
}

func handleFileInput(inPath string, outPath string, backupPath string) error {
	fstat, err := os.Stat(inPath)
	if err != nil {
		return err
	}
	var noBackup bool
	// 处理文件：初始化 outPath、backupPath
	if outPath == "" {
		// 省略输出，则默认输出覆盖输入文件，需要备份文件
		outPath = inPath
		// 处理备份路径
		if err := setBackupFilePath(fstat, inPath, outPath, &backupPath, &noBackup); err != nil {
			return err
		}
	} else {
		// 指定输出路径
		if inPath == outPath {
			// 当输出和输入一致，且未指定 backup 为 nobackup 时，还是要设置备份路径
			if err := setBackupFilePath(fstat, inPath, outPath, &backupPath, &noBackup); err != nil {
				return err
			}
		} else {
			// 输入和输出不一致，不需要备份
			noBackup = true
			backupPath = "--"
		}
		// 非 .md 结尾，默认输出路径为目录
		if !strings.HasSuffix(outPath, MarkdownSuffix) {
			// 传入输出目录
			if !existsDir(outPath) {
				return errors.New("输出目录不存在")
			}
			outPath = fmt.Sprintf("%s%c%s%s%s", outPath, os.PathSeparator,
				strings.TrimSuffix(fstat.Name(), MarkdownSuffix),
				DefaultOutputSuffix, MarkdownSuffix)
		}
		// 其他情况 outPath 不处理
	}

	bf, err := os.Open(inPath)
	if err != nil {
		return err
	}
	// 内容处理
	inContentBytes, err := ioutil.ReadAll(bf)
	if err != nil {
		return err
	}

	inContent := string(inContentBytes)
	// 手动关闭输入文件（因为可能后面是覆盖该文件，要写入）
	bf.Close()

	// 备份
	if !noBackup {
		os.Rename(inPath, backupPath)
	}

	// 文本写入
	of, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer of.Close()
	of.WriteString(FormatMarkdown(inContent))

	if !globalConfig.QuietMode {
		log.Printf("【输入文件】: %s 【输出文件】: %s 【备份文件】: %s", inPath, outPath, backupPath)
	}
	return nil
}

// 处理文件和目录输入
func handlePathInput(inPath string, outPath string, backupPath string) error {
	fstat, err := os.Stat(inPath)
	if err != nil {
		return err
	}
	if fstat.IsDir() {
		var allSuccess = true
		// 处理目录
		if backupPath == "" {
			// 备份路径为空，默认当前路径，inPath_bk
			backupPath = fmt.Sprintf("%s%s", inPath, DefaultBackupSuffix)
		}
		// step1. 备份目录到 backupPath
		if backupPath != NoBackupFlag {
			log.Printf("目录已备份，备份路径 %s", backupPath)
			copy.Copy(inPath, backupPath)
		}

		// step2. 遍历 inPath，逐个替换文件
		log.Println("开始处理目录中的 Markdown 文件...")
		filepath.Walk(inPath, func(path string, info os.FileInfo, err error) error {
			// 忽略目录
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(info.Name(), MarkdownSuffix) {
				if err := handleFileInput(path, "", NoBackupFlag); err != nil {
					allSuccess = false
					return err
				}
			}

			// 忽略非 Markdown 文件
			return nil
		})
		if !allSuccess {
			return errors.New("目录处理未成功，请从备份目录恢复")
		}
	} else {
		return handleFileInput(inPath, outPath, backupPath)
	}
	return nil
}

// 设置文件备份路径
func setBackupFilePath(fstat os.FileInfo, inPath string, outPath string, backupPath *string, noBackup *bool) error {
	// 备份路径为空，默认当前路径
	if *backupPath == "" {
		*backupPath = fmt.Sprintf("%s%c%s%s%s", filepath.Dir(inPath), os.PathSeparator,
			strings.TrimSuffix(fstat.Name(), MarkdownSuffix),
			DefaultBackupSuffix, MarkdownSuffix)
		return nil
	}
	// 要求不备份
	if *backupPath == NoBackupFlag {
		*noBackup = true
		*backupPath = "--"
		return nil
	}
	// 非 .md 结尾，默认备份路径为目录
	if !strings.HasSuffix(*backupPath, MarkdownSuffix) {
		// 判断目录是否存在
		if !existsDir(*backupPath) {
			return errors.New("备份目录不存在")
		}
		*backupPath = fmt.Sprintf("%s%c%s%s%s", backupPath, os.PathSeparator,
			strings.TrimSuffix(fstat.Name(), MarkdownSuffix),
			DefaultBackupSuffix, MarkdownSuffix)
	}
	// 文件路径与 inPath 或 outPath 相同
	if *backupPath == inPath || *backupPath == outPath {
		return errors.New("备份文件不能与输入文件或输出文件同名")
	}
	// 原 backup 路径不变
	return nil
}
