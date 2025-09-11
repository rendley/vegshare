package handler

import "github.com/go-chi/chi/v5"

func (h *TaskHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// GET /api/v1/admin/tasks/
	r.Get("/", h.GetAllTasks)

	// POST /api/v1/admin/tasks/{taskID}/accept
	r.Post("/{taskID}/accept", h.AcceptTask)

	// POST /api/v1/admin/tasks/{taskID}/complete
	r.Post("/{taskID}/complete", h.CompleteTask)

	// POST /api/v1/admin/tasks/{taskID}/fail
	r.Post("/{taskID}/fail", h.FailTask)

	return r
}
