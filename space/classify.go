package space

import (
	"unicode"
)

// 是否为中文字符
func isZh(c rune) bool {
	if c >= '\u4e00' && c <= '\u9fa5' {
		return true
	}
	return false
}

// 排除 * () ` 等特殊符号
func isGeneralEn(c rune) bool {
	if isDigit(c) {
		return true
	}
	if isAlpha(c) {
		return true
	}
	if isGeneralEnSymbol(c) {
		return true
	}
	return false
}

// 是否为数字
func isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

// 是否为英文字母
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// 是否为泛用英文符号
func isGeneralEnSymbol(c rune) bool {
	enSymbols := []rune{
		':', ';', '%', '!', '?', '°', '_',
		'<', '=', '>', '"', '$', '&', '\'', ',', '.',
		'/', '@', '\\', '^', '|',
	}
	for _, r := range enSymbols {
		if c == r {
			return true
		}
	}
	return false
}

// 是否为 Markdown 特殊的英文符号
func isMarkdownEnSymbol(c rune) bool {
	enSymbols := []rune{
		'*', '`',
	}
	for _, r := range enSymbols {
		if c == r {
			return true
		}
	}
	return false
}

// 是否为英文左括号
func isEnLeftBracket(c rune) bool {
	if c == '(' || c == '[' {
		return true
	}
	return false
}

// 是否为英文右括号
func isEnRightBracket(c rune) bool {
	if c == ')' || c == ']' {
		return true
	}
	return false
}
