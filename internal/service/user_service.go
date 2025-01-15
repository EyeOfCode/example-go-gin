package service

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	"example-go-project/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := u.userRepo.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) FindByID(ctx context.Context, userID string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.FindOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) Create(ctx context.Context, payload *dto.RegisterRequest) (*model.User,error) {
	now := time.Now()
	user := &model.User{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  payload.Password,
		Roles:     payload.Roles,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	res := &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return res, nil
}

func (u *UserService) Update(ctx context.Context, payload *dto.UpdateProfileRequest, id primitive.ObjectID) (*model.User, error) {
	now := time.Now()
	user := &model.User{
		ID: id,
		Name:      payload.Name,
		UpdatedAt: now,
	}
	
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	res := &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return res, nil
}

func (u *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *UserService) FindAll(ctx context.Context, filter dto.UserFilter, page, pageSize int) ([]model.User, int64, error) {
	mongoFilter := bson.D{}
	if filter.Name != "" {
		mongoFilter = append(mongoFilter, bson.E{
			Key: "name",
			Value: bson.D{{
				Key:   "$regex",
				Value: primitive.Regex{Pattern: filter.Name, Options: "i"},
			}},
		})
	}

	total, err := u.userRepo.Count(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	users, err := u.userRepo.FindAll(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}