package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
	"go/goserver.io/utils"
	"go/goserver.io/view"
)

/**
 * Main function
 */
func main() {
	var Props = constants.Props
	log.Printf("Precision Time Server Server: Console: %v, LogFile: %v, View: %v\n", Props.CONSOLE, Props.LOGFILE, Props.VIEW)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if Props.LOGFILE {
		f, err := os.OpenFile("goserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			utils.LogError(err)
		}
		defer f.Close()

		if Props.CONSOLE {
			log.SetOutput(io.MultiWriter(os.Stdout, f))
		} else {
			log.SetOutput(f)
		}
	}

	app := lib.CreateApp(Props.BASE_DIR)
	Router(app)
	log.Printf("ENV: %s\n", Props.ENV)
	log.Printf("API_BASEPATH: %s\n", constants.API_BASEPATH)

	var err error
	time.AfterFunc(time.Second, func() {
		if err == nil {
			log.Printf("Starting server port: %d\n", Props.PORT)
		}
	})

	//Catch Interrupt signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for sig := range sigc {
			fmt.Println("Signal: ", sig)
		}
	}()
	err = http.ListenAndServe(fmt.Sprintf(":%d", Props.PORT), nil)
	utils.LogError(err)
}

/**
 * Register URL routes and function handler
 */
func Router(app *lib.App) {
	app.UrlHandle(constants.URL_Static, httpStatic(app))
	app.UrlHandle(constants.URL_Login, view.Login)
	app.UrlHandle(constants.URL_Logout, view.Logout)
	app.UrlHandle(constants.URL_CreateAccount, view.CreateAccount)
	app.UrlHandle(constants.URL_LinkedInOauth, view.CreateAccountOauth)
	app.UrlHandle(constants.URL_VerifyEmail, view.CreateAccountVerifyEmail)
	app.UrlHandle(constants.URL_DeleteAccount, view.DeleteAccount)
	app.UrlHandle(constants.URL_ActivateAccount, view.ActivateAccount)
	app.UrlHandle(constants.URL_UserProfile, view.UserProfile)
	app.UrlHandle(constants.URL_LoadURL, view.LoadImage)
	app.UrlHandle(constants.URL_ForgotUsername, view.ForgotUsername)
	app.UrlHandle(constants.URL_ResetPassword, view.ResetPassword)
	app.UrlHandle(constants.URL_GetStarted, view.GetStarted)
	app.UrlHandle(constants.URL_AddClient, view.GetStarted)
	app.UrlHandle(constants.URL_ManageDevices, view.ManageDevices)
	app.UrlHandle(constants.URL_ManageLicenses, view.ManageLicenses)
	app.UrlHandle(constants.URL_Map, view.Map)
	app.UrlHandle(constants.URL_Install, view.Installation)
	app.UrlHandle(constants.URL_DeviceInfo, view.DeviceInfo)
	app.UrlHandle(constants.URL_AjaxIpdata, view.IpData)
	app.UrlHandle(constants.URL_AjaxMetrics, view.Metrics)
	app.UrlHandle(constants.URL_AjaxStatus, view.Status)
}

/**
 * Static recourses request handler
 */
func httpStatic(app *lib.App) func(http.ResponseWriter, *http.Request, *lib.App) {
	fs := http.FileServer(http.Dir(utils.Abs(constants.Props.STATIC_DIR)))
	h := http.StripPrefix(constants.URL_Static, fs)
	return func(w http.ResponseWriter, r *http.Request, app *lib.App) {
		h.ServeHTTP(w, r)
	}
}
