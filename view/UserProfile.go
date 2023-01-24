package view

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * Template handler for request: /userProfile
 */
func UserProfile(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:UserProfile: " + r.Method)

	ctx := lib.CreateContent(r)
	form := userProfileForm(r)
	ctx["form"] = form.GetFields

	if r.Method == "POST" {
		res, err := lib.ApiRequest(nil, r, constants.ProxyProfile, form.ToMap())
		error := app.GetError(err, res)
		if error != nil {
			ctx["error"] = error
		} else {
			token := lib.CreateJWT(res.Body)
			app.SetCookie(w, constants.TOKEN_COOKIE_NAME, token)
			app.Redirect(w, r, "/")
			return
		}
	} else {
		res, err := lib.ApiRequest(nil, r, constants.ProxyProfile)
		if app.ErrorHandler(w, err, res) {
			return
		}
		setUserProfileValues(form, res.Body)
	}
	lib.ExecTemplate(w, lib.Template("ping/userProfile.html"), ctx)
}

func userProfileForm(r *http.Request) *lib.Form {
	form := lib.InitForm(r, nil)
	form.CharField("UserName", "Username:", "search", `readonly="readonly"`)
	form.CharField("FirstName", "First Name:", "text", `required`)
	form.CharField("LastName", "Last Name:", "text", `required`)
	form.CharField("CellPhone", "Cell Phone:", "text", ``)
	form.CharField("WorkPhone", "Work Phone:", "text", ``)
	form.CharField("Title", "Title:", "text", ``)
	form.CharField("UserOrg__OrgName", "Company Name:", "text", ``)
	form.CharField("UserOrg__OrgKey", "", "hidden", ``)
	form.CharField("UserAddress__Country", "Country:", "text", ``)
	form.CharField("UserAddress__AddressLine2", "Location:", "text", ``)
	form.CharField("ExtendedUserAttribute__AVATAR", "", "hidden", ``)
	return form
}

func setUserProfileValues(form *lib.Form, body map[string]interface{}) {
	for key, value := range body {
		data, ok := value.(map[string]interface{})
		if ok {
			for subname, val := range data {
				f := form.Get(key + "__" + subname)
				if f != nil {
					field := f.(*lib.Field)
					field.SetValue(val.(string))
				}
			}
		} else {
			f := form.Get(key)
			if f != nil {
				field := f.(*lib.Field)
				field.SetValue(value.(string))
			}
		}
	}
}

/**
 * Upload image handler for request: /load
 */
func LoadImage(w http.ResponseWriter, r *http.Request, app *lib.App) {
	url := r.Header.Get("URL")
	req, err := http.NewRequest("GET", url, nil)
	for key, val := range r.Header {
		if !strings.Contains(constants.RESTRICT_KEYS, fmt.Sprintf("'%s'", key)) {
			req.Header.Set(key, fmt.Sprintf("%v", val))
		}
	}
	client := &http.Client{
		Timeout: time.Duration(constants.HTTP_TIMEOUT * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer res.Body.Close()
	for name, values := range res.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
