package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"siransbach/taskmanagementapi/auth"
	"siransbach/taskmanagementapi/fiberx"
	"siransbach/taskmanagementapi/tasks"
)

func (h *handlers) employerCreateTask(c *fiber.Ctx) error {
	var taskRequest struct {
		Title          string    `json:"title"`
		Description    string    `json:"description"`
		AssignedUserID int       `json:"assignedUserID"`
		DueDate        time.Time `json:"dueDate"`
	}
	if err := c.BodyParser(&taskRequest); err != nil {
		log.Err(err).Msg("could not parse task request")
		return fiberx.Err(c, fiber.StatusBadRequest)
	}

	if taskRequest.AssignedUserID > 0 {
		user, err := auth.NewDB(h.pg).FindOne(c.Context(), auth.FindOptions{
			IDs: []int{taskRequest.AssignedUserID},
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error().Msg(fmt.Sprintf("user %d does not exist", taskRequest.AssignedUserID))
				return fiberx.Err(c, fiber.StatusBadRequest)
			}
			log.Err(err).Msg("could not find user")
			return fiberx.Err(c, fiber.StatusInternalServerError)
		}
		if user.IsEmployer() {
			log.Error().Msg(fmt.Sprintf("user %d is not an employee", taskRequest.AssignedUserID))
			return fiberx.Err(c, fiber.StatusBadRequest, "assignedUserID is not an employee")
		}
	}

	task := tasks.Entry{
		Title:          taskRequest.Title,
		Description:    taskRequest.Description,
		AssignedUserID: taskRequest.AssignedUserID,
		DueDate:        taskRequest.DueDate,
	}
	id, err := tasks.NewDB(h.pg).Insert(c.Context(), task)
	if err != nil {
		log.Err(err).Msg("could not create task")
		return fiberx.Err(c, fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"taskId": id,
	})
}

func (h *handlers) employerGetTasks(c *fiber.Ctx) error {
	var (
		// search criteria
		assignedUserID int
		status         tasks.Status
		sortBy         tasks.DBColumn
		sortOrder      tasks.SortOrder

		err error
	)
	if v := c.Query("assignedUserId"); v != "" {
		assignedUserID, err = strconv.Atoi(v)
		if err != nil {
			log.Err(err).Msg("could not parse assignedUserID")
			return fiberx.Err(c, fiber.StatusBadRequest)
		}
	}
	if v := c.Query("status"); v != "" {
		status, err = tasks.ParseStatus(v)
		if err != nil {
			log.Err(err).Msg("could not parse status")
			return fiberx.Err(c, fiber.StatusBadRequest)
		}
	}
	if v := c.Query("sortBy"); v != "" {
		sortBy, err = tasks.ParseDBColumn(v)
		if err != nil {
			log.Err(err).Msg("could not parse sortBy")
			return fiberx.Err(c, fiber.StatusBadRequest)
		}
	}
	if v := c.Query("sortOrder"); v != "" {
		sortOrder, err = tasks.ParseSortOrder(v)
		if err != nil {
			log.Err(err).Msg("could not parse sortOrder")
			return fiberx.Err(c, fiber.StatusBadRequest)
		}
	}

	opts := tasks.FindOptions{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	if assignedUserID > 0 {
		opts.AssignedUserIDs = []int{assignedUserID}
	}
	if status != "" {
		opts.Statuses = []tasks.Status{status}
	}

	entries, err := tasks.NewDB(h.pg).Find(c.Context(), opts)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(fiber.Map{
				"tasks": []tasks.Entry{},
			})
		}
		log.Err(err).Msg("could not find tasks")
		return fiberx.Err(c, fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"tasks": entries,
	})
}

func (h *handlers) employerGetTaskSummary(c *fiber.Ctx) error {
	summaries, err := tasks.NewDB(h.pg).Summarize(c.Context())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg("could not summarize tasks")
			return fiberx.Err(c, fiber.StatusInternalServerError)
		}
	}
	return c.JSON(fiber.Map{
		"summaries": summaries,
	})
}
