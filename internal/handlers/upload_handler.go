package handlers

import (
	"context"
	"example-go-project/internal/service"
	"example-go-project/pkg/middleware"
	"example-go-project/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadHandler struct {
	fileService *service.FileService
	userService *service.UserService
}

func NewUploadHandler(fileService *service.FileService, userService *service.UserService) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
		userService: userService,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, ok := middleware.GetUserFromContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User not found")
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

	filesInfo, err := u.fileService.UploadFile(ctx, files, user)
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

	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	resFile, err := u.fileService.FindById(ctx, ObjID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = u.fileService.DeleteFile(ctx, resFile.ID.Hex())
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

	files, err := u.fileService.FindAll(ctx, bson.D{}, nil)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccess(c, http.StatusOK, files)
}
