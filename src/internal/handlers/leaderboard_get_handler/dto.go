package leaderboard_get_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

type Request struct {
	Language  string
	Mode      string
	SubMode   string
	PageIndex int64
	PageSize  int64
}

type LeaderboardItem struct {
	Rank     int64   `json:"rank" description:"Позиция в лидерборде" example:"1"`
	Username string  `json:"username" description:"Имя пользователя" example:"ffh"`
	WPM      float64 `json:"wpm" description:"Кол-во слов в минуту" example:"142.5"`
	Accuracy float64 `json:"accuracy" description:"Точность" example:"99.2"`
	PlayedAt string  `json:"playedAt" description:"Дата, когда тест был сыгран (RFC-3339)" example:"2020-01-01T00:00:00Z"`
} //@name LeaderboardGetHandler.LeaderboardItem

type ResponseBody struct {
	Leaderboard []LeaderboardItem      `json:"leaderboard" description:"Лидерборд"`
	Pagination  api.PaginationResponse `json:"pagination" description:"Информация о пагинации"`
} //@name LeaderboardGetHandler.ResponseBody

func newRequest(c *gin.Context) (*Request, error) {
	pageIndex, pageSize, err := api.GetPaginationQueryParams(c)
	if err != nil {
		return nil, err
	}

	language := c.Query("language")
	if language == "" {
		return nil, models.NewValidationError("language", "language is required")
	}

	mode := c.Query("mode")
	if mode == "" {
		return nil, models.NewValidationError("mode", "mode is required")
	}

	submode := c.Query("submode")
	if submode == "" {
		return nil, models.NewValidationError("submode", "submode is required")
	}

	return &Request{
		Language:  language,
		Mode:      mode,
		SubMode:   submode,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, nil
}

func newResponseBody(page *models.LeaderboardPage) ResponseBody {
	items := make([]LeaderboardItem, 0, len(page.Entries))
	for _, entry := range page.Entries {
		items = append(items, LeaderboardItem{
			Rank:     entry.Rank,
			Username: entry.Username,
			WPM:      entry.WPM,
			Accuracy: entry.Accuracy,
			PlayedAt: proto.MarshalTime(entry.PlayedAt),
		})
	}

	return ResponseBody{
		Leaderboard: items,
		Pagination: api.PaginationResponse{
			PageIndex:  page.PageIndex,
			PageSize:   page.PageSize,
			TotalPages: page.TotalPages,
		},
	}
}
