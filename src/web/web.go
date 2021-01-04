package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SetHandler(router *mux.Router) {
	router.HandleFunc("/transaction/new", HandlerNewTransaction)
}

func HandlerNewTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
