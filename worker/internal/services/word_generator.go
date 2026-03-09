package services

import (
	"fmt"
	"math"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
)

type WordGenerator struct {
	alphabet      []rune
	alphabetLen   int64
	startIndex    int64
	endIndex      int64
	currentIndex  int64
	lastLen       int64
	lastLenCount  int64
	lastWordCount int64
	curWordCount  int64
}

// NewWordGenerator создает новый генератор слов для заданной задачи
func NewWordGenerator(task *models.CrackTaskRequest) (*WordGenerator, error) {
	if task.Count <= 0 {
		return nil, fmt.Errorf("count must be positive")
	}
	if task.StartIndex < 0 {
		return nil, fmt.Errorf("start index must be non-negative")
	}
	if len(task.Alphabet) == 0 {
		return nil, fmt.Errorf("alphabet cannot be empty")
	}
	if task.MaxLen <= 0 {
		return nil, fmt.Errorf("max word length must be positive")
	}

	alphabet := []rune(task.Alphabet)
	alphabetLen := int64(len(alphabet))

	startIndex := task.StartIndex
	endIndex := startIndex + task.Count

	wordsCount := WordsCount(alphabetLen, task.MaxLen)
	if endIndex > wordsCount {
		endIndex = wordsCount
	}

	return &WordGenerator{
		alphabet:      alphabet,
		alphabetLen:   alphabetLen,
		startIndex:    startIndex,
		endIndex:      endIndex,
		currentIndex:  startIndex,
		lastLen:       1,
		lastLenCount:  alphabetLen,
		lastWordCount: 0,
		curWordCount:  alphabetLen,
	}, nil
}

func WordsCount(alphabetLen int64, maxLen int64) int64 {
	return int64(float64(alphabetLen) * (math.Pow(float64(alphabetLen), float64(maxLen)) - 1) / float64(alphabetLen-1))
}

func (wg *WordGenerator) indexToWord(globalIdx int64) string {
	N := int64(wg.alphabetLen)

	for globalIdx >= wg.curWordCount { // = потому что нумерация с 0
		wg.lastWordCount = wg.curWordCount
		wg.curWordCount = wg.curWordCount + wg.lastLenCount*N
		wg.lastLen += 1
	}

	// 2. Позиция внутри блока длины `length`
	pos := globalIdx - wg.lastWordCount

	// 3. Конвертируем pos в систему счисления с основанием N
	chars := make([]rune, wg.lastLen)
	for j := wg.lastLen - 1; j >= 0; j-- {
		chars[j] = wg.alphabet[pos%N]
		pos = pos / N
	}

	return string(chars)
}

// Next генерирует следующее слово в диапазоне
// Возвращает слово и флаг, есть ли еще слова
func (wg *WordGenerator) Next() (string, bool) {
	if wg.currentIndex >= wg.endIndex {
		return "", false
	}

	word := wg.indexToWord(wg.currentIndex)
	wg.currentIndex = wg.currentIndex + 1

	return word, true
}

// HasNext проверяет, есть ли еще слова для генерации
func (wg *WordGenerator) HasNext() bool {
	return wg.currentIndex < wg.endIndex
}
