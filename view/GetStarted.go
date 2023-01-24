package view

import (
	"log"
	"net/http"

	"go/goserver.io/constants"
	"go/goserver.io/lib"

	"github.com/dgofman/pongo2"
)

/**
 * Template handler for request: /
 */
func GetStarted(w http.ResponseWriter, r *http.Request, app *lib.App) {
	log.Printf("view:GetStarted: " + r.Method)

	resp, err := lib.ApiRequest(nil, r, constants.ApiStatus)
	if app.ErrorHandler(w, err, resp) {
		return
	}
	ctx := lib.CreateContent(r)
	ctx["_url_install"] = constants.URL_Install
	resp.Decode(&ctx)
	status := ctx["status"].(string)
	if r.Method == "POST" {
		if status == constants.STATUS_NEW {
			resp, err := lib.ApiRequest(nil, r, constants.ProxyInitial)
			if app.ErrorHandler(w, err, nil) {
				return
			}
			if resp.StatusCode == 400 && resp.Error != nil {
				ctx["error"] = resp.Error.Message
				lib.BuildTemplate(w, "getstarted_wrap.html", ctx)
				return
			}
		} else {
			resp, err := lib.ApiRequest(nil, r, constants.ApiCreateUser)
			if app.ErrorHandler(w, err, nil) {
				return
			}
			log.Printf("CREATE Status: %d - %s\n", resp.StatusCode, resp.Body)
			if resp.StatusCode == 400 && resp.Error != nil {
				ctx["error"] = resp.Error.Message
				lib.BuildTemplate(w, "getstarted_wrap.html", ctx)
				return
			}
		}
		lib.BuildTemplate(w, "setupService.html", ctx)
	} else {
		if status == constants.STATUS_NO_USER || status == constants.STATUS_NEW {
			lib.BuildTemplate(w, "getstarted_wrap.html", ctx)
		} else if status == constants.STATUS_FOREIGN_USER {
			app.Redirect(w, r, constants.URL_Login)
		} else if status == constants.STATUS_INITIAL || r.URL.Path == constants.URL_AddClient {
			lib.BuildTemplate(w, "setupService.html", ctx)
		} else {
			statuses := map[string]int{
				"present_linux_ntp":  toInt(ctx, "present_linux_ntp"),
				"present_linux_ptp":  toInt(ctx, "present_linux_ptp"),
				"presentwindows_ntp": toInt(ctx, "presentwindows_ntp"),
			}
			ctx["status"] = statuses
			if statuses["present_linux_ntp"] != 0 {
				ctx["os_protocol_label"] = "Linux NTP"
				ctx["ntp_linux_checked"] = "checked"
			} else if statuses["present_linux_ptp"] != 0 {
				ctx["os_protocol_label"] = "Linux PTP"
				ctx["ptp_linux_checked"] = "checked"
			} else {
				ctx["os_protocol_label"] = "Windows NTP"
				ctx["ntp_windows_checked"] = "checked"
			}
			lib.ExecTemplate(w, lib.Template("dashboard.html"), ctx)
		}
	}
}

func toInt(ctx pongo2.Context, key string) int {
	i := ctx[key]
	if i != nil {
		return int(i.(float64))
	}
	return 0
}
