package handlers

import (
	"fmt"
	"net/http"

	"example.com/taskservice/internal/domain"
	taskdomain "example.com/taskservice/internal/domain/task"
	domainTypes "example.com/taskservice/internal/domain/types"
)

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req taskMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		h.handleError(w, domain.NewBadRequest("invalid json", err))
		return
	}

	input, err := taskdomain.NewCreateInput(req.Title, req.Description, req.Status, req.ScheduledAt, req.RecurrenceType,
		req.Settings)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	created, err := h.usecase.CreateTask(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newTaskDTO(created))
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	positiveID, err := domainTypes.NewPositiveInt(id)
	if err != nil {
		err = fmt.Errorf("invalid id: %v", err)
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
	}

	task, rule, err := h.usecase.GetTaskByID(r.Context(), positiveID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := taskWithRuleDTO{
		Task: newTaskDTO(task),
		Rule: newRuleDTO(rule),
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	positiveID, err := domainTypes.NewPositiveInt(id)
	if err != nil {
		err = fmt.Errorf("invalid id: %v", err)
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
	}

	var req taskUpdateDTO
	if err := decodeJSON(r, &req); err != nil {
		h.handleError(w, domain.NewBadRequest("invalid json", err))
		return
	}

	input, err := taskdomain.NewUpdateInput(req.Title, req.Description, req.Status, req.ScheduledAt)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	updated, err := h.usecase.UpdateTask(r.Context(), positiveID, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newTaskDTO(updated))
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	positiveID, err := domainTypes.NewPositiveInt(id)
	if err != nil {
		err = fmt.Errorf("invalid id: %v", err)
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
	}

	if err := h.usecase.DeleteTask(r.Context(), positiveID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListTask(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.usecase.ListTask(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := make([]taskDTO, 0, len(tasks))
	for i := range tasks {
		response = append(response, newTaskDTO(&tasks[i]))
	}

	writeJSON(w, http.StatusOK, response)
}
