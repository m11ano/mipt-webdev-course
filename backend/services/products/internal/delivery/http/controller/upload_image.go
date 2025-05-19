package controller

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/m11ano/e"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/delivery/http/middleware"
	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/domain"
)

type UploadImageOut struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

var imageTypes = map[string]domain.FileTarget{
	"preview": domain.FileTargetProductPreview,
	"slider":  domain.FileTargetProductSlider,
}

// @Summary Загрузка изображения
// @Security BearerAuth
// @Tags products
// @Accept  multipart/form-data
// @Param file formData file true "Файл изображения"
// @Param image_type formData string true "Тип изображения, enum: preview, slider"
// @Success 200 {object} UploadImageOut
// @Failure 400 {object} middleware.ErrorJSON
// @Router /products/image [post]
func (ctrl *Controller) UploadImageHandler(c *fiber.Ctx) error {

	authData := middleware.ExtractAuthData(c)

	if !authData.IsAuth {
		return e.ErrUnauthorized
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("form field `file` is required")
	}

	imageType := c.FormValue("image_type")
	imageTypeValue, ok := imageTypes[imageType]
	if !ok {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("form field `imageType` is invalid")
	}

	f, err := fileHeader.Open()
	if err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot open uploaded file")
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return e.NewErrorFrom(e.ErrBadRequest).Wrap(err).SetMessage("cannot read uploaded file")
	}

	file, err := ctrl.fileUC.UploadImageFile(c.Context(), domain.FileTarget(imageTypeValue), fileHeader.Filename, fileData)
	if err != nil {
		return err
	}

	result := UploadImageOut{
		ID:  file.ID,
		URL: file.GetURL(&ctrl.cfg),
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
