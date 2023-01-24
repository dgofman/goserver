package view

import (
	"fmt"
	"net/http"

	"go/goserver.io/constants"
	"go/goserver.io/lib"
)

/**
 * AJAX handler for request: /device_info
 */
func DeviceInfo(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, fmt.Sprintf("%s?%s", constants.ApiDeviceInfo, r.URL.RawQuery))
}

/**
 * AJAX handler for request: /ajax_ipdata
 */
func IpData(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, fmt.Sprintf("%s?%s", constants.ApiIpData, r.URL.RawQuery))
}

/**
 * AJAX handler for request: /ajax_metrics
 */
func Metrics(w http.ResponseWriter, r *http.Request, app *lib.App) {
	lib.ApiRequest(w, r, fmt.Sprintf("%s?%s", constants.ApiMetrics, r.URL.RawQuery))
}
