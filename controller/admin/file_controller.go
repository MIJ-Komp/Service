package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/middleware"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FileController struct {
	UserService service.UserService
}

func NewFileController(userService *service.UserService) *FileController {
	return &FileController{UserService: *userService}
}

func (controller *FileController) Route(app *fiber.App) {
	file := app.Group("/api/admin/files", middleware.AuthMiddleware(controller.UserService))
	file.Post("/", controller.UploadPhoto)
	file.Get("/", controller.GetPhoto)
}

// UploadPhotos godoc
// @Summary      Upload Multiple Photos
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        photos formData file true "Multiple photo files (max 1MB each, .jpg/.jpeg/.png)"
// @Success      200 {object} response.WebResponse
// @Failure      400 {object} response.WebResponse
// @Failure      413 {object} response.WebResponse "File too large"
// @Security     ApiKeyAuth
// @Router       /api/admin/files [post]
func (controller *FileController) UploadPhoto(ctx *fiber.Ctx) error {
	const maxFileSize = 1 * 1024 * 1024 // 1MB
	uploadDir := "./uploads"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		panic(exception.NewValidationError("Invalid multipart form"))
	}

	files := form.File["photos"]
	if len(files) == 0 {
		panic(exception.NewValidationError("No files uploaded"))
	}

	var fileIDs []string

	for _, file := range files {
		if file.Size > maxFileSize {
			panic(exception.NewValidationError(fmt.Sprintf("File %s exceeds 1MB limit", file.Filename)))
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			panic(exception.NewValidationError(fmt.Sprintf("File %s has invalid extension", file.Filename)))
		}

		fileId := uuid.New()
		newFileName := fileId.String() + ext
		savePath := filepath.Join(uploadDir, newFileName)

		if err := ctx.SaveFile(file, savePath); err != nil {
			panic(exception.NewValidationError(fmt.Sprintf("Failed to save file %s", file.Filename)))
		}

		fileIDs = append(fileIDs, fileId.String())
	}

	return ctx.JSON(response.NewWebResponse(fileIDs, "Files uploaded successfully"))
}

// GetPhoto godoc
// @Summary      Get Photo by ID
// @Tags         File
// @Produce      image/jpeg
// @Param        id query string true "File UUID without extension"
// @Success      200 {string} binary "Photo content"
// @Failure      400 {object} response.WebResponse
// @Failure      404 {object} response.WebResponse
// @Security     ApiKeyAuth
// @Router       /api/admin/files [get]
func (controller *FileController) GetPhoto(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		panic(exception.NewValidationError("File ID is required"))
	}

	uploadDir := "./uploads"
	var filePath string
	var contentType string

	for _, ext := range []string{".jpg", ".jpeg", ".png"} {
		tempPath := filepath.Join(uploadDir, id+ext)
		if _, err := os.Stat(tempPath); err == nil {
			filePath = tempPath
			switch ext {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			}
			break
		}
	}

	if filePath == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(response.NewWebResponse(nil, "File not found"))
	}

	ctx.Set("Content-Type", contentType)
	return ctx.SendFile(filePath)
}
