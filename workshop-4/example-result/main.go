package main

import (
	"log"
	"temp-backend/internal/handler"
	"temp-backend/internal/repository"
	"temp-backend/internal/usecase"
	"temp-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize database
	db, err := database.NewDatabase("./users.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	transferRepo := repository.NewTransferRepository(db.DB)
	pointLedgerRepo := repository.NewPointLedgerRepository(db.DB)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, userRepo, pointLedgerRepo, db.DB)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	transferHandler := handler.NewTransferHandler(transferUsecase)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":      "User Management API with Transfer Feature",
			"version":      "2.1.0",
			"architecture": "Clean Architecture",
			"features":     []string{"Users CRUD", "Point Transfers", "Point Ledger"},
		})
	})

	// User routes
	userRoutes := app.Group("/users")
	userRoutes.Get("/", userHandler.GetUsers)
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Post("/", userHandler.CreateUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)

	// Transfer routes
	transferRoutes := app.Group("/transfers")
	transferRoutes.Post("/", transferHandler.CreateTransfer)
	transferRoutes.Get("/", transferHandler.GetTransfers)
	transferRoutes.Get("/:id", transferHandler.GetTransfer)

	log.Println("🚀 Server starting on port 3000 with Clean Architecture + Transfer Feature...")
	log.Fatal(app.Listen(":3000"))
}
