package view

import (
	"fmt"
	"log"
	"net/http"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * Template handler for request: /public/forgotUsername
 */
func ForgotUsername(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:ForgotUsername: " + r.Method)
	useremail := lib.FormValue(r, "email")

	ctx := lib.CreateContent(r)
	form := emailForm(r, useremail)
	ctx["form"] = form.GetFields
	if r.Method == "POST" {
		res, err := lib.ApiRequest(nil, r, fmt.Sprintf("%s?email=%s", constants.ProxyUsername, useremail))
		error := app.GetError(err, res)
		if error != nil {
			emailField := form.Get("email").(*lib.Field)
			emailField.AddError(*error)
		} else {
			lib.BuildTemplate(w, "ping/forgotUsernameSuccess.html", ctx)
			return
		}
	}
	lib.BuildTemplate(w, "ping/forgotUsername.html", ctx)
}

func emailForm(r *http.Request, value string) *lib.Form {
	form := lib.InitForm(r, nil)
	form.CharField("email", "", "email", `required`, value)
	return form
}
