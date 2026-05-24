package lesson_progress_service

import (
	"context"
	"time"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type SaveIn struct {
	UserID    models.ID
	Keymap    string
	LessonID  uint32
	StepID    uint32
	WPM       float64
	CPM       float64
	Accuracy  float64
	Duration  time.Duration
	CreatedAt time.Time
}

func (s *Service) Save(ctx context.Context, in *SaveIn) error {
	user, err := s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		UserID: &in.UserID,
	})
	if err != nil {
		return err
	}
	if user == nil {
		return UserNotFoundError{id: in.UserID}
	}

	keymaps := s.keymapRepository.Get()

	result, err := models.NewLessonStepResult(models.LessonStepResultOptions{
		UserID:    user.ID,
		Keymap:    in.Keymap,
		LessonID:  in.LessonID,
		StepID:    in.StepID,
		WPM:       in.WPM,
		CPM:       in.CPM,
		Accuracy:  in.Accuracy,
		Duration:  in.Duration,
		CreatedAt: in.CreatedAt,
	}, keymaps)
	if err != nil {
		return err
	}

	return s.lessonProgressRepository.Upsert(ctx, result)
}
