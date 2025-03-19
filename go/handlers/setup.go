package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"siransbach/taskmanagementapi/auth"
)

const apiVersion = "v1"

type handlers struct {
	pg *sql.DB
}

func Setup(app *fiber.App, pg *sql.DB) {
	h := handlers{pg: pg}

	app.Route("/api/"+apiVersion, func(api fiber.Router) {
		employeeRoutes := api.Group("/employee", userMustHaveRole(auth.RoleEmployee))
		employeeRoutes.Route("/tasks", func(tasks fiber.Router) {
			tasks.Get("/", h.employeeGetTasks)
			tasks.Put("/:id/status/:status", h.employeeUpdateTaskStatus)
		})

		employerRoutes := api.Group("/employer", userMustHaveRole(auth.RoleEmployer))
		employerRoutes.Route("/tasks", func(tasks fiber.Router) {
			tasks.Get("/", h.employerGetTasks)
			tasks.Post("/", h.employerCreateTask)
			tasks.Get("/summary", h.employerGetTaskSummary)
		})
	})
}
