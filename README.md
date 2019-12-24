# HeySpace

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2019-12-24-20-29-53_r93.png)

## 概述

`HeySpace` 是一款帮你优化排版，提升阅读体验的「**命令行工具**」。

- 核心功能：**中英文之间添加空格**
- 去除连续两个以上的空行
- 兼容 Markdown 格式
- 支持剪贴板输入输出，复制内容处理后可直接粘贴
- **支持文件、目录的输入和输出**；支持文件备份

> 以上内容均只在 Mac OS 下测试过。

### TODO

- [x] 剪贴板输入 / 输出
- [x] 文件目录输入 / 输出
- [ ] 服务监听模式
- [ ] 纯文本非 Markdown 处理
- [ ] PDF 模式空格、换行处理

### 效果展示

排版前：

> 因为 Go 的 `net/http` 包提供了基础的路由函数组合与丰富的功能函数。所以在社区里流行一种用 Go 编写 API 不需要框架的观点，在我们看来，如果你的项目的路由在个位数、URI 固定且不通过 URI 来传递参数，那么确实使用官方库也就足够。但在复杂场景下，官方的 http 库还是有些力有不逮。

排版后：

> 因为 Go 的 `net/http` 包提供了基础的路由函数组合与丰富的功能函数。所以在社区里流行一种用 Go 编写 API 不需要框架的观点，在我们看来，如果你的项目的路由在个位数、URI 固定且不通过 URI 来传递参数，那么确实使用官方库也就足够。但在复杂场景下，官方的 http 库还是有些力有不逮。

同时不影响 Markdown 的符号的正常使用，包括：



~~~markdown
# 标题

1. 有序列表一
2. 有序列表二

- 无序列表一
- 无序列表二

*斜体* abc

**粗体** abc

> 引用：**念奴娇·赤壁怀古**
>
> 大江东去，浪淘尽，千古风流人物。故垒西边，人道是：三国周郎赤壁。乱石穿空，惊涛拍岸，卷起千堆雪。江山如画，一时多少豪杰。
>
> 遥想公瑾当年，小乔初嫁了，雄姿英发。羽扇纶巾，谈笑间、樯橹灰飞烟灭。故国神游，多情应笑我，早生华发。人生如梦，一樽还酹江月。

`small code block` 小代码块

```go
fmt.Println("大代码块")
```
~~~

> 在各种复杂的应用场景下，Markdown 文本的各类符号都（也许）能够得到有效排版（欢迎提 bug），这也是我写这个工具的初衷，因为其他中英文排版工具没有对 Markdown 文本作特殊处理。
>
> -   比如粗体、斜体、小代码块的内容，与周边的词直接的关系；
>
> -   比如跳过「代码块」内容的处理

## 使用方式

```shell script
go get github.com/louisun/heyspace
```

```shell script
$ heyspace help
NAME:
   HeySpace - 在中英文之间添加空格

USAGE:
   heyspace [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

AUTHOR:
   Renzo <luyang.sun@outlook.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --in value, -i value      输入文件路径 (default: 默认剪贴板输入)
   --out value, -o value     输出文件路径 (default: 默认剪贴板输出)
   --backup value, -b value  备份目录路径
   --server, -s              服务器监听模式 (default: 关闭)
   --quiet, -q               不输出具体文件日志 (default: 关闭)
   --markdown, -m            Markdown 模式 (default: 开启)
   --pdf, -p                 PDF 模式 (default: 开启)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

### 无参数

不传入 `-i` 或 `--in` 参数，默认是从「**剪贴板**」获取输入，并输出到「**剪贴板**」便于粘贴。

```go
➜ heyspace

➜ heyspace --markdown # 等同于上面
```

-   步骤一：Copy 一段内容
-   步骤二：执行 `heyspace`
-   步骤三：Paste 处理之后的内容

> 而只要指定了 `-i` 参数，则以文件模式输出

### 输入为文件

`-i` 或 `--in` 参数如果是个文件，则会：

-   先备份文件，以防处理失败
-   若指定输出路径 `-o` 或 `-out`，则将处理后的内容写入输出路径

```
heyspace -i ~/Blog/Artical.md
```

> 上面的例子未指定输出路径，则默认覆盖 `~/Blog/Artical.md` 的内容，且由于未指定备份路径，会自动将原文件备份为  `~/Blog/Artical_bk.md` 文件。

```bash
heyspace -i ~/Blog/Artical.md -o ~/Blog/Artical_new.md
```

> 上面的例子指定了输出路径为文件，则不进行备份（原文件即备份），处理排版后输出的文件为指定路径

```bash
heyspace -i ~/Blog/Aritical.md -o ~/Blog
```

> 上面的子指定了输出路径为目录，则不进行备份，且默认输出文件名为原名字加 `_new`，即  `~/Blog/Artical_new.md`

### 输入为目录

`-i` 或 `--in` 参数如果是个路径，则会：

-   先备份整个目录，以防处理失败
-   遍历目录下的 Markdown 文件，直接替换每个文件为排版后的内容

```bash
heyspace -i ~/Blog
```

> 上面的例子未指定备份目录，自动备份目录到 `~/Blog_bk`，即原名字加后缀 `_bk` 的目录。然后对 `~/Blog` 目录下所有 `.md` 结尾的文本文件进行排版处理。

`-b` 或 `--backup` 参数表示备份路径

```shell
heyspace -i ~/Blog -b ~/MyBlog
```

> 上面的例子指定备份目录，则备份目录到 `~/MyBlog`，备份目录名可不存在。

```
heyspace -i ~/Blog -b nobackup
```

> 上面的例子表示不进行备份，原地替换 `~/Blog` 下的 Markdown 文件（谨慎使用）。

## 作者常用

下面介绍我是怎么用这个工具的

### 处理整个目录

```bash
# -q 或 --quiet 参数可以不输出一长串文件信息

# 暴力替换某个目录下的文件：
heyspace -i ~/Blog -b nobackup -q

# 做个备份
heyspace -i ~/Blog -b ~/somewhere -q
```

### 快捷键：剪贴板替换

> 再次提醒，我是在 Mac OS 下使用的，需要 Mac OS `Automator` 工具的帮助。

我用 `ctrl + command + z` 快捷键，执行这个 `heyspace` 脚本，可快速完成剪贴板内容的替换。

附上 shell 脚本：

```shell script
# 设置编码
export LANG=en_US.UTF-8;
export LC_ALL=en_US.UTF-8;

# 执行 heyspace
/Users/louisun/.local/bin/heyspace;

# Mac 执行完发送通知的命令
osascript -e "display notification \"${strPrompt}\" with title \"排版成功，请粘贴\" sound name \"default\"";
```

然后在 ` 设置 > 键盘 > 快捷键 > 服务 >` 中对该脚本设置快捷键

复制内容后执行快捷键，效果如下：

![](https://bucket-1255905387.cos.ap-shanghai.myqcloud.com/2019-12-12-19-31-49_r80.png)

## 待修复 bug

- [ ] `(中文)` 在「非链接」情况下，依然需要加空格