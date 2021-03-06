package style

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/ansi"
	"github.com/muesli/termenv"
	"tty-blog/cmd/view/renderer/webmedia"
)

const (
	Foreground int = iota // colorful.Color
	Background            // colorful.Color
	Bold                  // bool
	Italic                // bool
	CrossOut              // bool
	Underline             // bool
	Overline              // bool
	Link                  // string
	Media                 // *webmedia.MediaDesc
	MAX
)

type Style = map[int]interface{}

var profile = termenv.ColorProfile()

func toTerm(style Style) termenv.Style {
	ret := termenv.Style{}
	if v, ok := style[Foreground]; ok {
		ret = ret.Foreground(profile.Color(v.(colorful.Color).Hex()))
	}
	if v, ok := style[Background]; ok {
		ret = ret.Background(profile.Color(v.(colorful.Color).Hex()))
	}
	if v, ok := style[Bold]; ok && v.(bool) {
		ret = ret.Bold()
	}
	if v, ok := style[Italic]; ok && v.(bool) {
		ret = ret.Italic()
	}
	if v, ok := style[CrossOut]; ok && v.(bool) {
		ret = ret.CrossOut()
	}
	if v, ok := style[Underline]; ok && v.(bool) {
		ret = ret.Underline()
	}
	if v, ok := style[Overline]; ok && v.(bool) {
		ret = ret.Overline()
	}
	return ret
}

func merge(styles []Style) Style {
	ret := Style{}
	length := len(styles) - 1
	for i, _ := range styles {
		style := styles[length-i]
		for k := 0; k < MAX; k += 1 {
			if _, ok := ret[k]; !ok {
				if v, ok := style[k]; ok {
					ret[k] = v
				}
			}
		}
	}
	return ret
}

func Render(styles []Style, text string) string {
	merged := merge(styles)
	style := toTerm(merged)
	rendered := style.Styled(text)
	if meta, ok := merged[Link]; ok {
		link := meta.(string)
		if webmedia.InWebmediaTerm {
			rendered = webmedia.SetWebmediaLink(link, ansi.PrintableRuneWidth(text)) + rendered
		} else {
			rendered = webmedia.SetOSC8Link(link) + rendered + webmedia.SetOSC8Link("")
		}
	}
	if meta, ok := merged[Media]; ok {
		desc := meta.(*webmedia.MediaDesc)
		if webmedia.InWebmediaTerm {
			rendered += webmedia.SetWebmediaMedia(desc)
		}
	}
	return rendered
}
