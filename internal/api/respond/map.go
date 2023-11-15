package respond

import (
	"encoding/json"
	"net/http"
)

func WithMap[a comparable, b any](w http.ResponseWriter, m map[a]b) {
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		InternalServerError(w, "не могу размаршалить мапу %v", err)
	}
}
