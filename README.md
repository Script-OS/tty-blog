# TTY-BLOG

> 终端博客管理系统 - 在命令行中浏览、查看和管理你的博客

TTY-BLOG 提供类似 Unix Shell 的交互式环境，专为 Markdown 博客内容设计。支持终端美化的 Markdown 渲染、分页浏览和权限管理。

## 安装

```bash
# 克隆项目
git clone https://github.com/Script-OS/tty-blog.git
cd tty-blog

# 构建
go build -o tty-blog .

# 运行
./tty-blog
```

## 配置

配置文件路径：`~/.blogrc`

```yaml
# 编辑器命令及参数（可选，默认 nano -R）
editor: ["nano", "-R"]

# editor 模式密码（可选，默认无密码保护）
editorPassword: "your-secret-password"

# 博客根目录（可选，默认当前目录）
rootDir: "/path/to/your/blog"
```

### 配置示例

```yaml
editor: ["vim"]
editorPassword: "mypassword"
rootDir: "~/my-blog"
```

## 基本使用

启动后进入 REPL 环境，提示符格式：

```
guest@blog:/current/path> _
```

### 可用命令

| 命令 | 用法 | 说明 |
|------|------|------|
| `ls` | `ls [目录]` | 列出目录内容，按修改时间排序 |
| `cd` | `cd <目录>` | 切换当前目录 |
| `view` | `view <文件.md>` | 查看 Markdown 文件（美化渲染） |
| `edit` | `edit <文件>` | 编辑文件（需要 editor 权限） |
| `su` | `su [guest|editor]` | 切换用户身份 |
| `help` | `help [命令]` | 显示帮助信息 |
| `?` | `?` | 显示可用命令列表 |

### 命令详解

#### ls - 目录列表

```
guest@blog:/articles> ls
 TIME               │ NAME
────────────────────┼──────
2024-01-15 10:30:00 │ tech-guide.md
2024-01-10 14:20:00 │ life-notes
2024-01-05 09:00:00 │ drafts
```

- 目录显示为蓝色
- 隐藏文件（以 `.` 开头）不显示
- 按修改时间降序排列（最新在前）

#### cd - 目录切换

```
guest@blog:/articles> cd drafts
guest@blog:/articles/drafts> cd ..
guest@blog:/articles> cd /
guest@blog:/>
```

支持：
- 相对路径：`cd subdir`、`cd ..`
- 绝对路径：`cd /articles`
- 自动补全：输入部分路径后按 `Tab`

#### view - 查看 Markdown

```
guest@blog:/articles> view tech-guide.md
```

进入分页浏览模式，Markdown 内容美化渲染显示。

#### edit - 编辑文件

```
editor@blog:/articles> edit tech-guide.md
```

- 仅 `editor` 用户可执行
- 调用配置的编辑器程序
- 退出编辑器后返回 REPL

#### su - 用户切换

```
guest@blog:/> su editor
input password: ****
editor@blog:/>

editor@blog:/> su guest
guest@blog:/>
```

- 默认切换到 `editor`
- `editor` 身份需要密码验证
- `guest` 无需密码

#### help - 获取帮助

```
guest@blog:/> help
Usage of help:
  help <cmd>

guest@blog:/> help view
Usage of view:
  view <file>
```

## view 模式交互

在 `view` 命令的分页浏览模式中：

### 键盘操作

|按键 | 功能 |
|------|------|
| `↑` / `k` | 向上滚动一行 |
| `↓` / `j` | 向下滚动一行 |
| `PgUp` | 向上翻页 |
| `PgDown` | 向下翻页 |
| `q` / `ESC` / `Ctrl+C` | 退出返回 REPL |

### 鼠标操作

| 操作 | 功能 |
|------|------|
| 滚轮向上 | 向上滚动 |
| 滚轮向下 | 向下滚动 |

## Markdown 渲染特性

TTY-BLOG 支持丰富的 Markdown 元素渲染：

### 基本元素

| 元素 | 渲染效果 |
|------|----------|
| 一级标题 | 蓝底黄字高亮块 |
| 其他标题 | 蓝色粗体，带 `#` 前缀 |
| 引用块 | 左侧竖线 `│` 装饰 |
| 列表 | 项目符号 `•` 和缩进 |
| 任务列表 | `[✓]` / `[ ]` 复选框 |

### 代码块

支持语法高亮（chroma），使用 `monokai` 主题：

```
```go
func main() {
    fmt.Println("Hello")
}
```
```

### 表格

ASCII 边框表格渲染：

```
  Name        │ Age  │ City
──────────────┼──────┼──────
  Alice       │ 25   │ NYC
  Bob         │ 30   │ LA
```

### 强调样式

| 语法 | 效果 |
|------|------|
| `*文本*` | 斜体 |
| `**文本**` | 粗体 |
| `~~文本~~` | 删除线 |
| `` `代码` `` | 红底行内代码 |

### 链接与图片

- **链接**：OSC8 超链接协议，支持终端点击（兼容终端）
- **图片**：Webmedia 协议内嵌显示（需要兼容终端）

## 用户角色权限

| 角色 | 权限 | 说明 |
|------|------|------|
| `guest` | 查看、浏览 | 默认身份，可执行 `ls`、`cd`、`view` |
| `editor` | 编辑 | 需密码验证，可执行 `edit` 命令 |

## 终端兼容性

### 推荐终端

- iTerm2（macOS）
- Kitty
- WezTerm
- 支持 OSC8 的现代终端

### 终端功能支持

| 功能 | 终端要求 |
|------|----------|
| ANSI 颜色 | 支持 256 色 |
| OSC8 超链接 | 终端需支持 OSC8 协议 |
| 鼠标滚轮 | 终端需支持 X10 鼠标协议 |
| 图片嵌入 | 需要 Webmedia 兼容终端 |

## 自动补全

- 命令自动补全：输入命令前几个字母按 `Tab`
- 路径自动补全：输入部分路径按 `Tab`
- 补全会自动添加 `/`（目录）或空格（文件）

## 退出程序

在 REPL 中：
- `Ctrl+D` - 正常退出
- `Ctrl+C` - 中断退出

## 示例工作流

```
# 启动
$ ./tty-blog

# 浏览博客目录
guest@blog:/> ls
guest@blog:/> cd articles

# 查看文章
guest@blog:/articles> view hello-world.md

# 切换到编辑模式
guest@blog:/articles> su editor
input password: ****

# 编辑文章
editor@blog:/articles> edit hello-world.md

# 返回访客模式
editor@blog:/articles> su guest

# 退出
guest@blog:/articles> Ctrl+D
```

## 相关文档

- [架构文档](ARCHITECTURE.md) - 系统架构和模块设计