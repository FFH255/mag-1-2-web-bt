package lesson_progress_service

import (
	"context"
	"fmt"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/user_repository"
)

type UserNotFoundError struct {
	id models.ID
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user with id %d not found", e.id)
}

func IsUserNotFoundError(err error) bool {
	return models.IsCustomError[UserNotFoundError](err)
}

type lessonProgressRepository interface {
	Upsert(ctx context.Context, result *models.LessonStepResult) error
	GetAllByUserID(ctx context.Context, userID models.ID) ([]*models.LessonStepResult, error)
}

type userRepository interface {
	GetOne(ctx context.Context, in *user_repository.GetOneIn) (*models.User, error)
}

type keymapRepository interface {
	Get() models.Keymaps
}

type Service struct {
	lessonProgressRepository lessonProgressRepository
	userRepository           userRepository
	keymapRepository         keymapRepository
}

func New(
	lessonProgressRepository lessonProgressRepository,
	userRepository userRepository,
	keymapRepository keymapRepository,
) *Service {
	return &Service{
		lessonProgressRepository: lessonProgressRepository,
		userRepository:           userRepository,
		keymapRepository:         keymapRepository,
	}
}
