package utils

import "os"

// 判断所给路径是否存在
func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断路径是否是目录
func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// 判断路径是否是文件
func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// 判断是否存在该目录
func ExistsDir(path string) bool {
	return ExistsPath(path) && IsDir(path)
}

// 判断是否存在该文件
func ExistsFile(path string) bool {
	return ExistsPath(path) && IsFile(path)
}
