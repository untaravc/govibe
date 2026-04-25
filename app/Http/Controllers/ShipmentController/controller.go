package shipmentcontroller

import (
	"errors"
	"fmt"
	"strings"
	"time"

	authmiddleware "govibe/app/Http/Middleware/AuthMiddleware"
	"govibe/app/Http/Response"
	"govibe/app/Models"
	"govibe/app/Parser"
	appvalidator "govibe/app/Validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ShipmentController struct {
	db *gorm.DB
}

type shipmentStatusItem struct {
	Status uint   `json:"status"`
	Name   string `json:"name"`
}

const shipmentStatusCreated uint = 100
const shipmentLogStatusPlannedTransit uint = 0

func New(db *gorm.DB) *ShipmentController {
	return &ShipmentController{db: db}
}

func shipmentStatuses() []shipmentStatusItem {
	return []shipmentStatusItem{
		{Status: shipmentStatusCreated, Name: "Dibuat"},
		{Status: 200, Name: "Perjalanan"},
		{Status: 250, Name: "Tiba di Transit"},
		{Status: 251, Name: "Berangkat dari transit"},
		{Status: 300, Name: "Sampai Tujuan"},
	}
}

func shipmentStatusName(status uint) string {
	for _, s := range shipmentStatuses() {
		if s.Status == status {
			return s.Name
		}
	}
	return "Tidak diketahui"
}

func (ctl *ShipmentController) ShipmentStatusList(c *fiber.Ctx) error {
	return response.OK(c, "ok", fiber.Map{"shipment_statuses": shipmentStatuses()})
}

func (ctl *ShipmentController) Track(c *fiber.Ctx) error {
	code := strings.TrimSpace(c.Query("code"))
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "shipment code is required")
	}

	var s models.Shipment
	if err := ctl.db.Where("code = ?", code).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "shipment not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{
		"shipment": fiber.Map{
			"code":        s.Code,
			"status":      s.Status,
			"status_name": shipmentStatusName(s.Status),
			"created_at":  s.CreatedAt,
			"updated_at":  s.UpdatedAt,
		},
	})
}

func (ctl *ShipmentController) Index(c *fiber.Ctx) error {
	var shipments []models.Shipment
	q := ctl.db.Model(&models.Shipment{})

	shipmentType := strings.ToLower(strings.TrimSpace(c.Query("type")))
	if shipmentType != "" {
		statuses, ok := shipmentStatusesForType(shipmentType)
		if !ok {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "invalid shipment type")
		}
		q = q.Where("status IN ?", statuses)
	}

	if err := q.Order("id desc").Find(&shipments).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return response.OK(c, "ok", fiber.Map{"shipments": shipments})
}

func shipmentStatusesForType(shipmentType string) ([]uint, bool) {
	switch shipmentType {
	case "departure":
		return []uint{shipmentStatusCreated}, true
	case "transit":
		return []uint{200, 250, 251}, true
	case "arrive":
		return []uint{300}, true
	default:
		return nil, false
	}
}

func (ctl *ShipmentController) Show(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var s models.Shipment
	if err := ctl.db.
		Preload("Details").
		Preload("Logs", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", shipmentLogStatusPlannedTransit).Order("id asc")
		}).
		First(&s, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "shipment not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{"shipment": s})
}

