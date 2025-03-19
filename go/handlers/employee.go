package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"siransbach/taskmanagementapi/auth"
	"siransbach/taskmanagementapi/fiberx"
	"siransbach/taskmanagementapi/tasks"
)

func (h *handlers) employeeGetTasks(c *fiber.Ctx) error {
	user, err := auth.CurrentUser(c)
	if err != nil {
		log.Err(err).Msg("could not get current user")
		return fiberx.Err(c, fiber.StatusUnauthorized)
	}
	entries, err := tasks.NewDB(h.pg).Find(c.Context(),
		tasks.FindOptions{
			AssignedUserIDs: []int{user.ID},
		},
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Err(err).Msg("could not find tasks")
		return fiberx.Err(c, fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"tasks": entries,
	})
}

func (h *handlers) employeeUpdateTaskStatus(c *fiber.Ctx) error {
	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		log.Err(err).Msg("could not get current user")
		return fiberx.Err(c, fiber.StatusUnauthorized)
	}
	taskID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Err(err).Msg("could not parse task id")
		return fiberx.Err(c, fiber.StatusBadRequest)
	}
	status, err := tasks.ParseStatus(c.Params("status"))
	if err != nil {
		log.Err(err).Msg("could not parse task status")
		return fiberx.Err(c, fiber.StatusBadRequest)
	}
	if err := tasks.NewDB(h.pg).UpdateStatus(c.Context(), currentUser.ID, taskID, status); err != nil {
		log.Err(err).Msg("could not update task status")
		if errors.Is(err, sql.ErrNoRows) {
			return fiberx.Err(c, fiber.StatusNotFound)
		}
		return fiberx.Err(c, fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
