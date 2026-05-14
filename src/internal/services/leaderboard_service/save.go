package leaderboard_service

import (
	"context"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (s *Service) getRank(ctx context.Context, id models.LeaderboardID, userID models.ID) (*int64, error) {
	zeroRank, err := s.leaderboardRepo.GetRank(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if zeroRank == nil {
		return nil, nil
	}

	rank := *zeroRank + 1

	return &rank, nil
}

func (s *Service) UpdateRank(ctx context.Context, userID models.ID, language models.Language, mode models.Mode, submode models.Submode, wpm float64, accuracy float64, playedAt time.Time) (*models.RankChange, error) {
	id := models.LeaderboardID{Language: language, Mode: mode, SubMode: submode}

	if !isAllowedLeaderboardID(id) {
		return nil, nil
	}

	oldRank, err := s.getRank(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if err := s.leaderboardRepo.UpsertScore(ctx, id, userID, models.Score{
		WPM:      wpm,
		Accuracy: accuracy,
		PlayedAt: playedAt,
	}); err != nil {
		return nil, err
	}

	newRank, err := s.getRank(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if oldRank == nil && newRank == nil {
		return nil, nil
	}

	return &models.RankChange{
		OldRank: oldRank,
		NewRank: pointer.Get(newRank),
	}, nil
}
