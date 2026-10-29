package project

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"musashi/backend/internal/platform/httperr"
	"musashi/backend/utils"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// List godoc
// @id          listProjects
// @Summary     List projects
// @Tags        projects
// @Produce     json
// @Param       q query string false "Search query"
// @Param       page query int false "Page number" default(1)
// @Param       per_page query int false "Items per page" default(25)
// @Success     200 {object} ProjectListResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects [get]
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

// Create godoc
// @id          createProject
// @Summary     Create project
// @Tags        projects
// @Accept      json
// @Produce     json
// @Param       body body CreateProjectInput true "Project input"
// @Success     201 {object} Project
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects [post]
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

// Get godoc
// @id          getProject
// @Summary     Get project by ID
// @Tags        projects
// @Produce     json
// @Param       id path string true "Project ID"
// @Success     200 {object} Project
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id} [get]
func (h *Handler) Get(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Update godoc
// @id          updateProject
// @Summary     Update project
// @Tags        projects
// @Accept      json
// @Produce     json
// @Param       id path string true "Project ID"
// @Param       body body UpdateProjectInput true "Project update data"
// @Success     200 {object} Project
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     409 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id} [put]
func (h *Handler) Update(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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

// Delete godoc
// @id          deleteProject
// @Summary     Delete project
// @Tags        projects
// @Param       id path string true "Project ID"
// @Success     204 "No Content"
// @Failure     400 {object} httperr.ErrorResponse
// @Failure     404 {object} httperr.ErrorResponse
// @Failure     500 {object} httperr.ErrorResponse
// @Router      /projects/{id} [delete]
func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := utils.ParseID(c.Params("id"))
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
