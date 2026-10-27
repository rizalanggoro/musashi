package schedule

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

type ScheduleDTO struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Code      string `json:"code"`
	WorkItem  string `json:"work_item"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Duration  int    `json:"duration"`
	Weight    string `json:"weight"`
	Progress  string `json:"progress"`
	Status    string `json:"status"`
}

type CreateScheduleInput struct {
	Code      string `json:"code"`
	WorkItem  string `json:"work_item"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Weight    string `json:"weight"`
	Progress  string `json:"progress"`
	Status    string `json:"status"`
}

type UpdateScheduleInput struct {
	Code      *string `json:"code"`
	WorkItem  *string `json:"work_item"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Weight    *string `json:"weight"`
	Progress  *string `json:"progress"`
	Status    *string `json:"status"`
}

func toDTO(s *domain.Schedule) ScheduleDTO {
	return ScheduleDTO{
		ID:        s.ID.String(),
		ProjectID: s.ProjectID.String(),
		Code:      s.Code,
		WorkItem:  s.WorkItem,
		StartDate: s.StartDate.Format("2006-01-02"),
		EndDate:   s.EndDate.Format("2006-01-02"),
		Duration:  s.Duration,
		Weight:    fmt.Sprintf("%.2f", s.Weight),
		Progress:  fmt.Sprintf("%.2f", s.Progress),
		Status:    string(s.Status),
	}
}

func calcDuration(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		return 1
	}
	return days
}

func ListByProjectHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var schedules []domain.Schedule
		if err := db.Where("project_id = ?", projectID).Order("start_date ASC, code ASC").Find(&schedules).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch schedules")
		}

		data := make([]ScheduleDTO, len(schedules))
		for i := range schedules {
			data[i] = toDTO(&schedules[i])
		}
		return c.JSON(data)
	}
}

func CreateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var input CreateScheduleInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}
		if input.WorkItem == "" || input.StartDate == "" || input.EndDate == "" {
			return httperr.BadRequest(c, "Work item, start date, and end date are required", nil)
		}

		start, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			return httperr.BadRequest(c, "Invalid start date format (YYYY-MM-DD)", nil)
		}
		end, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			return httperr.BadRequest(c, "Invalid end date format (YYYY-MM-DD)", nil)
		}

		weight, _ := strconv.ParseFloat(input.Weight, 64)
		progress, _ := strconv.ParseFloat(input.Progress, 64)

		status := domain.SchedulePlanning
		if input.Status != "" {
			status = domain.ScheduleStatus(input.Status)
		}

		item := domain.Schedule{
			ProjectID: projectID,
			Code:      input.Code,
			WorkItem:  input.WorkItem,
			StartDate: start,
			EndDate:   end,
			Duration:  calcDuration(start, end),
			Weight:    weight,
			Progress:  progress,
			Status:    status,
		}

		if err := db.Create(&item).Error; err != nil {
			return httperr.Internal(c, "Failed to create schedule item")
		}

		return c.Status(fiber.StatusCreated).JSON(toDTO(&item))
	}
}

func UpdateHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid schedule ID", nil)
		}

		var s domain.Schedule
		if err := db.First(&s, id).Error; err != nil {
			return httperr.NotFound(c, "Schedule item not found")
		}

		var input UpdateScheduleInput
		if err := c.Bind().Body(&input); err != nil {
			return httperr.BadRequest(c, "Invalid request body", nil)
		}

		if input.Code != nil {
			s.Code = *input.Code
		}
		if input.WorkItem != nil {
			s.WorkItem = *input.WorkItem
		}
		if input.StartDate != nil {
			if t, err := time.Parse("2006-01-02", *input.StartDate); err == nil {
				s.StartDate = t
			}
		}
		if input.EndDate != nil {
			if t, err := time.Parse("2006-01-02", *input.EndDate); err == nil {
				s.EndDate = t
			}
		}
		s.Duration = calcDuration(s.StartDate, s.EndDate)

		if input.Weight != nil {
			w, _ := strconv.ParseFloat(*input.Weight, 64)
			s.Weight = w
		}
		if input.Progress != nil {
			p, _ := strconv.ParseFloat(*input.Progress, 64)
			s.Progress = p
		}
		if input.Status != nil {
			s.Status = domain.ScheduleStatus(*input.Status)
		}

		if err := db.Save(&s).Error; err != nil {
			return httperr.Internal(c, "Failed to update schedule item")
		}

		return c.JSON(toDTO(&s))
	}
}

func DeleteHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid schedule ID", nil)
		}

		res := db.Delete(&domain.Schedule{}, id)
		if res.Error != nil {
			return httperr.Internal(c, "Failed to delete schedule item")
		}
		if res.RowsAffected == 0 {
			return httperr.NotFound(c, "Schedule item not found")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

