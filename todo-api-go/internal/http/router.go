package httpapi

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Router(h Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/health", otelhttp.NewHandler(http.HandlerFunc(h.Health), "health"))

	// collection routes
	mux.Handle("/todos", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTodos(w, r)
		case http.MethodPost:
			h.CreateTodo(w, r)
		default:
			http.Error(w, "method not allowed", 405)
		}
	}), "todos"))

	// item routes
	mux.Handle("/todos/", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTodoByID(w, r)
		case http.MethodPut:
			h.UpdateTodo(w, r)
		default:
			http.Error(w, "method not allowed", 405)
		}
	}), "todo-by-id"))

	return mux
}
