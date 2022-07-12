package renderer

import (
	"bytes"
	"github.com/alecthomas/chroma/quick"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"strings"
	"tty-blog/cmd/view/renderer/common"
	"tty-blog/cmd/view/renderer/style"
)

func easyHex(s string) colorful.Color {
	ret, _ := colorful.Hex(s)
	return ret
}

func oncePrefix(special string, common string) func(lineNo int, lineNum int) string {
	trigger := false
	return func(lineNo int, lineNum int) string {
		if trigger {
			return common
		} else {
			if lineNo != lineNum && lineNo == 0 {
				trigger = true
			}
			return special
		}
	}
}

func initMarkdownStyle(renderer *TermRenderer) {
	firstHeadingStyle := style.Style{
		style.Foreground: easyHex("#ffff87"),
		style.Background: easyHex("#5f5fff"),
		style.Bold:       true,
	}
	headingStyle := style.Style{
		style.Foreground: easyHex("#7f7fff"),
		style.Bold:       true,
	}
	// heading
	renderer.BlockProc[ast.KindHeading] = BlockItem{
		Enter: func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
			level := node.(*ast.Heading).Level
			if level == 1 {
				return &common.BlockDecorator{
					Prefix: func(lineNo int, lineNum int) string {
						return style.Render([]style.Style{firstHeadingStyle}, " ")
					},
					Suffix: func(lineNo int, lineNum int) string {
						return style.Render([]style.Style{firstHeadingStyle}, " ")
					},
					InnerStyle: firstHeadingStyle,
					IsThin:     true,
					LinePrefix: " ",
					LineSuffix: " ",
				}
			} else {
				return &common.BlockDecorator{
					Prefix: oncePrefix(
						style.Render([]style.Style{headingStyle}, strings.Repeat("#", level)+" "),
						style.Render([]style.Style{headingStyle}, strings.Repeat(" ", level)+" "),
					),
					InnerStyle: headingStyle,
					LineSuffix: " ",
				}
			}
		},
	}
	// block-quote
	renderer.BlockProc[ast.KindBlockquote] = BlockItem{
		Enter: func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
			return &common.BlockDecorator{
				Prefix: func(lineNo int, lineNum int) string { return "│ " },
			}
		},
	}
	// list
	renderer.BlockProc[ast.KindList] = BlockItem{
		Enter: func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
			return &common.BlockDecorator{
				Prefix:     func(lineNo int, lineNum int) string { return "  " },
				LineSuffix: " ",
			}
		},
	}
	// list-item
	renderer.BlockProc[ast.KindListItem] = BlockItem{
		Enter: func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
			return &common.BlockDecorator{
				Prefix: oncePrefix("• ", "  "),
			}
		},
	}
	// fenced-code-block
	renderer.BlockProc[ast.KindFencedCodeBlock] = BlockItem{
		Enter: func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
			return &common.BlockDecorator{
				Prefix: func(lineNo int, lineNum int) string { return termenv.Style{}.Foreground(termenv.ANSIWhite).Styled("") },
				Suffix: func(lineNo int, lineNum int) string { return termenv.Style{}.Foreground(termenv.ANSIWhite).Styled("") },
			}
		},
		Render: func(ctx *RenderContext, node ast.Node, width int, content string, action []Action, source []byte) string {
			buf := &bytes.Buffer{}
			n := node.(*ast.FencedCodeBlock)
			lines := n.Lines().Len()
			code := ""
			for i := 0; i < lines; i += 1 {
				it := n.Lines().At(i)
				code += string(it.Value(source))
			}
			lang := string(node.(*ast.FencedCodeBlock).Language(source))
			quick.Highlight(buf, code, lang, "terminal256", "monokai")
			return DefaultBlockRender(ctx, node, width, buf.String(), action, source)
		},
	}
	// table
	renderer.BlockProc[astext.KindTable] = BlockItem{
		Enter:  EnterTable,
		Render: RenderTable,
	}
	// heading
	renderer.BlockProc[astext.KindTableHeader] = BlockItem{
		Enter:  EnterRow,
		Render: RenderHeading,
	}
	// row
	renderer.BlockProc[astext.KindTableRow] = BlockItem{
		Enter:  EnterRow,
		Render: RenderRow,
	}
	// cell
	renderer.BlockProc[astext.KindTableCell] = BlockItem{
		Render: RenderCell,
	}

	// text
	renderer.InlineProc[ast.KindText] = InlineItem{
		Render: func(ctx *RenderContext, node ast.Node, source []byte) string {
			//return string(node.Text(source))
			return strings.Replace(string(node.Text(source)), "\n", " ", -1)
		},
	}
	// emphasis
	renderer.InlineProc[ast.KindEmphasis] = InlineItem{
		Enter: func(ctx *RenderContext, node ast.Node) (*style.Style, string) {
			level := node.(*ast.Emphasis).Level
			if level == 1 {
				return &style.Style{
					style.Italic: true,
				}, ""
			} else {
				return &style.Style{
					style.Bold: true,
				}, ""
			}
		},
	}
	// code-span
	renderer.InlineProc[ast.KindCodeSpan] = InlineItem{
		Enter: func(ctx *RenderContext, node ast.Node) (*style.Style, string) {
			return &style.Style{
				style.Foreground: easyHex("#ff5f5f"),
				style.Background: easyHex("#303030"),
			}, " "
		},
		Render: func(ctx *RenderContext, node ast.Node, source []byte) string {
			return " "
		},
	}
	// task-check-box
	renderer.InlineProc[astext.KindTaskCheckBox] = InlineItem{
		Render: func(ctx *RenderContext, node ast.Node, source []byte) string {
			if node.(*astext.TaskCheckBox).IsChecked {
				return "[✓] "
			} else {
				return "[ ] "
			}
		},
	}
	// link
	renderer.InlineProc[ast.KindLink] = InlineItem{
		Render: func(ctx *RenderContext, node ast.Node, source []byte) string {
			href := string(node.(*ast.Link).Destination)
			return " →  " + href
		},
	}
	// media
	renderer.InlineProc[ast.KindImage] = InlineItem{
		Render: func(ctx *RenderContext, node ast.Node, source []byte) string {
			href := string(node.(*ast.Image).Destination)
			return " →  " + href
		},
	}
}

func New(width int) goldmark.Markdown {
	ret := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	termRenderer := NewTermRenderer()
	termRenderer.Width = width
	initMarkdownStyle(termRenderer)
	ret.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(termRenderer, 1000),
			),
		),
	)
	return ret
}

func EasyRender(renderer goldmark.Markdown, text []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := renderer.Convert(text, &buf)
	return buf.Bytes(), err
}
