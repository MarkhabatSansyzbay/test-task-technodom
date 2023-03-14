package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Redirecter interface {
	SaveDataset(values string) error
}

type RedirectRepository struct {
	db *sqlx.DB
}

func NewRedirecter(db *sqlx.DB) Redirecter {
	return &RedirectRepository{
		db: db,
	}
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
