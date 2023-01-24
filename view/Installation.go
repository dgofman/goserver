package view

import (
	"log"
	"net/http"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * Template handler for request: /installation
 */
func Installation(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:Installation: " + r.Method)

	ctx := lib.CreateContent(r)
	ctx["_url_api_status"] = constants.URL_AjaxStatus

	if r.Method == "POST" {
		ctx["type"] = r.URL.RawQuery
		lib.ExecTemplate(w, lib.Template("installationSuccess.html"), ctx)
		return
	} else {
		resp, err := lib.ApiRequest(nil, r, constants.ApiStatus)
		if app.ErrorHandler(w, err, resp) {
			return
		}
		status := resp.Body["status"].(string)
		if status == constants.STATUS_FOREIGN_USER ||
			status == constants.STATUS_NEW ||
			status == constants.STATUS_NO_USER {
			app.Redirect(w, r, "/")
			return
		}
		ctx["type"] = r.URL.RawQuery
		ctx["license_key"] = resp.Body["license_key"]
	}
	lib.ExecTemplate(w, lib.Template("installation.html"), ctx)
}

/**
 * AJAX handler for request: /status
 */
func Status(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, constants.ApiStatus)
}
