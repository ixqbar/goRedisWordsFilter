package filter

import (
	"fmt"
	"testing"
)

func TestWords(t *testing.T) {
	wordsMapping := NewWordsMapping()

	wordsMapping.AddNewWord([]rune("毛泽东"), false)
	wordsMapping.AddNewWord([]rune("毛片"), false)

	result := wordsMapping.FilterSentence([]rune("我看毛泽东的照片"))

	fmt.Println(string(result))

	if string(result) != "我看***的照片" {
		t.Error("filter fail")
	}
}
