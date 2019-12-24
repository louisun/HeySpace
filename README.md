# HeySpace

HeySpace 专治空格强迫症，空格党 High 起来！


## 概述

> HeySpace 能帮你优化排版，提升阅读体验。

- 核心功能：**中英文之间添加空格**
- 去除连续两个以上的空行
- 兼容 Markdown 格式
  - 支持跳过代码块
  - 不影响 Markdown 特殊符号的使用
- 支持剪贴板输入输出，复制内容处理后可直接粘贴
- 支持文件、目录输入输出，安全备份

## TODO

- [x] 剪贴板输入 / 输出
- [x] 文件目录输入 / 输出
- [ ] 服务监听模式
- [ ] 纯文本非 Markdown 处理
- [ ] PDF 模式空格、换行处理


## 使用方式

```shell script
go get github.com/louisun/heyspace
```

```shell script
$ heyspace help
NAME:
   HeySpace - 在中英文之间添加空格

USAGE:
   HeySpace [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

AUTHOR:
   Renzo <luyang.sun@outlook.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value    输入文件路径 (default: 默认剪贴板输入)
   --out value, -o value     输出文件路径 (default: 默认剪贴板输出)
   --backup value, -b value  备份目录路径
   --server, -s              服务器监听模式 (default: 关闭)
   --markdown, -m            Markdown 模式 (default: 开启)
   --pdf, -p                 PDF 模式 (default: 开启)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```



```shell
# 这样就把剪贴板的内容加空格了，默认输入和输出都是 Markdown 格式
➜ heyspace

➜ heyspace --markdown # 等同于上面
```

> 效果请看 exapmples



## 说明

Shell 执行的话需要是 UTF-8 编码

```shell script
export LANG=en_US.UTF-8;
export LC_ALL=en_US.UTF-8;
```

附上我在 Mac OS 上 Automator 的脚本：

```shell script
export LANG=en_US.UTF-8;
export LC_ALL=en_US.UTF-8;
/Users/louisun/.local/bin/heyspace;
osascript -e "display notification \"${strPrompt}\" with title \"排版成功，请粘贴\" sound name \"default\"";
```

然后在 `设置>键盘>快捷键>服务>` 中对该脚本设置快捷键

复制内容后执行快捷键，效果如下：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2019-12-12-19-31-49_r80.png)


## 待修复 bug

- [ ] `(中文)` 在「非链接」情况下，依然需要加空格

