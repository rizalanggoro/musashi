package export

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"musashi/backend/internal/domain"
	"musashi/backend/internal/platform/httperr"
)

func VolumeExcelHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var project domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		var rows []struct {
			domain.Volume
			WorkItem      string  `gorm:"column:work_item"`
			Unit          string  `gorm:"column:unit"`
			PlannedVolume float64 `gorm:"column:planned_volume"`
		}

		query := db.Table("volumes").
			Select("volumes.*, boq_items.work_item, boq_items.unit, boq_items.planned_volume").
			Joins("JOIN boq_items ON boq_items.id = volumes.boq_item_id").
			Where("boq_items.project_id = ?", projectID).
			Order("boq_items.created_at ASC, volumes.created_at ASC")

		if err := query.Scan(&rows).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch volume data")
		}

		f := excelize.NewFile()
		sheet := "Volume Take-off"
		f.SetSheetName("Sheet1", sheet)

		// Header project info
		f.SetCellValue(sheet, "A1", "PROJECT:")
		f.SetCellValue(sheet, "B1", project.Name)
		f.SetCellValue(sheet, "A2", "SECTION:")
		f.SetCellValue(sheet, "B2", project.Section)
		f.SetCellValue(sheet, "A3", "PLANT:")
		f.SetCellValue(sheet, "B3", project.Plant)
		f.SetCellValue(sheet, "A4", "DATE:")
		f.SetCellValue(sheet, "B4", time.Now().Format("2006-01-02"))

		// Table headers
		headers := []string{"NO", "WORK ITEM", "VOLUME", "UNIT", "LENGTH (m)", "WIDTH (m)", "HEIGHT (m)", "QUANTITY", "UNIT WEIGHT", "TOTAL QUANTITY"}
		for colIdx, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 6)
			f.SetCellValue(sheet, cell, h)
		}

		// Data rows
		for i, r := range rows {
			rowNum := 7 + i
			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), r.WorkItem)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), r.PlannedVolume)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), r.Unit)
			if r.Length != nil {
				f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), *r.Length)
			}
			if r.Width != nil {
				f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), *r.Width)
			}
			if r.Height != nil {
				f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), *r.Height)
			}
			f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), r.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", rowNum), r.UnitWeight)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", rowNum), r.TotalQuantity)
		}

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			return httperr.Internal(c, "Failed to generate Excel file")
		}

		filename := fmt.Sprintf("Volume-%s-%s.xlsx", project.Name, time.Now().Format("20060102"))
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		return c.Send(buf.Bytes())
	}
}

func RabExcelHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var project domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		var items []domain.RabItem
		if err := db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch RAB data")
		}

		f := excelize.NewFile()
		sheet := "RAB"
		f.SetSheetName("Sheet1", sheet)

		// Header project info
		f.SetCellValue(sheet, "A1", "RENCANA ANGGARAN BIAYA (RAB)")
		f.SetCellValue(sheet, "A2", "PROJECT:")
		f.SetCellValue(sheet, "B2", project.Name)
		f.SetCellValue(sheet, "A3", "SECTION:")
		f.SetCellValue(sheet, "B3", project.Section)
		f.SetCellValue(sheet, "A4", "PLANT:")
		f.SetCellValue(sheet, "B4", project.Plant)
		f.SetCellValue(sheet, "A5", "DATE:")
		f.SetCellValue(sheet, "B5", time.Now().Format("2006-01-02"))

		// Table headers
		headers := []string{"NO", "WORK ITEM", "VOLUME", "UNIT", "UNIT PRICE", "TOTAL COST"}
		for colIdx, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 7)
			f.SetCellValue(sheet, cell, h)
		}

		var totalRab int64
		for i, r := range items {
			rowNum := 8 + i
			totalRab += r.TotalCost
			f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), r.WorkItem)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), r.Volume)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), r.Unit)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), r.UnitPrice)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), r.TotalCost)
		}

		totalRow := 8 + len(items)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", totalRow), "TOTAL RAB")
		f.SetCellValue(sheet, fmt.Sprintf("F%d", totalRow), totalRab)

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			return httperr.Internal(c, "Failed to generate Excel file")
		}

		filename := fmt.Sprintf("RAB-%s-%s.xlsx", project.Name, time.Now().Format("20060102"))
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		return c.Send(buf.Bytes())
	}
}

func RabPdfHandler(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		projectID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return httperr.BadRequest(c, "Invalid project ID", nil)
		}

		var project domain.Project
		if err := db.Where("id = ? AND deleted_at IS NULL", projectID).First(&project).Error; err != nil {
			return httperr.NotFound(c, "Project not found")
		}

		var items []domain.RabItem
		if err := db.Where("project_id = ?", projectID).Order("created_at ASC").Find(&items).Error; err != nil {
			return httperr.Internal(c, "Failed to fetch RAB data")
		}

		pdf := fpdf.New("L", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(0, 10, "MUSASHI - RENCANA ANGGARAN BIAYA (RAB)")
		pdf.Ln(12)

		pdf.SetFont("Arial", "", 10)
		pdf.Cell(30, 6, "PROJECT:")
		pdf.Cell(100, 6, project.Name)
		pdf.Cell(30, 6, "SECTION:")
		pdf.Cell(0, 6, project.Section)
		pdf.Ln(6)
		pdf.Cell(30, 6, "PLANT:")
		pdf.Cell(100, 6, project.Plant)
		pdf.Cell(30, 6, "DATE:")
		pdf.Cell(0, 6, time.Now().Format("2006-01-02"))
		pdf.Ln(10)

		// Table Header
		pdf.SetFillColor(250, 220, 220) // Pink #FADCDC
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(12, 8, "NO", "1", 0, "C", true, 0, "")
		pdf.CellFormat(110, 8, "WORK ITEM", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "VOLUME", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 8, "UNIT", "1", 0, "C", true, 0, "")
		pdf.CellFormat(45, 8, "UNIT PRICE", "1", 0, "C", true, 0, "")
		pdf.CellFormat(50, 8, "TOTAL COST", "1", 1, "C", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		var totalRab int64
		for i, r := range items {
			totalRab += r.TotalCost
			pdf.CellFormat(12, 7, strconv.Itoa(i+1), "1", 0, "C", false, 0, "")
			pdf.CellFormat(110, 7, r.WorkItem, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", r.Volume), "1", 0, "R", false, 0, "")
			pdf.CellFormat(25, 7, r.Unit, "1", 0, "C", false, 0, "")
			pdf.CellFormat(45, 7, fmt.Sprintf("Rp %s", formatNumber(r.UnitPrice)), "1", 0, "R", false, 0, "")
			pdf.CellFormat(50, 7, fmt.Sprintf("Rp %s", formatNumber(r.TotalCost)), "1", 1, "R", false, 0, "")
		}

		// Total Row
		pdf.SetFont("Arial", "B", 10)
		pdf.SetFillColor(228, 0, 15) // Red #E4000F
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(177, 8, "TOTAL RAB", "1", 0, "R", true, 0, "")
		pdf.CellFormat(95, 8, fmt.Sprintf("Rp %s", formatNumber(totalRab)), "1", 1, "R", true, 0, "")

		var buf bytes.Buffer
		if err := pdf.Output(&buf); err != nil {
			return httperr.Internal(c, "Failed to generate PDF file")
		}

		filename := fmt.Sprintf("RAB-%s-%s.pdf", project.Name, time.Now().Format("20060102"))
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		return c.Send(buf.Bytes())
	}
}

func formatNumber(n int64) string {
	in := strconv.FormatInt(n, 10)
	var out []byte
	l := len(in)
	for i, c := range in {
		if i > 0 && (l-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

