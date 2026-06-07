package userusecase

import (
	"context"
	"fmt"
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/google/uuid"
)

const (
	presignedURLTTL = 5 * time.Minute
)

type GenerateProfileImagePresignedURLInput struct {
	UserID      valueobject.UserID
	ContentType valueobject.ImageContentType
}

type GenerateProfileImagePresignedURLOutput struct {
	URL  valueobject.URL
	Path string
}

type GenerateProfileImagePresignedURL struct {
	storageClient shareddomain.StorageClient
}

func (uc *GenerateProfileImagePresignedURL) Execute(ctx context.Context, input GenerateProfileImagePresignedURLInput) (*GenerateProfileImagePresignedURLOutput, error) {
	path := fmt.Sprintf("profiles/%s/%s.%s", input.UserID.Value(), uuid.New().String(), input.ContentType.Extension())
	putURL, err := uc.storageClient.GeneratePutURL(ctx, path, input.ContentType.Value(), presignedURLTTL)
	if err != nil {
		return nil, err
	}
	return &GenerateProfileImagePresignedURLOutput{
		URL:  putURL,
		Path: path,
	}, nil
}

func NewGenerateProfileImagePresignedURL(storageClient shareddomain.StorageClient) *GenerateProfileImagePresignedURL {
	return &GenerateProfileImagePresignedURL{storageClient: storageClient}
}
