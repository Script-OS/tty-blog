package renderer

import (
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/wrap"
	"github.com/yuin/goldmark/ast"
	"strings"
	"tty-blog/cmd/view/renderer/style"
)

type Action struct {
	Position int
	Style    *style.Style
}

type BlockEnterType func(ctx *RenderContext, node ast.Node, source []byte) BlockDecorator
type BlockRenderType func(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string
type InlineEnterType func(ctx *RenderContext, node ast.Node, source []byte) (*style.Style, string)
type InlineRenderType func(ctx *RenderContext, node ast.Node, source []byte) (string, string)

type BlockItem struct {
	Enter  BlockEnterType
	Render BlockRenderType
}

type InlineItem struct {
	Enter  InlineEnterType
	Render InlineRenderType
}

func DecoratorWidth(ctx *RenderContext) int {
	width := 0
	for _, deco := range ctx.Deco {
		width += deco.Width()
	}
	return width
}

func Wrap(s string, limit int) string {
	wrapper := wrap.NewWriter(limit)
	wrapper.PreserveSpace = true
	_, _ = wrapper.Write([]byte(s))
	return wrapper.String()
}

func DefaultBlockRender(ctx *RenderContext, node ast.Node, width int, content string, actions []Action, source []byte) string {
	margin := DecoratorWidth(ctx)
	realWidth := width - margin
	styles := []style.Style{}
	for _, deco := range ctx.Deco {
		styles = append(styles, deco.Style())
	}
	thin := ctx.Deco[len(ctx.Deco)-1].Thin()
	decoLines := []string{}

	//head := ctx.Deco[len(ctx.Deco)-1].Push()
	//if head != "" {
	//	lines := strings.Split(head, "\n")
	//	for lineIndex, line := range lines {
	//		rendered := line + strings.Repeat(" ", realWidth+ctx.Deco[len(ctx.Deco)-1].Width()-ansi.PrintableRuneWidth(line))
	//		length := len(ctx.Deco) - 1
	//		for i, _ := range ctx.Deco {
	//			if i == 0 {
	//				continue
	//			}
	//			rendered = ctx.Deco[length-i].Deco(rendered, lineIndex, len(lines))
	//		}
	//		decoLines = append(decoLines, rendered)
	//	}
	//}

	//
	lines := strings.Split(Wrap(content, realWidth), "\n")
	index := 0
	actionIndex := 0
	for lineIndex, line := range lines {
		it := []rune(line)

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

		if !(len(lines) == lineIndex+1) || !thin {
			rendered += style.Render(styles, strings.Repeat(" ", realWidth-ansi.PrintableRuneWidth(line)))
		}

		length := len(ctx.Deco) - 1
		for i, _ := range ctx.Deco {
			rendered = ctx.Deco[length-i].Deco(rendered, lineIndex, len(line))
		}
		decoLines = append(decoLines, rendered)
	}

	tail := ctx.Deco[len(ctx.Deco)-1].Pop()
	if tail != "" {
		lines := strings.Split(tail, "\n")
		for lineIndex, line := range lines {
			rendered := line + strings.Repeat(" ", realWidth+ctx.Deco[len(ctx.Deco)-1].Width()-ansi.PrintableRuneWidth(line))
			length := len(ctx.Deco) - 1
			for i, _ := range ctx.Deco {
				if i == 0 {
					continue
				}
				rendered = ctx.Deco[length-i].Deco(rendered, lineIndex, len(lines))
			}
			decoLines = append(decoLines, rendered)
		}
	}

	if len(decoLines) == 0 {
		return ""
	}
	return strings.Join(decoLines, "\n") + "\n"
}
