package services

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
	"github.com/TKaterinna/CrackHash/worker/internal/repo"
	"github.com/google/uuid"
)

type CalcService struct {
	repo         *repo.CalcRepo
	resultSender *ResultSender
}

func NewCalcService(repo *repo.CalcRepo, resultSender *ResultSender) *CalcService {
	return &CalcService{
		repo:         repo,
		resultSender: resultSender,
	}
}

func (s *CalcService) Save(req *models.CrackTaskRequest) error {
	// if err = s.repo.SaveTask(req.RequestId, req.PartNumber, req.PartCount, req.MaxLen, req.CheckHash); err != nil {
	// 	return err
	// }

	print("Save")

	go s.work(req)

	return nil
}

func (s *CalcService) GetResult(requestId uuid.UUID) (string, error) {
	var data string
	var err error

	if data, err = s.repo.GetResult(requestId); err != nil {
		return "", err
	}

	return data, nil
}

func (s *CalcService) getMD5Hash(word string) string {
	hash := md5.Sum([]byte(word))
	return hex.EncodeToString(hash[:])
}

func (s *CalcService) checkWord(word string, checkHash string) bool {
	curHash := s.getMD5Hash(word)

	if curHash == checkHash {
		return true
	}

	return false
}

func (s *CalcService) work(req *models.CrackTaskRequest) {
	var err error
	var words []string
	var wg *WordGenerator

	if wg, err = NewWordGenerator(req); err != nil {
		res := &models.CrackTaskResult{
			RequestId: req.RequestId,
			Results:   nil,
			Status:    models.StatusERROR,
		}
		s.resultSender.Send(res)
		return
	}

	log.Println("START work")

	for {
		var word string
		var isNotEnd bool
		if word, isNotEnd = wg.Next(); isNotEnd == false {
			break
		}

		log.Println("Check word = ", word)

		if s.checkWord(word, req.CheckHash) {
			log.Println("CORRECT WORD: ", word)
			words = append(words, word)
			break
		}
	}

	res := &models.CrackTaskResult{
		RequestId: req.RequestId,
		Results:   words,
		Status:    models.StatusDONE,
	}
	s.resultSender.Send(res)
}
