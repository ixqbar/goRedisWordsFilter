package filter

import (
	"bufio"
	cuckoo "github.com/seiflotfy/cuckoofilter"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type WordsItem struct {
	total  int
	detail map[int][]string
}

func NewWordsItem() *WordsItem {
	return &WordsItem{
		total:  0,
		detail: make(map[int][]string, 0),
	}
}

func (obj *WordsItem) AddNewWord(word []rune) bool {
	wordLen := len(word)
	if _, ok := obj.detail[wordLen]; !ok {
		obj.detail[wordLen] = make([]string, 0)
	}

	if InStringArray(string(word), obj.detail[wordLen]) {
		return false
	}

	obj.total += 1
	obj.detail[wordLen] = append(obj.detail[wordLen], string(word))

	return true
}

func (obj *WordsItem) DeleteWord(word []rune) bool {
	wordLen := len(word)

	if _, ok := obj.detail[wordLen]; !ok {
		return false
	}

	foundIndex := -1
	for i, v := range obj.detail[wordLen] {
		if v == string(word) {
			foundIndex = i
			break
		}
	}

	if foundIndex >= 0 {
		obj.detail[wordLen][foundIndex] = obj.detail[wordLen][len(obj.detail[wordLen])-1]
		obj.detail[wordLen][len(obj.detail[wordLen])-1] = ""
		obj.detail[wordLen] = obj.detail[wordLen][:len(obj.detail[wordLen])-1]
		obj.total--
		return true
	}

	return false
}

type WordsMapping struct {
	sync.Mutex
	total  int
	detail map[rune]*WordsItem
	cf     *cuckoo.Filter
}

func NewWordsMapping() *WordsMapping {
	return &WordsMapping{
		total:  0,
		detail: make(map[rune]*WordsItem, 0),
		cf:     cuckoo.NewFilter(1000),
	}
}

func (obj *WordsMapping) Clear() {
	obj.Lock()
	defer obj.Unlock()

	if obj.total == 0 {
		return
	}

	obj.total = 0
	obj.detail = make(map[rune]*WordsItem, 0)
	obj.cf.Reset()
}

func (obj *WordsMapping) Load() {
	startTime := time.Now()
	file, err := os.Open(GConfig.DictWordsPath)
	if err != nil {
		Logger.Print(err)
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		lineContent, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			Logger.Print(err)
			continue
		}

		lineWord := strings.TrimSpace(lineContent)
		if len(lineWord) == 0 {
			continue
		}

		obj.AddNewWord([]rune(lineWord), false)
	}

	Logger.Printf("load found %d words time elapsed %v", obj.TotalWords(), time.Since(startTime))
}

func (obj *WordsMapping) ExistsWord(word string) bool {
	obj.Lock()
	defer obj.Unlock()

	return obj.cf.Lookup([]byte(word))
}

func (obj *WordsMapping) AddNewWord(word []rune, autoSave bool) bool {
	obj.Lock()
	defer obj.Unlock()

	if autoSave {
		if obj.cf.Lookup([]byte(string(word))) == false {
			go func() {
				file, err := os.OpenFile(GConfig.DictWordsPath, os.O_WRONLY|os.O_APPEND, 0644);
				if err != nil {
					Logger.Print(err)
					return
				}

				defer file.Close()
				file.WriteString(string(word) + "\n")
			}()
		}
	}

	firstWordRune := word[0]

	if _, ok := obj.detail[firstWordRune]; !ok {
		obj.detail[firstWordRune] = NewWordsItem()
	}

	if obj.detail[firstWordRune].AddNewWord(word) {
		obj.total += 1
		if autoSave == false {
			obj.cf.InsertUnique([]byte(string(word)))
		}
		return true
	}

	return false
}

func (obj *WordsMapping) DeleteWord(word []rune) bool {
	obj.Lock()
	defer obj.Unlock()

	firstWordRune := word[0]

	if _, ok := obj.detail[firstWordRune]; !ok {
		return true
	}

	if obj.detail[firstWordRune].DeleteWord(word) {
		obj.cf.Delete([]byte(string(word)))
		obj.total--
		return true
	}

	return false
}

func (obj *WordsMapping) TotalWords() int {
	obj.Lock()
	defer obj.Unlock()

	return obj.total
}

func (obj *WordsMapping) FilterSentence(sentence []rune) []rune {
	obj.Lock()
	defer obj.Unlock()

	if obj.total == 0 {
		return sentence
	}

	sentenceLen := len(sentence)
	result := make([]rune, sentenceLen)
	copy(result, sentence)

	for i := 0; i < sentenceLen; i++ {
		startWord := sentence[i]

		if _, ok := obj.detail[startWord]; !ok {
			continue
		}

		maxSubLen := sentenceLen - i

		for maxSubLen > 0 {
			if _, ok := obj.detail[startWord].detail[maxSubLen]; !ok {
				maxSubLen--
				continue
			}

			found := false
			tempWords := string(sentence[i : i+maxSubLen])

			for _, val := range obj.detail[startWord].detail[maxSubLen] {
				if tempWords == val {
					found = true
					break
				}
			}

			if found {
				for j := i; j < i+maxSubLen; j++ {
					result[j] = rune('*')
				}
				i += maxSubLen
				maxSubLen = 0
			} else {
				maxSubLen--
			}
		}
	}

	return result
}
