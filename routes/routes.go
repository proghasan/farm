package routes

import (
	"farm/config"
	"farm/database"
	"farm/handlers"
	"farm/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, cfg *config.Config) {
	app.Use(cors.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Auth routes (public)
	auth := app.Group("/api/auth")
	auth.Post("/login", middleware.Login(cfg))

	// Protected routes
	api := app.Group("/api", middleware.JWTAuth(cfg.JWTSecret))

	// Users
	users := api.Group("/users")
	users.Get("/", wrapH(handlers.ListUsers))
	users.Get("/profile", wrapH(handlers.GetProfile))
	users.Get("/:id", wrapH(handlers.GetUser))
	users.Post("/", wrapH(handlers.CreateUser))
	users.Put("/:id", wrapH(handlers.UpdateUser))
	users.Delete("/:id", wrapH(handlers.DeleteUser))

	// Species
	species := api.Group("/species")
	species.Get("/", wrapH(handlers.ListSpecies))
	species.Get("/:id", wrapH(handlers.GetSpecies))
	species.Post("/", wrapH(handlers.CreateSpecies))
	species.Put("/:id", wrapH(handlers.UpdateSpecies))
	species.Delete("/:id", wrapH(handlers.DeleteSpecies))

	// Breeds
	breeds := api.Group("/breeds")
	breeds.Get("/", wrapH(handlers.ListBreeds))
	breeds.Get("/:id", wrapH(handlers.GetBreed))
	breeds.Post("/", wrapH(handlers.CreateBreed))
	breeds.Put("/:id", wrapH(handlers.UpdateBreed))
	breeds.Delete("/:id", wrapH(handlers.DeleteBreed))

	// Animals
	animals := api.Group("/animals")
	animals.Get("/", wrapH(handlers.ListAnimals))
	animals.Get("/:id", wrapH(handlers.GetAnimal))
	animals.Get("/:id/profile", wrapH(handlers.GetAnimalProfile))
	animals.Post("/", wrapH(handlers.CreateAnimal))
	animals.Put("/:id", wrapH(handlers.UpdateAnimal))
	animals.Delete("/:id", wrapH(handlers.DeleteAnimal))

	// Weight Histories
	weights := api.Group("/weight-histories")
	weights.Get("/", wrapH(handlers.ListWeightHistories))
	weights.Get("/:id", wrapH(handlers.GetWeightHistory))
	weights.Post("/", wrapH(handlers.CreateWeightHistory))
	weights.Put("/:id", wrapH(handlers.UpdateWeightHistory))
	weights.Delete("/:id", wrapH(handlers.DeleteWeightHistory))

	// Vaccines
	vaccines := api.Group("/vaccines")
	vaccines.Get("/", wrapH(handlers.ListVaccines))
	vaccines.Get("/:id", wrapH(handlers.GetVaccine))
	vaccines.Post("/", wrapH(handlers.CreateVaccine))
	vaccines.Put("/:id", wrapH(handlers.UpdateVaccine))
	vaccines.Delete("/:id", wrapH(handlers.DeleteVaccine))

	// Vaccinations
	vaccinations := api.Group("/vaccinations")
	vaccinations.Get("/", wrapH(handlers.ListVaccinations))
	vaccinations.Get("/:id", wrapH(handlers.GetVaccination))
	vaccinations.Post("/", wrapH(handlers.CreateVaccination))
	vaccinations.Put("/:id", wrapH(handlers.UpdateVaccination))
	vaccinations.Delete("/:id", wrapH(handlers.DeleteVaccination))

	// Inventory Categories
	invCats := api.Group("/inventory-categories")
	invCats.Get("/", wrapH(handlers.ListInventoryCategories))
	invCats.Get("/:id", wrapH(handlers.GetInventoryCategory))
	invCats.Post("/", wrapH(handlers.CreateInventoryCategory))
	invCats.Put("/:id", wrapH(handlers.UpdateInventoryCategory))
	invCats.Delete("/:id", wrapH(handlers.DeleteInventoryCategory))

	// Inventory Items
	invItems := api.Group("/inventory-items")
	invItems.Get("/", wrapH(handlers.ListInventoryItems))
	invItems.Get("/:id", wrapH(handlers.GetInventoryItem))
	invItems.Post("/", wrapH(handlers.CreateInventoryItem))
	invItems.Put("/:id", wrapH(handlers.UpdateInventoryItem))
	invItems.Delete("/:id", wrapH(handlers.DeleteInventoryItem))

	// Inventory Transactions
	invTxns := api.Group("/inventory-transactions")
	invTxns.Get("/", wrapH(handlers.ListInventoryTransactions))
	invTxns.Get("/:id", wrapH(handlers.GetInventoryTransaction))
	invTxns.Post("/", wrapH(handlers.CreateInventoryTransaction))
	invTxns.Put("/:id", wrapH(handlers.UpdateInventoryTransaction))
	invTxns.Delete("/:id", wrapH(handlers.DeleteInventoryTransaction))

	// Account Heads
	acctHeads := api.Group("/account-heads")
	acctHeads.Get("/", wrapH(handlers.ListAccountHeads))
	acctHeads.Get("/:id", wrapH(handlers.GetAccountHead))
	acctHeads.Post("/", wrapH(handlers.CreateAccountHead))
	acctHeads.Put("/:id", wrapH(handlers.UpdateAccountHead))
	acctHeads.Delete("/:id", wrapH(handlers.DeleteAccountHead))

	// Account Transactions
	acctTxns := api.Group("/account-transactions")
	acctTxns.Get("/", wrapH(handlers.ListAccountTransactions))
	acctTxns.Get("/:id", wrapH(handlers.GetAccountTransaction))
	acctTxns.Post("/", wrapH(handlers.CreateAccountTransaction))
	acctTxns.Put("/:id", wrapH(handlers.UpdateAccountTransaction))
	acctTxns.Delete("/:id", wrapH(handlers.DeleteAccountTransaction))

	// Pregnancies
	pregnancies := api.Group("/pregnancies")
	pregnancies.Get("/", wrapH(handlers.ListPregnancies))
	pregnancies.Get("/:id", wrapH(handlers.GetPregnancy))
	pregnancies.Post("/", wrapH(handlers.CreatePregnancy))
	pregnancies.Put("/:id", wrapH(handlers.UpdatePregnancy))
	pregnancies.Delete("/:id", wrapH(handlers.DeletePregnancy))
}

func wrapH(fn func(c fiber.Ctx, db *gorm.DB) error) fiber.Handler {
	return func(c fiber.Ctx) error {
		return fn(c, database.DB)
	}
}
