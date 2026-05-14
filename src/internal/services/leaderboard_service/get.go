package leaderboard_service

import (
	"context"
	"math"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (s *Service) GetLeaderboard(ctx context.Context, language models.Language, mode models.Mode, submode models.Submode, pageIndex, pageSize int64) (*models.LeaderboardPage, error) {
	id := models.LeaderboardID{Language: language, Mode: mode, SubMode: submode}

	if !isAllowedLeaderboardID(id) {
		return nil, InvalidLeaderboardIDError{ID: id}
	}

	start := pageIndex * pageSize
	stop := start + pageSize - 1

	entries, err := s.leaderboardRepo.GetPage(ctx, id, start, stop)
	if err != nil {
		return nil, err
	}

	total, err := s.leaderboardRepo.Count(ctx, id)
	if err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(pageSize)))
	}

	userIDs := make([]models.ID, len(entries))
	for i, e := range entries {
		userIDs[i] = e.UserID
	}

	usernames, err := s.userRepo.GetUsernamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		entries[i].Rank = start + int64(i) + 1
		entries[i].Username = usernames[entries[i].UserID]
	}

	return &models.LeaderboardPage{
		Entries:    entries,
		PageIndex:  pageIndex,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
