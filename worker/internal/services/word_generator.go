package services

import (
	"fmt"
	"math/big"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
)

type WordGenerator struct {
	alphabet      []rune
	alphabetLen   int
	maxWordLen    int
	startIndex    *big.Int
	endIndex      *big.Int
	currentIndex  *big.Int
	lastLen       int
	lastLenCount  *big.Int
	lastWordCount *big.Int
	curWordCount  *big.Int
}

// NewWordGenerator создает новый генератор слов для заданной задачи
func NewWordGenerator(task *models.CrackTaskRequest) (*WordGenerator, error) {
	if task.PartNumber < 0 || task.PartNumber >= task.PartCount {
		return nil, fmt.Errorf("invalid part number: %d (must be 0 <= partNumber < partCount)", task.PartNumber)
	}
	if task.PartCount <= 0 {
		return nil, fmt.Errorf("part count must be positive")
	}
	if len(task.Alphabet) == 0 {
		return nil, fmt.Errorf("alphabet cannot be empty")
	}
	if task.MaxLen <= 0 {
		return nil, fmt.Errorf("max word length must be positive")
	}

	alphabet := []rune(task.Alphabet)
	alphabetLen := len(alphabet)

	// Вычисляем общее количество слов
	totalWords := calculateTotalWords(alphabetLen, task.MaxLen)

	// Вычисляем диапазон для этого воркера
	partSize := new(big.Int).Div(totalWords, big.NewInt(int64(task.PartCount)))
	remainder := new(big.Int).Mod(totalWords, big.NewInt(int64(task.PartCount)))

	startIndex := new(big.Int).Mul(partSize, big.NewInt(int64(task.PartNumber)))
	if int64(task.PartNumber) < remainder.Int64() {
		startIndex.Add(startIndex, big.NewInt(int64(task.PartNumber)))
	} else {
		startIndex.Add(startIndex, remainder)
	}

	endIndex := new(big.Int).Mul(partSize, big.NewInt(int64(task.PartNumber+1)))
	if int64(task.PartNumber+1) < remainder.Int64() {
		endIndex.Add(endIndex, big.NewInt(int64(task.PartNumber+1)))
	} else {
		endIndex.Add(endIndex, remainder)
	}

	// Последний воркер получает остаток
	if task.PartNumber == task.PartCount-1 {
		endIndex.Set(totalWords)
	}

	return &WordGenerator{
		alphabet:      alphabet,
		alphabetLen:   alphabetLen,
		maxWordLen:    task.MaxLen,
		startIndex:    startIndex,
		endIndex:      endIndex,
		currentIndex:  new(big.Int).Set(startIndex),
		lastLen:       1,
		lastLenCount:  big.NewInt(int64(alphabetLen)),
		lastWordCount: big.NewInt(0),
		curWordCount:  big.NewInt(int64(alphabetLen)),
	}, nil
}

// calculateTotalWords вычисляет общее количество слов
// Формула: sum_{i=1}^{maxLen} alphabetLen^i = (alphabetLen^(maxLen+1) - alphabetLen) / (alphabetLen - 1)
func calculateTotalWords(alphabetLen, maxWordLen int) *big.Int {
	if alphabetLen == 1 {
		return big.NewInt(int64(maxWordLen))
	}

	base := big.NewInt(int64(alphabetLen))
	exp := big.NewInt(int64(maxWordLen + 1))
	numerator := new(big.Int).Exp(base, exp, nil)
	numerator.Sub(numerator, base)

	denominator := big.NewInt(int64(alphabetLen - 1))
	result := new(big.Int).Div(numerator, denominator)

	return result
}

func (wg *WordGenerator) indexToWord(globalIdx *big.Int) string {
	N := big.NewInt(int64(wg.alphabetLen))

	for globalIdx.Cmp(wg.curWordCount) >= 0 { // = потому что нумерация с 0
		wg.lastWordCount.Set(wg.curWordCount)
		wg.curWordCount.Add(wg.curWordCount, wg.lastLenCount.Mul(wg.lastLenCount, N))
		wg.lastLen += 1
	}

	// 2. Позиция внутри блока длины `length`
	pos := new(big.Int).Sub(globalIdx, wg.lastWordCount)

	// 3. Конвертируем pos в систему счисления с основанием N
	chars := make([]rune, wg.lastLen)
	for j := wg.lastLen - 1; j >= 0; j-- {
		chars[j] = wg.alphabet[new(big.Int).Mod(pos, N).Int64()]
		pos.Div(pos, N)
	}

	return string(chars)
}

// Next генерирует следующее слово в диапазоне
// Возвращает слово и флаг, есть ли еще слова
func (wg *WordGenerator) Next() (string, bool) {
	if wg.currentIndex.Cmp(wg.endIndex) >= 0 {
		return "", false
	}

	word := wg.indexToWord(wg.currentIndex)
	wg.currentIndex.Add(wg.currentIndex, big.NewInt(1))

	return word, true
}

// HasNext проверяет, есть ли еще слова для генерации
func (wg *WordGenerator) HasNext() bool {
	return wg.currentIndex.Cmp(wg.endIndex) < 0
}
