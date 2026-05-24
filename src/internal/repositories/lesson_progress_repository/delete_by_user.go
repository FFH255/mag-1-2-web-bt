package lesson_progress_repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID models.ID) error {
	const query = `
		DELETE FROM lesson_step_results
		WHERE user_id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	_, err := r.db.Exec(ctx, query, args)
	return err
}
