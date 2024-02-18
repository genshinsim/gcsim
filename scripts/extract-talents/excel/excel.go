package excel

var textMap map[TextMapHash]string

func init() {
	load("TextMap/TextMapEN.json", &textMap)
}

type TextMapHash uint32

func (h TextMapHash) String() string {
	return textMap[h]
}
