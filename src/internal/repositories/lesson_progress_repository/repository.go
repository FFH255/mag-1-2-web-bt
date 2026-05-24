package lesson_progress_repository

import "github.com/ruslanonly/blindtyping/src/internal/shared/postgres"

type Repository struct {
	db *postgres.Database
}

func New(db *postgres.Database) *Repository {
	return &Repository{
		db: db,
	}
}
