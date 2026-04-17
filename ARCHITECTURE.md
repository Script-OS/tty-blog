# TTY-BLOG 架构文档

## 一、项目概述

TTY-BLOG 是一个终端博客管理系统，提供类似 Shell 的交互式环境，用于浏览、查看和编辑博客内容。项目模拟类 Unix Shell 的命令行体验，专注于 Markdown 文件的终端渲染展示。

## 二、预期功能

### 2.1 核心功能

| 功能 | 命令 | 描述 |
|------|------|------|
| 目录浏览 | `ls [dir]` | 查看博客目录结构，显示文件修改时间和名称 |
| 目录切换 | `cd <dir>` | 在博客目录树中导航 |
| 文件查看 | `view <file>` | 以终端美化方式渲染 Markdown 文件，支持分页和滚动 |
| 文件编辑 | `edit <file>` | 调用外部编辑器编辑文件（需 editor 权限） |
| 用户切换 | `su [guest|editor]` | 切换用户身份，editor 模式需密码验证 |
| 帮助信息 | `help <cmd>` | 查看命令帮助 |

### 2.2 用户角色

| 角色 | 权限 | 说明 |
|------|------|------|
| guest | 查看、浏览 | 默认身份，可执行 ls、cd、view |
| editor | 编辑 | 需密码验证，可执行 edit 命令 |

### 2.3 Markdown 渲染特性

- 标题：分级样式渲染（一级标题特殊背景色）
- 引用块：左侧竖线装饰
- 列表：项目符号和嵌套支持
- 代码块：语法高亮（chroma）
- 表格：ASCII 边框绘制
- 链接：OSC8 超链接协议支持
- 图片：Webmedia 终端协议内嵌显示
- 强调：粗体、斜体样式
- 行内代码：背景色高亮

### 2.4 交互特性

- 分页浏览：PageUp/PageDown、上下箭头、鼠标滚轮
- 自动补全：路径和命令自动补全
- 终端适配：自动检测终端宽度

## 三、技术架构

### 3.1 技术栈

| 组件 | 库 | 用途 |
|------|-----|------|
| Markdown 解析 | `github.com/yuin/goldmark` | 解析 Markdown 文本为 AST |
| Markdown 扩展 | `github.com/yuin/goldmark-emoji` | Emoji 支持 |
| 代码高亮 | `github.com/alecthomas/chroma` | 代码块语法高亮 |
| 终端样式 | `github.com/muesli/termenv` | ANSI 样式渲染 |
| 终端宽度 | `github.com/muesli/ansi` | 可打印字符宽度计算 |
| 文本换行 | `github.com/muesli/reflow` | 终端宽度自适应换行 |
| 命令行输入 | `github.com/chzyer/readline` | REPL 循环、自动补全 |
| 颜色处理 | `github.com/lucasb-eyer/go-colorful` | 颜色空间转换 |
| 字符宽度 | `github.com/mattn/go-runewidth` | 字符显示宽度计算 |
| 配置解析 | `gopkg.in/yaml.v3` | YAML 配置文件解析 |

### 3.2 模块结构

```
tty-blog/
├── main.go                  # 程序入口、REPL 循环
├── dispatch.go              # 命令分发器
├── help.go                  # 帮助命令实现
├── go.mod / go.sum          # Go 模块定义
├── README.md                # 用户手册
├── ARCHITECTURE.md          # 架构文档
│
├── global/                  # 全局状态模块
│   ├── vars.go              # 全局变量定义
│   ├── config.go            # 配置加载
│   └── completer.go         # 路径自动补全
│
└── cmd/                     # 命令实现模块
    ├── ls/cmd.go            # ls 命令
    ├── cd/cmd.go            # cd 命令
    ├── su/cmd.go            # su 命令
    ├── edit/cmd.go          # edit 命令
    │
    └── view/                 # view 命令模块
        ├── cmd.go           # view 命令入口
        ├── viewport.go      # 分页显示控制器
        │
        └── renderer/         # Markdown 渲染器
            ├── renderer.go  # 核心渲染器框架
            ├── basic.go     # 默认块渲染逻辑
            ├── utils.go     # Markdown 样式定义
            ├── table.go     # 表格渲染逻辑
            │
            ├── style/       # 样式处理
            │   └── style.go # 样式合并与渲染
            │
            ├── common/      # 通用组件
            │   └── base.go  # BlockDecorator 实现
            │
            ├── input/       # 输入处理
            │   ├── key.go   # 键盘事件
            │   └── mouse.go # 鼠标事件
            │
            └── webmedia/    # Webmedia 协议
                └── webmedia.go # OSC 序列生成
```

