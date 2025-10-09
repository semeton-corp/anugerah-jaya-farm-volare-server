package storage

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
)

const (
	ContentDisposition string = "inline"
)

type IStorage interface {
	UploadFile(fileData multipart.File, metadata map[string]string) (*s3.PutObjectOutput, error)
	DownloadFile(key string) (*s3.GetObjectOutput, error)
	DeleteFile(key string) (*s3.DeleteObjectOutput, error)
	GetPresignedUrl(key string) (string, error)
}

type Storage struct {
	client          *s3.Client
	presignedClient *s3.PresignClient
}

func New() IStorage {
	creds := credentials.NewStaticCredentialsProvider(viper.GetString("s3.access_key"), viper.GetString("s3.secret_key"), "")

	client := s3.New(s3.Options{
		Credentials:  creds,
		AppID:        viper.GetString("s3.app_id"),
		BaseEndpoint: aws.String(viper.GetString("s3.endpoint")),
		Region:       viper.GetString("s3.region"),
		UsePathStyle: true,
	}, func(o *s3.Options) {
		o.ClientLogMode = aws.LogSigning | aws.LogRequest | aws.LogResponseWithBody
	})

	presignedClient := s3.NewPresignClient(client)

	return &Storage{
		client:          client,
		presignedClient: presignedClient,
	}
}

func (s *Storage) UploadFile(fileData multipart.File, metadata map[string]string) (*s3.PutObjectOutput, error) {
	out, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:             aws.String(viper.GetString("s3.bucket")),
		Key:                aws.String(metadata["key"]),
		Body:               fileData,
		ACL:                types.ObjectCannedACLPublicRead,
		ContentDisposition: aws.String(ContentDisposition),
		ContentType:        aws.String(metadata["content-type"]),
		Metadata:           metadata,
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Storage) DownloadFile(key string) (*s3.GetObjectOutput, error) {
	out, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket:       aws.String(viper.GetString("s3.bucket")),
		Key:          aws.String(key),
		ChecksumMode: types.ChecksumModeEnabled,
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Storage) DeleteFile(key string) (*s3.DeleteObjectOutput, error) {
	out, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("s3.bucket")),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Storage) GetPresignedUrl(key string) (string, error) {
	out, err := s.presignedClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("s3.bucket")),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(10*time.Minute))
	if err != nil {
		return "", err
	}

	return out.URL, nil
}
