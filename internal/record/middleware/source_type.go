package middleware

import (
	"log"
	"net/http"

	"golang.org/x/exp/slices"
)

type SourceType struct {
	sourceTypes []string
}

func (st *SourceType) Populate(sourceTypes []string) {
	st.sourceTypes = sourceTypes
}

func (st *SourceType) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sourceType := r.Header.Get("Source-Type")
		if slices.Contains(st.sourceTypes, sourceType) {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Invalid source type: %s\n", sourceType)
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	})
}
