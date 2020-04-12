package errorpage

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	rest *httprouter.Router
)

func Init(r *httprouter.Router) {
	rest = r
}

func Start() {

	rest.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(405)
		w.Write([]byte("405"))
		return

	})
	// ERROR BRAK STRONY 404
	rest.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(404)
		w.Write([]byte("404"))
		return

	})

}
