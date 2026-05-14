package leaderboard_service

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type leaderboardRepository interface {
	UpsertScore(ctx context.Context, id models.LeaderboardID, userID models.ID, score models.Score) error
	UpsertScores(ctx context.Context, records []models.ScoreRecord) error
	GetRank(ctx context.Context, id models.LeaderboardID, userID models.ID) (*int64, error)
	GetPage(ctx context.Context, id models.LeaderboardID, start, stop int64) ([]models.LeaderboardEntry, error)
	Count(ctx context.Context, id models.LeaderboardID) (int64, error)
}

type statisticsRepository interface {
	GetPersonalBestsPaginated(ctx context.Context, limit, offset int) ([]statistics_repository.PersonalBestRecord, error)
}

type userRepository interface {
	GetUsernamesByIDs(ctx context.Context, ids []models.ID) (map[models.ID]string, error)
}

// LeaderboardKeys defines all active leaderboard combinations.
var LeaderboardKeys = []models.LeaderboardID{
	{Language: "russian", Mode: models.ModeTime, SubMode: models.SubmodeTime15s},
	{Language: "russian", Mode: models.ModeTime, SubMode: models.SubmodeTime1m},
	{Language: "english", Mode: models.ModeTime, SubMode: models.SubmodeTime15s},
	{Language: "english", Mode: models.ModeTime, SubMode: models.SubmodeTime1m},
}

type InvalidLeaderboardIDError struct {
	ID models.LeaderboardID
}

func (e InvalidLeaderboardIDError) Error() string {
	return "invalid leaderboard id"
}

func IsInvalidLeaderboardIDError(err error) bool {
	return models.IsCustomError[InvalidLeaderboardIDError](err)
}

func isAllowedLeaderboardID(id models.LeaderboardID) bool {
	for _, key := range LeaderboardKeys {
		if key == id {
			return true
		}
	}
	return false
}

type Service struct {
	leaderboardRepo leaderboardRepository
	statsRepo       statisticsRepository
	userRepo        userRepository
}

func New(
	leaderboardRepo leaderboardRepository,
	statsRepo statisticsRepository,
	userRepo userRepository,
) *Service {
	return &Service{
		leaderboardRepo: leaderboardRepo,
		statsRepo:       statsRepo,
		userRepo:        userRepo,
	}
}
