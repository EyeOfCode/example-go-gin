package handlers

import (
	"context"
	fileRepository "example-go-project/internal/repository/files"
	userRepository "example-go-project/internal/repository/user"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadHandler struct {
	fileRepo fileRepository.LocalFileRepository
	userRepo userRepository.UserRepository
}

func NewUploadHandler(fileRepo fileRepository.LocalFileRepository, userRepo userRepository.UserRepository) *UploadHandler {
	return &UploadHandler{
		fileRepo: fileRepo,
		userRepo: userRepo,
	}
}

// @Summary     Upload multiple files
// @Description Upload multiple files to the server
// @Tags        uploads
// @Accept      multipart/form-data
// @Produce     json
// @Security    Bearer
// @Param       files formData []file true "Multiple files to upload"
// @Router      /local_upload [post]
func (u *UploadHandler) UploadMultipleLocalFiles(c *gin.Context) {
		userID, _ := c.Get("userID")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userIDStr, ok := userID.(string)
		if !ok {
			utils.SendError(c, http.StatusInternalServerError, "Failed to get user ID")
			return
		}

		user, err := u.userRepo.FindByID(ctx, userIDStr)
		if err != nil && err != mongo.ErrNoDocuments {
			utils.SendError(c, http.StatusInternalServerError, "User not found")
			return
		}

    form, err := c.MultipartForm()
    if err != nil {
				utils.SendError(c, http.StatusBadRequest, "Failed to parse form")
        return
    }

    files := form.File["files"]
    if len(files) == 0 {
				utils.SendError(c, http.StatusBadRequest, "No files received")
        return
    }

    filesInfo, err := u.fileRepo.Uploads(ctx, files, user)
    if err != nil {
				utils.SendError(c, http.StatusInternalServerError, err.Error())
				return
    }

    res := gin.H{"files": filesInfo}

    utils.SendSuccess(c, http.StatusOK, res, "Files uploaded successfully")
}

// @Summary     Delete a file
// @Description Delete a file from the server
// @Tags        uploads
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Param       id path string true "File ID"
// @Router      /local_upload/{id} [delete]
func (u *UploadHandler) DeleteFile(c *gin.Context) {
		id := c.Param("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := u.fileRepo.Delete(ctx, id)
		if err != nil {
				utils.SendError(c, http.StatusInternalServerError, err.Error())
				return
		}

		utils.SendSuccess(c, http.StatusOK, nil, "File deleted successfully")
}

// @Summary     Get all files
// @Description Get all files from the server
// @Tags        uploads
// @Accept      json
// @Produce     json
// @Security    Bearer
// @Router      /local_upload [get]
func (u *UploadHandler) GetFileAll(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		files, err := u.fileRepo.FindAll(ctx, bson.D{}, nil)
		if err != nil {
				utils.SendError(c, http.StatusInternalServerError, err.Error())
				return
		}

		utils.SendSuccess(c, http.StatusOK, files)
}