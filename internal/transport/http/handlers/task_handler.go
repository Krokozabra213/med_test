package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Krokozabra213/common/types"
	"github.com/gorilla/mux"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Usecase interface {
	Create(ctx context.Context, input *taskdomain.CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id types.PositiveInt[int64]) (*taskdomain.Task, error)
	Update(ctx context.Context, id types.PositiveInt[int64], input *taskdomain.UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id types.PositiveInt[int64]) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type TaskHandler struct {
	usecase Usecase
}

func NewTaskHandler(usecase Usecase) *TaskHandler {
	return &TaskHandler{usecase: usecase}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req taskMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	input, err := taskdomain.NewCreateInput(req.Title, req.Description, req.Status, req.ScheduledAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	created, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newTaskDTO(created))
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	positiveID, err := types.NewPositiveInt(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))
	}

	task, err := h.usecase.GetByID(r.Context(), positiveID)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newTaskDTO(task))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	positiveID, err := types.NewPositiveInt(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))
	}

	var req taskMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	input, err := taskdomain.NewUpdateInput(req.Title, req.Description, req.Status, req.ScheduledAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := h.usecase.Update(r.Context(), positiveID, input)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newTaskDTO(updated))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIntIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	positiveID, err := types.NewPositiveInt(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid id: %v", err))
	}

	if err := h.usecase.Delete(r.Context(), positiveID); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.usecase.List(r.Context())
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	response := make([]taskDTO, 0, len(tasks))
	for i := range tasks {
		response = append(response, newTaskDTO(&tasks[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func getIntIDFromRequest(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]
	if rawID == "" {
		return 0, errors.New("missing task id")
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("wrong id format")
	}

	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, taskdomain.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}
