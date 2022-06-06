package config

import (
	"github.com/NoCLin/douyin-backend-go/utils/Sensitivefilter"
	"log"
)

func initSensitiveTree() *Sensitivefilter.TrieFilter {

	fs, err := Sensitivefilter.NewFileStock("utils/Sensitivefilter/senstiveStock.txt")
	if err != nil {
		log.Println(err)
	}
	var filter = Sensitivefilter.NewTrieFilter(fs)
	filter.Excludes('-', ' ')

	return filter
}
