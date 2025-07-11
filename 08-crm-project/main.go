package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/leminkhoa/go-crm-basic/database"
	"github.com/leminkhoa/go-crm-basic/lead"
)

func main() {
	app := fiber.New()
	initDatabase()
	setupRoutes(app)
	app.Listen(3000)
	defer database.DBConn.Close()

}

func initDatabase() {
	database.Connect()
	DBConn := database.GetDB()
	fmt.Println("Connection opened to database")
	DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}
