package view

import (
	"fmt"
	"log"
	"net/http"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * Template handler for request: /public/activateAccount
 */
func ActivateAccount(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:ActivateAccount: " + r.Method)

	token := r.URL.Query().Get("token")

	ctx := lib.CreateContent(r)
	form := userProfileForm(r)
	ctx["form"] = form.GetFields

	if r.Method == "POST" {
		res, err := lib.ApiRequest(nil, r, fmt.Sprintf("%s?token=%s", constants.ProxyProfile, token), form.ToMap())
		error := app.GetError(err, res)
		if error != nil {
			ctx["error"] = error
		} else {
			lib.BuildTemplate(w, "ping/activateAccountSuccess.html", ctx)
			return
		}
	} else {
		if token == "" {
			app.Redirect(w, r, "/")
			return
		}
		res, err := lib.ApiRequest(nil, r, fmt.Sprintf("%s?token=%s", constants.ProxyActivate, token))
		if app.ErrorHandler(w, err, res) {
			return
		}
		setUserProfileValues(form, res.Body)
	}
	lib.BuildTemplate(w, "ping/userProfile.html", ctx)
}
