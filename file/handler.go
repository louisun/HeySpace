package file

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/louisun/heyspace/config"
	"github.com/louisun/heyspace/space"
	"github.com/louisun/heyspace/utils"

	"github.com/otiai10/copy"
)

// 处理文件和目录输入
func HandlePathInput(inPath string, outPath string, backupPath string) error {
	fstat, err := os.Stat(inPath)
	if err != nil {
		return err
	}
	if fstat.IsDir() {
		var allSuccess = true
		// 处理目录
		if backupPath == "" {
			// 备份路径为空，默认当前路径，inPath_bk
			backupPath = fmt.Sprintf("%s%s", inPath, config.DEFAULT_BACKUP_SUFFIX)
		}
		// step1. 备份目录到 backupPath
		if backupPath != config.NO_BACKUP_FLAG {
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

			if strings.HasSuffix(info.Name(), config.MARKDOWN_SUFFIX) {
				if err := HandleFileInput(path, "", config.NO_BACKUP_FLAG); err != nil {
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
		return HandleFileInput(inPath, outPath, backupPath)
	}
	return nil
}

// 处理文件输出
func HandleFileInput(inPath string, outPath string, backupPath string) error {
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
		if !strings.HasSuffix(outPath, config.MARKDOWN_SUFFIX) {
			// 传入输出目录
			if !utils.ExistsDir(outPath) {
				return errors.New("输出目录不存在")
			}
			outPath = fmt.Sprintf("%s%c%s%s%s", outPath, os.PathSeparator,
				strings.TrimSuffix(fstat.Name(), config.MARKDOWN_SUFFIX),
				config.DEFAULT_OUTPUT_SUFFIX, config.MARKDOWN_SUFFIX)
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

	var handler space.TextHandler
	inContent := string(inContentBytes)
	if config.GlobalConfig.MarkdownMode {
		handler = space.NewMarkdownHandler(&inContent)
	} else {
		handler = space.NewPlainTextHander(&inContent)
	}
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
	of.WriteString(handler.HandleText())

	if !config.GlobalConfig.QuietMode {
		log.Printf("【输入文件】: %s 【输出文件】: %s 【备份文件】: %s", inPath, outPath, backupPath)
	}
	return nil
}

// 设置文件备份路径
func setBackupFilePath(fstat os.FileInfo, inPath string, outPath string, backupPath *string, noBackup *bool) error {
	// 备份路径为空，默认当前路径
	if *backupPath == "" {
		*backupPath = fmt.Sprintf("%s%c%s%s%s", filepath.Dir(inPath), os.PathSeparator,
			strings.TrimSuffix(fstat.Name(), config.MARKDOWN_SUFFIX),
			config.DEFAULT_BACKUP_SUFFIX, config.MARKDOWN_SUFFIX)
		return nil
	}
	// 要求不备份
	if *backupPath == config.NO_BACKUP_FLAG {
		*noBackup = true
		*backupPath = "--"
		return nil
	}
	// 非 .md 结尾，默认备份路径为目录
	if !strings.HasSuffix(*backupPath, config.MARKDOWN_SUFFIX) {
		// 判断目录是否存在
		if !utils.ExistsDir(*backupPath) {
			return errors.New("备份目录不存在")
		}
		*backupPath = fmt.Sprintf("%s%c%s%s%s", backupPath, os.PathSeparator,
			strings.TrimSuffix(fstat.Name(), config.MARKDOWN_SUFFIX),
			config.DEFAULT_BACKUP_SUFFIX, config.MARKDOWN_SUFFIX)
	}
	// 文件路径与 inPath 或 outPath 相同
	if *backupPath == inPath || *backupPath == outPath {
		return errors.New("备份文件不能与输入文件或输出文件同名")
	}
	// 原 backup 路径不变
	return nil
}
