package notificationinfra

import (
	"context"

	notificationdomain "github.com/Watari995/musclead/internal/notification/internal/domain"
	"github.com/Watari995/musclead/internal/shared/dbtx"
	"github.com/Watari995/musclead/internal/shared/sqlconv"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/go-gorp/gorp/v3"
)

type deviceTokenRepository struct {
	dbmap *gorp.DbMap
}

func NewDeviceTokenRepository(dbmap *gorp.DbMap) notificationdomain.DeviceTokenRepository {
	return &deviceTokenRepository{dbmap: dbmap}
}

func NewDeviceTokenQuery(dbmap *gorp.DbMap) notificationdomain.DeviceTokenQuery {
	return &deviceTokenRepository{dbmap: dbmap}
}

const upsertDeviceTokenSQL = `
INSERT INTO device_tokens (id, user_id, token, platform, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    user_id = VALUES(user_id),
    platform = VALUES(platform),
    updated_at = VALUES(updated_at)
`

func (r *deviceTokenRepository) Save(ctx context.Context, deviceToken *notificationdomain.DeviceToken) error {
	q := dbtx.Querier(ctx, r.dbmap)
	params, err := buildUpsertDeviceTokenParams(deviceToken)
	if err != nil {
		return err
	}
	if _, err := q.Exec(upsertDeviceTokenSQL, params...); err != nil {
		return err
	}
	return nil
}

func (r *deviceTokenRepository) FindAllByUserID(ctx context.Context, userID valueobject.UserID) ([]notificationdomain.DeviceTokenView, error) {
	q := dbtx.Querier(ctx, r.dbmap)
	userIDBytes, err := userID.Bytes()
	if err != nil {
		return nil, err
	}
	var rows []DeviceTokenModel
	if _, err := q.Select(&rows, `
		SELECT id, user_id, token, platform, created_at, updated_at
		FROM device_tokens
		WHERE user_id = ?
	`, userIDBytes); err != nil {
		return nil, err
	}
	views := make([]notificationdomain.DeviceTokenView, 0, len(rows))
	for _, row := range rows {
		view, err := toDeviceTokenView(row)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func buildUpsertDeviceTokenParams(deviceToken *notificationdomain.DeviceToken) ([]any, error) {
	idBytes, err := deviceToken.ID().Bytes()
	if err != nil {
		return nil, err
	}
	userIDBytes, err := deviceToken.UserID().Bytes()
	if err != nil {
		return nil, err
	}
	return []any{
		idBytes,
		userIDBytes,
		deviceToken.Token(),
		deviceToken.Platform().Value(),
		deviceToken.CreatedAt(),
		deviceToken.UpdatedAt(),
	}, nil
}

func toDeviceTokenView(row DeviceTokenModel) (notificationdomain.DeviceTokenView, error) {
	id, err := sqlconv.NewPrimaryIDFromBytes[valueobject.DeviceTokenID](row.ID)
	if err != nil {
		return notificationdomain.DeviceTokenView{}, err
	}
	platform, err := valueobject.NewNotificationPlatformFromString(row.Platform)
	if err != nil {
		return notificationdomain.DeviceTokenView{}, err
	}
	return notificationdomain.DeviceTokenView{
		ID:       *id,
		Token:    row.Token,
		Platform: *platform,
	}, nil
}
