package handler

import (
	"strconv"
	"temp-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// TransferHandler handles HTTP requests for transfers
type TransferHandler struct {
	transferUsecase domain.TransferUsecase
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(transferUsecase domain.TransferUsecase) *TransferHandler {
	return &TransferHandler{
		transferUsecase: transferUsecase,
	}
}

// CreateTransfer handles POST /transfers
func (h *TransferHandler) CreateTransfer(c *fiber.Ctx) error {
	var req domain.TransferCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid request body",
		})
	}

	transfer, err := h.transferUsecase.CreateTransfer(&req)
	if err != nil {
		// Determine error type and status code
		statusCode := 400
		errorCode := "VALIDATION_ERROR"

		switch err.Error() {
		case "insufficient points":
			statusCode = 409
			errorCode = "INSUFFICIENT_POINTS"
		case "cannot transfer to yourself":
			statusCode = 422
			errorCode = "INVALID_OPERATION"
		case "sender user not found", "receiver user not found":
			statusCode = 400
			errorCode = "USER_NOT_FOUND"
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"error":   errorCode,
			"message": err.Error(),
		})
	}

	// Set the idempotency key in response header
	c.Set("Idempotency-Key", transfer.IdemKey)

	return c.Status(201).JSON(&domain.TransferCreateResponse{
		Transfer: transfer,
	})
}

// GetTransfer handles GET /transfers/:id
func (h *TransferHandler) GetTransfer(c *fiber.Ctx) error {
	idemKey := c.Params("id")
	if idemKey == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Transfer ID is required",
		})
	}

	transfer, err := h.transferUsecase.GetTransferByIdemKey(idemKey)
	if err != nil {
		if err.Error() == "transfer not found" {
			return c.Status(404).JSON(fiber.Map{
				"error":   "NOT_FOUND",
				"message": "Transfer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfer",
		})
	}

	return c.JSON(&domain.TransferGetResponse{
		Transfer: transfer,
	})
}

// GetTransfers handles GET /transfers?userId={id}
func (h *TransferHandler) GetTransfers(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId query parameter is required",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid userId",
		})
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 200 {
			pageSize = ps
		}
	}

	response, err := h.transferUsecase.GetTransfersByUserID(userID, page, pageSize)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(400).JSON(fiber.Map{
				"error":   "USER_NOT_FOUND",
				"message": "User not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfers",
		})
	}

	return c.JSON(response)
}
