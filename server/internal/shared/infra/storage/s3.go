package sharedstorage

import (
	"context"
	"time"

	shareddomain "github.com/Watari995/musclead/internal/shared/domain"
	"github.com/Watari995/musclead/internal/valueobject"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3Client(client *s3.Client, bucket string) shareddomain.StorageClient {
	return &s3Client{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
	}
}

func (c *s3Client) GeneratePutURL(ctx context.Context,
	key, contentType string, ttl time.Duration) (valueobject.URL, error) {
	req, err := c.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return valueobject.URL{}, err
	}
	url, err := c.buildURL(req.URL)
	if err != nil {
		return valueobject.URL{}, err
	}
	return *url, nil
}

func (c *s3Client) GenerateGetURL(ctx context.Context, key string, ttl time.Duration) (valueobject.URL, error) {
	req, err := c.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return valueobject.URL{}, err
	}
	url, err := c.buildURL(req.URL)
	if err != nil {
		return valueobject.URL{}, err
	}
	return *url, nil
}

func (c *s3Client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

// buildURLはURLをvalueobject.URLに変換する
func (c *s3Client) buildURL(urlValue string) (*valueobject.URL, error) {
	url, err := valueobject.NewURL(urlValue)
	if err != nil {
		return nil, err
	}
	return url, nil
}
