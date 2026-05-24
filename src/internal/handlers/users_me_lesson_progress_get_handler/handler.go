package users_me_lesson_progress_get_handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruslanonly/blindtyping/src/internal"
	"github.com/ruslanonly/blindtyping/src/internal/api"
	"github.com/ruslanonly/blindtyping/src/internal/api/middleware"
	"github.com/ruslanonly/blindtyping/src/internal/models"
	"github.com/ruslanonly/blindtyping/src/internal/services/lesson_progress_service"
	"github.com/ruslanonly/blindtyping/src/internal/shared/proto"
)

const handlerName = "users_me_lesson_progress_get_handler"

type progressGetter interface {
	GetByUser(ctx context.Context, userID models.ID) (lesson_progress_service.GetOut, error)
}

type StepResult struct {
	Accuracy  float64 `json:"accuracy" example:"0.95"`
	CPM       float64 `json:"cpm" example:"250.5"`
	CreatedAt string  `json:"createdAt" example:"2024-01-15T10:30:00Z"`
	Time      int64   `json:"time" example:"45000"`
	WPM       float64 `json:"wpm" example:"50.0"`
} //@name UsersMeLessonProgressGetHandler.StepResult

// LessonStepProgressMap — ключи: stepId
type LessonStepProgressMap map[string]StepResult //@name UsersMeLessonProgressGetHandler.LessonStepProgressMap

// LessonProgressMap — ключи: lessonId
type LessonProgressMap map[string]LessonStepProgressMap //@name UsersMeLessonProgressGetHandler.LessonProgressMap

// ResponseBody mirrors TLessonKeymapsStats from the frontend.
// Ключи верхнего уровня — название раскладки (keymap), например "qwerty", "dvorak".
type ResponseBody map[string]LessonProgressMap //@name UsersMeLessonProgressGetHandler.ResponseBody

func newResponseBody(out lesson_progress_service.GetOut) ResponseBody {
	body := make(ResponseBody, len(out))
	for keymap, lessons := range out {
		body[keymap] = make(LessonProgressMap, len(lessons))
		for lessonID, steps := range lessons {
			lessonKey := fmt.Sprint(lessonID)
			body[keymap][lessonKey] = make(LessonStepProgressMap, len(steps))
			for stepID, r := range steps {
				body[keymap][lessonKey][fmt.Sprint(stepID)] = StepResult{
					Accuracy:  r.Accuracy,
					CPM:       r.CPM,
					CreatedAt: r.CreatedAt,
					Time:      r.Time,
					WPM:       r.WPM,
				}
			}
		}
	}
	return body
}

type Handler struct {
	progressGetter progressGetter
	logger         internal.Logger
}

func (h *Handler) handleError(ctx context.Context, c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "something went wrong serverside"

	ctx = h.logger.WithError(h.logger.WithStatusCode(ctx, status), err)
	h.logger.Error(ctx)

	proto.WriteError(c, status, message)
}

// Handle godoc
// @Summary Получить прогресс уроков
// @Description Получить весь прогресс уроков текущего пользователя
// @Tags Lessons
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} ResponseBody
// @Failure 500 {object} proto.Error
// @Router /users/me/lessons/progress [get]
func (h *Handler) Handle(c *gin.Context) {
	ctx := h.logger.WithHandlerName(c.Request.Context(), handlerName)

	userID := models.ID(api.GetUserID(c))

	result, err := h.progressGetter.GetByUser(ctx, userID)
	if err != nil {
		h.handleError(ctx, c, err)
		return
	}

	proto.WriteJSON(c, http.StatusOK, newResponseBody(result))
}

func (h *Handler) Method() string {
	return http.MethodGet
}

func (h *Handler) Path() string {
	return "/users/me/lessons/progress"
}

func (h *Handler) Middleware() []string {
	return []string{middleware.Auth}
}

func New(progressGetter progressGetter, logger internal.Logger) *Handler {
	return &Handler{
		progressGetter: progressGetter,
		logger:         logger,
	}
}
