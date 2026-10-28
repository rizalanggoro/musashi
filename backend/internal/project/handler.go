package project

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/platform/httperr"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(c fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "25"))

	data, total, err := h.svc.List(q, page, perPage)
	if err != nil {
		return httperr.Internal(c, "Failed to fetch projects")
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 25
	}

	return c.JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input CreateProjectInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Create(input)
	if err != nil {
		if errors.Is(err, ErrProjectNameRequired) {
			return httperr.BadRequest(c, "Project name is required", map[string]string{"name": "Project name is required"})
		}
		if errors.Is(err, ErrProjectAlreadyExists) {
			return httperr.Conflict(c, err.Error())
		}
		return httperr.Internal(c, "Failed to create project")
	}

	return c.Status(fiber.StatusCreated).JSON(dto)
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	dto, err := h.svc.Get(id)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to fetch project")
	}

	return c.JSON(dto)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	var input UpdateProjectInput
	if err := c.Bind().Body(&input); err != nil {
		return httperr.BadRequest(c, "Invalid request body", nil)
	}

	dto, err := h.svc.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		if errors.Is(err, ErrProjectAlreadyExists) {
			return httperr.Conflict(c, err.Error())
		}
		return httperr.Internal(c, "Failed to update project")
	}

	return c.JSON(dto)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httperr.BadRequest(c, "Invalid project ID", nil)
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, ErrProjectNotFound) {
			return httperr.NotFound(c, "Project not found")
		}
		return httperr.Internal(c, "Failed to delete project")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Backward-compatible helpers
func ListHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).List
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Create
}

func GetHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Get
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Update
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return NewHandler(NewService(NewRepository(db))).Delete
}
