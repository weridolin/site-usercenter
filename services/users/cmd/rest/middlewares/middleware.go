package middlewares

import "net/http"

type SuperAdminRequiredMiddleware struct {
}

func NewUserAgentMiddleware() *SuperAdminRequiredMiddleware {
	return &SuperAdminRequiredMiddleware{}
}

func (m *SuperAdminRequiredMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
