package service

import (
	"context"
	"example-go-project/internal/model"
	repository "example-go-project/internal/repository"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (f *FileService) UploadFile(ctx context.Context, files []*multipart.FileHeader, user *model.User) ([]*model.FileStorage,error) {
	res, err := f.fileStoreRepo.Uploads(ctx, files, user)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (f *FileService) DeleteFile(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	return f.fileStoreRepo.Delete(ctx, objectID)
}

func (f *FileService) FindAll(ctx context.Context, query bson.D, opts *options.FindOptions) ([]*model.FileStorage, error) {
	return f.fileStoreRepo.FindAll(ctx, query, opts)
}

func (f *FileService) FindById(ctx context.Context, id primitive.ObjectID) (*model.FileStorage, error) {
	return f.fileStoreRepo.FindOne(ctx, bson.M{"_id": id})
}