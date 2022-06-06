package Sensitivefilter

//定义一个敏感词仓库
type WordStock interface {
	// ReadAll 获取所有的词语
	ReadAll() []string
}
