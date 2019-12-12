# HeySpace

HeySpace 专治空格强迫症，空格党 High 起来！

`[]~(￣▽￣)~*`

## 使用方式

```shell script
➜ heyspace help
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


- [x] 剪贴板输入 / 输出
- [ ] 文件目录输入 / 输出
- [ ] 服务监听模式
- [ ] 纯文本非 Markdown 处理
- [ ] PDF 模式空格、换行处理