### 3.3 模块职责

#### 入口层

| 文件 | 职责 |
|------|------|
| `main.go` | 初始化配置、注册命令、运行 REPL 循环、解析用户输入 |
| `dispatch.go` | 命令注册表管理、命令路由分发 |
| `help.go` | 帮助命令实现、动态调用其他命令的 `-help` |

#### 全局状态层（global）

| 文件 | 职责 |
|------|------|
| `vars.go` | 定义全局状态：用户身份、工作目录、文件系统根 |
| `config.go` | 加载 `~/.blogrc` 配置文件、配置默认值合并 |
| `completer.go` | 基于 fs.FS 的路径自动补全实现 |

#### 命令层（cmd）

| 命令 | 职责 | 数据流向 |
|------|------|----------|
| `ls` | 读取目录、按修改时间排序、格式化输出 | fs.FS → stdout |
| `cd` | 验证路径、更新工作目录 | WorkDir 更新 |
| `su` | 密码验证、切换用户身份 | User 更新 |
| `edit` | 权限检查、调用外部编辑器 | 外部进程 |
| `view` | 读取文件、渲染 Markdown、分页显示 | fs.FS → 渲染器 → viewport |

#### 渲染层（cmd/view/renderer）

| 文件 | 职责 |
|------|------|
| `renderer.go` | goldmark AST 遍历、渲染上下文管理、Block/Inline 处理器注册 |
| `basic.go` | 默认块渲染逻辑：换行、样式应用、装饰器嵌套 |
| `utils.go` | Markdown 各元素样式定义、goldmark 初始化 |
| `table.go` | 表格特殊渲染：列宽计算、边框绘制 |
| `style/style.go` | 样式栈合并、ANSI 序列生成、OSC8/Webmedia 链接 |
| `common/base.go` | BlockDecorator 接口实现 |
| `input/key.go` | 键盘事件解析、键码映射 |
| `input/mouse.go` | 鼠标事件解析（X10 协议） |
| `webmedia/webmedia.go` | OSC 序列生成：超链接、图片嵌入 |

### 3.4 数据流

#### REPL 主循环

```
┌─────────────────────────────────────────────────────────┐
│                      main.go                             │
├─────────────────────────────────────────────────────────┤
│  1. 初始化配置 (global.Config)                           │
│  2. 初始化文件系统 (global.Root = os.DirFS(rootDir))    │
│  3. 注册命令 (RegisterCommand)                          │
│  4. 进入 REPL 循环                                      │
│     ┌────────────────────────────────────────────┐      │
│     │  readline.Readline() → 用户输入            │      │
│     │  strings.Split → 解析命令和参数            │      │
│     │  Dispatch(args) → 命令分发                 │      │
│     └────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────┘
```

#### view 命令渲染流程

```
Markdown 文件
    │
    ▼
fs.ReadFile(global.Root, path)
    │
    ▼
goldmark.Parser → AST (Document)
    │
    ▼
TermRenderer.render (AST 遍历)
    │  ┌─────────────────────────────────────────┐
    │  │ Block 元素: Heading, Blockquote, List...│
    │  │   Enter → 创建 BlockDecorator           │
    │  │   Render → 应用样式 + 装饰器            │
    │  ├─────────────────────────────────────────┤
    │  │ Inline 元素: Text, Link, CodeSpan...    │
    │  │   Enter → 创建 Style Action             │
    │  │   Render → 追加文本                     │
    │  └─────────────────────────────────────────┘
    │
    ▼
ANSI 文本输出 (带样式序列)
    │
    ▼
RenderInPage (viewport.go)
    │  ┌─────────────────────────────────────────┐
    │  │ 进入 Alternate Screen                   │
    │  │ 进入 Raw Mode                           │
    │  │ 渲染可见行                              │
    │  │ 循环处理输入事件                        │
    │  │   ↑/↓/PgUp/PgDown → 滚动               │
    │  │   鼠标滚轮 → 滚动                       │
    │  │   q/Ctrl+C/ESC → 退出                  │
    │  │ 退出 Alternate Screen                   │
    │  └─────────────────────────────────────────┘
    │
    ▼
返回主 REPL
```

### 3.5 关键接口

#### 命令接口

```go
type CMD = func(args []string)
```

所有命令实现为 `func(args []string)`，通过 `RegisterCommand` 注册到分发器。

#### Block 装饰器接口

