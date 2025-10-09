package service

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/storage"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type FileService struct {
	log     *zap.Logger
	storage storage.IStorage
}

type IFileService interface {
	UploadFile(file *multipart.FileHeader) (dto.FileResponse, error)
	DownloadFile(key string) (dto.FileDetailResponse, error)
	DeleteFile(key string) error
	GetPresignedUrl(key string) (dto.FileResponse, error)
}

func NewFileService(
	log *zap.Logger,
	storage storage.IStorage,
) IFileService {
	return &FileService{
		log:     log,
		storage: storage,
	}
}

func (s *FileService) UploadFile(file *multipart.FileHeader) (dto.FileResponse, error) {
	fileData, err := file.Open()
	if err != nil {
		s.log.Error("failed open file data", zap.Error(err))
		return dto.FileResponse{}, err
	}

	key := fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().Unix())
	fileContentType := file.Header.Get("Content-Type")

	metadata := make(map[string]string)
	metadata["key"] = key
	metadata["real-name"] = file.Filename
	metadata["content-type"] = fileContentType

	_, err = s.storage.UploadFile(fileData, metadata)
	if err != nil {
		s.log.Error("failed upload file", zap.Error(err))
		return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		URL: util.ConstructURL(key),
	}, nil
}

func (s *FileService) DownloadFile(key string) (dto.FileDetailResponse, error) {
	out, err := s.storage.DownloadFile(key)
	if err != nil {
		s.log.Error("failed to dowload file", zap.Error(err))
	}

	return dto.FileDetailResponse{
		Metadata: out.Metadata,
		Body:     out.Body,
	}, nil
}

func (s *FileService) DeleteFile(key string) error {
	_, err := s.storage.DeleteFile(key)
	if err != nil {
		s.log.Error("failed delete file", zap.Error(err))
		return err
	}

	return err
}

func (s *FileService) GetPresignedUrl(key string) (dto.FileResponse, error) {
	out, err := s.storage.GetPresignedUrl(key)
	if err != nil {
		s.log.Error("failed get presigned url", zap.Error(err))
		return dto.FileResponse{}, err
	}

	return dto.FileResponse{
		URL: out,
	}, nil
}
