package common

import (
	"github.com/go-ego/gse"
)

var seg gse.Segmenter

func Load() {
	if !seg.Load {
		seg, _ = gse.New("./data/dict/zh/s_1.txt,./data/dict/zh/user.txt")
	}
}

func Segment(keywords string) []string {
	Load()
	//删除重复
	allKeys := make(map[string]bool)
	list := []string{}
	segments := seg.Segment([]byte(keywords))
	for _, s := range segments {
		//过滤字符
		text := gse.FilterSymbol(s.Token().Text())
		if text != "" {
			// pos := s.Token().Pos()
			// if pos != "j" {
			if _, value := allKeys[text]; !value {
				allKeys[text] = true
				list = append(list, text)
			}
			// }
		}
	}
	return list
}
