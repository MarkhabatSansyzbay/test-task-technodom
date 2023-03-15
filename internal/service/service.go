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
	Redirects() ([]models.Link, error)
	RedirectByID(id int) (models.Link, error)
	CreateRedirect(redirect models.Link) error
	UpdateRedirect(id int, newActiveLink string) error
	DeleteRedirect(id int) error
}

type RedirectService struct {
	repo repository.Redirecter
}

func NewRedirecter(repo repository.Redirecter) Redirecter {
	return &RedirectService{
		repo: repo,
	}
}

func (s *RedirectService) DeleteRedirect(id int) error {
	if err := s.repo.DeleteRedirect(id); err != nil {
		return fmt.Errorf("repo.DeleteRedirect(): %s", err)
	}

	return nil
}

func (s *RedirectService) UpdateRedirect(id int, newActiveLink string) error {
	if err := s.repo.UpdateRedirect(id, newActiveLink); err != nil {
		return fmt.Errorf("repo.UpdateRedirect(): %s", err)
	}

	return nil
}

func (s *RedirectService) CreateRedirect(redirect models.Link) error {
	if err := s.repo.CreateRedirect(redirect); err != nil {
		return fmt.Errorf("repo.CreateRedirect(): %s", err)
	}

	return nil
}

func (s *RedirectService) RedirectByID(id int) (models.Link, error) {
	var res models.Link
	res, err := s.repo.RedirectByID(id)
	if err != nil {
		return models.Link{}, fmt.Errorf("repo.RedirectByID(): %s", err)
	}

	return res, nil
}

func (s *RedirectService) Redirects() ([]models.Link, error) {
	rows, err := s.repo.AllRedirects()
	if err != nil {
		return nil, fmt.Errorf("repo.AllRedirects(): %s", err)
	}
	defer rows.Close()

	var redirects []models.Link
	for rows.Next() {
		var redirect models.Link

		if err := rows.Scan(&redirect.ID, &redirect.ActiveLink, &redirect.HistoryLink); err != nil {
			return nil, fmt.Errorf("rows.Scan(): %s", err)
		}

		redirects = append(redirects, redirect)
	}

	return redirects, nil
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
