package notificationinfra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type notificationRepository struct {
	dbmap *gorp.DbMap
}

func NewNotificationRepository(dbmap *gorp.DbMap) notificationdomain.NotificationRepository {
	return &notificationRepository{dbmap: dbmap}
}

func (r *notificationRepository) FindByID(ctx context.Context, id valueobject.NotificationID, userID valueobject.UserID) (*notificationdomain.Notification, error) {
	// TODO: implement
	return nil, nil
}

func (r *notificationRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*notificationdomain.Notification, error) {
	// TODO: implement
	return nil, nil
}

func (r *notificationRepository) Save(ctx context.Context, notification *notificationdomain.Notification) error {
	// TODO: implement
	return nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id valueobject.NotificationID) error {
	// TODO: implement
	return nil
}

func toNotificationEntity(row NotificationModel) (*notificationdomain.Notification, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.NotificationID](row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sqlconv.NewPrimaryIDFromBytes[valueobject.UserID](row.UserID)
	if err != nil {
		return nil, err
	}
	return notificationdomain.NewNotification(
		*id,
		*userID,
		row.NotificationType,
		row.Metadata,
		row.ReadAt,
		row.CreatedAt,
	), nil
}

// 未使用import対策（実装時に削除）
var _ = errors.Is
var _ = sql.ErrNoRows
var _ = dbtx.Querier
