package view

import (
	"log"
	"net/http"

	"go/goserver.io/lib"
)

/**
 * Template handler for request: /manage_devices
 */
func ManageDevices(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:ManageDevices: " + r.Method)

	ctx := lib.CreateContent(r)

	if r.Method == "POST" {

	} else {

	}
	lib.BuildTemplate(w, "manage_devices.html", ctx)
}
