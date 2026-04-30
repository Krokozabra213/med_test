package transporthttp

import (
	"net/http"

	"github.com/gorilla/mux"

	swaggerdocs "example.com/taskservice/internal/transport/http/docs"
	httphandlers "example.com/taskservice/internal/transport/http/handlers"
)

func NewRouter(handler *httphandlers.Handler, docsHandler *swaggerdocs.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/tasks", handler.CreateTask).Methods(http.MethodPost)
	api.HandleFunc("/tasks", handler.ListTask).Methods(http.MethodGet)
	api.HandleFunc("/tasks/{id:[0-9]+}", handler.GetTaskByID).Methods(http.MethodGet)
	api.HandleFunc("/tasks/{id:[0-9]+}", handler.UpdateTask).Methods(http.MethodPut)
	api.HandleFunc("/tasks/{id:[0-9]+}", handler.DeleteTask).Methods(http.MethodDelete)

	api.HandleFunc("/rules", handler.ListRule).Methods(http.MethodGet)
	api.HandleFunc("/rules/{id:[0-9]+}", handler.GetRuleByID).Methods(http.MethodGet)
	api.HandleFunc("/rules/{id:[0-9]+}", handler.UpdateRule).Methods(http.MethodPut)
	api.HandleFunc("/rules/{id:[0-9]+}", handler.DeleteRule).Methods(http.MethodDelete)

	return router
}
