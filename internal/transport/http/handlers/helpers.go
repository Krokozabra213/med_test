package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"example.com/taskservice/internal/domain"
	"github.com/gorilla/mux"
)

func (h *Handler) handleError(w http.ResponseWriter, err error) {

	appErr := getAppError(err)
	if appErr == nil {
		h.log.Error("internal error", "error", err.Error())
		writeError(w, http.StatusInternalServerError, "sorry, internal error")
		return
	}

	h.log.Error(appErr.Message, "error", appErr.Error(), "code", appErr.Code)

	httpStatus := mapCodeToHTTP(appErr.Code)
	writeError(w, httpStatus, appErr.Message)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func mapCodeToHTTP(code domain.Code) int {
	switch code {
	case domain.CodeBadRequest:
		return http.StatusBadRequest
	case domain.CodeNotFound:
		return http.StatusNotFound
	case domain.CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func getAppError(err error) *domain.AppError {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return nil
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
