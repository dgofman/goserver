package view

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * Template handler for request: /public/resetPassword
 */
func ResetPassword(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:ResetPassword: " + r.Method)

	token := r.URL.Query().Get("token")

	ctx := lib.CreateContent(r)
	form := usernameForm(r)

	if r.Method == "POST" {
		username := lib.FormValue(r, "username")
		email := lib.FormValue(r, "email")  //First page reset by username and password
		if strings.Contains(token, email) { //Step two Reset password using token
			form = resetPasswordForm(r)
			ctx["form"] = form.GetFields
			ctx["token"] = token
			err := ValidatePassword(r.FormValue("Password"), r.FormValue("Confirm"))
			if err != nil {
				ctx["error"] = err.Error
			} else {
				res, err := lib.ApiRequestImpl(nil, r, fmt.Sprintf("%s?%s", constants.ProxyPassword, r.URL.RawQuery), false, map[string]interface{}{
					"username": username,
					"password": lib.Encode(r.FormValue("Password"), constants.SECRET_KEY+username),
				})
				error := app.GetError(err, res)
				if error != nil {
					ctx["error"] = error
				} else {
					lib.ExecTemplate(w, lib.Template("ping/resetPasswordSuccess.html"), ctx)
					return
				}
			}
		} else {
			res, err := lib.ApiRequest(nil, r, fmt.Sprintf("%s", constants.ProxyPassword), map[string]interface{}{
				"username": username,
				"email":    email,
			})
			error := app.GetError(err, res)
			if error != nil {
				ctx["error"] = error
			} else {
				lib.ExecTemplate(w, lib.Template("ping/resetPasswordSuccess.html"), ctx)
				return
			}
		}
	} else {
		if token != "" {
			res, err := lib.ApiRequest(nil, r, fmt.Sprintf("%s?%s", constants.ProxyPassword, r.URL.RawQuery))
			error := app.GetError(err, res)
			if error != nil {
				ctx["error"] = error
			} else {
				ctx["token"] = token
				usernameField := form.Get("username").(*lib.Field)
				usernameField.SetValue(res.Body["UserName"].(string))
			}
		}
	}
	ctx["form"] = form.GetFields
	lib.ExecTemplate(w, lib.Template("ping/resetPassword.html"), ctx)
}

func usernameForm(r *http.Request) *lib.Form {
	form := lib.InitForm(r, nil)
	form.CharField("username", "", "text", `required`)
	form.CharField("email", "", "email", `required`)
	return form
}

func resetPasswordForm(r *http.Request) *lib.Form {
	form := lib.InitForm(r, nil)
	form.CharField("username", "", "text", `required`)
	form.CharField("Password", "Create a password:", "password", `minlength="8" required`, nil)
	form.CharField("Confirm", "Confirm your password:", "password", `minlength="8" required`, nil)
	return form
}
