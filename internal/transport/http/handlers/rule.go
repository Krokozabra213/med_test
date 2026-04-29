package handlers

import (
	"fmt"
	"net/http"

	"example.com/taskservice/internal/domain"
	taskdomain "example.com/taskservice/internal/domain/task"
	domainTypes "example.com/taskservice/internal/domain/types"
)

func (h *Handler) UpdateRule(w http.ResponseWriter, r *http.Request) {
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

	var req ruleMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		h.handleError(w, domain.NewBadRequest("invalid json", err))
		return
	}

	input, err := taskdomain.NewUpdateRuleInput(req.Title, req.Description, req.ScheduledAt, req.RecurrenceType, req.Settings)
	if err != nil {
		h.handleError(w, domain.NewBadRequest(err.Error(), err))
		return
	}

	updated, err := h.usecase.UpdateRule(r.Context(), positiveID, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newRuleDTO(updated))
}

func (h *Handler) ListRule(w http.ResponseWriter, r *http.Request) {
	rules, err := h.usecase.ListRule(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := make([]ruleDTO, 0, len(rules))
	for i := range rules {
		response = append(response, newRuleDTO(&rules[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetRuleByID(w http.ResponseWriter, r *http.Request) {
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

	rule, err := h.usecase.GetRuleByID(r.Context(), positiveID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := newRuleDTO(rule)

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteRule(w http.ResponseWriter, r *http.Request) {
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

	if err := h.usecase.DeleteRule(r.Context(), positiveID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
