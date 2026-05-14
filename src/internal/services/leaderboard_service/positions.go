package leaderboard_service

import (
	"context"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

func (s *Service) GetUserPositions(ctx context.Context, userID models.ID) ([]models.LeaderboardPosition, error) {
	var positions []models.LeaderboardPosition

	for _, key := range LeaderboardKeys {
		rank, err := s.leaderboardRepo.GetRank(ctx, key, userID)
		if err != nil {
			return nil, err
		}

		if rank == nil {
			continue
		}

		positions = append(positions, models.LeaderboardPosition{
			Language: key.Language,
			Mode:     key.Mode,
			SubMode:  key.SubMode,
			Rank:     *rank + 1, // 1-based
		})
	}

	if positions == nil {
		positions = []models.LeaderboardPosition{}
	}

	return positions, nil
}
