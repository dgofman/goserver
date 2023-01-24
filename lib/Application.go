package lib

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgofman/pongo2"

	"go/goserver.io/constants"
	"go/goserver.io/utils"
)

/**
 * This function wraps http.HandlerFunc calls plus reference to this *App object
 */
type Handler func(http.ResponseWriter, *http.Request, *App)

/**
 * App definition
 */
type App struct {
	routes []Route
}

/**
 * Route definition
 */
type Route struct {
	pattern *regexp.Regexp
	handler Handler
}

/**
 * Load error.html template file
 */
var errorTmpl = Template("error.html")

/**
 * Register a function handler for HTTP URL
 */
func (app *App) UrlHandle(pattern string, handler Handler) {
	re := regexp.MustCompile("^" + pattern)
	route := Route{pattern: re, handler: handler}
	app.routes = append(app.routes, route)
}

/**
 * Get user cookies
 */
func (app *App) GetCookie(r *http.Request, name string) string {
	var cookie, err = r.Cookie(name)
	if err == nil {
		return cookie.Value
	} else {
		return ""
	}
}

/**
 * Add user cookies
 */
func (app *App) AddCookie(w http.ResponseWriter, name string, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Expires: expires, HttpOnly: true})
}

/**
 * Update user cookies or set cookie only for session
 */
func (app *App) SetCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, HttpOnly: true})
}

/**
 * Delete user cookies value by name
 */
func (app *App) DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: "", MaxAge: -1})
}

/**
 * Common redirect URL handler
 */
func (app *App) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 301)
}

/**
 * Generate an error page by using "templates/error.html" file
 */
func (app *App) Error(w http.ResponseWriter, err string) {
	ExecTemplate(w, errorTmpl, pongo2.Context{"error": fmt.Sprintf("%v", err)})
}

/**
 * Validate if response body is valid return false, otherwise generate an error page
 */
func (app *App) ErrorHandler(w http.ResponseWriter, err error, resp *constants.Response) bool {
	error := app.GetError(err, resp)
	if error != nil {
		app.Error(w, *error)
		return true
	}
	return false
}

/**
 * Validation of the HTTP calls and JSON result
 * @param error - the second argument of return function
 * @param resp - reference the the Response object
 * @return either NIL not errors or error message
 */
func (app *App) GetError(err error, resp *constants.Response) *string {
	var error string
	if err != nil {
		error = err.Error()
	} else if resp != nil && resp.Error != nil {
		error = resp.Error.Message
	} else {
		return nil
	}
	return &error
}

/**
 * Create application and refine default route for / URL
 */
func CreateApp(basedir string) *App {
	app := &App{}
	http.HandleFunc("/", httpMiddleware(app, http.FileServer(http.Dir(basedir))))
	return app
}

/**
 * All URL requests function-handler
 */
func httpMiddleware(app *App, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authTokenVerification(app, r) {
			app.Redirect(w, r, constants.URL_Login)
			return
		}
		for _, rt := range app.routes {
			matches := rt.pattern.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				rt.handler(w, r, app)
				return
			}
		}
		//w.WriteHeader(404)
		err := "File not found: " + r.URL.Path
		utils.LogError(err)
		app.Error(w, err)
	})
}

/**
 * Verify access to the URL without authentication JWT token
 */
func authTokenVerification(app *App, r *http.Request) bool {
	if constants.Whitelist_urls.MatchString(r.URL.Path) ||
		constants.Whitelist_reg_urls.MatchString(r.URL.Path) {
		return true
	}
	if r.Header.Get("X-Auth-Token") != "" {
		if validateHeader(r) {
			return true
		} else {
			log.Printf("Header: %s\n", r.Header)
			setTokenHeader(r, r.Header.Get("X-Auth-Token"))
			if validateHeader(r) {
				return true
			}
		}
	}
	setTokenHeader(r, app.GetCookie(r, constants.TOKEN_COOKIE_NAME))
	return validateHeader(r)
}

/**
 * Validate the request header
 */
func validateHeader(r *http.Request) bool {
	return r.Header.Get("X-Auth-Token") != "" &&
		r.Header.Get("X-Auth-Subject") != "" &&
		r.Header.Get("X-Auth-User-Email") != "" &&
		r.Header.Get("X-Auth-User-Fn") != "" &&
		r.Header.Get("X-Auth-User-Ln") != "" &&
		r.Header.Get("X-Auth-User-Name") != ""
}

/**
 * Set an authentication JWT token in the header
 */
func setTokenHeader(r *http.Request, token string) {
	if token != "" {
		json, err := JwtDecode(token)
		if err != nil {
			utils.LogError(err)
			return
		}
		r.Header.Set("X-Auth-Token", token)
		r.Header.Set("X-Auth-Subject", toString(json.Subject))
		r.Header.Set("X-Auth-User-Email", toString(json.Email))
		r.Header.Set("X-Auth-User-Fn", toString(json.FirstName))
		r.Header.Set("X-Auth-User-Ln", toString(json.LastName))
		r.Header.Set("X-Auth-User-Name", toString(json.UserName))
		r.Header.Set("X-Auth-User-Avatar", toString(json.Avatar))
	}
}

/**
 * Convert JSON value to the string if not NIL
 */
func toString(val interface{}) string {
	if val != nil {
		return val.(string)
	}
	return ""
}
