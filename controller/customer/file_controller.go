package customer

import (
	"os"
	"path/filepath"

	"api.mijkomp.com/exception"
	"api.mijkomp.com/models/response"
	"api.mijkomp.com/service"
	"github.com/gofiber/fiber/v2"
)

type FileController struct {
	UserService service.UserService
}

func NewFileController(userService *service.UserService) *FileController {
	return &FileController{UserService: *userService}
}

func (controller *FileController) Route(app *fiber.App) {
	file := app.Group("/api/files")
	file.Get("/", controller.GetPhoto)
}

// GetPhoto godoc
// @Summary      Get Photo by ID
// @Tags         File
// @Produce      image/jpeg
// @Param        id query string true "File UUID without extension"
// @Success      200 {string} binary "Photo content"
// @Failure      400 {object} response.WebResponse
// @Failure      404 {object} response.WebResponse
// @Router       /api/files [get]
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
