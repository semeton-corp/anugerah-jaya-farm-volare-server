package rest

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/middleware"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type FileHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IFileService
}

func (h *FileHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/files")
	v1.Use(middleware.Authentication())
	v1.Post("/upload", h.UploadFile)
	v1.Get("/presigned/:id", h.GetPresignedUrl)
	v1.Get("/:id", h.DownloadFile)
	v1.Delete("/:id", h.DeleteFile)
}

func NewFileHandler(log *zap.Logger, service service.IFileService, validator *validator.Validate) *FileHandler {
	return &FileHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		h.log.Error("failed parse retrieve file from form")
		return err
	}

	res, err := h.service.UploadFile(file)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success upload file")
}

func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	res, err := h.service.DownloadFile(id)
	if err != nil {
		return err
	}

	c.Attachment(fmt.Sprintf("%s.%s", res.Metadata["real-name"], strings.SplitN(res.Metadata["content-type"], "/", 2)[1]))
	return c.SendStream(res.Body)
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	err := h.service.DeleteFile(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *FileHandler) GetPresignedUrl(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	res, err := h.service.GetPresignedUrl(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get presigned url")
}
