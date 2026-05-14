package leaderboard_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/ruslanonly/blindtyping/src/internal/models"
)

type Repository struct {
	client *redis.Client
}

func leaderboardKey(id models.LeaderboardID) string {
	return fmt.Sprintf("leaderboard:%s:%s:%s", id.Language, id.Mode, id.SubMode)
}

func indexKey(id models.LeaderboardID) string {
	return fmt.Sprintf("leaderboard_index:%s:%s:%s", id.Language, id.Mode, id.SubMode)
}

func (r *Repository) UpsertScore(ctx context.Context, id models.LeaderboardID, userID models.ID, score models.Score) error {
	key := leaderboardKey(id)
	idxKey := indexKey(id)
	userIDStr := userID.String()

	composite := score.Composite()

	data, err := json.Marshal(score)
	if err != nil {
		return err
	}

	if err := r.client.ZAdd(ctx, key, redis.Z{
		Score:  composite,
		Member: userIDStr,
	}).Err(); err != nil {
		return err
	}

	return r.client.HSet(ctx, idxKey, userIDStr, data).Err()
}

func (r *Repository) GetRank(ctx context.Context, id models.LeaderboardID, userID models.ID) (*int64, error) {
	key := leaderboardKey(id)
	rank, err := r.client.ZRevRank(ctx, key, userID.String()).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

func (r *Repository) GetPage(ctx context.Context, id models.LeaderboardID, start, stop int64) ([]models.LeaderboardEntry, error) {
	key := leaderboardKey(id)
	idxKey := indexKey(id)

	userIDStrs, err := r.client.ZRevRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	if len(userIDStrs) == 0 {
		return []models.LeaderboardEntry{}, nil
	}

	fields := make([]string, len(userIDStrs))
	copy(fields, userIDStrs)
	vals, err := r.client.HMGet(ctx, idxKey, fields...).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]models.LeaderboardEntry, 0, len(userIDStrs))
	for i, userIDStr := range userIDStrs {
		uid, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			continue
		}

		var score models.Score
		if vals[i] != nil {
			raw, ok := vals[i].(string)
			if !ok {
				continue
			}
			if err := json.Unmarshal([]byte(raw), &score); err != nil {
				continue
			}
		}

		entries = append(entries, models.LeaderboardEntry{
			UserID:   models.ID(uid),
			WPM:      score.WPM,
			Accuracy: score.Accuracy,
			PlayedAt: score.PlayedAt,
		})
	}

	return entries, nil
}

func (r *Repository) UpsertScores(ctx context.Context, records []models.ScoreRecord) error {
	pipe := r.client.Pipeline()

	for _, rec := range records {
		key := leaderboardKey(rec.ID)
		idxKey := indexKey(rec.ID)
		userIDStr := rec.UserID.String()
		composite := rec.Score.Composite()

		data, err := json.Marshal(rec.Score)
		if err != nil {
			return err
		}

		pipe.ZAdd(ctx, key, redis.Z{
			Score:  composite,
			Member: userIDStr,
		})
		pipe.HSet(ctx, idxKey, userIDStr, data)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) Count(ctx context.Context, id models.LeaderboardID) (int64, error) {
	key := leaderboardKey(id)
	return r.client.ZCard(ctx, key).Result()
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}
