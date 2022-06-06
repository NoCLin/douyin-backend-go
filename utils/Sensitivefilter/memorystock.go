package Sensitivefilter

//从内存入的方式获取敏感词汇
type MemoryStock struct {
	words []string
}

func NewMemoryStock(words ...string) (*MemoryStock, error) {
	var s = &MemoryStock{}
	s.words = words
	return s, nil
}

func (this *MemoryStock) ReadAll() []string {
	return this.words
}
