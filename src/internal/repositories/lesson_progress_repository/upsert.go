package lesson_progress_repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (r *Repository) Upsert(ctx context.Context, result *models.LessonStepResult) error {
	const query = `
		INSERT INTO lesson_step_results (
			user_id, keymap, lesson_id, step_id,
			wpm, cpm, accuracy, duration_ms, created_at
		) VALUES (
			@user_id, @keymap, @lesson_id, @step_id,
			@wpm, @cpm, @accuracy, @duration_ms, @created_at
		)
		ON CONFLICT ON CONSTRAINT uq_lesson_step_result DO UPDATE SET
			wpm        = EXCLUDED.wpm,
			cpm        = EXCLUDED.cpm,
			accuracy   = EXCLUDED.accuracy,
			duration_ms = EXCLUDED.duration_ms,
			created_at = EXCLUDED.created_at
		RETURNING id
	`

	args := pgx.NamedArgs{
		"user_id":     result.UserID,
		"keymap":      result.Keymap,
		"lesson_id":   result.LessonID,
		"step_id":     result.StepID,
		"wpm":         result.WPM,
		"cpm":         result.CPM,
		"accuracy":    result.Accuracy,
		"duration_ms": result.Duration.Milliseconds(),
		"created_at":  result.CreatedAt,
	}

	return r.db.QueryRow(ctx, query, args).Scan(&result.ID)
}
