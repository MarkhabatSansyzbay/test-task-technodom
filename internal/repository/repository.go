package repository

import (
	"database/sql"
	"fmt"
	"task/internal/models"

	"github.com/jmoiron/sqlx"
)

type Redirecter interface {
	SaveDataset(values string) error
	AllRedirects() (*sql.Rows, error)
	RedirectByID(id int) (models.Link, error)
	CreateRedirect(redirect models.Link) error
	UpdateRedirect(id int, newActiveLink string) error
}

type RedirectRepository struct {
	db *sqlx.DB
}

func NewRedirecter(db *sqlx.DB) Redirecter {
	return &RedirectRepository{
		db: db,
	}
}

func (r *RedirectRepository) UpdateRedirect(id int, newActiveLink string) error {
	query := `
		UPDATE links SET active_link = $2, history_link = active_link
		WHERE ID = $1;
	`

	if _, err := r.db.Exec(query, id, newActiveLink); err != nil {
		return err
	}

	return nil
}

func (r *RedirectRepository) CreateRedirect(redirect models.Link) error {
	query := `
		INSERT INTO links (active_link, history_link)
		VALUES ($1, $2);
	`

	if _, err := r.db.Exec(query, redirect.ActiveLink, redirect.HistoryLink); err != nil {
		return err
	}

	return nil
}

func (r *RedirectRepository) RedirectByID(id int) (models.Link, error) {
	query := `
		SELECT ID, active_link, history_link
		FROM links WHERE ID=$1;
	`

	var redirect models.Link
	if err := r.db.QueryRow(query, id).Scan(&redirect.ID, &redirect.ActiveLink, &redirect.HistoryLink); err != nil {
		return models.Link{}, err
	}

	return redirect, nil
}

func (r *RedirectRepository) AllRedirects() (*sql.Rows, error) {
	query := `
		SELECT ID, active_link, history_link
		FROM links;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *RedirectRepository) SaveDataset(values string) error {
	query := fmt.Sprintf(`
		INSERT INTO links(active_link, history_link)
		VALUES %s
	`, values)

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}
