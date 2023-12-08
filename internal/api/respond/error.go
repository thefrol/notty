package respond

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	"gitlab.com/thefrol/notty/internal/entity"
)

// Error отвечает клиенту ошибкой
func Errorf(w http.ResponseWriter, code int, format string, args ...any) {
	w.WriteHeader(code)

	msg := fmt.Sprintf(format, args...)
	e := entity.Error{
		Code: strconv.Itoa(code),
		Desc: msg,
	}

	_, _, err := easyjson.MarshalToHTTPResponseWriter(e, w)
	if err != nil {
		http.Error(w, "не удалось демаршалить ошибку:"+msg+" в "+err.Error(), http.StatusInternalServerError) // todo отвечать структурой
		return
	}
}

func InternalServerError(w http.ResponseWriter, format string, args ...any) {
	Errorf(w, http.StatusInternalServerError, format, args...)
}

func BadRequest(w http.ResponseWriter, format string, args ...any) {
	Errorf(w, http.StatusBadRequest, format, args...)
}

func NotFound(w http.ResponseWriter, format string, args ...any) {
	Errorf(w, http.StatusNotFound, format, args...)
}
