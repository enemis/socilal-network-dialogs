package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"social-network-dialogs/internal/dialog"
	"social-network-dialogs/internal/logger"
	"social-network-dialogs/internal/validator"
)

type DialogHandler struct {
	dialogService *dialog.DialogService
	logger        logger.LoggerInterface
}

type MessageInput struct {
	Recipient uuid.UUID `json:"recipient" binding:"required"`
	Message   string    `json:"message" binding:"required"`
}

func NewDialogHandler(service *dialog.DialogService, logger logger.LoggerInterface) *DialogHandler {
	return &DialogHandler{
		dialogService: service,
		logger:        logger,
	}
}

func (h *DialogHandler) SendMessage(c *gin.Context) {
	userUuid, ok := h.fetchUser(c)
	if ok != true {
		return
	}

	var input MessageInput

	validator := validator.NewValidator(input)

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, validator.DecryptErrors(err).(map[string]interface{}))
		return
	}

	_, err := h.dialogService.SendDirectMessage(*userUuid, input.Recipient, input.Message)

	if err != nil {
		h.logger.Error("error send message", err, nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"message": "internal server error"})
	}

	c.Status(http.StatusOK)
}

func (h *DialogHandler) fetchUser(c *gin.Context) (*uuid.UUID, bool) {
	user, exist := c.Get("user")

	if exist != true {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"message": "unauthorized"})
		return nil, false
	}
	userUuid, ok := user.(uuid.UUID)
	if ok != true {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"message": "unauthorized"})
		return nil, false
	}
	return &userUuid, true
}

func (h *DialogHandler) GetDialog(c *gin.Context) {
	userUuid, ok := h.fetchUser(c)
	if ok != true {
		return
	}

	collocutorId := c.Param("user_id")
	collocutorUuid, err := uuid.Parse(collocutorId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"user_id": "must be valid user id"})
		return
	}

	messages, err := h.dialogService.ListLastDirectMessages(*userUuid, collocutorUuid, 0, 50)

	if err != nil {
		h.logger.Error("error fetch dialog", err, nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"messages": messages})
}
