package renderer

import (
	"fmt"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"tty-blog/cmd/view/renderer/common"
	"tty-blog/cmd/view/renderer/style"
)

type BlockDecorator interface {
	Deco(line string, lineNo int, lineNum int) string
	Style() style.Style
	Push() string
	Pop() string
	Thin() bool
	Width() int
}

type BlankDeco struct{}

func (*BlankDeco) Deco(line string, lineNo int, lineNum int) string { return line }
func (*BlankDeco) Style() style.Style                               { return style.Style{} }
func (*BlankDeco) Push() string                                     { return "" }
func (*BlankDeco) Pop() string                                      { return "" }
func (*BlankDeco) Thin() bool                                       { return false }
func (*BlankDeco) Width() int                                       { return 0 }

type RenderContext struct {
	Deco      []BlockDecorator
	localLine string
	width     int
	actions   []Action
	Meta      map[string]interface{}
}

func (ctx *RenderContext) Reset() {
	ctx.Meta = map[string]interface{}{}
	ctx.Deco = []BlockDecorator{}
	ctx.localLine = ""
	ctx.width = 0
	ctx.actions = []Action{}
}

type TermRenderer struct {
	Width      int
	ctx        *RenderContext
	BlockProc  map[ast.NodeKind]BlockItem
	InlineProc map[ast.NodeKind]InlineItem
}

func NewTermRenderer() *TermRenderer {
	return &TermRenderer{
		Width:      80,
		ctx:        &RenderContext{},
		BlockProc:  map[ast.NodeKind]BlockItem{},
		InlineProc: map[ast.NodeKind]InlineItem{},
	}
}

func (r *TermRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.render)
	reg.Register(ast.KindHeading, r.render)
	reg.Register(ast.KindBlockquote, r.render)
	reg.Register(ast.KindCodeBlock, r.render)
	reg.Register(ast.KindFencedCodeBlock, r.render)
	reg.Register(ast.KindHTMLBlock, r.render)
	reg.Register(ast.KindList, r.render)
	reg.Register(ast.KindListItem, r.render)
	reg.Register(ast.KindParagraph, r.render)
	reg.Register(ast.KindTextBlock, r.render)
	reg.Register(ast.KindThematicBreak, r.render)

	// inlines
	reg.Register(ast.KindAutoLink, r.render)
	reg.Register(ast.KindCodeSpan, r.render)
	reg.Register(ast.KindEmphasis, r.render)
	reg.Register(ast.KindImage, r.render)
	reg.Register(ast.KindLink, r.render)
	reg.Register(ast.KindAutoLink, r.render)
	reg.Register(ast.KindRawHTML, r.render)
	reg.Register(ast.KindText, r.render)
	reg.Register(ast.KindString, r.render)

	// tables
	reg.Register(astext.KindTable, r.render)
	reg.Register(astext.KindTableHeader, r.render)
	reg.Register(astext.KindTableRow, r.render)
	reg.Register(astext.KindTableCell, r.render)

	// definitions
	reg.Register(astext.KindDefinitionList, r.render)
	reg.Register(astext.KindDefinitionTerm, r.render)
	reg.Register(astext.KindDefinitionDescription, r.render)

	// footnotes
	reg.Register(astext.KindFootnote, r.render)
	reg.Register(astext.KindFootnoteList, r.render)
	reg.Register(astext.KindFootnoteLink, r.render)
	reg.Register(astext.KindFootnoteBacklink, r.render)

	// checkboxes
	reg.Register(astext.KindTaskCheckBox, r.render)

	// strikethrough
	reg.Register(astext.KindStrikethrough, r.render)

	// emoji
	reg.Register(east.KindEmoji, r.render)
}

func (r *TermRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	debug := false
	if !debug {
		if entering {
			if node.Type() == ast.TypeBlock {
				if proc, ok := r.BlockProc[node.Kind()]; ok && proc.Enter != nil {
					deco := proc.Enter(r.ctx, node, source)
					r.ctx.Deco = append(r.ctx.Deco, deco)
				} else {
					r.ctx.Deco = append(r.ctx.Deco, &common.BlockDecorator{})
				}
				r.ctx.localLine = ""
				r.ctx.width = 0
				r.ctx.actions = []Action{}
			} else if node.Type() == ast.TypeInline {
				action := &style.Style{}
				s := ""
				if proc, ok := r.InlineProc[node.Kind()]; ok && proc.Enter != nil {
					action, s = proc.Enter(r.ctx, node, source)
				}
				r.ctx.actions = append(r.ctx.actions, Action{
					Position: r.ctx.width,
					Style:    action,
				})
				r.ctx.width += len([]rune(s))
				r.ctx.localLine += s
			} else if node.Type() == ast.TypeDocument {
				r.ctx.Reset()
			}
		} else {
			if node.Type() == ast.TypeBlock {
				line := r.ctx.localLine
				if proc, ok := r.BlockProc[node.Kind()]; ok && proc.Render != nil {
					line = proc.Render(r.ctx, node, r.Width, line, r.ctx.actions, source)
				} else {
					line = DefaultBlockRender(r.ctx, node, r.Width, line, r.ctx.actions, source)
				}
				_, err := w.WriteString(line)
				if err != nil {
					return ast.WalkStop, err
				}
				r.ctx.Deco = r.ctx.Deco[:len(r.ctx.Deco)-1]
				r.ctx.localLine = ""
				r.ctx.width = 0
				r.ctx.actions = []Action{}
			} else if node.Type() == ast.TypeInline {
				line := ""
				suffix := ""
				if proc, ok := r.InlineProc[node.Kind()]; ok && proc.Render != nil {
					line, suffix = proc.Render(r.ctx, node, source)
				}
				r.ctx.width += len([]rune(line))
				r.ctx.localLine += line
				r.ctx.actions = append(r.ctx.actions, Action{
					Position: r.ctx.width,
					Style:    nil,
				})
				r.ctx.width += len([]rune(suffix))
				r.ctx.localLine += suffix
			} else if node.Type() == ast.TypeDocument {
				r.ctx.Reset()
			}
		}
	} else {
		if node.Type() == ast.TypeInline {
			//_, err := w.Write([]byte(fmt.Sprintf("<%s>%s</%s>", node.Kind().String(), string(node.Text(source)), node.Kind().String())))
			//if err != nil {
			//	return ast.WalkStop, err
			//}
			if entering {
				_, err := w.Write([]byte(fmt.Sprintf("<%s>", node.Kind().String())))
				if err != nil {
					return ast.WalkStop, err
				}
			} else {
				if node.Kind() == ast.KindText {
					_, err := w.Write(node.Text(source))
					if err != nil {
						return ast.WalkStop, err
					}
				}
				_, err := w.Write([]byte(fmt.Sprintf("</%s>", node.Kind().String())))
				if err != nil {
					return ast.WalkStop, err
				}
			}
		} else {
			if entering {
				_, err := w.Write([]byte(fmt.Sprintf("<%s>", node.Kind().String())))
				if err != nil {
					return ast.WalkStop, err
				}
			} else {
				_, err := w.Write([]byte(fmt.Sprintf("\n</%s>\n", node.Kind().String())))
				if err != nil {
					return ast.WalkStop, err
				}
			}
		}
	}
	return ast.WalkContinue, nil
}
