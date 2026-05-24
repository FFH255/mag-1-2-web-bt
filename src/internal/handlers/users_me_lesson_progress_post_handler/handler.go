package users_me_lesson_progress_post_handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/lesson_progress_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_lesson_progress_post_handler"

type progressSaver interface {
	Save(ctx context.Context, in *lesson_progress_service.SaveIn) error
}

type RequestBody struct {
	Keymap     string  `json:"keymap" example:"qwerty"`
	LessonID   uint32  `json:"lessonId" example:"0"`
	StepID     uint32  `json:"stepId" example:"0"`
	WPM        float64 `json:"wpm" example:"42.5"`
	CPM        float64 `json:"cpm" example:"210.3"`
	Accuracy   float64 `json:"accuracy" example:"98.7"`
	DurationMs uint64  `json:"durationMs" example:"60000"`
	CreatedAt  string  `json:"createdAt" example:"2025-10-19T19:02:29+03:00"`
} //@name UsersMeLessonProgressPostHandler.RequestBody

type Handler struct {
	progressSaver progressSaver
	logger        internal.Logger
}

func (h *Handler) newSaveIn(c *gin.Context) (*lesson_progress_service.SaveIn, error) {
	body := new(RequestBody)
	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		return nil, err
	}

	createdAt, err := proto.UnmarshalTime(body.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &lesson_progress_service.SaveIn{
		UserID:    models.ID(api.GetUserID(c)),
		Keymap:    body.Keymap,
		LessonID:  body.LessonID,
		StepID:    body.StepID,
		WPM:       body.WPM,
		CPM:       body.CPM,
		Accuracy:  body.Accuracy,
		Duration:  proto.ParseMilliseconds(&body.DurationMs),
		CreatedAt: createdAt,
	}, nil
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	var (
		status  = http.StatusInternalServerError
		message = "something went wrong serverside"
	)

	switch {
	case lesson_progress_service.IsUserNotFoundError(err):
		status = http.StatusNotFound
		message = "user not found"
	case models.IsValidationError(err):
		status = http.StatusBadRequest
		message = err.Error()
	}

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)

	switch status {
	case http.StatusInternalServerError:
		h.logger.Error(ctx)
	default:
		h.logger.Warning(ctx)
	}

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary Сохранить результат шага урока
// @Description Upsert результата шага урока для текущего пользователя
// @Tags Lessons
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body RequestBody true "Тело запроса"
// @Success 204
// @Failure 400 {object} proto.Error
// @Failure 404 {object} proto.Error
// @Failure 500 {object} proto.Error
// @Router /users/me/lessons/progress [post]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	in, err := h.newSaveIn(c)
	if err != nil {
		ctx = h.logger.WithStatusCode(ctx, http.StatusBadRequest)
		h.logger.Warning(h.logger.WithError(ctx, err))
		proto.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.progressSaver.Save(ctx, in); err != nil {
		h.handleError(ctx, c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Method() string {
	return http.MethodPost
}

func (h *Handler) Path() string {
	return "/users/me/lessons/progress"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(progressSaver progressSaver, logger internal.Logger) *Handler {
	return &Handler{
		progressSaver: progressSaver,
		logger:        logger,
	}
}
