package service

import (
	"encoding/json"
	"fmt"
	"os"

	"task/internal/models"
	"task/internal/repository"
)

type Redirecter interface {
	SaveDataset(fileName string) error
}

type RedirectService struct {
	repo repository.Redirecter
}

func NewRedirecter(repo repository.Redirecter) Redirecter {
	return &RedirectService{
		repo: repo,
	}
}

func (s *RedirectService) SaveDataset(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("ReadFile(): %s", err)
	}

	var links []models.Link
	if err := json.Unmarshal(data, &links); err != nil {
		return fmt.Errorf("Unmarshal(): %s", err)
	}

	values := linksToQueryStr(links)
	if err := s.repo.SaveDataset(values); err != nil {
		return fmt.Errorf("error saving dataset to the DB: %s", err)
	}
	return nil
}

func linksToQueryStr(links []models.Link) string {
	var res string
	for _, pair := range links {
		res += fmt.Sprintf("('%s','%s'),", pair.ActiveLink, pair.HistoryLink)
	}

	return res[:len(res)-1] + ";"
}
