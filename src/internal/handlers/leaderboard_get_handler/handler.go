package leaderboard_get_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/leaderboard_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "/leaderboard"

type leaderboardGetter interface {
	GetLeaderboard(ctx context.Context, language models.Language, mode models.Mode, submode models.Submode, pageIndex, pageSize int64) (*models.LeaderboardPage, error)
}

type Handler struct {
	log               internal.Logger
	leaderboardGetter leaderboardGetter
}

// Handle godoc
// @Summary Получить лидерборд
// @Description Получить лидерборд по языку, режиму и подрежиму
// @Tags Leaderboards
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param language query string true "Язык" Example(english)
// @Param mode query string true "Режим" Example(time)
// @Param submode query string true "Подрежим" Example(15s)
// @Param pageIndex query int true "Индекс запрашиваемой страницы (начиная с нуля)" Example(0)
// @Param pageSize query int true "Размер запрашиваемой страницы" Example(50)
// @Success 200 {object} ResponseBody
// @Failure 400 {object} proto.Error "Неверные параметры запроса"
// @Failure 401 {object} proto.Error "Пользователь не авторизован"
// @Router /leaderboard [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.log.WithHandlerName(c.Request.Context(), handlerName)

	request, err := newRequest(c)
	if err != nil {
		proto.WriteError(c, http.StatusBadRequest, err)
		h.log.Warning(h.log.WithError(h.log.WithStatusCode(ctx, http.StatusBadRequest), err))
		return
	}

	page, err := h.leaderboardGetter.GetLeaderboard(
		ctx,
		models.Language(request.Language),
		models.Mode(request.Mode),
		models.Submode(request.SubMode),
		request.PageIndex,
		request.PageSize,
	)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	responseBody := newResponseBody(page)
	h.log.Info(ctx, "get leaderboard success")
	proto.WriteJSON(c, http.StatusOK, responseBody)
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "something went wrong"

	switch {
	case leaderboard_service.IsInvalidLeaderboardIDError(err):
		status = http.StatusBadRequest
		message = "invalid leaderboard id"
	}

	ctx = h.log.WithError(h.log.WithStatusCode(ctx, status), err)
	if status == http.StatusInternalServerError {
		h.log.Error(ctx)
	} else {
		h.log.Warning(ctx)
	}

	proto.WriteError(c, status, message)
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

func New(leaderboardGetter leaderboardGetter, log internal.Logger) *Handler {
	return &Handler{
		log:               log,
		leaderboardGetter: leaderboardGetter,
	}
}
