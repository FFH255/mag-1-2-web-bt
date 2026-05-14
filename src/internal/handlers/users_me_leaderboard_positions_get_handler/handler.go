package users_me_leaderboard_positions_get_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "/users/me/leaderboard-positions"

type positionsGetter interface {
	GetUserPositions(ctx context.Context, userID models.ID) ([]models.LeaderboardPosition, error)
}

type Handler struct {
	log             internal.Logger
	positionsGetter positionsGetter
}

// Handle godoc
// @Summary Получить местоположение текущего пользователя в лидербордах
// @Description Получить местоположение текущего пользователя в лидербордах
// @Tags Leaderboards
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} ResponseBody
// @Failure 401 {object} proto.Error "Пользователь не авторизован"
// @Router /users/me/leaderboard-positions [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.log.WithHandlerName(c.Request.Context(), handlerName)

	request := newRequest(c)

	positions, err := h.positionsGetter.GetUserPositions(ctx, request.UserID)
	if err != nil {
		proto.WriteError(c, http.StatusInternalServerError, "something went wrong")
		h.log.Error(h.log.WithError(h.log.WithStatusCode(ctx, http.StatusInternalServerError), err))
		return
	}

	res := newResponseBody(positions)
	proto.WriteJSON(c, http.StatusOK, res)
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return handlerName
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(positionsGetter positionsGetter, logger internal.Logger) *Handler {
	return &Handler{
		log:             logger,
		positionsGetter: positionsGetter,
	}
}
