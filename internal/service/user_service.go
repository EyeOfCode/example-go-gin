package service

import (
	"context"
	"example-go-project/internal/dto"
	"example-go-project/internal/model"
	"example-go-project/internal/repository"
	"example-go-project/pkg/config"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
	config      *config.Config
}

func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client, config *config.Config) *UserService {
	return &UserService{
		userRepo:    userRepo,
		redisClient: redisClient,
		config:      config,
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

func (u *UserService) Create(ctx context.Context, payload *dto.RegisterRequest) (*model.User, error) {
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
	req := bson.M{
		"name": payload.Name,
	}
	user, err := u.userRepo.Update(ctx, req, id)
	if err != nil {
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

func (u *UserService) Login(ctx context.Context, password string, user *model.User) (*utils.TokenPair, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	auth := utils.NewAuthHandler(u.config.JWTSecretKey, u.config.JWTRefreshKey, u.config.JWTExpiresIn, u.config.JWTRefreshIn)
	tokenPair, err := auth.GenerateTokenPair(user.ID.Hex(), user.Roles)
	if err != nil {
		return nil, err
	}

	expires, _ := time.ParseDuration(u.config.JWTExpiresIn)
	if err := u.redisClient.Set(ctx,
		tokenPair.AccessToken,
		user.ID.Hex(),
		expires).Err(); err != nil {
		return nil, err
	}
	return tokenPair, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*utils.TokenPair, error) {
	auth := utils.NewAuthHandler(s.config.JWTSecretKey, s.config.JWTRefreshKey, s.config.JWTExpiresIn, s.config.JWTRefreshIn)
	claims, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	blacklisted, err := s.redisClient.Get(ctx, "blacklist:"+refreshToken).Result()
	if blacklisted != "" || err != redis.Nil {
		return nil, err
	}

	user, err := s.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	tokenPair, err := auth.GenerateTokenPair(user.ID.Hex(), user.Roles)
	if err != nil {
		return nil, err
	}

	expires, _ := time.ParseDuration(s.config.JWTExpiresIn)
	if err := s.redisClient.Set(ctx,
		tokenPair.AccessToken,
		user.ID.Hex(),
		expires).Err(); err != nil {
		return nil, err
	}

	// Keep blacklist for 48h to prevent reuse
	if err := s.redisClient.Set(ctx,
		"blacklist:"+refreshToken,
		"true",
		48*time.Hour).Err(); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *UserService) ValidateTokenWithRedis(ctx context.Context, token string) error {
	// Check blacklist with 24h window
	blacklisted, err := s.redisClient.Get(ctx, "blacklist:"+token).Result()
	if err != redis.Nil || blacklisted != "" {
		return gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{
				"status": http.StatusUnauthorized,
			},
		}
	}

	// Check active tokens
	_, err = s.redisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
			Meta: gin.H{
				"status": http.StatusUnauthorized,
			},
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	pipe := s.redisClient.Pipeline()

	// Blacklist access token for 24h
	expires, _ := time.ParseDuration(s.config.JWTExpiresIn)
	pipe.Set(ctx,
		"blacklist:"+accessToken,
		"true",
		expires)

	// Blacklist refresh token for 48h
	pipe.Set(ctx,
		"blacklist:"+refreshToken,
		"true",
		48*time.Hour)

	// Remove active access token
	pipe.Del(ctx, accessToken)

	_, err := pipe.Exec(ctx)
	return err
}
