package lesson_progress_repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (r *Repository) GetAllByUserID(ctx context.Context, userID models.ID) ([]*models.LessonStepResult, error) {
	const query = `
		SELECT id, user_id, keymap, lesson_id, step_id,
		       wpm, cpm, accuracy, duration_ms, created_at
		FROM lesson_step_results
		WHERE user_id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.LessonStepResult

	for rows.Next() {
		var (
			res        models.LessonStepResult
			durationMs int64
		)

		err := rows.Scan(
			&res.ID,
			&res.UserID,
			&res.Keymap,
			&res.LessonID,
			&res.StepID,
			&res.WPM,
			&res.CPM,
			&res.Accuracy,
			&durationMs,
			&res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		res.Duration = time.Duration(durationMs) * time.Millisecond
		results = append(results, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
