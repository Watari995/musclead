package mealusecase

import (
	"context"
	"fmt"
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/google/uuid"
)

const (
	presignedURLTTL = 3 * time.Minute
)

type GenerateMealPhotoImagePresignedURLInput struct {
	UserID      valueobject.UserID
	ContentType valueobject.ImageContentType
}

type GenerateMealPhotoImagePresignedURLOutput struct {
	URL  valueobject.URL
	Path string
}

type GenerateMealPhotoImagePresignedURL struct {
	storageClient shareddomain.StorageClient
}

func (uc *GenerateMealPhotoImagePresignedURL) Execute(ctx context.Context, input GenerateMealPhotoImagePresignedURLInput) (*GenerateMealPhotoImagePresignedURLOutput, error) {
	path := fmt.Sprintf("meals/%s/%s.%s", input.UserID.Value(), uuid.New().String(), input.ContentType.Extension())
	putURL, err := uc.storageClient.GeneratePutURL(ctx, path, input.ContentType.Value(), presignedURLTTL)
	if err != nil {
		return nil, err
	}
	return &GenerateMealPhotoImagePresignedURLOutput{
		URL:  putURL,
		Path: path,
	}, nil
}

func NewGenerateMealPhotoImagePresignedURL(storageClient shareddomain.StorageClient) *GenerateMealPhotoImagePresignedURL {
	return &GenerateMealPhotoImagePresignedURL{storageClient: storageClient}
}
