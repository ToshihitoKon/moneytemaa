package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ToshihitoKon/moneytemaa/src/db"
	"github.com/ToshihitoKon/moneytemaa/src/slack"
	"github.com/ToshihitoKon/moneytemaa/src/web"
	"github.com/gorilla/mux"
)

var DB *sql.DB

func main() {
	db.NewDB()
	fmt.Println("hello")

	router := mux.NewRouter()
	web.SetHandler(router.PathPrefix("/api").Subrouter())
	slack.SetHandler(router.PathPrefix("/slack").Subrouter())

	router.Walk(func(route *mux.Route, rtr *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("path:", pathTemplate)
		}
		return nil
	})

	log.Fatal(http.ListenAndServe(":5000", router))
}
