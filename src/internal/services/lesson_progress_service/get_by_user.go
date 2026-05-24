package lesson_progress_service

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type StepResult struct {
	WPM       float64 `json:"wpm"`
	CPM       float64 `json:"cpm"`
	Accuracy  float64 `json:"accuracy"`
	Time      int64   `json:"time"`
	CreatedAt string  `json:"createdAt"`
}

// GetOut mirrors the frontend TLessonKeymapsStats structure:
// keymap -> lessonId -> stepId -> StepResult
type GetOut map[string]map[uint32]map[uint32]StepResult

func (s *Service) GetByUser(ctx context.Context, userID models.ID) (GetOut, error) {
	user, err := s.userRepository.GetOne(ctx, &user_repository.GetOneIn{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, UserNotFoundError{id: userID}
	}

	results, err := s.lessonProgressRepository.GetAllByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	out := make(GetOut)

	for _, r := range results {
		keymap := string(r.Keymap)

		if out[keymap] == nil {
			out[keymap] = make(map[uint32]map[uint32]StepResult)
		}
		if out[keymap][r.LessonID] == nil {
			out[keymap][r.LessonID] = make(map[uint32]StepResult)
		}

		out[keymap][r.LessonID][r.StepID] = StepResult{
			WPM:       float64(r.WPM),
			CPM:       float64(r.CPM),
			Accuracy:  float64(r.Accuracy),
			Time:      r.Duration.Milliseconds(),
			CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return out, nil
}
