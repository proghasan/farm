package routes

import (
	"farm/internal/config"
	"farm/internal/database"
	"farm/internal/handlers"
	"farm/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func Setup(app *fiber.App, cfg *config.Config) {
	app.Use(cors.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := app.Group("/api/auth")
	auth.Post("/login", middleware.Login(cfg))

	api := app.Group("/api", middleware.JWTAuth(cfg.JWTSecret))

	userH := handlers.NewUserHandler(database.DB)
	speciesH := handlers.NewSpeciesHandler(database.DB)
	breedH := handlers.NewBreedHandler(database.DB)
	animalH := handlers.NewAnimalHandler(database.DB)
	weightH := handlers.NewWeightHandler(database.DB)
	vaccineH := handlers.NewVaccineHandler(database.DB)
	vaccinationH := handlers.NewVaccinationHandler(database.DB)
	categoryH := handlers.NewCategoryHandler(database.DB)
	invItemH := handlers.NewInventoryItemHandler(database.DB)
	invTxnH := handlers.NewInventoryTransactionHandler(database.DB)
	acctHeadH := handlers.NewAccountHeadHandler(database.DB)
	acctTxnH := handlers.NewAccountTransactionHandler(database.DB)
	pregnancyH := handlers.NewPregnancyHandler(database.DB)

	users := api.Group("/users")
	users.Get("/", userH.List)
	users.Get("/profile", userH.Profile)
	users.Get("/:id", userH.Get)
	users.Post("/", userH.Create)
	users.Put("/:id", userH.Update)
	users.Delete("/:id", userH.Delete)

	species := api.Group("/species")
	species.Get("/", speciesH.List)
	species.Get("/:id", speciesH.Get)
	species.Post("/", speciesH.Create)
	species.Put("/:id", speciesH.Update)
	species.Delete("/:id", speciesH.Delete)

	breeds := api.Group("/breeds")
	breeds.Get("/", breedH.List)
	breeds.Get("/:id", breedH.Get)
	breeds.Post("/", breedH.Create)
	breeds.Put("/:id", breedH.Update)
	breeds.Delete("/:id", breedH.Delete)

	animals := api.Group("/animals")
	animals.Get("/", animalH.List)
	animals.Get("/:id", animalH.Get)
	animals.Get("/:id/profile", animalH.Profile)
	animals.Post("/", animalH.Create)
	animals.Put("/:id", animalH.Update)
	animals.Delete("/:id", animalH.Delete)

	weights := api.Group("/weight-histories")
	weights.Get("/", weightH.List)
	weights.Get("/:id", weightH.Get)
	weights.Post("/", weightH.Create)
	weights.Put("/:id", weightH.Update)
	weights.Delete("/:id", weightH.Delete)

	vaccines := api.Group("/vaccines")
	vaccines.Get("/", vaccineH.List)
	vaccines.Get("/:id", vaccineH.Get)
	vaccines.Post("/", vaccineH.Create)
	vaccines.Put("/:id", vaccineH.Update)
	vaccines.Delete("/:id", vaccineH.Delete)

	vaccinations := api.Group("/vaccinations")
	vaccinations.Get("/", vaccinationH.List)
	vaccinations.Get("/:id", vaccinationH.Get)
	vaccinations.Post("/", vaccinationH.Create)
	vaccinations.Put("/:id", vaccinationH.Update)
	vaccinations.Delete("/:id", vaccinationH.Delete)

	invCats := api.Group("/inventory-categories")
	invCats.Get("/", categoryH.List)
	invCats.Get("/:id", categoryH.Get)
	invCats.Post("/", categoryH.Create)
	invCats.Put("/:id", categoryH.Update)
	invCats.Delete("/:id", categoryH.Delete)

	invItems := api.Group("/inventory-items")
	invItems.Get("/", invItemH.List)
	invItems.Get("/:id", invItemH.Get)
	invItems.Post("/", invItemH.Create)
	invItems.Put("/:id", invItemH.Update)
	invItems.Delete("/:id", invItemH.Delete)

	invTxns := api.Group("/inventory-transactions")
	invTxns.Get("/", invTxnH.List)
	invTxns.Get("/:id", invTxnH.Get)
	invTxns.Post("/", invTxnH.Create)
	invTxns.Put("/:id", invTxnH.Update)
	invTxns.Delete("/:id", invTxnH.Delete)

	acctHeads := api.Group("/account-heads")
	acctHeads.Get("/", acctHeadH.List)
	acctHeads.Get("/:id", acctHeadH.Get)
	acctHeads.Post("/", acctHeadH.Create)
	acctHeads.Put("/:id", acctHeadH.Update)
	acctHeads.Delete("/:id", acctHeadH.Delete)

	acctTxns := api.Group("/account-transactions")
	acctTxns.Get("/", acctTxnH.List)
	acctTxns.Get("/:id", acctTxnH.Get)
	acctTxns.Post("/", acctTxnH.Create)
	acctTxns.Put("/:id", acctTxnH.Update)
	acctTxns.Delete("/:id", acctTxnH.Delete)

	pregnancies := api.Group("/pregnancies")
	pregnancies.Get("/", pregnancyH.List)
	pregnancies.Get("/:id", pregnancyH.Get)
	pregnancies.Post("/", pregnancyH.Create)
	pregnancies.Put("/:id", pregnancyH.Update)
	pregnancies.Delete("/:id", pregnancyH.Delete)
}
