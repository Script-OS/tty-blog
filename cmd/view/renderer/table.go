package renderer

import (
	"github.com/muesli/ansi"
	"github.com/yuin/goldmark/ast"
	"strings"
	"tty-blog/cmd/view/renderer/common"
	"tty-blog/cmd/view/renderer/style"
)

func EnterTable(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
	meta := []int{}
	for row := node.FirstChild(); row != nil; row = row.NextSibling() {
		index := 0
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			width := ansi.PrintableRuneWidth(strings.Replace(string(cell.Text(source)), "\n", " ", -1))
			if len(meta) <= index {
				meta = append(meta, width)
			} else {
				if meta[index] < width {
					meta[index] = width
				}
			}
			index += 1
		}
	}
	ctx.Meta["table"] = meta
	return &common.BlockDecorator{
		Prefix:     func(lineNo int, lineNum int) string { return "  " },
		LinePrefix: " ",
	}
}

func RenderTable(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string {
	delete(ctx.Meta, "table")
	delete(ctx.Meta, "row")
	delete(ctx.Meta, "col")
	return DefaultBlockRender(ctx, node, width, content, actions, source)
}

func RenderHeading(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string {
	line := ""
	for _, width := range ctx.Meta["table"].([]int) {
		line += "┼" + strings.Repeat("─", width)
	}
	return DefaultBlockRender(ctx, node, width, ctx.Meta["row"].(string), actions, source) + DefaultBlockRender(ctx, node, width, string([]rune(line)[1:]), actions, source)
}

func EnterRow(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator {
	ctx.Meta["row"] = ""
	ctx.Meta["col"] = 0
	return &BlankDeco{}
}

func RenderRow(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string {
	return DefaultBlockRender(ctx, node, width, ctx.Meta["row"].(string), actions, source)
}

func RenderCell(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string {
	index := 0
	actionIndex := 0
	styles := []style.Style{}
	for _, deco := range ctx.Deco {
		styles = append(styles, deco.Style())
	}

	it := []rune(content)

	rendered := ""
	for len(it) > 0 {
		part := len(it)
		var actionRef *Action = nil
		if actionIndex < len(actions) {
			action := actions[actionIndex]
			if action.Position <= index+len(it) {
				part = action.Position - index
				//index = action.Position
				actionRef = &action
				actionIndex += 1
			}
		}
		index += part

		if part != 0 {
			rendered += style.Render(styles, string(it[0:part]))
			it = it[part:]
		}

		if actionRef != nil {
			if actionRef.Style != nil {
				styles = append(styles, *actionRef.Style)
			} else {
				styles = styles[:len(styles)-1]
			}
		}
	}

	for ; actionIndex < len(actions); actionIndex += 1 {
		action := actions[actionIndex]
		if action.Position <= index {
			if action.Style != nil {
				styles = append(styles, *action.Style)
			} else {
				styles = styles[:len(styles)-1]
			}
		} else {
			break
		}
	}

	rendered += style.Render(styles, strings.Repeat(" ", ctx.Meta["table"].([]int)[ctx.Meta["col"].(int)]-ansi.PrintableRuneWidth(rendered)))

	prefix := ""
	if ctx.Meta["col"].(int) != 0 {
		prefix = "│"
	}
	ctx.Meta["row"] = ctx.Meta["row"].(string) + prefix + rendered
	ctx.Meta["col"] = ctx.Meta["col"].(int) + 1
	return ""
}
