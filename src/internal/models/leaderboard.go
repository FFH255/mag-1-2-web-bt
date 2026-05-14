package models

import "time"

type Score struct {
	WPM      float64   `json:"wpm"`
	Accuracy float64   `json:"accuracy"`
	PlayedAt time.Time `json:"played_at"`
}

// MaxTimestamp is a far-future Unix timestamp (year 2100) used to invert time ordering in composite scores.
const MaxTimestamp int64 = 4102444800

// Composite encodes WPM, Accuracy, and PlayedAt into a single float64 for Redis sorted set ranking.
// Primary: WPM descending, Secondary: Accuracy descending, Tertiary: PlayedAt ascending (earlier is better).
func (s Score) Composite() float64 {
	return s.WPM*1_000_000 + s.Accuracy*1_000 + float64(MaxTimestamp-s.PlayedAt.Unix())*0.001
}

type LeaderboardEntry struct {
	Rank     int64
	UserID   ID
	Username string
	WPM      float64
	Accuracy float64
	PlayedAt time.Time
}

type LeaderboardPosition struct {
	Language Language
	Mode     Mode
	SubMode  Submode
	Rank     int64
}

type LeaderboardPage struct {
	Entries    []LeaderboardEntry
	PageIndex  int64
	PageSize   int64
	TotalPages int64
}

type RankChange struct {
	OldRank *int64
	NewRank int64
}

type LeaderboardID struct {
	Language Language
	Mode     Mode
	SubMode  Submode
}

type ScoreRecord struct {
	ID     LeaderboardID
	UserID ID
	Score  Score
}
