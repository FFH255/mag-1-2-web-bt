package leaderboard_service

import (
	"context"
	"fmt"

	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/repositories/statistics_repository"
)

const warmUpPageSize = 100

func (s *Service) WarmUp(ctx context.Context) error {
	offset := 0
	for {
		records, err := s.statsRepo.GetPersonalBestsPaginated(ctx, warmUpPageSize, offset)
		if err != nil {
			return fmt.Errorf("leaderboard warmup: failed to get personal bests (offset=%d): %w", offset, err)
		}
		if len(records) == 0 {
			break
		}

		if err := s.leaderboardRepo.UpsertScores(ctx, newScoreRecordBatch(records)); err != nil {
			return fmt.Errorf("leaderboard warmup: failed to upsert batch at offset=%d: %w", offset, err)
		}

		if len(records) < warmUpPageSize {
			break
		}
		offset += warmUpPageSize
	}

	return nil
}

func newScoreRecordBatch(records []statistics_repository.PersonalBestRecord) []models.ScoreRecord {
	batch := make([]models.ScoreRecord, 0, len(records))
	for _, rec := range records {
		leaderboardID := models.LeaderboardID{Language: rec.Language, Mode: rec.Mode, SubMode: rec.SubMode}
		if !isAllowedLeaderboardID(leaderboardID) {
			continue
		}

		batch = append(batch, models.ScoreRecord{
			ID:     leaderboardID,
			UserID: rec.UserID,
			Score:  models.Score{WPM: rec.WPM, Accuracy: rec.Accuracy, PlayedAt: rec.PlayedAt},
		})
	}
	return batch
}
