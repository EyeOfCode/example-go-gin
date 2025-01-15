package service

import (
	"context"
	"example-go-project/internal/model"
	repository "example-go-project/internal/repository"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileService struct {
	fileStoreRepo repository.LocalFileRepository
}

func NewFileService(fileStoreRepo repository.LocalFileRepository) *FileService {
	return &FileService{
		fileStoreRepo: fileStoreRepo,
	}
}

func (f *FileService) UploadFile(ctx context.Context, files []*multipart.FileHeader, user *model.User) error {
	_, err := f.fileStoreRepo.Uploads(ctx, files, user)
	return err
}

func (f *FileService) DeleteFile(ctx context.Context, id string) error {
	return f.fileStoreRepo.Delete(ctx, id)
}

func (f *FileService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]model.FileStorage, error) {
	return f.fileStoreRepo.FindAll(ctx, query, opts)
}