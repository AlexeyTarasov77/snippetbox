package decorators

import (
	"net/http"
)

func RequirePost(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			errStatus := http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(errStatus), errStatus)
			return
		}
		f(w, r)
	}
}
