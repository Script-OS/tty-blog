package webmedia

import (
	"encoding/base64"
	"strconv"
	"strings"
)

var InWebmediaTerm = false

type MediaDesc struct {
	Text  string
	Lines int
	Url   string
}

func init() {
	//if os.Getenv("TERM") == "xterm-webmedia-256color" {
	InWebmediaTerm = true
	//}
}

func makeOSCSeq(id int, payload string) string {
	return "\x1b]" + strconv.FormatInt(int64(id), 10) + ";" + payload + "\x1b\\"
}

func SetOSC8Link(link string) string {
	return makeOSCSeq(8, ";"+link)
}

func ResetWebmedia() string {
	return makeOSCSeq(9999, "")
}

func SetWebmediaLink(link string, textLen int) string {
	params := []string{
		"link",
		strconv.FormatInt(int64(textLen), 10),
		link,
	}
	return makeOSCSeq(9999, strings.Join(params, ";"))
}

func SetWebmediaMedia(desc *MediaDesc) string {
	params := []string{
		"media",
		base64.StdEncoding.EncodeToString([]byte(desc.Text)),
		strconv.FormatInt(int64(desc.Lines), 10),
		desc.Url,
	}
	return makeOSCSeq(9999, strings.Join(params, ";"))
}
