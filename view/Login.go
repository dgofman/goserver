package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dgofman/pongo2"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
	"go/goserver.io/utils"
)

var errorTemplate = `
<div class='error-messages'>
	<div>
		<span class='alert-eq-icon'></span>
		<span class='text-error'><pre>%s</pre></span>
	</div>
</div>
`
var defaultTemplate = createTemplate()

/**
 * Template handler for request: /login
 */
func Login(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:Login: " + r.Method)
	var _err string
	var username, password, rememberUsername string = "", "", ""

	if r.Method == "POST" {
		username = lib.FormValue(r, "username")
		password = lib.FormValue(r, "password")
		rememberUsername = r.FormValue("rememberUsername")

		if rememberUsername != "" {
			app.AddCookie(w, "username", lib.Encode(username), time.Now().AddDate(0, 1, 0))
		} else {
			app.DeleteCookie(w, "username")
		}

		res, err := lib.ApiRequestImpl(nil, r, constants.Authenticate, false, map[string]interface{}{
			"username": username,
			"password": password,
		})
		if err != nil {
			_err = err.Error()
		} else if res.Error != nil {
			if res.Error.Code == constants.EQX_ESE_011_EMAIL_MISSING_IN_PROFILE {
				app.Redirect(w, r, constants.URL_UpdateEmail+"?username="+username)
				return
			} else {
				_err = res.Error.Message
			}
		} else if res.StatusCode == 200 {
			token := lib.CreateJWT(res.Body)
			app.SetCookie(w, constants.TOKEN_COOKIE_NAME, token)
			app.Redirect(w, r, "/")
			return
		} else {
			log.Println(res.Body)
			app.Error(w, constants.DEFAULT_HTTP_ERROR)
			return
		}
	} else {
		app.DeleteCookie(w, constants.TOKEN_COOKIE_NAME)
		username = lib.Decode(app.GetCookie(r, "username"))
		rememberUsername = username
	}
	html := strings.Replace(defaultTemplate, "$username", username, -1)
	html = strings.Replace(html, "$focusField", map[bool]string{true: "password", false: "username"}[username != ""], -1)
	html = strings.Replace(html, "$rememberUsernameChecked", map[bool]string{true: "checked", false: ""}[rememberUsername != ""], -1)
	html = strings.Replace(html, "$rememberUsername", "rememberUsername", -1)

	if _err != "" {
		html = regexp.MustCompile("<div class='error-messages'>[\\s\\S]*?</div>").ReplaceAllString(html, fmt.Sprintf(errorTemplate, _err))
	} else {
		html = regexp.MustCompile("<div class='error-messages'>[\\s\\S]*?</div>").ReplaceAllString(html, "")
	}
	var t, err = pongo2.FromString(html)
	if err != nil {
		utils.LogError(err)
		app.Error(w, err.Error())
	} else {
		lib.ExecTemplate(w, t, pongo2.Context{})
	}
}

func Logout(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:Logout: " + r.Method)
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	app.DeleteCookie(w, constants.TOKEN_COOKIE_NAME)
	app.Redirect(w, r, "/")
}

func createTemplate() string {
	content, _ := ioutil.ReadFile(filepath.Join(lib.TemplatePath(), "/ping/login.html"))
	out := string(content)
	out = strings.Replace(out, "${basePath}", "/", -1)
	out = strings.Replace(out, "$basePath", "/", -1)
	out = strings.Replace(out, "$PingFedBaseURL", "", -1)
	out = strings.Replace(out, "$locale.getLanguage()", "en", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'title')", "GOService", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'header')", "GO SERVICE", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'loginTitle')", "Login", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'headerMessage')", "Welcome to Go Service. Please sign-in to continue.", -1)
	out = regexp.MustCompile("\\$templateMessages.getMessage\\(\\$messageKeyPrefix, .usernameTitle.\\)").ReplaceAllString(out, "Username")
	out = regexp.MustCompile("\\$templateMessages.getMessage\\(\\$messageKeyPrefix, .passwordTitle.\\)").ReplaceAllString(out, "Password")
	out = strings.Replace(out, "$templateMessages.getMessage($messageKeyPrefix, 'rememberUsernameTitle')", "Remember my username", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'forgotYour')", "Forgot your", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($goserverKeyPrefix, 'forgotOr')", "or", -1)
	out = strings.Replace(out, "$templateMessages.getMessage($messageKeyPrefix, 'signInButtonTitle')", "Log In", -1)
	out = strings.Replace(out, "$name", "username", -1)
	out = strings.Replace(out, "$pass", "password", -1)
	out = strings.Replace(out, "$url", "", -1)
	out = regexp.MustCompile("#if\\(\\$authnMessageKey\\)[\\s\\S]*?#end").ReplaceAllString(out, "")
	out = regexp.MustCompile("#if\\(\\$errorMessageKey\\)[\\s\\S]*?#end").ReplaceAllString(out, "")
	out = regexp.MustCompile("#if\\(\\$serverError\\)[\\s\\S]*?#end").ReplaceAllString(out, "")
	out = regexp.MustCompile("#if\\(\\$loginFailed \\|\\| \\(\\$rememberUsernameCookieExists && \\$enableRememberUsername\\) \\|\\| \\$isChainedUsernameAvailable\\)[\\s\\S]*?#end").ReplaceAllString(out, "\tdocument.getElementById('$$focusField').focus();")
	out = regexp.MustCompile("<!--[\\s\\S]*?-->").ReplaceAllString(out, "")
	return out
}
