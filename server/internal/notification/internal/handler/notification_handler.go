package notificationhandler

import (
	"net/http"

	"github.com/Watari995/musclead/internal/myerror"
	"github.com/Watari995/musclead/internal/shared/httpx"
	notificationdto "github.com/Watari995/musclead/internal/notification/dto"
	notificationusecase "github.com/Watari995/musclead/internal/notification/internal/usecase"
	"github.com/Watari995/musclead/internal/valueobject"
)

type NotificationHandler struct {
	getNotifications *notificationusecase.GetNotifications
	getNotification  *notificationusecase.GetNotification
	readNotification *notificationusecase.ReadNotification
}

func New(
	getNotifications *notificationusecase.GetNotifications,
	getNotification *notificationusecase.GetNotification,
	readNotification *notificationusecase.ReadNotification,
) http.Handler {
	h := &NotificationHandler{
		getNotifications: getNotifications,
		getNotification:  getNotification,
		readNotification: readNotification,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notifications", h.GetNotifications)
	mux.HandleFunc("GET /notifications/{id}", h.GetNotification)
	mux.HandleFunc("PUT /notifications/{id}/read", h.ReadNotification)
	return mux
}

// GetNotifications godoc
//
// @Summary 通知一覧取得
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} notificationdto.GetNotificationsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	output, err := h.getNotifications.Execute(r.Context(), notificationusecase.GetNotificationsInput{UserID: userID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	dtos := make([]notificationdto.NotificationDTO, 0, len(output.Notifications))
	for _, n := range output.Notifications {
		dtos = append(dtos, notificationdto.NotificationFromEntity(n))
	}
	httpx.WriteJSON(w, http.StatusOK, notificationdto.GetNotificationsResponse{
		Notifications: dtos,
		UnreadCount:   output.UnreadCount,
	})
}

// GetNotification godoc
//
// @Summary 通知詳細取得
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "notification id"
// @Success 200 {object} notificationdto.NotificationDTO
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /notifications/{id} [get]
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	id, err := valueobject.NewPrimaryIDFromString[valueobject.NotificationID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid notification id"))
		return
	}
	output, err := h.getNotification.Execute(r.Context(), notificationusecase.GetNotificationInput{ID: *id, UserID: userID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, notificationdto.NotificationFromEntity(output.Notification))
}

// ReadNotification godoc
//
// @Summary 通知既読化
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "notification id"
// @Success 204
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) ReadNotification(w http.ResponseWriter, r *http.Request) {
	userID, err := httpx.UserIDFromContext(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	id, err := valueobject.NewPrimaryIDFromString[valueobject.NotificationID](r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, myerror.NewBadRequestError().SetMessage("invalid notification id"))
		return
	}
	if _, err := h.readNotification.Execute(r.Context(), notificationusecase.ReadNotificationInput{ID: *id, UserID: userID}); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
