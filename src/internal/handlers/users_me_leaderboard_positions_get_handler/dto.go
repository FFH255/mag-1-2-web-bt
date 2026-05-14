package users_me_leaderboard_positions_get_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/models"
)

type LeaderboardPosition struct {
	Language string `json:"language" example:"english"`
	Mode     string `json:"mode" example:"time"`
	SubMode  string `json:"submode" example:"15s"`
	Rank     int64  `json:"rank" example:"42"`
} //@name UsersMeLeaderboardPositionsGetHandler.LeaderboardPosition

type ResponseBody struct {
	LeaderboardPositions []LeaderboardPosition `json:"leaderboardPositions"`
} //@name UsersMeLeaderboardPositionsGetHandler.ResponseBody

type Request struct {
	UserID models.ID
}

func newResponseBody(positions []models.LeaderboardPosition) *ResponseBody {
	items := make([]LeaderboardPosition, 0, len(positions))
	for _, p := range positions {
		items = append(items, LeaderboardPosition{
			Language: string(p.Language),
			Mode:     string(p.Mode),
			SubMode:  string(p.SubMode),
			Rank:     p.Rank,
		})
	}

	return &ResponseBody{
		LeaderboardPositions: items,
	}
}

func newRequest(c *gin.Context) *Request {
	return &Request{
		UserID: models.ID(api.GetUserID(c)),
	}
}