func (ctl *ShipmentController) Store(c *fiber.Ctx) error {
	var req createShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	req.Code = strings.TrimSpace(req.Code)
	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerPhone = strings.TrimSpace(req.CustomerPhone)
	req.PriceType = strings.TrimSpace(req.PriceType)
	if req.CustomerEmail != nil {
		trimmed := strings.TrimSpace(*req.CustomerEmail)
		if trimmed == "" {
			req.CustomerEmail = nil
		} else {
			req.CustomerEmail = &trimmed
		}
	}
	for i := range req.Details {
		req.Details[i].ItemName = strings.TrimSpace(req.Details[i].ItemName)
	}
	if req.Code == "" {
		req.Code = generateShipmentCode()
	}
	req.Status = shipmentStatusCreated

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{"errors": errs})
	}

	userID, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	var created models.Shipment
	if err := ctl.db.Transaction(func(tx *gorm.DB) error {
		s := models.Shipment{
			Code:                req.Code,
			CustomerName:        req.CustomerName,
			OfficeOriginID:      req.OfficeOriginID,
			OfficeDestinationID: req.OfficeDestinationID,
			CustomerPhone:       req.CustomerPhone,
			CustomerEmail:       req.CustomerEmail,
			Price:               req.Price,
			UserID:              userID,
			Wight:               req.Wight,
			Length:              req.Length,
			Width:               req.Width,
			Height:              req.Height,
			PriceType:           req.PriceType,
			Status:              req.Status,
		}

		if err := tx.Create(&s).Error; err != nil {
			return err
		}

		if len(req.Details) > 0 {
			details := make([]models.ShipmentDetail, 0, len(req.Details))
			for _, d := range req.Details {
				details = append(details, models.ShipmentDetail{
					ShipmentID: s.ID,
					ItemName:   d.ItemName,
					ItemPrice:  d.ItemPrice,
					CategoryID: d.CategoryID,
				})
			}
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		if err := createTransitLogs(tx, s.ID, userID, req.Transits); err != nil {
			return err
		}

		if err := tx.
			Preload("Details").
			Preload("Logs", func(db *gorm.DB) *gorm.DB {
				return db.Where("status = ?", shipmentLogStatusPlannedTransit).Order("id asc")
			}).
			First(&created, s.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.Created(c, "created", fiber.Map{"shipment": created})
}

func (ctl *ShipmentController) Update(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req updateShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}

	if req.Code != nil {
		trimmed := strings.TrimSpace(*req.Code)
		req.Code = &trimmed
	}
	if req.CustomerName != nil {
		trimmed := strings.TrimSpace(*req.CustomerName)
		req.CustomerName = &trimmed
	}
	if req.CustomerPhone != nil {
		trimmed := strings.TrimSpace(*req.CustomerPhone)
		req.CustomerPhone = &trimmed
	}
	if req.PriceType != nil {
		trimmed := strings.TrimSpace(*req.PriceType)
		req.PriceType = &trimmed
	}
	if req.CustomerEmail != nil {
		trimmed := strings.TrimSpace(*req.CustomerEmail)
		if trimmed == "" {
			req.CustomerEmail = nil
		} else {
			req.CustomerEmail = &trimmed
		}
	}
	if req.Details != nil {
		for i := range *req.Details {
			(*req.Details)[i].ItemName = strings.TrimSpace((*req.Details)[i].ItemName)
		}
	}

	if errs := appvalidator.Validate(req); len(errs) > 0 {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation error", fiber.Map{"errors": errs})
	}

	userID, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	var updated models.Shipment
	if err := ctl.db.Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]any, 16)
		if req.Code != nil {
			updates["code"] = *req.Code
		}
		if req.CustomerName != nil {
			updates["customer_name"] = *req.CustomerName
		}
		if req.OfficeOriginID != nil {
			updates["office_origin_id"] = *req.OfficeOriginID
		}
		if req.OfficeDestinationID != nil {
			updates["office_destination_id"] = *req.OfficeDestinationID
		}
		if req.CustomerPhone != nil {
			updates["customer_phone"] = *req.CustomerPhone
		}
		if req.CustomerEmail != nil {
			updates["customer_email"] = *req.CustomerEmail
		}
		if req.Price != nil {
			updates["price"] = *req.Price
		}
		if req.Wight != nil {
			updates["wight"] = *req.Wight
		}
		if req.Length != nil {
			updates["length"] = *req.Length
		}
		if req.Width != nil {
			updates["width"] = *req.Width
		}
		if req.Height != nil {
			updates["height"] = *req.Height
		}
		if req.PriceType != nil {
			updates["price_type"] = *req.PriceType
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}

		if len(updates) > 0 {
			tx2 := tx.Model(&models.Shipment{}).Where("id = ?", uint(id)).Updates(updates)
			if tx2.Error != nil {
				return tx2.Error
			}
			if tx2.RowsAffected == 0 {
				return fiber.NewError(fiber.StatusNotFound, "shipment not found")
			}
		} else {
			var exists models.Shipment
			if err := tx.Select("id").First(&exists, uint(id)).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fiber.NewError(fiber.StatusNotFound, "shipment not found")
				}
				return err
			}
		}

		if req.Details != nil {
			if err := tx.Where("shipment_id = ?", uint(id)).Delete(&models.ShipmentDetail{}).Error; err != nil {
				return err
			}
			if len(*req.Details) > 0 {
				details := make([]models.ShipmentDetail, 0, len(*req.Details))
				for _, d := range *req.Details {
					details = append(details, models.ShipmentDetail{
						ShipmentID: uint(id),
						ItemName:   d.ItemName,
						ItemPrice:  d.ItemPrice,
						CategoryID: d.CategoryID,
					})
				}
				if err := tx.Create(&details).Error; err != nil {
					return err
				}
			}
		}

		if req.Transits != nil {
			if err := tx.
				Where("shipment_id = ?", uint(id)).
				Where("status = ?", shipmentLogStatusPlannedTransit).
				Delete(&models.ShipmentLog{}).Error; err != nil {
				return err
			}
			if err := createTransitLogs(tx, uint(id), userID, *req.Transits); err != nil {
				return err
			}
		}

		return tx.
			Preload("Details").
			Preload("Logs", func(db *gorm.DB) *gorm.DB {
				return db.Where("status = ?", shipmentLogStatusPlannedTransit).Order("id asc")
			}).
			First(&updated, uint(id)).Error
	}); err != nil {
		if fe, ok := err.(*fiber.Error); ok && fe != nil {
			return fe
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.OK(c, "ok", fiber.Map{"shipment": updated})
}

func (ctl *ShipmentController) Destroy(c *fiber.Ctx) error {
	id, err := parser.UintParam(c, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := ctl.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shipment_id = ?", uint(id)).Delete(&models.ShipmentDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("shipment_id = ?", uint(id)).Delete(&models.ShipmentLog{}).Error; err != nil {
			return err
		}
		res := tx.Delete(&models.Shipment{}, uint(id))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "shipment not found")
		}
		return nil
	}); err != nil {
		if fe, ok := err.(*fiber.Error); ok && fe != nil {
			return fe
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.OK(c, "deleted", fiber.Map{"deleted": true})
}

func authenticatedUserID(c *fiber.Ctx) (uint, error) {
	userID, ok := c.Locals(authmiddleware.LocalUserID).(uint)
	if !ok || userID == 0 {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "authenticated user not found")
	}
	return userID, nil
}

func generateShipmentCode() string {
	now := time.Now().UTC()
	return fmt.Sprintf("SHP-%s-%06d", now.Format("20060102150405"), now.Nanosecond()/1000)
}

func createTransitLogs(tx *gorm.DB, shipmentID uint, userID uint, transits []shipmentTransitInput) error {
	if len(transits) == 0 {
		return nil
	}

	logs := make([]models.ShipmentLog, 0, len(transits))
	for _, t := range transits {
		logs = append(logs, models.ShipmentLog{
			ShipmentID: shipmentID,
			OfficeID:   t.OfficeID,
			UserID:     userID,
			Status:     shipmentLogStatusPlannedTransit,
		})
	}

	return tx.Create(&logs).Error
}