```go
type BlockDecorator interface {
    Deco(line string, lineNo int, lineNum int) string  // 行装饰
    Style() style.Style                                // 内部样式
    Push() string                                      // 块开始额外内容
    Pop() string                                       // 块结束额外内容
    Thin() bool                                        // 是否精简模式（末行无填充）
    Width() int                                        // 装饰占用宽度
}
```

用于为块级元素（标题、引用、列表、代码块）添加边框、前缀、后缀等视觉装饰。

#### Block/Inline 处理器

```go
type BlockItem struct {
    Enter  BlockEnterType  // 进入节点时创建装饰器
    Render BlockRenderType // 离开节点时渲染内容
}

type InlineItem struct {
    Enter  InlineEnterType  // 进入节点时创建样式
    Render InlineRenderType // 离开节点时追加文本
}
```

处理器注册到 `TermRenderer.BlockProc` 和 `TermRenderer.InlineProc`，由 AST 遍历器自动调用。

#### 样式系统

```go
type Style = map[int]interface{}

const (
    Foreground int = iota  // colorful.Color
    Background             // colorful.Color
    Bold                   // bool
    Italic                 // bool
    CrossOut               // bool
    Underline              // bool
    Overline               // bool
    Link                   // string (OSC8 URL)
    Media                  // *webmedia.MediaDesc
)
```

样式以栈形式管理，渲染时合并栈中所有样式并转换为 ANSI 序列。

### 3.6 配置系统

配置文件路径：`~/.blogrc`

```yaml
editor: ["nano", "-R"]    # 编辑器命令及参数
editorPassword: ""        # editor 模式密码
rootDir: "."              # 博客根目录
```

配置加载流程：
1. 解析 `~/.blogrc` → `ConfigType`
2. 合合默认值（空字段使用默认）
3. 存储至 `global.Config`

### 3.7 终端协议支持

#### OSC8 超链接

标准终端超链接协议，支持在支持的终端中点击链接：

```
ESC ] 8 ; ; <url> ESC \  <text>  ESC ] 8 ; ; ESC \
```

#### Webmedia 协议

自定义 OSC 序列（ID 9999），用于增强终端功能：

| 功能 | 序列格式 |
|------|----------|
| 链接设置 | `ESC ] 9999 ; link ; <len> ; <url> ESC \` |
| 图片嵌入 | `ESC ] 9999 ; media ; <id> ; <base64> ; <lines> ; <url> ESC \` |
| 清理图片 | `ESC ] 9999 ; cleanMedia ESC \` |
| 重置 | `ESC ] 9999 ; ESC \` |

## 四、运行机制

### 4.1 程序启动

1. 加载配置文件
2. 初始化文件系统（`global.Root = os.DirFS(rootDir)`）
3. 注册命令到分发器
4. 显示 Banner
5. 进入 REPL 循环

### 4.2 REPL 循环

```
循环:
  1. 设置窗口标题
  2. 计算提示符路径
  3. 等待用户输入 (readline)
  4. 空行 → 继续
  5. 解析命令和参数
  6. Dispatch → 执行命令
  7. 返回步骤 1

退出条件:
  - EOF (Ctrl+D)
  - ErrInterrupt (Ctrl+C)
```

### 4.3 分页显示机制

viewport 使用 Alternate Screen Buffer，独立于主 REPL 屏幕：

1. 进入 Alternate Screen
2. 设置终端 Raw Mode（禁用行缓冲）
3. 启用鼠标事件报告
4. 渲染可见区域
5. 循环处理输入事件，更新滚动位置
6. 退出时恢复终端状态

### 4.4 路径处理

所有路径通过 `global.CalcPath` 统一处理：

```go
func CalcPath(path string) string {
    // 相对路径 → 基于 WorkDir 转绝对路径
    // 清理路径（去除冗余分隔符）
    // 转换为 fs.FS 可用格式
}
```

## 五、扩展点

### 5.1 添加新命令

1. 在 `cmd/<name>/cmd.go` 创建命令模块
2. 定义 `Name` 常量和 `Run` 函数
3. 定义 `Completer`（可选）
4. 在 `main.go` 中 `RegisterCommand`

### 5.2 扩展 Markdown 渲染

1. 在 `renderer/utils.go` 的 `initMarkdownStyle` 中注册处理器
2. 或在 `renderer/` 下创建独立处理器文件

### 5.3 自定义样式

修改 `renderer/utils.go` 中的样式定义：

```go
headingStyle := style.Style{
    style.Foreground: easyHex("#7f7fff"),
    style.Bold:       true,
}
```