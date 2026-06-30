package notificationinfra

import (
	"context"
	"database/sql"
	"errors"
	"time"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type notificationRepository struct {
	dbmap *gorp.DbMap
}

func NewNotificationRepository(dbmap *gorp.DbMap) notificationdomain.NotificationRepository {
	return &notificationRepository{dbmap: dbmap}
}

func (r *notificationRepository) FindByIDAndUserID(ctx context.Context, id valueobject.NotificationID, userID valueobject.UserID) (*notificationdomain.Notification, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var row NotificationModel
	if err := q.SelectOne(&row, `
		SELECT id, user_id, notification_type, metadata, read_at, created_at
		FROM notifications
		WHERE id = ? AND user_id = ?
	`, idBytes, userIDBytes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var result *notificationdomain.Notification
	result, err = toNotificationEntity(row)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *notificationRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]*notificationdomain.Notification, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []NotificationModel
	if _, err := q.Select(&rows, `
		SELECT id, user_id, notification_type, metadata, read_at, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userIDBytes); err != nil {
		return nil, err
	}
	result := make([]*notificationdomain.Notification, 0, len(rows))
	for _, row := range rows {
		notification, err := toNotificationEntity(row)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}
	return result, nil
}

func (r *notificationRepository) Save(ctx context.Context, notification *notificationdomain.Notification) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildInsertNotificationParams(notification)
	if err != nil {
		return err
	}
	if _, err := q.Exec(`
		INSERT INTO notifications (id, user_id, notification_type, metadata, read_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, params...); err != nil {
		return err
	}
	return nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id valueobject.NotificationID, userID valueobject.UserID) error {
	q := dbtx.Querier(ctx, r.dbmap)
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	if _, err := q.Exec("UPDATE notifications SET read_at = ? WHERE id = ?", time.Now(), idBytes); err != nil {
		return err
	}
	return nil
}

func buildInsertNotificationParams(notification *notificationdomain.Notification) ([]any, error) {
	idBytes, err := notification.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := notification.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	return []any{
		idBytes,
		userIDBytes,
		notification.NotificationType().String(),
		notification.Metadata(),
		notification.ReadAt(),
		notification.CreatedAt(),
	}, nil
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
	notificationType, err := valueobject.NewNotificationTypeFromString(row.NotificationType)
	if err != nil {
		return nil, err
	}
	return notificationdomain.NewNotification(
		*id,
		*userID,
		*notificationType,
		row.Metadata,
		row.ReadAt,
		row.CreatedAt,
	), nil
}
