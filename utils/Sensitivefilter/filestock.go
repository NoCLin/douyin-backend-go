package Sensitivefilter

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type FileStock struct {
	words []string
}

//从文件中读取出所有敏感词汇
func NewFileStock(p string) (*FileStock, error) {
	var s = &FileStock{}
	var rFile, err = os.OpenFile(p, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer rFile.Close()
	var reader = bufio.NewReader(rFile)
	var line []byte

	for {
		if line, _, err = reader.ReadLine(); err != nil {
			if err == io.EOF {
				break
			}
			log.Println("读取敏感词汇文件出错")
			return nil, err
		}
		var sLine = strings.TrimSpace(string(line))
		if sLine == "" {
			continue
		}
		s.words = append(s.words, sLine)
	}

	return s, nil
}

//实现WordStock接口
func (this *FileStock) ReadAll() []string {

	log.Printf("****** words ****** %d", len(this.words))

	return this.words
}
