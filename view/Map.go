package view

import (
	"log"
	"net/http"

	"go/goserver.io/lib"
)

/**
 * Template handler for request: /map
 */
func Map(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:Map: " + r.Method)

	ctx := lib.CreateContent(r)

	if r.Method == "POST" {

	} else {

	}
	lib.BuildTemplate(w, "map.html", ctx)
}
